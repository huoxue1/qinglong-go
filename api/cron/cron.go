package cron

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/cron"
	"github.com/huoxue1/qinglong-go/utils/res"
	"os"
	"strconv"
	"time"
)

func Api(group *gin.RouterGroup) {
	group.GET("", get())
	group.POST("", post())
	group.DELETE("", del())
	group.PUT("", put())
	group.PUT("/disable", disable())
	group.PUT("/enable", enable())
	group.PUT("/pin", pin())
	group.PUT("/unpin", unpin())
	group.GET("/:id/log", log1())
	group.PUT("/run", run())
	group.PUT("/stop", stop())
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sorter, ok := ctx.GetQueryMap("sorter")
		if !ok {
			sorter = map[string]string{"field": "name", "type": "ASC"}
		}
		filters := ctx.QueryMap("queryString")["filters"]
		page, _ := strconv.Atoi(ctx.Query("page"))
		size, _ := strconv.Atoi(ctx.Query("size"))
		if size == 0 {
			size = 1000
		}
		crons, err := cron.GetCrons(page, size, ctx.Query("searchValue"), sorter, filters)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(gin.H{
			"data":  crons,
			"total": models.Count(ctx.Query("searchValue")),
		}))
	}
}

func post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := new(models.Crontabs)
		err := ctx.ShouldBindJSON(c)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		c.Status = 1
		c.Createdat = time.Now().Format(time.RFC3339)
		c.Updatedat = time.Now().Format(time.RFC3339)
		c.Timestamp = time.Now().Format("Mon Jan 02 2006 15:04:05 MST")
		id, err := cron.AddCron(c)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		c.Id = id
		ctx.JSON(200, res.Ok(c))
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
		err = cron.DeleteCron(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := new(models.Crontabs)
		err := ctx.ShouldBindJSON(c)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = cron.UpdateCron(c)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(cron.GetCron(c.Id)))
	}
}

func disable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = cron.DisableCron(ids)
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
		err = cron.EnableCron(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func pin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = cron.PinCron(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func unpin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = cron.UnPinCron(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func run() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = cron.RunCron(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func stop() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(502, err))
			return
		}
		err = cron.StopCron(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func log1() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		c := cron.GetCron(id)
		data, _ := os.ReadFile(c.LogPath)
		ctx.JSON(200, res.Ok(string(data)))
	}
}
