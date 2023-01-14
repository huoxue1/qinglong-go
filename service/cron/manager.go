package cron

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/config"
	"github.com/huoxue1/qinglong-go/service/env"
	"github.com/huoxue1/qinglong-go/utils"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	manager     sync.Map
	execManager sync.Map
)

func init() {
	initTask()
	models.SetAllCornStop()
}

type task struct {
	cmd   string
	isNow bool
	isCon bool
	envs  map[string][]int
	dir   string
}

func initTask() {
	enableCrons := models.FindAllEnableCron()
	for _, c := range enableCrons {
		AddTask(c)
	}
}

func stopCron(crontabs *models.Crontabs) {
	defer func() {
		_ = recover()
	}()
	value, ok := execManager.Load(crontabs.Id)
	if !ok {
		return
	}
	cancel := value.(func())
	cancel()
}

func runCron(crontabs *models.Crontabs) {
	envFromDb := env.LoadEnvFromDb()
	envfromFile := env.LoadEnvFromFile("data/config/config.sh")
	for s, s2 := range envfromFile {
		if _, ok := envFromDb[s]; ok {
			envFromDb[s] = envFromDb[s] + "&" + s2
		} else {
			envFromDb[s] = s2
		}
	}
	ta := handCommand(crontabs.Command)

	if ta.isCon {
		for e, indexs := range ta.envs {
			value, ok := envFromDb[e]
			if !ok {
				break
			}
			values := strings.Split(value, "&")
			for _, index := range indexs {
				logFile, _ := os.OpenFile("data/log/"+time.Now().Format("2006-01-02")+"/"+crontabs.Name+"_"+uuid.New().String()+".log", os.O_RDWR|os.O_CREATE, 0666)
				e2 := exec.Command(strings.Split(ta.cmd, " ")[0], strings.Split(ta.cmd, " ")[1:]...)
				e2.Env = []string{e + "=" + values[index]}
				stdoutPipe, _ := e2.StdoutPipe()

				for s, s2 := range envFromDb {
					if s != e {
						e2.Env = append(e2.Env, s+"="+s2)
					}
				}
				go func() {
					err := e2.Start()
					if err != nil {
						log.Errorln(err.Error())
					}
					go io.Copy(logFile, stdoutPipe)
					e2.Wait()
				}()

			}
		}
	} else {
		if _, ok := execManager.Load(crontabs.Id); ok {
			log.Warningln(fmt.Sprintf("the task %s is running,skip the task", crontabs.Name))
			return
		}
		crontabs.LogPath = "data/log/" + time.Now().Format("2006-01-02") + "/" + crontabs.Name + "_" + uuid.New().String() + ".log"
		crontabs.Status = 0
		models.UpdateCron(crontabs)
		cancelChan := make(chan int, 1)
		ctx := context.WithValue(context.Background(), "cancel", cancelChan)
		execManager.Store(crontabs.Id, func() {
			cancelChan <- 1
		})
		now := time.Now()
		_ = os.Mkdir("data/log/"+time.Now().Format("2006-01-02"), 0666)
		file, _ := os.OpenFile(crontabs.LogPath, os.O_RDWR|os.O_CREATE, 0666)
		cmdDir := "./data/scripts/"
		if strings.HasPrefix(ta.cmd, "go") {
			cmdDir = ta.dir
		}
		go utils.RunWithOption(ctx, &utils.RunOption{
			Command: ta.cmd,
			Env:     envFromDb,
			OnStart: func(ctx context.Context) {
				writer := ctx.Value("log").(io.Writer)
				writer.Write([]byte(fmt.Sprintf("##开始执行..  %s\n\n", now.Format("2006-01-02 15:04:05"))))
			},
			OnEnd: func(ctx context.Context) {
				writer := ctx.Value("log").(io.Writer)
				writer.Write([]byte(fmt.Sprintf("\n##执行结束..  %s，耗时%.1f秒\n\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().Sub(now).Seconds())))
				crontabs.Status = 1
				crontabs.LastExecutionTime = now.Unix()
				crontabs.LastRunningTime = int64(time.Now().Sub(now).Seconds())
				models.UpdateCron(crontabs)
				execManager.LoadAndDelete(crontabs.Id)
				file.Close()
			},
			LogFile: file,
			CmdDir:  cmdDir,
		})
	}
}

func AddTask(crontabs *models.Crontabs) {
	crons := strings.Split(crontabs.Schedule, " ")
	var c *cron.Cron
	if len(crons) == 5 {
		c = cron.New()

	} else if len(crons) == 6 {
		c = cron.New(cron.WithParser(
			cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
	} else {
		log.Errorf("the task %s cron %s is error", crontabs.Name, crontabs.Command)
		return
	}
	_, err := c.AddFunc(crontabs.Schedule, func() {
		runCron(crontabs)
	})
	if err != nil {
		log.Errorln("添加task错误" + crontabs.Schedule + err.Error())
		return
	}
	c.Start()
	manager.Store(crontabs.Id, c)

}

func handCommand(command string) *task {
	ta := new(task)
	commands := strings.Split(command, " ")

	pythonCmd := config.GetKey("PythonCmd", "python")
	JsCmd := config.GetKey("JsCmd", "node")
	ShCmd := config.GetKey("ShCmd", "bash")

	if commands[0] == "task" {
		if strings.HasSuffix(commands[1], ".py") {
			ta.cmd = pythonCmd + " " + commands[1]
		} else if strings.HasSuffix(commands[1], ".js") {
			ta.cmd = JsCmd + " " + commands[1]
		} else if strings.HasSuffix(commands[1], ".sh") {
			ta.cmd = ShCmd + " " + commands[1]
		} else if strings.HasSuffix(commands[1], ".ts") {
			ta.cmd = "ts-node-transpile-only " + commands[1]
		} else if strings.HasSuffix(commands[1], ".go") {

			log.Infoln(filepath.Base(commands[1]))
			ta.cmd = fmt.Sprintf(`go run %s`, filepath.Base(commands[1]))
			ta.dir = path.Join("./data", "scripts", filepath.Dir(commands[1]))
		}
		if len(commands) > 2 {
			if commands[2] == "now" {
				ta.isNow = true
			} else if commands[2] == "desi" {
				if len(commands) >= 3 {
					var envIndex []int
					for _, i := range commands[4:] {
						index, _ := strconv.Atoi(i)
						envIndex = append(envIndex, index)
					}
					ta.envs = map[string][]int{
						commands[3]: envIndex,
					}
				}

			} else if commands[2] == "conc" {
				ta.isCon = true
				if len(commands) >= 3 {
					var envIndex []int
					for _, i := range commands[4:] {
						index, _ := strconv.Atoi(i)
						envIndex = append(envIndex, index)
					}
					ta.envs = map[string][]int{
						commands[3]: envIndex,
					}
				}
			}
		}

	} else {
		ta.cmd = command
	}
	return ta
}

func getGoModule(filePath string) string {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorln("not get the go module name")
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(strings.Split(string(file), "\n")[0], "module"))
}
