package cron

import (
	"github.com/huoxue1/qinglong-go/utils"
	"os"

	log "github.com/huoxue1/go-utils/base/log"
	"github.com/panjf2000/ants/v2"
)

func init() {
	initPool()
}

var (
	pool *ants.PoolWithFunc
)

func run(task *utils.RunOption) error {
	return pool.Invoke(task)
}

func initPool() {

	pool1, err := ants.NewPoolWithFunc(5, func(i2 interface{}) {
		option := i2.(*utils.RunOption)
		utils.RunWithOption(option.Ctx, option)
	})
	if err != nil {
		log.Errorln("创建定时任务协程池失败" + err.Error())
		os.Exit(2)
	}
	pool = pool1
}
