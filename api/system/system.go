package system

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/service/system"
	"github.com/huoxue1/qinglong-go/utils/res"
	"os"
	"path"
)

func Api(group *gin.RouterGroup) {
	group.GET("/", get())
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := os.Stat(path.Join("data", "config", "auth.json"))
		exist := os.IsExist(err)
		ctx.JSON(200, res.Ok(system.System{
			IsInitialized:  exist,
			Version:        "2.0.14",
			LastCommitTime: "",
			LastCommitId:   "",
			Branch:         "master",
		}))
	}
}
