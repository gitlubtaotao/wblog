package crouter

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/controllers"
	"github.com/gitlubtaotao/wblog/controllers/admin"
	"github.com/gitlubtaotao/wblog/controllers/client"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"net/http"
)

//初始化路由
func InitRouter(engine *gin.Engine) {
	indexInit(engine)
	signUp(engine)
	signInAndOut(engine)
	visitorRouter(engine)
	subscriberRouter(engine)
	otherRouter(engine)
	adminRouter(engine)
	captchaRoute(engine)
}

//captchaRoute
func captchaRoute(engine *gin.Engine) {
	controller := new(controllers.CaptchaController)
	engine.GET("/getCaptcha", controller.GetCaptcha)
	engine.GET("/verifyCaptcha", controller.VerifyCaptcha)
	engine.GET("/captcha/:source", controller.GetCaptchaPng)
	engine.GET("/captcha", controller.Captcha)
}

func indexInit(engine *gin.Engine) {
	engine.NoRoute(controllers.Handle404)
	engine.GET("/", client.Index)
	engine.GET("/index", client.Index)
	engine.GET("/rss", controllers.RssGet)
}

func signUp(engine *gin.Engine) {
	if system.GetConfiguration().SignupEnabled {
		engine.GET("/admin/signup", new(admin.RegisterController).SignUpGet)
		engine.POST("admin/signup", new(admin.RegisterController).SignUpPost)
	}
}

//登录和退出
func signInAndOut(engine *gin.Engine) {
	session := admin.SessionController{}
	engine.GET("/admin/signin", session.GetSignIn)
	engine.POST("/admin/signin", session.PostSignIn)
	engine.GET("/logout", controllers.LogoutGet)
	engine.GET("/oauth2callback", controllers.Oauth2Callback)
	engine.GET("/auth/:authType", controllers.AuthGet)
	engine.GET("/password/index", session.GetPassword)
	engine.GET("/password/modifyPassword/:hash",session.ModifyPassword)
	engine.POST("/password/updatePassword",session.UpdatePassword)
	engine.POST("/passwords", session.UpdatePassword)
	engine.POST("/password/sendNotice", session.SendNotice)
}

func visitorRouter(engine *gin.Engine) {
	visitor := engine.Group("/visitor")
	visitor.Use(authRequired())
	{
		visitor.POST("/new_comment", controllers.CommentPost)
		visitor.POST("/comment/:id/delete", controllers.CommentDelete)
	}
}

//订阅者访问
func subscriberRouter(engine *gin.Engine) {
	engine.GET("/subscribe", controllers.SubscribeGet)
	engine.POST("/subscribe", controllers.Subscribe)
	engine.GET("/active", controllers.ActiveSubscriber)
	engine.GET("/unsubscribe", controllers.UnSubscribe)
}

func otherRouter(engine *gin.Engine) {
	engine.GET("/page/:id", controllers.PageGet)
	engine.GET("/post/:id", admin.PostGet)
	engine.GET("/tag/:tag", controllers.TagGet)
	engine.GET("/archives/:year/:month", controllers.ArchiveGet)
	engine.GET("/link/:id", controllers.LinkGet)
}

//后台路由
func adminRouter(engine *gin.Engine) {
	authorized := engine.Group("/admin")
	authorized.Use(AdminScopeRequired())
	{
		authorized.GET("", admin.Home)
		authorized.GET("/index", admin.Home)
		authorized.POST("/upload", controllers.Upload)
		authorized.GET("/page", controllers.PageIndex)
		authorized.GET("/new_page", controllers.PageNew)
		authorized.POST("/new_page", controllers.PageCreate)
		authorized.GET("/page/:id/edit", controllers.PageEdit)
		authorized.POST("/page/:id/edit", controllers.PageUpdate)
		authorized.POST("/page/:id/publish", controllers.PagePublish)
		authorized.POST("/page/:id/delete", controllers.PageDelete)
		
		// post
		post := new(admin.PostController)
		authorized.GET("/posts", post.Index)
		authorized.GET("/posts/new", post.New)
		authorized.POST("/posts", post.Create)
		authorized.GET("/post/:id/edit", post.Edit)
		authorized.POST("/post/:id/edit", post.Update)
		authorized.POST("/post/:id/publish", admin.PostPublish)
		authorized.POST("/post/:id/delete", post.Delete)
		// tag
		authorized.POST("/new_tag", controllers.TagCreate)
		authorized.GET("/user", controllers.UserIndex)
		authorized.POST("/user/:id/lock", controllers.UserLock)
		// profile
		authorized.GET("/profile", controllers.ProfileGet)
		authorized.POST("/profile", controllers.ProfileUpdate)
		authorized.POST("/profile/email/bind", controllers.BindEmail)
		authorized.POST("/profile/email/unbind", controllers.UnbindEmail)
		authorized.POST("/profile/github/unbind", controllers.UnbindGithub)
		
		// subscriber
		authorized.GET("/subscriber", controllers.SubscriberIndex)
		authorized.POST("/subscriber", controllers.SubscriberPost)
		
		// link
		authorized.GET("/link", controllers.LinkIndex)
		authorized.POST("/new_link", controllers.LinkCreate)
		authorized.POST("/link/:id/edit", controllers.LinkUpdate)
		authorized.POST("/link/:id/delete", controllers.LinkDelete)
		// comment
		authorized.POST("/comment/:id", controllers.CommentRead)
		authorized.POST("/read_all", controllers.CommentReadAll)
		
		// backup
		authorized.POST("/backup", controllers.BackupPost)
		authorized.POST("/restore", controllers.RestorePost)
		
		// mail
		authorized.POST("/new_mail", controllers.SendMail)
		authorized.POST("/new_batchmail", controllers.SendBatchMail)
	}
	
}
func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(controllers.CONTEXT_USER_KEY); user != nil {
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
		if user, _ := c.Get(controllers.CONTEXT_USER_KEY); user != nil {
			if u, ok := user.(*models.User); ok && u.IsAdmin {
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
