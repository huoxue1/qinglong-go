package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/internal/res"
	"github.com/huoxue1/qinglong-go/service/notification"
	"os"
	"path"
)

func getNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := os.ReadFile(path.Join("data", "config", "push.json"))
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		m := make(map[string]interface{}, 5)
		_ = json.Unmarshal(data, &m)
		ctx.JSON(200, res.Ok(m))
	}
}

func putNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := ctx.GetRawData()
		if err != nil {
			ctx.JSON(403, res.Err(403, err))
			return
		}
		_ = os.MkdirAll(path.Join("data", "config"), 0666)
		err = notification.HandlePush(string(data))
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}
