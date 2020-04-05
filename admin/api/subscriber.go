package admin

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
)

type SubscriberApi struct {
	*api.BaseApi
}

func (s *SubscriberApi) Index(ctx *gin.Context) {
	repository := repositories.NewSubscriberRepository(ctx)
	attr := map[string]interface{}{}
	columns := []string{"email", "verify_state", "subscribe_state", "created_at"}
	subscribers, err := repository.AllListSubscriber(attr, columns)
	if err != nil {
		_ = seelog.Critical(err)
		s.HandleMessage(ctx, err.Error())
		ctx.Abort()
	}
	user, _ := s.AdminUser(ctx)
	ctx.HTML(http.StatusOK, "subscriber/index.html",
		s.RenderComments(gin.H{
			"subscribers": subscribers,
			"user":        user,
		}))
}
