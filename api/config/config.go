package config

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/internal/res"
	"os"
)

const (
	Dir = "data/config/"
)

func Api(group *gin.RouterGroup) {
	group.GET("/files", files())
	group.GET("/:name", getFile())
	group.POST("/save", save())
}

func files() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		entries, err := os.ReadDir(Dir)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
		}
		var result []map[string]string

		for _, entry := range entries {
			if !entry.IsDir() && entry.Name() != "auth.json" {
				result = append(result, map[string]string{
					"title": entry.Name(),
					"value": entry.Name(),
				})
			}
		}
		ctx.JSON(200, res.Ok(result))
	}
}

func getFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		file, err := os.ReadFile(Dir + name)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(string(file)))
	}
}

func save() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type req struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}
		r := new(req)
		err := ctx.ShouldBindJSON(r)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		err = os.WriteFile(Dir+r.Name, []byte(r.Content), 0666)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		ctx.JSON(200, res.OkMessage(true, "保存成功"))
	}
}
