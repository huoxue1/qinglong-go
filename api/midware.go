package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/utils"
	"github.com/huoxue1/qinglong-go/utils/res"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	unExcludedPath = []string{
		"/api/system",
		"/api/user/login",
		"/api/user/init",
	}
)

func OpenJwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.URL.Path, "/open/auth/token") {
			ctx.Next()
			return
		} else {
			tokenHeader := ctx.GetHeader("Authorization")
			if tokenHeader == "" {
				ctx.JSON(401, res.Err(401, errors.New("no authorization token was found")))
				ctx.Abort()
				return
			}
			authToken := strings.Split(tokenHeader, " ")[1]
			claims, _ := utils.ParseToken(authToken)
			if claims.ExpiresAt < time.Now().Unix() {
				ctx.JSON(401, res.Err(401, errors.New("the authorization token is expired")))
				ctx.Abort()
				return
			}
			userId, _ := strconv.Atoi(claims.UserID)
			app, err := models.GetApp(userId)
			if err != nil {
				ctx.JSON(401, res.Err(401, errors.New("the authorization token is invalid")))
				ctx.Abort()
				return
			}
			for _, scope := range app.Scopes {
				if strings.HasPrefix(ctx.Request.URL.Path, "/open/"+scope) {
					ctx.Next()
				}
			}
			ctx.Abort()
		}
	}
}

func Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, s := range unExcludedPath {
			if strings.HasPrefix(ctx.Request.URL.Path, s) {
				ctx.Next()
				return
			}
		}

		data, err := os.ReadFile(path.Join("data", "config", "auth.json"))
		if err != nil {
			ctx.Abort()
			return
		}
		auth := new(models.AuthFile)
		_ = json.Unmarshal(data, auth)
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			queryToken, b := ctx.GetQuery("token")
			if b {
				tokenHeader = "Bearer " + queryToken
			} else {
				ctx.JSON(401, res.Err(401, errors.New("no authorization token was found")))
				ctx.Abort()
				return
			}

		}
		authToken := strings.Split(tokenHeader, " ")[1]

		mobile := utils.IsMobile(ctx.GetHeader("User-Agent"))
		claims, err := utils.ParseToken(authToken)
		if err != nil {
			ctx.JSON(401, res.Err(401, errors.New("the authorization fail")))
			ctx.Abort()
			return
		}
		if claims.ExpiresAt < time.Now().Unix() {
			ctx.JSON(401, res.Err(401, errors.New("the authorization token is expired")))
			ctx.Abort()
			return
		}
		if mobile {
			if authToken != auth.Tokens.Mobile && authToken != auth.Token {
				ctx.JSON(401, res.Err(401, errors.New("the authorization token is error")))
				ctx.Abort()
				return
			} else {
				ctx.Next()
				return
			}
		} else {
			if authToken != auth.Tokens.Desktop && authToken != auth.Token {
				ctx.JSON(401, res.Err(401, errors.New("the authorization token is error")))
				ctx.Abort()
				return
			} else {
				ctx.Next()
				return
			}
		}

	}
}
