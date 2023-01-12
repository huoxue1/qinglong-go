package main

import (
	nested "github.com/Lyrics-you/sail-logrus-formatter/sailor"
	"github.com/huoxue1/qinglong-go/controller"
	"github.com/huoxue1/qinglong-go/service"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
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
}

func main() {
	service.AppInit()
	engine := controller.Router()
	_ = engine.Run(":5700")
}
