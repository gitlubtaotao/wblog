package tools

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/api/admin"
	"github.com/gitlubtaotao/wblog/api/client"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"net/http"
)

type Routes struct {
	engine *gin.Engine
}

func NewRoutes(engine *gin.Engine) *Routes {
	return &Routes{engine: engine}
}

//初始化路由
func (r *Routes) InitRouter() {
	r.indexInit()
	r.signUp()
	r.signInAndOut()
	r.visitorRouter()
	r.subscriberRouter()
	r.otherRouter()
	r.adminRouter()
	r.captchaRoute()
}

//captchaRoute
func (r *Routes) captchaRoute() {
	controller := new(api.CaptchaController)
	r.engine.GET("/getCaptcha", controller.GetCaptcha)
	r.engine.GET("/verifyCaptcha", controller.VerifyCaptcha)
	r.engine.GET("/captcha/:source", controller.GetCaptchaPng)
	r.engine.GET("/captcha", controller.Captcha)
}

func (r *Routes) indexInit() {
	r.engine.NoRoute(api.Handle404)
	r.engine.GET("/", client.Index)
	r.engine.GET("/index", client.Index)
	r.engine.GET("/rss", api.RssGet)
}

func (r *Routes) signUp() {
	if system.GetConfiguration().SignupEnabled {
		r.engine.GET("/admin/signup", new(admin.RegisterApi).SignUpGet)
		r.engine.POST("admin/signup", new(admin.RegisterApi).SignUpPost)
	}
}

//登录和退出
func (r *Routes) signInAndOut() {
	session := admin.SessionApi{}
	auth := api.AuthApi{}
	r.engine.GET("/admin/signin", session.GetSignIn)
	r.engine.POST("/admin/signin", session.PostSignIn)
	r.engine.GET("/logout", session.LogoutGet)
	r.engine.GET("/githubCallback", auth.GithubCallback)
	r.engine.GET("/auth/:authType", auth.AuthGet)
	r.engine.GET("/password/index", session.GetPassword)
	r.engine.GET("/password/modifyPassword/:hash", session.ModifyPassword)
	r.engine.POST("/password/updatePassword", session.UpdatePassword)
	r.engine.POST("/passwords", session.UpdatePassword)
	r.engine.POST("/password/sendNotice", session.SendNotice)
}

func (r *Routes) visitorRouter() {
	visitor := r.engine.Group("/visitor")
	visitor.Use(authRequired())
	{
		visitor.POST("/new_comment", api.CommentPost)
		visitor.POST("/comment/:id/delete", api.CommentDelete)
	}
}

//订阅者访问
func (r *Routes) subscriberRouter() {
	r.engine.GET("/subscribe", api.SubscribeGet)
	r.engine.POST("/subscribe", api.Subscribe)
	r.engine.GET("/active", api.ActiveSubscriber)
	r.engine.GET("/unsubscribe", api.UnSubscribe)
}

func (r *Routes) otherRouter() {
	r.engine.GET("/page/:id", api.PageGet)
	r.engine.GET("/post/:id", admin.PostGet)
	r.engine.GET("/tag/:tag", api.TagGet)
	r.engine.GET("/archives/:year/:month", api.ArchiveGet)
	r.engine.GET("/link/:id", api.LinkGet)
}

//后台路由
func (r *Routes) adminRouter() {
	authorized := r.engine.Group("/admin")
	authorized.Use(AdminScopeRequired())
	{
		authorized.GET("", admin.Home)
		authorized.GET("/index", admin.Home)
		authorized.POST("/upload", api.Upload)
		authorized.GET("/page", api.PageIndex)
		authorized.GET("/new_page", api.PageNew)
		authorized.POST("/new_page", api.PageCreate)
		authorized.GET("/page/:id/edit", api.PageEdit)
		authorized.POST("/page/:id/edit", api.PageUpdate)
		authorized.POST("/page/:id/publish", api.PagePublish)
		authorized.POST("/page/:id/delete", api.PageDelete)
		
		// post
		post := new(admin.PostApi)
		authorized.GET("/posts", post.Index)
		authorized.GET("/posts/new", post.New)
		authorized.POST("/posts", post.Create)
		authorized.GET("/post/:id/edit", post.Edit)
		authorized.POST("/post/:id/edit", post.Update)
		authorized.POST("/post/:id/publish", admin.PostPublish)
		authorized.POST("/post/:id/delete", post.Delete)
		// tag
		authorized.POST("/new_tag", api.TagCreate)
		authorized.GET("/user", admin.UserIndex)
		authorized.POST("/user/:id/lock", admin.UserLock)
		// profile
		user := &admin.UserApi{}
		authorized.GET("/user/profile", user.ProfileGet)
		authorized.POST("/user/:id/profile", user.ProfileUpdate)
		auth := api.AuthApi{}
		authorized.POST("/profile/email/bind", auth.BindEmail)
		authorized.POST("/profile/email/unbind", auth.UnbindEmail)
		authorized.POST("/profile/github/unbind", auth.UnbindGithub)
		
		// subscriber
		authorized.GET("/subscriber", api.SubscriberIndex)
		authorized.POST("/subscriber", api.SubscriberPost)
		
		// link
		authorized.GET("/link", api.LinkIndex)
		authorized.POST("/new_link", api.LinkCreate)
		authorized.POST("/link/:id/edit", api.LinkUpdate)
		authorized.POST("/link/:id/delete", api.LinkDelete)
		// comment
		authorized.POST("/comment/:id", api.CommentRead)
		authorized.POST("/read_all", api.CommentReadAll)
		
		// backup
		authorized.POST("/backup", api.BackupPost)
		authorized.POST("/restore", api.RestorePost)
		
		// mail
		authorized.POST("/new_mail", api.SendMail)
		authorized.POST("/new_batchmail", api.SendBatchMail)
	}
	
}
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(api.CONTEXT_USER_KEY); user != nil {
			if _, ok := user.(*models.User); ok {
				c.Next()
				return
			}
		}
		seelog.Warnf("User not authorized to visit %s", c.Request.RequestURI)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "Forbidden!",
		})
		c.Abort()
	}
}

//AuthRequired grants access to authenticated users, requires SharedData middleware
func AdminScopeRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(api.CONTEXT_USER_KEY); user != nil {
			if u, ok := user.(*models.User); ok && u.IsAdmin {
				c.Next()
				return
			}
		}
		
		_ = seelog.Warnf("User not authorized to visit %s", c.Request.RequestURI)
		c.Redirect(http.StatusSeeOther, "/admin/signin")
		c.Abort()
	}
}
