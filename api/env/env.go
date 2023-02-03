package env

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/env"
	"github.com/huoxue1/qinglong-go/utils/res"
	"io"
	"time"
)

func Api(group *gin.RouterGroup) {
	group.GET("", get())
	group.POST("", post())
	group.PUT("", put())
	group.DELETE("", del())
	group.PUT("/enable", enable())
	group.PUT("/disable", disable())
	group.POST("/upload", upload())
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		envs, err := env.QueryEnv(ctx.Query("searchValue"))
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
		}
		ctx.JSON(200, res.Ok(envs))
	}
}

func post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		envs := make([]*models.Envs, 0)
		err := ctx.ShouldBindJSON(&envs)
		if err != nil {
			ctx.JSON(200, res.Err(502, err))
			return
		}
		status := 0
		for _, e := range envs {
			e.Createdat = time.Now()
			e.Updatedat = time.Now()
			e.Timestamp = time.Now().Format("Mon Jan 02 2006 15:04:05 MST")
			e.Status = status
			id, err := env.AddEnv(e)
			if err != nil {
				ctx.JSON(200, res.Err(503, err))
				return
			}
			e.Id = id
		}
		ctx.JSON(200, res.Ok(envs))
	}
}

func put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		e := new(models.Envs)
		err := ctx.ShouldBindJSON(e)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = env.UpdateEnv(e)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(env.GetEnv(e.Id)))
	}
}

func del() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = env.DeleteEnv(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}
func enable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = env.EnableEnv(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}
func disable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		err = env.DisableEnv(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func upload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("env")
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		data, _ := file.Open()
		defer data.Close()
		content, _ := io.ReadAll(data)
		envs := make([]*models.Envs, 0)
		err = json.Unmarshal(content, &envs)
		if err != nil {
			ctx.JSON(200, res.Err(502, err))
			return
		}
		status := 0
		for _, e := range envs {
			e.Createdat = time.Now()
			e.Updatedat = time.Now()
			e.Timestamp = time.Now().Format("Mon Jan 02 2006 15:04:05 MST")
			e.Status = status
			id, err := env.AddEnv(e)
			if err != nil {
				ctx.JSON(200, res.Err(503, err))
				return
			}
			e.Id = id
		}
		ctx.JSON(200, res.Ok(envs))
	}
}
