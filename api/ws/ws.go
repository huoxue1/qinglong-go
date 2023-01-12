package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/huoxue1/qinglong-go/service/client"
	"github.com/huoxue1/qinglong-go/utils/res"
)

func Api(group *gin.RouterGroup) {
	group.GET("/info", info())
	group.GET("/:id/:name/websocket", wsHandle())
}

func info() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"websocket": true, "origins": []string{"*:*"}, "cookie_needed": false, "entropy": int64(3563341155)})
	}
}

func wsHandle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conn, _, _, err := ws.UpgradeHTTP(ctx.Request, ctx.Writer)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}

		client.AddWs(ctx.Param("id")+"_"+ctx.Param("name"), conn)
	}
}
