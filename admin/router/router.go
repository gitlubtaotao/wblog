package admin

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/admin/api"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"net/http"
)

type IRoutes interface {
	Register()
}
type Routes struct {
	engine *gin.Engine
	group  *gin.RouterGroup
}

func NewRoutes(router *gin.Engine) IRoutes {
	return &Routes{engine: router, group: router.Group("/admin")}
}

func (r *Routes) Register() {
	r.sessionRoute()
	r.registerRoute()
	r.passwordRoute()
	r.captchaRoute()
	r.group.Use(r.AdminScopeRequired())
	{
		r.homeRoute()
		r.user()
	}
	
}

func (r *Routes) sessionRoute() {
	session := admin.SessionApi{}
	r.group.GET("/login", session.New)
	r.group.POST("/login", session.Create)
	r.group.DELETE("/destroy", session.Destroy)
}
func (r *Routes) registerRoute() {
	if system.GetConfiguration().SignupEnabled {
		register := admin.RegisterApi{}
		r.group.GET("/register", register.New)
		r.group.POST("/register", register.Create)
	}
}
func (r *Routes) passwordRoute() {
	password := admin.PasswordApi{}
	r.group.GET("/password", password.New)
	r.group.POST("/create", password.Create)
	r.group.POST("/modify", password.Modify)
	r.group.GET("/send_notice", password.SendEmail)
}

func (r *Routes) captchaRoute() {
	captcha := &admin.CaptchaApi{}
	r.group.GET("/verify_captcha", captcha.Verify)
	r.group.GET("/captcha/:source", captcha.Image)
	r.group.GET("/captcha", captcha.Get)
	
}

/*
@title: 首页注册路由
*/
func (r *Routes) homeRoute() {
	home := admin.HomeApi{}
	{
		r.group.GET("", home.Index)
		r.group.GET("/index", home.Index)
	}
}
func (r *Routes) user()  {
	user := &admin.UserApi{}
	r.group.GET("/user/profile", user.Get)
	r.group.POST("/user/:id/profile", user.Update)
	r.group.GET("/user", user.Index)
	r.group.POST("/user/:id/lock", user.Lock)
}

//AuthRequired grants access to authenticated users, requires SharedData middleware
func (r *Routes) AdminScopeRequired() gin.HandlerFunc {
	config := system.GetConfiguration()
	return func(c *gin.Context) {
		if user, _ := c.Get(config.AdminUser); user != nil {
			if u, ok := user.(*models.User); ok && u.IsAdmin {
				c.Next()
				return
			}
		}
		_ = seelog.Warnf("User not authorized to visit %s", c.Request.RequestURI)
		c.Redirect(http.StatusSeeOther, "/admin/login")
		c.Abort()
	}
}
