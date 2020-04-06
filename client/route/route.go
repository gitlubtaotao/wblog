package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	client2 "github.com/gitlubtaotao/wblog/client/api"
)

type IRoute interface {
	Register()
}
type Route struct {
	engine *gin.Engine
	group  *gin.RouterGroup
}

func NewRoute(engine *gin.Engine) IRoute {
	return &Route{engine: engine, group: engine.Group("/client")}
}
func (r *Route) Register() {
	r.engine.NoRoute(new(api.BaseApi).Handle404)
	home := client2.HomeApi{}
	r.engine.GET("/", home.Index)
	post := r.engine.Group("/post")
	{
		post.GET("/:id", new(client2.PostApi).Show)
		post.GET("", new(client2.PostApi).Index)
	}
	page := r.engine.Group("/page")
	{
		
		page.GET("/:id", new(client2.PageApi).Show)
		page.GET("", new(client2.PageApi).Index)
	}
	
	r.engine.GET("/rss", new(client2.RssApi).RssGet)
	
}
