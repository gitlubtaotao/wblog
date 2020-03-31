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
	Create(link models.Link) (models.Link, error)
	Update(link *models.Link) error
	UpdateAttr(link *models.Link, attr map[string]interface{}) error
	UpdateLink() error
	Delete(id uint) error
	MaxSort() int
	FirstLink(id uint) (models.Link, error)
}

//
type LinkService struct {
	Models *models.Link
}

//删除
func (l *LinkService) Delete(id uint) error {
	return database.DBCon.Where("id = ?", id).Delete(&models.Link{}).Error
}

func (l *LinkService) Update(link *models.Link) error {
	return database.DBCon.Save(&link).Error
}

func (l *LinkService) UpdateAttr(link *models.Link, attr map[string]interface{}) error {
	return database.DBCon.Model(&link).Update(attr).Error
}

func (l *LinkService) UpdateLink() error {
	return database.DBCon.Save(&l.Models).Error
}

func (l *LinkService) FirstLink(id uint) (models.Link, error) {
	var link models.Link
	err := database.DBCon.First(&link, id).Error
	return link, err
}

func NewLinkService() ILinkService {
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

//查询当前最大的排序
func (l *LinkService) MaxSort() int {
	var count int
	database.DBCon.Model(&l.Models).Select("sum(sort) as max_sort").Scan(&count)
	return count
}

//创建链接
func (l *LinkService) Create(link models.Link) (models.Link, error) {
	err := database.DBCon.Create(&link).Error
	return link, err
}
