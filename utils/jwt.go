package utils

// jwt身份验证demo

import (
	"github.com/dgrijalva/jwt-go"
	"math/rand"
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
func GenerateToken(userid string, hour int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(hour) * time.Hour)

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

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-?")
	rand.Seed(time.Now().UnixMilli())
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
