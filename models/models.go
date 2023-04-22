package models

import (
	"github.com/huoxue1/go-utils/base/log"
	xLog "github.com/huoxue1/go-utils/base/log/xorm"
	_ "modernc.org/sqlite"
	"os"
	"xorm.io/xorm"
)

var (
	engine *xorm.Engine
)

func InitModels() {
	_ = os.MkdirAll("data/db", 0666)
	en, err := xorm.NewEngine("sqlite", "data/db/database.sqlite")
	if err != nil {
		log.Errorln("[sql] " + err.Error())
		return
	}
	_ = en.Sync2(new(Apps), new(Auths), new(Crontabs), new(Crontabviews), new(Dependences), new(Envs), new(Subscriptions))
	en.SetLogger(xLog.GetXormLogger(log.StandardLogger(), "info", false))
	engine = en
}
