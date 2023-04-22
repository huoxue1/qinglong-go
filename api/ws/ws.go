package ws

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/internal/res"
	"github.com/huoxue1/qinglong-go/service/client"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
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
		conn, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{
			Subprotocols:         nil,
			InsecureSkipVerify:   false,
			OriginPatterns:       []string{"*"},
			CompressionMode:      0,
			CompressionThreshold: 0,
		})
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}

		wsjson.Write(context.Background(), conn, map[string]string{"123": "11"})
		wsjson.Write(context.Background(), conn, map[string]string{"123": "11"})
		wsjson.Write(context.Background(), conn, map[string]string{"123": "11"})

		c := make(chan any, 100)
		client.AddChan(c)

		for true {
			wsjson.Write(context.Background(), conn, map[string]string{"123": "11"})
			time.Sleep(1000)
			wsjson.Write(context.Background(), conn, map[string]string{"123": "11"})
			data := <-c
			err := wsjson.Write(context.Background(), conn, data)
			if err != nil {
				break
			}
		}
	}

}
