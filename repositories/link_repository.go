package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gitlubtaotao/wblog/models"
	service2 "github.com/gitlubtaotao/wblog/service"
	"strconv"
)

type ILinkRepository interface {
	ListAllLink(columns []string) (links []*models.Link, err error)
	ListLink(per, page int, columns []string) (links []*models.Link, err error)
	Create() (link models.Link, err error)
	UpdateAttr() (models.Link, error)
	Update(link *models.Link) error
	Delete() (uint, error)
	Show() (link models.Link, err error)
}

type LinkRepository struct {
	Ctx     *gin.Context
	service service2.ILinkService
}

func (l *LinkRepository) Delete() (uint, error) {
	id := l.getId()
	return id, l.service.Delete(id)
}

func (l *LinkRepository) Update(link *models.Link) error {
	return l.service.Update(link)
}

//更新link
func (l *LinkRepository) UpdateAttr() (link models.Link, err error) {
	id := l.getId()
	link, err = l.service.FirstLink(id)
	if err != nil {
		return
	}
	var attr = map[string]interface{}{
		"name": l.Ctx.PostForm("name"),
		"url":  l.Ctx.PostForm("url"),
	}
	err = l.service.UpdateAttr(&link, attr)
	return
}

func NewLinkRepository(ctx *gin.Context) ILinkRepository {
	return &LinkRepository{Ctx: ctx, service: service2.NewLinkService()}
}

//查询link没有进行分页
func (l *LinkRepository) ListAllLink(columns []string) (links []*models.Link, err error) {
	return l.ListLink(0, 0, columns)
}

//查询link
func (l *LinkRepository) ListLink(per, page int, columns []string) (links []*models.Link, err error) {
	return l.service.ListLink(per, page, columns)
}

//创建链接
func (l *LinkRepository) Create() (link models.Link, err error) {
	err = l.Ctx.ShouldBindWith(&link, binding.Form)
	if err != nil {
		return
	}
	//获取系统当前的sort
	link.Sort = l.service.MaxSort() + 1
	valid := ValidatorRepository{model: link}
	err = valid.HandlerError()
	if err != nil {
		return
	}
	return l.service.Create(link)
}

func (l *LinkRepository) Show() (link models.Link, err error) {
	id := l.getId()
	return l.service.FirstLink(id)
}

//获取当前最大的排序字段
func (l *LinkRepository) maxSort() int {
	return l.service.MaxSort()
}

func (l *LinkRepository) getId() uint {
	id, _ := strconv.ParseUint(l.Ctx.Param("id"), 10, 64)
	return uint(id)
}
