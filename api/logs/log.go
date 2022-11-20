package logs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/service/scripts"
	"github.com/huoxue1/qinglong-go/utils/res"
	"os"
	"path"
)

func APi(group *gin.RouterGroup) {
	group.GET("", get())
	group.GET("/:name", getFile())
}

func get() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		files := scripts.GetFiles(path.Join("data", "log"), "")
		ctx.JSON(200, res.Ok(files))
	}
}

func getFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileName := ctx.Param("name")
		path := ctx.Query("path")
		data, _ := os.ReadFile(fmt.Sprintf("data/log/%s/%s", path, fileName))
		ctx.JSON(200, res.Ok(string(data)))
	}
}
