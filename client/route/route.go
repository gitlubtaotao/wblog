package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	client2 "github.com/gitlubtaotao/wblog/client/api"
)

type  IRoute interface {
	Register()
}
type Route struct {
	engine *gin.Engine
	group *gin.RouterGroup
}

func NewRoute(engine *gin.Engine) IRoute {
	return &Route{engine:engine,group: engine.Group("/client")}
}
func (r *Route) Register()  {
	r.engine.NoRoute(new(api.BaseApi).Handle404)
	home := client2.HomeApi{}
	r.engine.GET("/", home.Index)
}
