package service

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/jinzhu/gorm"
)

type ILinkService interface {
	ListAllLink(columns []string) (links []*models.Link, err error)
	ListLink(per, page int, columns []string) (links []*models.Link, err error)
}

//
type LinkService struct {
	Models *models.Link
}

func NewLinkService() ILinkService  {
	return &LinkService{}
}
func (l *LinkService) ListAllLink(columns []string) (links []*models.Link, err error) {
	return l.ListLink(0, 0, columns)
}

func (l *LinkService) ListLink(per, page int, columns []string) (links []*models.Link, err error) {
	if per == 0 {
		per = system.GetConfiguration().PageSize
	}
	var temp *gorm.DB
	temp = database.DBCon.Find(&links)
	if page != 0 {
		temp = temp.Limit(per).Offset((page - 1) * per)
	}
	if len(columns) > 0 {
		temp = temp.Select(columns)
	}
	err = temp.Error
	return
}
