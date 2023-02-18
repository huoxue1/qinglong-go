package env_check

import (
	"fmt"
	"github.com/dablelv/go-huge-util/zip"
	"github.com/huoxue1/qinglong-go/service/config"
	"github.com/huoxue1/qinglong-go/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func CheckStatic() {
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
