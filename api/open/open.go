package open

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/open"
	"github.com/huoxue1/qinglong-go/utils/res"
	"strconv"
)

func Api(group *gin.RouterGroup) {
	group.GET("", get())
	group.POST("", post())
	group.PUT("", put())
	group.DELETE("", del())
	group.PUT("/:id/reset-secret", reset())
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apps, err := models.QueryApp()
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(apps))
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
		err = open.DeleteApp(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m := new(models.Apps)
		err := ctx.ShouldBindJSON(m)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		err = open.UpdateApp(m)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(m))
	}
}

func post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m := new(models.Apps)
		err := ctx.ShouldBindJSON(m)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		id, err := open.AddApp(m)
		if err != nil {
			return
		}
		app, _ := models.GetApp(id)
		ctx.JSON(200, res.Ok(app))
	}
}

func reset() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		app, err := models.GetApp(id)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		err = open.ResetApp(app)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(app))
	}
}
