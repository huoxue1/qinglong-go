package utils

import (
	"github.com/imroc/req/v3"
	"net/http"
)

var (
	client *req.Client
)

func init() {
	client = req.C().SetProxy(http.ProxyFromEnvironment)
}

func GetClient() *req.Client {
	return client
}
