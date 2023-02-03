package controller

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/api"
	"io/ioutil"
	"strings"
)

func Router() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	engine.Use(static.Serve("/", static.LocalFile("static/dist/", false)))
	engine.NoRoute(func(ctx *gin.Context) {
		accept := ctx.Request.Header.Get("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			content, err := ioutil.ReadFile("static/dist/index.html")
			if (err) != nil {
				ctx.Writer.WriteHeader(404)
				_, _ = ctx.Writer.WriteString("Not Found")
				return
			}
			ctx.Writer.WriteHeader(200)
			ctx.Writer.Header().Add("Accept", "text/html")
			_, _ = ctx.Writer.Write(content)
			ctx.Writer.Flush()
		}
	})
	api.Api(engine.Group("/api", api.Jwt()))
	api.Open(engine.Group("/open", api.OpenJwt()))

	return engine
}
