package controller

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/api"
	"net/http"
)

func Router() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	engine.Use(static.Serve("/", static.LocalFile("static/dist/", false)))
	engine.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodGet {
			ctx.Redirect(301, "/")
			return
		}
		ctx.Next()
	})
	api.Api(engine.Group("/api", api.Jwt()))
	api.Open(engine.Group("/open", api.OpenJwt()))

	return engine
}
