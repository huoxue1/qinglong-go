package subscription

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/config"
	"github.com/huoxue1/qinglong-go/service/cron"
	"github.com/huoxue1/qinglong-go/utils"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	manager sync.Map
)

func stopSubscription(sub *models.Subscriptions) {
	defer func() {
		_ = recover()
	}()
	sub.Status = 1
	_ = models.UpdateSubscription(sub)
	value, ok := manager.Load(sub.Id)
	if !ok {
		return
	}
	cancel := value.(func())
	cancel()
}

func downloadFiles(subscriptions *models.Subscriptions) {
	if subscriptions.Type == "public-repo" {
		os.RemoveAll(path.Join("data", "scripts", subscriptions.Alias))
		os.RemoveAll(path.Join("data", "repo", subscriptions.Alias))
		err := downloadPublicRepo(subscriptions)
		if err != nil {
			return
		}
		addScripts(subscriptions)
		file, _ := os.OpenFile(subscriptions.LogPath, os.O_APPEND|os.O_RDWR, 0666)
		file.WriteString(fmt.Sprintf("\n##执行结束..  %s，耗时0秒\n\n", time.Now().Format("2006-01-02 15:04:05")))
		_ = file.Close()
		subscriptions.Status = 1
		models.UpdateSubscription(subscriptions)
	}
}

func downloadPublicRepo(subscriptions *models.Subscriptions) error {
	subscriptions.LogPath = "data/log/" + time.Now().Format("2006-01-02") + "/" + subscriptions.Alias + "_" + uuid.New().String() + ".log"
	_ = os.MkdirAll(filepath.Dir(subscriptions.LogPath), 0666)
	cmd := fmt.Sprintf("clone -b %s --single-branch %s %s", subscriptions.Branch, subscriptions.Url, path.Join("data", "repo", subscriptions.Alias))
	command := exec.Command("git", strings.Split(cmd, " ")...)
	pipe, err := command.StdoutPipe()
	stderrPipe, _ := command.StderrPipe()
	if err != nil {
		return err
	}
	subscriptions.Status = 0
	err = models.UpdateSubscription(subscriptions)
	if err != nil {
		return err
	}
	file, _ := os.OpenFile(subscriptions.LogPath, os.O_CREATE|os.O_RDWR, 0666)
	file.Write([]byte(fmt.Sprintf("##开始执行..  %s\n\n", time.Now().Format("2006-01-02 15:04:05"))))
	err = command.Start()
	if err != nil {
		return err
	}
	manager.Store(subscriptions.Id, func() {
		command.Process.Kill()
	})
	defer manager.LoadAndDelete(subscriptions.Id)
	go io.Copy(io.MultiWriter(file, os.Stdout), pipe)
	go io.Copy(file, stderrPipe)
	command.Wait()
	return err
}

func addScripts(subscriptions *models.Subscriptions) {
	file, _ := os.OpenFile(subscriptions.LogPath, os.O_CREATE|os.O_RDWR, 0666)
	defer file.Close()
	var extensions []string
	if subscriptions.Extensions != "" {
		extensions = strings.Split(subscriptions.Extensions, " ")
	} else {
		extensions = strings.Split(config.GetKey("RepoFileExtensions"), " ")
	}
	dir, err := os.ReadDir(path.Join("data", "repo", subscriptions.Alias))
	if err != nil {
		return
	}
	for _, entry := range dir {
		// 判断文件后缀
		if !utils.In(strings.TrimPrefix(filepath.Ext(entry.Name()), "."), extensions) {
			if !entry.IsDir() {
				continue
			}
		}
		// 判断黑名单
		if utils.In(entry.Name(), strings.Split(subscriptions.Blacklist, "|")) {
			continue
		}
		compile := regexp.MustCompile(`(` + subscriptions.Whitelist + `)`)
		if compile.MatchString(entry.Name()) {
			name, c, _ := getSubCron(path.Join("data", "repo", subscriptions.Alias, entry.Name()))
			if c != "" {
				command, err := models.GetCronByCommand(fmt.Sprintf("task %s", path.Join(subscriptions.Alias, entry.Name())))
				if err != nil {
					file.WriteString("已添加新的定时任务  " + name + "\n")
					_, _ = cron.AddCron(&models.Crontabs{
						Name:      name,
						Command:   fmt.Sprintf("task %s", path.Join(subscriptions.Alias, entry.Name())),
						Schedule:  c,
						Timestamp: time.Now().Format("Mon Jan 02 2006 15:04:05 MST"),
						Status:    1,
						Labels:    []string{},
					})
				} else {
					command.Name = name
					command.Schedule = c
					_ = cron.UpdateCron(command)
				}

			}
			utils.Copy(path.Join("data", "repo", subscriptions.Alias, entry.Name()), path.Join("data", "scripts", subscriptions.Alias, entry.Name()))
		} else {
			depen := regexp.MustCompile(`(` + subscriptions.Dependences + `)`)
			if depen.MatchString(entry.Name()) {
				utils.Copy(path.Join("data", "repo", subscriptions.Alias, entry.Name()), path.Join("data", "scripts", subscriptions.Alias, entry.Name()))
			}
		}
	}
}

func getSubCron(filePath string) (name string, cron string, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}
	cronReg := regexp.MustCompile(`([0-9\-*/,]{1,} ){4,5}([0-9\-*/,]){1,}`)
	nameEnv := regexp.MustCompile(`new\sEnv\(['|"](.*?)['|"]\)`)
	if cronReg.Match(data) {
		cron = string(cronReg.FindAll(data, 1)[0])
		cron = strings.TrimPrefix(cron, "//")
		if nameEnv.Match(data) {
			name = string(nameEnv.FindAllSubmatch(data, 1)[0][1])
		} else {
			name = path.Base(filePath)
		}
	} else {
		return "", "", errors.New("not found cron")
	}
	return
}
