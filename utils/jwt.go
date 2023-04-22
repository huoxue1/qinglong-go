package utils

// jwt身份验证demo

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"strings"
	"time"
)

// 设置jwt密钥secret
var jwtSecret = []byte("qinglong")

// GenerateToken 生成token的函数
func GenerateToken(userid string, hour int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(hour) * time.Hour)

	// 生成token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userid,
		"exp":    expireTime.Unix(),
		"iat":    nowTime.Unix(),
		"issuer": "qinglong-go",
	})
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken 验证token的函数
func ParseToken(token string) (string, int64, error) {
	// 对token的密钥进行验证
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSecret, nil
	})
	if err != nil {
		return "", 0, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok && tokenClaims.Valid {
			return claims["userid"].(string), int64(claims["exp"].(float64)), nil
		}
	}

	return "", 0, err
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
