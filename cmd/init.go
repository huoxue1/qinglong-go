package cmd

import (
	"github.com/huoxue1/go-utils/base/log"
	"github.com/huoxue1/go-utils/base/log/hook/file"
	"github.com/huoxue1/go-utils/base/log/hook/std"
	fileInit "github.com/huoxue1/qinglong-go/internal/init"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/config"
	"github.com/huoxue1/qinglong-go/service/cron"
	"github.com/huoxue1/qinglong-go/service/subscription"
	"github.com/huoxue1/qinglong-go/utils"
	"path"
)

func InitLog() {
	if !utils.FileExist(path.Join("data", "config", "config.sh")) {
		fileInit.InitConfig()
	}
	level := config.GetKey("logLevel", "info")
	fileHook, _ := file.NewFileHook(file.WithLevel(level), file.WithDir(path.Join("data", "log")))
	log.AddHook(std.NewStdHook(std.WithLevel(level)))
	log.AddHook(fileHook)
}

func initCron() {
	cron.InitTask()
	subscription.InitSub()
	models.SetAllCornStop()
}
