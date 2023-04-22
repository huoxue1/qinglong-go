package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/internal/auth"
	"github.com/huoxue1/qinglong-go/internal/res"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/utils"
	"strconv"
	"strings"
	"time"
)

var (
	unExcludedPath = []string{
		"/api/system",
		"/api/user/login",
		"/api/user/init",
		"api/public/panel/log",
		"/api/user/notification/init",
		"/api/user/two-factor/login",
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
			userId, exp, _ := utils.ParseToken(authToken)
			if exp < time.Now().Unix() {
				ctx.JSON(401, res.Err(401, errors.New("the authorization token is expired")))
				ctx.Abort()
				return
			}
			id, _ := strconv.Atoi(userId)
			app, err := models.GetApp(id)
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

		authFile, err := auth.GetAuthInfo()
		if err != nil {
			ctx.JSON(401, res.Err(401, errors.New("the authorization fail")))
			ctx.Abort()
			return
		}
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
		_, exp, err := utils.ParseToken(authToken)
		if err != nil {
			ctx.JSON(401, res.Err(401, errors.New("the authorization fail")))
			ctx.Abort()
			return
		}
		if exp < time.Now().Unix() {
			ctx.JSON(401, res.Err(401, errors.New("the authorization token is expired")))
			ctx.Abort()
			return
		}
		if mobile {
			if authToken != authFile.Tokens.Mobile && authToken != authFile.Token {
				ctx.JSON(401, res.Err(401, errors.New("the authorization token is error")))
				ctx.Abort()
				return
			} else {
				ctx.Next()
				return
			}
		} else {
			if authToken != authFile.Tokens.Desktop && authToken != authFile.Token {
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
