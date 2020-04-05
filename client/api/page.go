package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/gitlubtaotao/wblog/system"
	"math"
	"net/http"
)

type PageApi struct {
	*UtilApi
}

func (p *PageApi) Index(ctx *gin.Context) {
	var (
		pageIndex, _ = p.PageIndex(ctx)
		pageSize     = system.GetConfiguration().PageSize
	)
	pages, err := p.listPublishPage(ctx)
	if err != nil {
		return
	}
	total, err := p.TotalByCategory(ctx)
	if err != nil {
		return
	}
	ctx.HTML(http.StatusOK, "page/index.html", gin.H{
		"pages":     pages,
		"pageIndex": pageIndex,
		"totalPage": int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func (p *PageApi) Show(ctx *gin.Context) {

}

func (p *PageApi) TotalByCategory(ctx *gin.Context) (total int, err error) {
	repository := repositories.NewPageRepository(ctx)
	total, err = repository.TotalPage(map[string]interface{}{
		"is_published": true,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	return
}
func (p *PageApi) listPublishPage(ctx *gin.Context) (pages []*models.Page, err error) {
	repository := repositories.NewPageRepository(ctx)
	page, _ := p.PageIndex(ctx)
	pages, err = repository.PublishPage(0, uint(page), map[string]interface{}{}, []string{})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	return
}
