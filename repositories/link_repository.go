package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	service2 "github.com/gitlubtaotao/wblog/service"
)

type ILinkRepository interface {
	ListAllLink(columns []string)(links []*models.Link,err error)
	ListLink(per,page int,columns []string)(links []*models.Link,err error)
}

type LinkRepository struct {
	Ctx *gin.Context
	service service2.ILinkService
}

func NewLinkRepository(ctx *gin.Context) ILinkRepository {
	return &LinkRepository{Ctx: ctx, service: service2.NewLinkService()}
}

func (l *LinkRepository)  ListAllLink(columns []string)(links []*models.Link,err error) {
	return l.ListLink(0,0,columns)
}

func (l *LinkRepository) ListLink(per,page int, columns []string)(links []*models.Link,err error)  {
	return l.service.ListLink(per,page, columns)
}