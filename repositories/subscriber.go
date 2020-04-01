package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/service"
)

type ISubscriberRepository interface {
	AllListSubscriber(attr map[string]interface{}, columns []string) ([]*models.Subscriber, error)
	ListSubscriber(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Subscriber, error)
}

type SubscriberRepository struct {
	Ctx     *gin.Context
	service service.ISubscriberService
}

func (s SubscriberRepository) AllListSubscriber(attr map[string]interface{}, columns []string) ([]*models.Subscriber, error) {
	return s.service.AllListSubscriber(attr, columns)
}

func (s SubscriberRepository) ListSubscriber(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Subscriber, error) {
	return s.service.ListSubscriber(per,page,attr,columns)
}

func NewSubscriberRepository(ctx *gin.Context) ISubscriberRepository {
	return SubscriberRepository{Ctx: ctx,service: service.NewSubscriberService()}
}
