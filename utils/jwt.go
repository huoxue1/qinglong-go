package utils

// jwt身份验证demo

import (
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/utils/res"
	"os"
	"path"
	"strings"
	"time"
)

// 设置jwt密钥secret
var jwtSecret = []byte("qinglong")

type Claims struct {
	UserID string `json:"userid"`
	jwt.StandardClaims
}

// GenerateToken 生成token的函数
func GenerateToken(userid string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(48 * time.Hour)

	claims := Claims{
		userid, // 自行添加的信息
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 设置token过期时间
			Issuer:    "gin-blog",        // 设置jwt签发者
		},
	}
	// 生成token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 验证token的函数
func ParseToken(token string) (*Claims, error) {
	// 对token的密钥进行验证
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	// 判断token是否过期
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

var (
	unExcludedPath = []string{
		"/api/system",
		"/api/user/login",
		"/api/user/init",
	}
)

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
			ctx.JSON(401, res.Err(401, errors.New("no authorization token was found")))
			ctx.Abort()
			return
		}
		authToken := strings.Split(tokenHeader, " ")[1]

		mobile := IsMobile(ctx.GetHeader("User-Agent"))
		claims, _ := ParseToken(authToken)
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

func IsMobile(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}

	isMobile := false
	mobileKeywords := []string{"Mobile", "Android", "Silk/", "Kindle",
		"BlackBerry", "Opera Mini", "Opera Mobi", "app"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			isMobile = true
			break
		}
	}

	return isMobile
}
