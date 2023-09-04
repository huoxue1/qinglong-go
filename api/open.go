package api

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/api/config"
	"github.com/huoxue1/qinglong-go/api/cron"
	"github.com/huoxue1/qinglong-go/api/dependencies"
	"github.com/huoxue1/qinglong-go/api/env"
	"github.com/huoxue1/qinglong-go/api/logs"
	"github.com/huoxue1/qinglong-go/api/open"
	"github.com/huoxue1/qinglong-go/api/scripts"
	"github.com/huoxue1/qinglong-go/api/server"
	"github.com/huoxue1/qinglong-go/api/subscription"
	"github.com/huoxue1/qinglong-go/api/system"
)

func Open(group *gin.RouterGroup) {
	group.GET("/auth/token", open.Auth())
	system.Api(group.Group("/system"))
	cron.Api(group.Group("/crons"))
	env.Api(group.Group("/envs"))
	config.Api(group.Group("/configs"))
	scripts.Api(group.Group("/scripts"))
	open.Api(group.Group("/apps"))
	subscription.Api(group.Group("/subscriptions"))
	logs.APi(group.Group("/logs"))
	dependencies.Api(group.Group("/dependencies"))
	server.Api(group.Group("/server"))
}
