package subscription

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/internal/res"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/subscription"
	"os"
	"strconv"
)

func Api(group *gin.RouterGroup) {
	group.GET("", get())
	group.POST("", post())
	group.PUT("", put())
	group.PUT("/disable", disable())
	group.PUT("/enable", enable())
	group.DELETE("", del())
	group.PUT("/run", run())
	group.GET("/:id/log", log1())
	group.PUT("/stop", stop())
}

func run() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		err = subscription.RunSubscription(ids)
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
			ctx.JSON(503, res.Err(503, err))
			return
		}
		err = subscription.StopSubscription(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		subs, err := models.QuerySubscription(ctx.Query("searchValue"))
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(subs))

	}
}

func post() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		sub := new(models.Subscriptions)
		err := ctx.ShouldBindJSON(sub)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		id, err := subscription.AddSubscription(sub)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		sub.Id = id
		ctx.JSON(200, res.Ok(sub))
	}
}

func enable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		err = subscription.EnableSubscription(ids)
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
			ctx.JSON(503, res.Err(503, err))
			return
		}
		err = subscription.DisableSubscription(ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := new(models.Subscriptions)
		err := ctx.ShouldBindJSON(s)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		err = subscription.UpdateSubscription(s)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(s))
	}
}

func del() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ids []int
		err := ctx.ShouldBindJSON(&ids)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		err = subscription.DeleteSubscription(ids)
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
		s, _ := models.GetSubscription(id)
		data, _ := os.ReadFile(s.LogPath)
		ctx.JSON(200, res.Ok(string(data)))
	}
}
