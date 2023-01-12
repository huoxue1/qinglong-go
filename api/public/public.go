package public

import (
	"github.com/gin-gonic/gin"
	"path"
	"time"
)

func Api(group *gin.RouterGroup) {
	group.GET("/panel/log", log())
}
func log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.File(path.Join("data", "log", "qinglong-go", time.Now().Format("2006-01-02")+".log"))
	}
}
