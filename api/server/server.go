package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/service/server"
)

func Api(group *gin.RouterGroup) {
	group.Match([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPut,
		http.MethodOptions,
	}, "/:path", handle())
}

func handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Param("path")
		query := make(map[string]string)
		ctx.ShouldBindQuery(&query)
		headers := make(map[string]string)
		ctx.ShouldBindHeader(&headers)

		body := make(map[string]any)
		ctx.ShouldBind(&body)
		data := server.Run(ctx, path, query, body, headers, false)
		ctx.Writer.WriteHeader(200)
		ctx.Writer.Header().Add("Content-Type", "application/json")
		_, _ = ctx.Writer.WriteString(data)
	}
}
