package user

import (
	"encoding/json"
	"github.com/huoxue1/qinglong-go/utils"
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

type Ip struct {
	Ip          string `json:"ip"`
	Pro         string `json:"pro"`
	ProCode     string `json:"proCode"`
	City        string `json:"city"`
	CityCode    string `json:"cityCode"`
	Region      string `json:"region"`
	RegionCode  string `json:"regionCode"`
	Addr        string `json:"addr"`
	RegionNames string `json:"regionNames"`
	Err         string `json:"err"`
}

func GetNetIp(ip string) (*Ip, error) {
	i := new(Ip)
	_, err := utils.GetClient().R().SetQueryParams(map[string]string{"ip": ip, "json": "true"}).SetResult(i).Get("https://whois.pconline.com.cn/ipJson.jsp")
	if err != nil {
		return nil, err
	}
	return i, nil
}
