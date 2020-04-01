package service

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/jinzhu/gorm"
)

type ISubscriberService interface {
	AllListSubscriber(attr map[string]interface{}, columns []string) ([]*models.Subscriber, error)
	ListSubscriber(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Subscriber, error)
}

type SubscriberService struct {
	Model *models.Subscriber
}

/*
@title  查询所有的订阅者
@description   根据不同的查询条件和查询列，查询订阅者
@auth   Xutaotao   2020.4.1
@param attr  需要过滤的条件
@param columns 需要查询的字段
@return
*/
func (s SubscriberService) AllListSubscriber(attr map[string]interface{}, columns []string) ([]*models.Subscriber, error) {
	return s.ListSubscriber(0, 0, attr, columns)
}

/*
@title  查询所有的订阅者
@description   根据不同的查询条件和查询列，查询订阅者
@auth   Xutaotao   2020.4.1
@param per limit
@param page offset
@param attr  需要过滤的条件
@param columns 需要查询的字段
@return
*/
func (s SubscriberService) ListSubscriber(per, page uint, attr map[string]interface{}, columns []string) (subscribers []*models.Subscriber, err error) {
	if per == 0 {
		per = uint(system.GetConfiguration().PageSize)
	}
	var temp *gorm.DB
	temp = database.DBCon.Find(&subscribers)
	if page != 0 {
		temp = temp.Limit(per).Offset((page - 1) * per)
	}
	if len(attr) > 0 {
		temp = temp.Where(attr)
	}
	if len(columns) > 0 {
		temp = temp.Select(columns)
	}
	err = temp.Error
	return subscribers, err
}

func NewSubscriberService() ISubscriberService {
	return SubscriberService{Model: &models.Subscriber{}}
}
