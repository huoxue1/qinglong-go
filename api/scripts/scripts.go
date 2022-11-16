package scripts

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/service/scripts"
	"github.com/huoxue1/qinglong-go/utils/res"
	"os"
	path2 "path"
)

func Api(group *gin.RouterGroup) {
	group.GET("", get())
	group.PUT("", put())
	group.POST("", post())
	group.DELETE("", del())
	group.GET("/:name", getFile())
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		files := scripts.GetFiles()
		ctx.JSON(200, res.Ok(files))
	}
}

func put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type Req struct {
			Path     string `json:"path"`
			FileName string `json:"filename"`
			Content  string `json:"content"`
		}
		r := new(Req)
		err := ctx.ShouldBindJSON(r)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		err = os.WriteFile(fmt.Sprintf("data/scripts/%s/%s", r.Path, r.FileName), []byte(r.Content), 0666)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		ctx.JSON(200, res.Ok(true))
	}
}

func post() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type Req struct {
			Path      string `json:"path" form:"path"`
			FileName  string `json:"filename" form:"filename"`
			Content   string `json:"content" form:"content"`
			Directory string `json:"directory" form:"directory"`
		}
		r := new(Req)
		err := ctx.ShouldBind(r)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		if r.Directory != "" {
			err := os.MkdirAll(path2.Join("data", "scripts", r.Path, r.Directory), 0666)
			if err != nil {
				ctx.JSON(503, res.Err(503, err))
				return
			}
		} else {
			if r.FileName == "undefined" {
				file, err := ctx.FormFile("file")
				if err != nil {
					ctx.JSON(503, res.Err(503, err))
					return
				}
				err = ctx.SaveUploadedFile(file, path2.Join("data", "scripts", r.Path, file.Filename))
				if err != nil {
					ctx.JSON(503, res.Err(503, err))
					return
				}

			} else {
				_, err := os.Create(path2.Join("data", "scripts", r.Path, r.FileName))
				if err != nil {
					ctx.JSON(503, res.Err(503, err))
					return
				}
			}
			ctx.JSON(200, res.Ok(true))

		}
	}
}

func getFile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileName := ctx.Param("name")
		path := ctx.Query("path")
		data, _ := os.ReadFile(fmt.Sprintf("data/scripts/%s/%s", path, fileName))
		ctx.JSON(200, res.Ok(string(data)))
	}
}

func del() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type Req struct {
			Path     string `json:"path"`
			FileName string `json:"filename"`
			Type     string `json:"type"`
		}
		r := new(Req)
		err := ctx.ShouldBindJSON(r)
		if err != nil {
			ctx.JSON(502, res.Err(502, err))
			return
		}
		if r.Type == "file" {
			err := os.Remove(fmt.Sprintf("data/scripts/%s/%s", r.Path, r.FileName))
			if err != nil {
				ctx.JSON(503, res.Err(503, err))
				return
			}
		} else {
			err := os.RemoveAll(fmt.Sprintf("data/scripts/%s/%s", r.Path, r.FileName))
			if err != nil {
				ctx.JSON(503, res.Err(503, err))
				return
			}
		}
		ctx.JSON(200, res.Ok(true))
	}

}
