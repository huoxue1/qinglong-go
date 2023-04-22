package system

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/internal/auth"
	"github.com/huoxue1/qinglong-go/internal/res"
	"github.com/huoxue1/qinglong-go/service/config"
	"github.com/huoxue1/qinglong-go/service/system"
)

func Api(group *gin.RouterGroup) {
	group.GET("", get())
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.JSON(200, res.Ok(system.System{
			IsInitialized:  auth.IsInit(),
			Version:        config.GetVersion(),
			LastCommitTime: "",
			LastCommitId:   "",
			Branch:         "qinglong-go",
		}))
	}
}
