package cron

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/env"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
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
	cmd := value.(*exec.Cmd)
	_ = cmd.Process.Kill()
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
				logFile, _ := os.OpenFile("data/log/"+time.Now().Format("2006-01-02")+"/"+crontabs.Name+".log", os.O_RDWR, 0666)
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
					go syncLog(stdoutPipe, logFile)
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
		_ = os.Mkdir("data/log/"+time.Now().Format("2006-01-02"), 0666)
		logFile := &myWriter{crontabs.LogPath}

		e2 := exec.Command(strings.Split(ta.cmd, " ")[0], strings.Split(ta.cmd, " ")[1:]...)
		execManager.Store(crontabs.Id, e2)
		stdoutPipe, _ := e2.StdoutPipe()
		for s, s2 := range envFromDb {
			e2.Env = append(e2.Env, s+"="+s2)
		}
		startTime := time.Now()
		logFile.Write([]byte(fmt.Sprintf("##开始执行..  %s\n\n", startTime.Format("2006-01-02 15:04:05"))))
		go func() {
			err := e2.Start()
			if err != nil {
				log.Errorln(err.Error())
			}
			go syncLog(stdoutPipe, logFile)
			e2.Wait()
			logFile.Write([]byte(fmt.Sprintf("##执行结束..  %s，耗时%.1f秒\n\n", time.Now().Format("2006-01-02 15:04:05"), time.Now().Sub(startTime).Seconds())))
			crontabs.Status = 1
			models.UpdateCron(crontabs)
			execManager.LoadAndDelete(crontabs.Id)
		}()
	}
}

func AddTask(crontabs *models.Crontabs) {
	c := cron.New()
	_, err := c.AddFunc(crontabs.Schedule, func() {
		runCron(crontabs)
	})
	if err != nil {
		log.Errorln("添加task错误" + err.Error())
		return
	}
	c.Start()
	manager.Store(crontabs.Id, c)
}

func handCommand(command string) *task {
	ta := new(task)
	commands := strings.Split(command, " ")
	if commands[0] == "task" {
		if strings.HasSuffix(commands[1], ".py") {
			ta.cmd = "python data/scripts/" + commands[1]
		} else if strings.HasSuffix(commands[1], ".js") {
			ta.cmd = "node data/scripts/" + commands[1]
		} else if strings.HasSuffix(commands[1], ".sh") {
			ta.cmd = "bash data/scripts/" + commands[1]
		} else if strings.HasSuffix(commands[1], ".ts") {
			ta.cmd = "ts-node-transpile-only data/scripts/" + commands[1]
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

type myWriter struct {
	fileName string
}

func (m *myWriter) Write(p []byte) (n int, err error) {
	file, _ := os.OpenFile(m.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	n, err = file.Write(p)
	file.Close()
	return n, err
}

//通过管道同步获取日志的函数
func syncLog(reader io.ReadCloser, writer io.Writer) {
	buf := make([]byte, 1)
	for {
		strNum, err := reader.Read(buf)
		if strNum > 0 {
			outputByte := buf[:strNum]
			writer.Write(outputByte)
		}
		if err != nil {
			//读到结尾
			if err == io.EOF || strings.Contains(err.Error(), "file already closed") {
				return
			}
		}
	}
}
