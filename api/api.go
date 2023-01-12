package api

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/api/config"
	"github.com/huoxue1/qinglong-go/api/cron"
	"github.com/huoxue1/qinglong-go/api/dependencies"
	"github.com/huoxue1/qinglong-go/api/env"
	"github.com/huoxue1/qinglong-go/api/logs"
	"github.com/huoxue1/qinglong-go/api/open"
	"github.com/huoxue1/qinglong-go/api/public"
	"github.com/huoxue1/qinglong-go/api/scripts"
	"github.com/huoxue1/qinglong-go/api/subscription"
	"github.com/huoxue1/qinglong-go/api/system"
	"github.com/huoxue1/qinglong-go/api/user"
	"github.com/huoxue1/qinglong-go/api/ws"
)

func Api(group *gin.RouterGroup) {
	system.Api(group.Group("/system"))
	cron.Api(group.Group("/crons"))
	user.Api(group.Group("/user"))
	env.Api(group.Group("/envs"))
	config.Api(group.Group("/configs"))
	scripts.Api(group.Group("/scripts"))
	open.Api(group.Group("/apps"))
	subscription.Api(group.Group("/subscriptions"))
	logs.APi(group.Group("/logs"))
	dependencies.Api(group.Group("/dependencies"))
	ws.Api(group.Group("/ws"))
	public.Api(group.Group("/public"))
}
