package user

import (
	"github.com/huoxue1/qinglong-go/utils"
)

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
