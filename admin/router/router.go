package admin

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/admin/api"
	"github.com/gitlubtaotao/wblog/api"
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
	base := api.BaseApi{}
	r.engine.NoRoute(base.Handle404)
	r.sessionRoute()
	r.registerRoute()
	r.passwordRoute()
	r.captchaRoute()
	r.group.Use(r.AdminScopeRequired())
	{
		r.auth()
		r.homeRoute()
		r.user()
		r.post()
		r.tag()
		r.page()
		link := &admin.LinkApi{}
		r.group.GET("/link", link.Index)
		r.group.POST("/link", link.Create)
		r.group.GET("/link/:id/show", link.Get)
		r.group.POST("/link/:id/update", link.Update)
		r.group.DELETE("/link/:id/delete", link.Delete)
		
		adminSub := admin.SubscriberApi{}
		r.group.GET("/subscriber", adminSub.Index)
		
		upload := admin.UploadApi{}
		r.group.POST("/upload", upload.Upload)
		
		// comment
		comment := admin.CommentApi{}
		r.group.POST("/comment/:id", comment.CommentRead)
		r.group.POST("/read_all", comment.CommentReadAll)
		
		mail := admin.MailApi{}
		r.group.POST("/mail/send", mail.Send)
		r.group.POST("/mail/batch/send", mail.SendBatch)
		// backup
		backup := admin.BackUpApi{}
		r.group.POST("/backup", backup.BackupPost)
		r.group.POST("/restore", backup.RestorePost)
	}
}

func (r *Routes) sessionRoute() {
	session := admin.SessionApi{}
	r.group.GET("/login", session.New)
	r.group.POST("/login", session.Create)
	r.group.GET("/session/destroy", session.Destroy)
	auth := admin.AuthApi{}
	r.engine.GET("/githubCallback", auth.GithubCallback)
	r.engine.GET("/auth/:authType", auth.AuthGet)
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
	r.group.POST("/password", password.Create)
	r.group.GET("/password/modify/:hash", password.Modify)
	r.group.POST("/password/send_notice", password.SendEmail)
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

func (r *Routes) user() {
	user := &admin.UserApi{}
	r.group.GET("/user/profile", user.Get)
	r.group.POST("/user/:id/profile", user.Update)
	r.group.GET("/user", user.Index)
	r.group.GET("/user/lock/:id", user.Lock)
}

func (r *Routes) auth() {
	auth := admin.AuthApi{}
	r.group.POST("/profile/email/bind", auth.BindEmail)
	r.group.POST("/profile/email/unbind", auth.UnbindEmail)
	r.group.POST("/profile/github/unbind", auth.UnbindGithub)
}

func (r *Routes) post() {
	post := new(admin.PostApi)
	r.group.GET("/posts", post.Index)
	r.group.GET("/posts/new", post.New)
	r.group.POST("/posts", post.Create)
	r.group.GET("/post/:id/edit", post.Edit)
	r.group.POST("/post/:id/update", post.Update)
	r.group.POST("/post/:id/publish", post.PostPublish)
	r.group.POST("/post/:id/delete", post.Delete)
}
func (r *Routes) page() {
	page := admin.PageApi{}
	r.group.GET("/page", page.Index)
	r.group.POST("/page", page.Create)
	r.group.GET("/page/new", page.New)
	r.group.GET("/page/edit/:id", page.Edit)
	r.group.GET("/page/get/:id", page.Get)
	r.group.POST("/page/update/:id", page.Update)
	r.group.DELETE("/page/delete/:id", page.Delete)
	r.group.POST("/page/publish/:id", page.Publish)
}

func (r *Routes) tag() {
	// tag
	tag := admin.TagApi{}
	r.group.GET("/tag/:format", tag.Index)
	r.group.POST("/tag", tag.Create)
	r.group.DELETE("/tag/:id", tag.Delete)
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
