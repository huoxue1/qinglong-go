package api

import (
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/api/config"
	"github.com/huoxue1/qinglong-go/api/cron"
	"github.com/huoxue1/qinglong-go/api/env"
	"github.com/huoxue1/qinglong-go/api/scripts"
	"github.com/huoxue1/qinglong-go/api/system"
	"github.com/huoxue1/qinglong-go/api/user"
)

func Api(group *gin.RouterGroup) {
	system.Api(group.Group("/system"))
	cron.Api(group.Group("/crons"))
	user.Api(group.Group("/user"))
	env.Api(group.Group("/envs"))
	config.Api(group.Group("/configs"))
	scripts.Api(group.Group("/scripts"))
}
