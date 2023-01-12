package user

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/models"
	"github.com/huoxue1/qinglong-go/service/notification"
	"github.com/huoxue1/qinglong-go/service/user"
	"github.com/huoxue1/qinglong-go/utils"
	"github.com/huoxue1/qinglong-go/utils/res"
	"os"
	"path"
	"time"
)

//go:embed config_sample.sh
var sample []byte

//go:embed package_sample.json
var pack []byte

func Api(group *gin.RouterGroup) {
	group.GET("/", get())
	group.PUT("/init", appInit())
	group.POST("/login", login())
	group.POST("/logout", logout())
	group.PUT("/notification/init", putNotification())
	group.PUT("/notification", putNotification())
	group.GET("/notification", getNotification())
}

func logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, res.Ok(true))
	}
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		info, err := user.GetUserInfo()
		if err != nil {

			return
		}
		ctx.JSON(200, res.Ok(gin.H{"username": info.Username}))
	}
}

func appInit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := os.Stat(path.Join("data", "config", "auth.json"))
		exist := os.IsExist(err)
		if exist {
			ctx.JSON(400, res.Err(400, err))
			return
		}
		_ = os.MkdirAll(path.Join("data", "config"), 0666)
		_ = os.MkdirAll(path.Join("data", "log"), 0666)
		_ = os.MkdirAll(path.Join("data", "repo"), 0666)
		_ = os.MkdirAll(path.Join("data", "scripts"), 0666)
		_ = os.MkdirAll(path.Join("data", "deps"), 0666)
		_ = os.MkdirAll(path.Join("data", "raw"), 0666)
		_ = os.WriteFile(path.Join("data", "config", "config.sh"), sample, 0666)
		_ = os.WriteFile(path.Join("data", "scripts", "package.json"), pack, 0666)
		_ = os.WriteFile(path.Join("data", "config", "config.sample.sh"), sample, 0666)
		type Req struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		}
		r := new(Req)
		err = ctx.ShouldBindJSON(r)
		if err != nil {
			ctx.JSON(503, res.ErrMessage(503, err.Error()))
			return
		}
		m := new(models.AuthFile)
		m.Username = r.UserName
		m.Password = r.Password
		m.Tokens.Mobile = ""
		m.Tokens.Desktop = ""
		data, _ := json.Marshal(m)
		_ = os.WriteFile(path.Join("data", "config", "auth.json"), data, 0666)
		ctx.JSON(200, res.Ok(true))
	}
}

func login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type Req struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		}
		r := new(Req)
		err := ctx.ShouldBindJSON(r)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
		data, err := os.ReadFile(path.Join("data", "config", "auth.json"))
		if err != nil {
			ctx.Abort()
			return
		}
		auth := new(models.AuthFile)
		_ = json.Unmarshal(data, auth)
		if auth.Username == r.UserName && auth.Password == r.Password {
			token, err := utils.GenerateToken(r.UserName, 48)
			if err != nil {
				ctx.JSON(503, res.Err(503, err))
				return
			}
			ip, err := user.GetNetIp(ctx.RemoteIP())
			if err != nil {
				ip = new(user.Ip)
				err = nil
			}
			mobile := utils.IsMobile(ctx.GetHeader("User-Agent"))
			if mobile {
				auth.Tokens.Mobile = token
				auth.Token = token
				file, _ := json.Marshal(auth)
				err := os.WriteFile(path.Join("data", "config", "auth.json"), file, 0666)
				if err != nil {
					ctx.JSON(503, res.Err(503, err))
					return
				}
				notification.Push.Send("登录通知", fmt.Sprintf("你于%s登录mobile端登录成功，ip地址 %s", time.Now().Format("2006-01-02 15:04:05"), ctx.ClientIP()))
				ctx.JSON(200, res.Ok(gin.H{
					"token":     token,
					"platform":  "mobile",
					"retries":   0,
					"lastip":    ctx.RemoteIP(),
					"lastaddr":  ip.Addr,
					"lastlogon": time.Now().UnixNano(),
				}))
			} else {
				auth.Tokens.Desktop = token
				auth.Token = token
				file, _ := json.Marshal(auth)
				err := os.WriteFile(path.Join("data", "config", "auth.json"), file, 0666)
				if err != nil {
					ctx.JSON(503, res.Err(503, err))
					return
				}
				notification.Push.Send("登录通知", fmt.Sprintf("你于%s登录pc端登录成功，ip地址 %s", time.Now().Format("2006-01-02 15:04:05"), ctx.ClientIP()))
				ctx.JSON(200, res.Ok(gin.H{
					"token":     token,
					"platform":  "desktop",
					"retries":   0,
					"lastip":    ctx.RemoteIP(),
					"lastaddr":  ip.Addr,
					"lastlogon": time.Now().Unix(),
				}))
			}
		} else {
			ctx.JSON(400, res.ErrMessage(400, "账号密码错误！"))
		}

	}
}
