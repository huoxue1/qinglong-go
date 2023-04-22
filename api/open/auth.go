package open

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/internal/res"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/utils"
	"strconv"
	"time"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientId := ctx.Query("client_id")
		clientSecret := ctx.Query("client_secret")
		app, err := models.GetAppById(clientId)
		if err != nil {
			ctx.JSON(401, res.Err(401, err))
			return
		}
		if app.ClientSecret != clientSecret {
			ctx.JSON(401, res.Err(401, errors.New("the auth fail")))
			return
		}
		token, err := utils.GenerateToken(strconv.Itoa(app.Id), 720)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		app.Tokens = append(app.Tokens, token)
		models.UpdateApp(app)
		ctx.JSON(200, res.Ok(map[string]any{
			"token":      token,
			"token_type": "Bearer",
			"expiration": time.Now().Add(time.Hour * 720).Unix(),
		}))

	}
}
