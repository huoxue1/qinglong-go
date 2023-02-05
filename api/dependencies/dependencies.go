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

	group.DELETE("", del())
}

var (
	typMap = map[string]int{
		"nodejs":  0,
		"python3": 1,
		"linux":   2,
	}
)

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		dependences, err := models.QueryDependences(ctx.Query("searchValue"), typMap[ctx.Query("type")])
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(dependences))
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
		dependencies.DelDep(ids)
		ctx.JSON(200, res.Ok(true))
	}
}

func post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var deps []*models.Dependences
		err := ctx.ShouldBindJSON(&deps)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
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
