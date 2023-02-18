package main

import (
	"flag"
	nested "github.com/Lyrics-you/sail-logrus-formatter/sailor"
	"github.com/huoxue1/qinglong-go/controller"
	"github.com/huoxue1/qinglong-go/service"
	"github.com/huoxue1/qinglong-go/service/config"
	env_check "github.com/huoxue1/qinglong-go/utils/env-check"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

var (
	address string
)

func init() {
	w, err := rotates.New(path.Join("data", "log", "qinglong-go", "%Y-%m-%d.log"), rotates.WithRotationTime(time.Hour*24))
	if err != nil {
		log.Errorf("rotates init err: %v", err)
		panic(err)
	}
	log.SetOutput(io.MultiWriter(w, os.Stdout))
	log.SetFormatter(&nested.Formatter{
		FieldsOrder:           nil,
		TimeStampFormat:       "2006-01-02 15:04:05",
		CharStampFormat:       "",
		HideKeys:              false,
		Position:              true,
		Colors:                true,
		FieldsColors:          true,
		FieldsSpace:           true,
		ShowFullLevel:         false,
		LowerCaseLevel:        true,
		TrimMessages:          true,
		CallerFirst:           false,
		CustomCallerFormatter: nil,
	})
	flag.StringVar(&address, "add", "0.0.0.0:5700", "the ql listen address!")
	flag.Parse()
	config.SetAddress(address)
}

func main() {
	env_check.CheckStatic()
	service.AppInit()
	engine := controller.Router()
	_ = engine.Run(address)
}
