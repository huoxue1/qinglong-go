package auth

import (
	"encoding/json"
	"github.com/huoxue1/qinglong-go/internal/cache"
	"github.com/huoxue1/qinglong-go/models"
)

const (
	CacheAuthKey = "auth_info"
)

func InitAuthInfo(username, password string) error {
	m := new(models.AuthFile)
	m.Username = username
	m.Password = password
	m.IsTwoFactorChecking = false
	m.Tokens.Mobile = ""
	m.Tokens.Desktop = ""
	data, _ := json.Marshal(m)
	return cache.SetBytes(CacheAuthKey, data)
}

func GetAuthInfo() (*models.AuthFile, error) {
	m := new(models.AuthFile)
	data, err := cache.GetBytes(CacheAuthKey)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func UpdateAuthInfo(file *models.AuthFile) error {
	data, _ := json.Marshal(file)
	return cache.SetBytes(CacheAuthKey, data)
}

func IsInit() bool {
	return cache.Exists(CacheAuthKey)
}
