package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/service"
)

type IPageRepository interface {
	New() (models.Page, error)
	GinCreate() error
	GinUpdate(id uint) error
	Create(page models.Page) (models.Page, error)
	UpdateAttr(page *models.Page, attr map[string]interface{}) error
	Update(page *models.Page) error
	UpdatePage(page *models.Page) error
	Delete(id uint) error
	Publish(id uint) error
	FindPage(id uint) (models.Page, error)
	ListAllPage(attr map[string]interface{}) ([]*models.Page, error)
	ListPage(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Page, error)
}

type PageRepository struct {
	service service.IPageService
	Ctx     *gin.Context
}

func (p *PageRepository) Delete(id uint) error {
	page, err := p.FindPage(id)
	if err != nil {
		return err
	}
	return p.service.Delete(page)
}

func (p *PageRepository) Publish(id uint) error {
	page, err := p.FindPage(id)
	if err != nil {
		return err
	}
	err = p.UpdateAttr(&page, map[string]interface{}{
		"is_published": !page.IsPublished,
	})
	return err
}

func (p *PageRepository) GinUpdate(id uint) error {
	page, err := p.FindPage(id)
	if err != nil {
		return err
	}
	isPublished := p.Ctx.PostForm("isPublished")
	var attr = map[string]interface{}{
		"title":        p.Ctx.PostForm("title"),
		"body":         p.Ctx.PostForm("body"),
		"is_published": "on" == isPublished,
	}
	return p.UpdateAttr(&page, attr)
}

func (p *PageRepository) GinCreate() (err error) {
	var page models.Page
	err = p.Ctx.ShouldBindWith(&page, binding.Form)
	if err != nil {
		return
	}
	isPublished := p.Ctx.PostForm("isPublished")
	page.IsPublished = "on" == isPublished
	valid := ValidatorRepository{model: page}
	err = valid.HandlerError()
	if err != nil {
		return
	}
	_, err = p.Create(page)
	return
}

func (p *PageRepository) New() (models.Page, error) {
	return p.service.New()
}

func (p *PageRepository) Create(page models.Page) (models.Page, error) {
	return p.service.Create(page)
}

func (p *PageRepository) UpdateAttr(page *models.Page, attr map[string]interface{}) error {
	return p.service.UpdateAttr(page, attr)
}

func (p *PageRepository) Update(page *models.Page) error {
	return p.service.Update(page)
}

func (p *PageRepository) UpdatePage(page *models.Page) error {
	_ = p.service.SetModel(page)
	return p.service.UpdatePage()
}

func (p *PageRepository) FindPage(id uint) (models.Page, error) {
	return p.service.FindPage(id)
}

func (p *PageRepository) ListAllPage(attr map[string]interface{}) ([]*models.Page, error) {
	return p.ListPage(0, 0, attr, []string{})
}

func (p *PageRepository) ListPage(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Page, error) {
	return p.service.ListPage(per, page, attr, columns)
}

func NewPageRepository(ctx *gin.Context) IPageRepository {
	return &PageRepository{Ctx: ctx, service: service.NewPageService()}
}
