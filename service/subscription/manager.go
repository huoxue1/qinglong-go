package subscription

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/huoxue1/go-utils/base/log"
	"github.com/huoxue1/qinglong-go/internal/cron-manager"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/config"
	"github.com/huoxue1/qinglong-go/service/cron"
	"github.com/huoxue1/qinglong-go/utils"
)

var (
	manager sync.Map
)

func InitSub() {
	log.Infoln("开始初始化订阅任务定时！")
	subscriptions, err := models.QuerySubscription("")
	if err != nil {
		return
	}
	for _, subscription := range subscriptions {
		cron_manager.AddCron(fmt.Sprintf("sub_%d", subscription.Id), subscription.GetCron(), func() {
			downloadFiles(subscription)
		})
	}
}

func getDepFiles() []string {
	var files []string
	dir, err := os.ReadDir(path.Join("data", "deps"))
	if err != nil {
		return []string{}
	}
	for _, entry := range dir {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}

	}
	return files
}

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

		os.RemoveAll(path.Join("data", "repo", subscriptions.Alias))
		err := downloadPublicRepo(subscriptions)
		if err != nil {
			return
		}
		os.RemoveAll(path.Join("data", "scripts", subscriptions.Alias))
		if config.GetKey("AutoAddCron", "true") == "true" {
			addScripts(subscriptions)
		} else {
			log.Infoln("未配置自动添加定时任务，不添加任务！")
		}

		file, _ := os.OpenFile(subscriptions.LogPath, os.O_APPEND|os.O_RDWR, 0666)
		file.WriteString(fmt.Sprintf("\n##执行结束..  %s，耗时0秒\n\n", time.Now().Format("2006-01-02 15:04:05")))
		_ = file.Close()
		subscriptions.Status = 1
		models.UpdateSubscription(subscriptions)
		manager.LoadAndDelete(subscriptions.Id)
	} else if subscriptions.Type == "file" {
		addRawFiles(subscriptions)
	}
}

func addRawFiles(subscriptions *models.Subscriptions) {
	_ = models.UpdateSubscription(subscriptions)
	defer func() {
		subscriptions.Status = 1
		_ = models.UpdateSubscription(subscriptions)
	}()
	err := utils.DownloadFile(subscriptions.Url, path.Join("data", "raw", subscriptions.Alias))
	if err != nil {
		_, _ = subscriptions.WriteString(err.Error() + "\n")
		return
	}
	name, c, err := getSubCron(path.Join("data", "raw", subscriptions.Alias))
	if err != nil {
		_, _ = subscriptions.WriteString(err.Error() + "\n")
		return
	}
	utils.Copy(path.Join("data", "raw", subscriptions.Alias), path.Join("data", "scripts", subscriptions.Alias))
	if c != "" {
		command, err := models.GetCronByCommand(fmt.Sprintf("task %s", subscriptions.Alias))
		if err != nil {
			subscriptions.WriteString("已添加新的定时任务  " + name + "\n")
			_, _ = cron.AddCron(&models.Crontabs{
				Name:      name,
				Command:   fmt.Sprintf("task %s", subscriptions.Alias),
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
}

func downloadPublicRepo(subscriptions *models.Subscriptions) error {

	_, _ = subscriptions.Write([]byte(fmt.Sprintf("##开始执行..  %s\n\n", time.Now().Format("2006-01-02 15:04:05"))))

	_, err := git.PlainClone(path.Join("data", "repo", subscriptions.Alias), false, &git.CloneOptions{
		URL:           subscriptions.Url,
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", subscriptions.Branch)),
		Progress:      subscriptions,
		ProxyOptions:  transport.ProxyOptions{URL: config.GetKey("ProxyUrl", "")},
	})
	if err != nil {
		_, _ = subscriptions.Write([]byte(fmt.Sprintf("err:  %s", err.Error())))
		return err
	}

	return err
}

func addScripts(subscriptions *models.Subscriptions) {
	depFiles := getDepFiles()
	var extensions []string
	if subscriptions.Extensions != "" {
		extensions = strings.Split(subscriptions.Extensions, " ")
	} else {
		extensions = strings.Split(config.GetKey("RepoFileExtensions", "js py sh"), " ")
	}
	dir, err := os.ReadDir(path.Join("data", "repo", subscriptions.Alias))
	if err != nil {
		return
	}
	crontabs, _ := models.QueryCronByDir(subscriptions.Alias)
	cronMap := make(map[string]*models.Crontabs, len(crontabs))
	for _, crontab := range crontabs {
		cronMap[crontab.Command] = crontab
	}
	isGoMod := false
	for _, entry := range dir {

		if entry.Name() == "go.mod" {
			isGoMod = true
		}
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
					_, err1 := cron.AddCron(&models.Crontabs{
						Name:      name,
						Command:   fmt.Sprintf("task %s", path.Join(subscriptions.Alias, entry.Name())),
						Schedule:  c,
						Timestamp: time.Now().Format("Mon Jan 02 2006 15:04:05 MST"),
						Status:    1,
						Labels:    []string{},
					})
					if err1 != nil {
						_, _ = subscriptions.WriteString("定时任务添加失败： " + name + " " + err1.Error())

					} else {
						_, _ = subscriptions.WriteString("已添加新的定时任务  " + name + "\n")
					}
				} else {
					command.Name = name
					command.Schedule = c
					_ = cron.UpdateCron(command)
					delete(cronMap, command.Command)
				}

			}
			utils.Copy(path.Join("data", "repo", subscriptions.Alias, entry.Name()), path.Join("data", "scripts", subscriptions.Alias, entry.Name()))
		} else {
			depen := regexp.MustCompile(`(` + subscriptions.Dependences + `)`)
			if depen.MatchString(entry.Name()) {
				utils.Copy(path.Join("data", "repo", subscriptions.Alias, entry.Name()), path.Join("data", "scripts", subscriptions.Alias, entry.Name()))
			}
		}
		if utils.In(entry.Name(), depFiles) {
			subscriptions.WriteString("已替换依赖文件： " + entry.Name() + "\n")
			utils.Copy(path.Join("data", "deps", entry.Name()), path.Join("data", "scripts", subscriptions.Alias, entry.Name()))
		}
	}
	if config.GetKey("AutoDelCron", "true") == "true" {
		for _, m := range cronMap {
			subscriptions.WriteString("已删除失效的任务 " + m.Name + "\n")
			models.DeleteCron(m.Id)
		}
	}
	if isGoMod {
		subscriptions.WriteString("检测到go模块，开始自动下载golang依赖!!")
		cancelChan := make(chan int, 1)
		ctx := context.WithValue(context.Background(), "cancel", cancelChan)
		utils.RunWithOption(ctx, &utils.RunOption{
			Command: "go mod tidy",
			Env:     map[string]string{},
			OnStart: func(ctx context.Context) {

			},
			OnEnd: func(ctx context.Context) {

			},
			LogFile: subscriptions,
			CmdDir:  path.Join("data", "scripts", subscriptions.Alias),
		})
	}
	for _, depFile := range depFiles {
		utils.Copy(path.Join("data", "deps", depFile), path.Join("data", "scripts", subscriptions.Alias, depFile))

	}
}

func getSubCron(filePath string) (name string, cron string, err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}
	cronReg := regexp.MustCompile(`([0-9\-*/,]{1,} ){4,5}([0-9\-*/,]){1,}`)
	nameEnv := regexp.MustCompile(`new\sEnv\(['|"](.*?)['|"]\)`)
	if nameEnv.Match(data) {
		if cronReg.Match(data) {
			cron = strings.TrimPrefix(strings.TrimPrefix(string(cronReg.FindAll(data, 1)[0]), "//"), " ")
		} else {
			key := config.GetKey("DefaultCronRule", "0 9 * * *")
			if key == "" {
				key = "0 9 * * *"
			}
			cron = key
		}
		name = string(nameEnv.FindAllSubmatch(data, 1)[0][1])
		return
	} else {
		return "", "", errors.New("not found cron")
	}
}
