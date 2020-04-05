package tools

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	admin2 "github.com/gitlubtaotao/wblog/admin/api"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/api/admin"
	"github.com/gitlubtaotao/wblog/api/client"
	"github.com/gitlubtaotao/wblog/models"
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
	r.signInAndOut()
	r.visitorRouter()
	r.subscriberRouter()
	r.otherRouter()
	r.adminRouter()
	r.captchaRoute()
}

//captchaRoute
func (r *Routes) captchaRoute() {

}

func (r *Routes) indexInit() {
	r.engine.NoRoute(api.Handle404)
	r.engine.GET("/", client.Index)
	r.engine.GET("/index", client.Index)
	r.engine.GET("/rss", api.RssGet)
}



//登录和退出
func (r *Routes) signInAndOut() {
	session := admin.SessionApi{}
	auth := admin2.AuthApi{}
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
		comment := &api.CommentApi{}
		visitor.POST("/new_comment", comment.CommentPost)
		visitor.POST("/comment/:id/delete", comment.CommentDelete)
	}
}

//订阅者访问
func (r *Routes) subscriberRouter() {
	subscriber := api.SubscribeApi{}
	r.engine.GET("/subscribe", subscriber.SubscribeGet)
	r.engine.POST("/subscribe", subscriber.Subscribe)
	r.engine.GET("/active", api.ActiveSubscriber)
	r.engine.GET("/unsubscribe", api.UnSubscribe)
}

func (r *Routes) otherRouter() {
	r.engine.GET("/post/:id", admin2.PostGet)
	r.engine.GET("/tag/:tag", admin2.TagGet)
	r.engine.GET("/archives/:year/:month", api.ArchiveGet)
	//link := admin.LinkApi{}
	//r.engine.GET("/link/:id", link.LinkGet)
}

//后台路由
func (r *Routes) adminRouter() {
	authorized := r.engine.Group("/admin")
	authorized.Use(AdminScopeRequired())
	{
		
		
		
		
		
		// mail
		
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
