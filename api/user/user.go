package user

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huoxue1/qinglong-go/internal/auth"
	"github.com/huoxue1/qinglong-go/internal/res"
	"github.com/huoxue1/qinglong-go/service/notification"
	"github.com/huoxue1/qinglong-go/service/user"
	"github.com/huoxue1/qinglong-go/utils"
	"os"
	"path"
	"time"
)

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
	group.GET("/two-factor/init", twoFactorInit())
	group.PUT("/two-factor/active", twoFactorActive())
	group.PUT("/two-factor/login", twoFactorLogin())
	group.PUT("/two-factor/deactive", twoFactorDeactive())
}

func twoFactorDeactive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		info, err := auth.GetAuthInfo()
		if err != nil {
			ctx.JSON(401, res.Err(401, err))
			return
		}
		info.TwoFactorSecret = ""
		info.IsTwoFactorChecking = false
		err = auth.UpdateAuthInfo(info)
		if err != nil {
			ctx.JSON(500, res.Err(500, err))
			return
		}
		ctx.JSON(200, res.Ok(nil))
	}
}

func twoFactorLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		info, err := auth.GetAuthInfo()
		if err != nil {
			ctx.JSON(401, res.Err(401, err))
			return
		}
		if info.TwoFactorSecret == "" {
			ctx.JSON(400, res.Err(400, errors.New("the two factor is not initialized")))
			return
		}
		type req struct {
			UserName string `json:"username"`
			Password string `json:"password"`
			Code     string `json:"code"`
		}
		var r req
		err = ctx.ShouldBindJSON(&r)
		if err != nil {
			ctx.JSON(400, res.Err(400, err))
			return
		}

		if info.Username != r.UserName || info.Password != r.Password {
			ctx.JSON(400, res.Err(400, errors.New("the username or password is invalid")))
			return
		}

		if !auth.VerifyTOTP(info.TwoFactorSecret, r.Code) {
			ctx.JSON(400, res.Err(400, errors.New("the two factor code is invalid")))
			return
		}
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
			info.Tokens.Mobile = token
			info.Token = token
			err := auth.UpdateAuthInfo(info)
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
			info.Tokens.Desktop = token
			info.Token = token
			err := auth.UpdateAuthInfo(info)
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
	}
}

func twoFactorActive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		info, err := auth.GetAuthInfo()
		if err != nil {
			ctx.JSON(401, res.Err(401, err))
			return
		}
		if info.TwoFactorSecret == "" {
			ctx.JSON(400, res.Err(400, errors.New("the two factor is not initialized")))
			return
		}
		type req struct {
			Code string `json:"code"`
		}
		var r req
		err = ctx.ShouldBindJSON(&r)
		if err != nil {
			ctx.JSON(400, res.Err(400, err))
			return
		}
		if !auth.VerifyTOTP(info.TwoFactorSecret, r.Code) {
			ctx.JSON(400, res.Err(400, errors.New("the two factor code is invalid")))
			return
		}
		info.IsTwoFactorChecking = true
		_ = auth.UpdateAuthInfo(info)
		ctx.JSON(200, res.Ok(true))
	}
}

func twoFactorInit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		info, err := auth.GetAuthInfo()
		if err != nil {
			ctx.JSON(401, res.Err(401, err))
			return
		}
		//if info.TwoFactorSecret != "" {
		//	ctx.JSON(400, res.Err(400, errors.New("the two factor is initialized")))
		//	return
		//}
		secret, qrcode, err := auth.GenerateTOTP(info.Username, "qinglong-go")
		if err != nil {
			ctx.JSON(500, res.Err(500, err))
			return
		}
		info.TwoFactorSecret = secret
		_ = auth.UpdateAuthInfo(info)
		ctx.JSON(200, res.Ok(gin.H{"secret": secret, "url": qrcode}))

	}
}

func logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, res.Ok(true))
	}
}

func get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		info, err := auth.GetAuthInfo()
		if err != nil {
			ctx.JSON(401, res.Err(401, err))
			return
		}
		ctx.JSON(200, res.Ok(gin.H{"username": info.Username, "twoFactorActivated": info.IsTwoFactorChecking}))
	}
}

func appInit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isInit := auth.IsInit()
		if isInit {
			ctx.JSON(400, res.Err(400, errors.New("the app is initialized")))
			return
		}
		_ = os.MkdirAll(path.Join("data", "deps"), 0666)
		_ = os.Link(path.Join("data", "deps"), path.Join("data", "scripts", "deps"))
		_ = os.MkdirAll(path.Join("data", "log"), 0666)
		_ = os.MkdirAll(path.Join("data", "repo"), 0666)
		_ = os.MkdirAll(path.Join("data", "scripts"), 0666)
		_ = os.MkdirAll(path.Join("data", "deps"), 0666)
		_ = os.MkdirAll(path.Join("data", "raw"), 0666)
		_ = os.WriteFile(path.Join("data", "scripts", "package.json"), pack, 0666)

		type Req struct {
			UserName string `json:"username"`
			Password string `json:"password"`
		}
		r := new(Req)
		err := ctx.ShouldBindJSON(r)
		if err != nil {
			ctx.JSON(503, res.ErrMessage(503, err.Error()))
			return
		}
		err = auth.InitAuthInfo(r.UserName, r.Password)
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}
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
		authInfo, err := auth.GetAuthInfo()
		if err != nil {
			ctx.JSON(503, res.Err(503, err))
			return
		}

		if authInfo.IsTwoFactorChecking {
			ctx.JSON(200, res.Err(420, errors.New("")))
			return
		}

		if authInfo.Username == r.UserName && authInfo.Password == r.Password {
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
				authInfo.Tokens.Mobile = token
				authInfo.Token = token
				err := auth.UpdateAuthInfo(authInfo)
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
				authInfo.Tokens.Desktop = token
				authInfo.Token = token
				err := auth.UpdateAuthInfo(authInfo)
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
