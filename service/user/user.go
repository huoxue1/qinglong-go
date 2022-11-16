package user

import (
	"encoding/json"
	"os"
)

type Info struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Tokens   struct {
		Desktop string `json:"desktop"`
		Mobile  string `json:"mobile"`
	} `json:"tokens"`
	Lastlogon           int64  `json:"lastlogon"`
	Retries             int    `json:"retries"`
	Lastip              string `json:"lastip"`
	Lastaddr            string `json:"lastaddr"`
	Platform            string `json:"platform"`
	IsTwoFactorChecking bool   `json:"isTwoFactorChecking"`
}

func GetUserInfo() (*Info, error) {
	i := new(Info)
	file, err := os.ReadFile("./data/config/auth.json")
	if err != nil {
		return i, err
	}
	err = json.Unmarshal(file, i)
	return i, err

}
