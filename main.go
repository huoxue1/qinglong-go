package main

import (
	"fmt"
	nested "github.com/Lyrics-you/sail-logrus-formatter/sailor"
	"github.com/dablelv/go-huge-util/zip"
	"github.com/huoxue1/qinglong-go/controller"
	"github.com/huoxue1/qinglong-go/service"
	"github.com/huoxue1/qinglong-go/service/config"
	"github.com/huoxue1/qinglong-go/utils"
	rotates "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"strings"
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
	checkStatic()
	service.AppInit()
	engine := controller.Router()
	_ = engine.Run(":5700")
}

func checkStatic() {
	if !utils.FileExist("./static/") {
		version := config.GetVersion()
		if !strings.HasPrefix(version, "v") {
			version = "v1.0.0"
		}
		log.Warningln("检测到静态文件资源不存在，即将自动下载文件！")
		log.Infoln("downloading file from ", fmt.Sprintf("https://github.com/huoxue1/qinglong/releases/download/%s/static.zip", version))
		response, err := utils.GetClient().R().Get(fmt.Sprintf("https://github.com/huoxue1/qinglong/releases/download/%s/static.zip", version))
		if err != nil {
			log.Errorln("下载静态资源文件失败 " + err.Error())
			return
		}
		err = os.WriteFile("static.zip", response.Bytes(), 0666)
		if err != nil {
			log.Errorln("写入压缩文件错误 " + err.Error())
			return
		}
		err = zip.Unzip("static.zip", ".")
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
}
