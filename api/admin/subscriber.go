package admin

//TODO-TAO 查询前端订阅者，需要将admin和client进行拆分，分布启用两个不同的服务

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/models"
	"net/http"
)

type SubscriberApi struct {
	*api.BaseApi
}


func (s *SubscriberApi) Index(ctx *gin.Context) {
	subscribers, _ := models.ListSubscriber(false)
	user, _ := s.CurrentUser(ctx)
	ctx.HTML(http.StatusOK, "subscriber/index.html",
		s.RenderComments(gin.H{
			"subscribers": subscribers,
			"user":        user,
		}))
}

