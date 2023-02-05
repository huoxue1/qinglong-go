package service

import (
	"context"
	"github.com/huoxue1/qinglong-go/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func AppInit() {
	go runYarn()
}

func runYarn() {
	defer func() {
		recover()
	}()
	_, err := os.Stat(path.Join("data", "scripts", "package.json"))
	if os.IsNotExist(err) {
		return
	}
	ch := make(chan int, 1)
	utils.RunTask(context.WithValue(context.Background(), "cancel", ch), "pnpm install", map[string]string{}, func(ctx context.Context) {
		log.Infoln("开始执行pnpm初始化！")
	}, func(ctx context.Context) {
		log.Infoln("pnpm初始化执行完成！")
	}, os.Stdout)
}
