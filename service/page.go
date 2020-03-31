package service

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/jinzhu/gorm"
)

type IPageService interface {
	New() (models.Page, error)
	Create(page models.Page) (models.Page, error)
	UpdateAttr(page *models.Page, attr map[string]interface{}) error
	Update(page *models.Page) error
	UpdatePage() error
	Delete(page models.Page) error
	FindPage(id uint) (models.Page, error)
	ListPage(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Page, error)
	SetModel(model *models.Page) error
	GetModel() (*models.Page, error)
}

//
type PageService struct {
	Model *models.Page
}

func (p *PageService) Delete(page models.Page) error {
	return database.DBCon.Delete(page).Error
}

func (p *PageService) SetModel(model *models.Page) error {
	p.Model = model
	return nil
}

func (p *PageService) GetModel() (*models.Page, error) {
	return p.Model, nil
}

func (p *PageService) New() (models.Page, error) {
	var page models.Page
	database.DBCon.NewRecord(&page)
	return page, nil
}

func (p *PageService) Create(page models.Page) (models.Page, error) {
	err := database.DBCon.Create(&page).Error
	return page, err
}

func (p *PageService) UpdateAttr(page *models.Page, attr map[string]interface{}) error {
	return database.DBCon.Model(&page).Update(attr).Error
}

func (p *PageService) Update(page *models.Page) error {
	return database.DBCon.Save(&page).Error
}

func (p *PageService) UpdatePage() error {
	return database.DBCon.Save(&p.Model).Error
}

func (p *PageService) FindPage(id uint) (models.Page, error) {
	var page models.Page
	err := database.DBCon.First(&page, id).Error
	return page, err
}

func (p *PageService) ListPage(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Page, error) {
	var pages []*models.Page
	if per == 0 {
		per = uint(system.GetConfiguration().PageSize)
	}
	var temp *gorm.DB
	temp = database.DBCon.Find(&pages)
	if page != 0 {
		temp = temp.Limit(per).Offset((page - 1) * per)
	}
	if len(attr) > 0 {
		temp = temp.Where(attr)
	}
	if len(columns) > 0 {
		temp = temp.Select(columns)
	}
	err := temp.Error
	return pages, err
}

func NewPageService() IPageService {
	return &PageService{Model: &models.Page{}}
}
