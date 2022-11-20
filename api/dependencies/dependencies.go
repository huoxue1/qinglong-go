package dependencies

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/dependencies"
	"github.com/huoxue1/qinglong-go/utils/res"
	"strconv"
)

func Api(group *gin.RouterGroup) {
	group.POST("", post())
	group.GET("", get())
	group.GET("/:id", getDep())
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dependences, err := models.QueryDependences(ctx.Query("searchValue"))
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		ctx.JSON(200, res.Ok(dependences))
	}
}

func post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var deps []*models.Dependences
		err := ctx.ShouldBindJSON(&deps)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		for _, dep := range deps {
			dependencies.AddDep(dep)
		}

		ctx.JSON(200, res.Ok(deps))
	}
}

func getDep() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		dep, _ := models.GetDependences(id)
		ctx.JSON(200, res.Ok(dep))
	}
}
