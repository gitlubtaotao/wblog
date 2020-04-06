package tools

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	admin2 "github.com/gitlubtaotao/wblog/admin/api"
	"github.com/gitlubtaotao/wblog/api"
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
	r.visitorRouter()
	r.subscriberRouter()
	r.otherRouter()
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
	r.engine.GET("/tag/:tag", admin2.TagGet)
	r.engine.GET("/archives/:year/:month", api.ArchiveGet)
	//link := admin.LinkApi{}
	//r.engine.GET("/link/:id", link.LinkGet)
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


