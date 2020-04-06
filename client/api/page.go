package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/gitlubtaotao/wblog/system"
	"math"
	"net/http"
	"strconv"
	"sync"
)

type PageApi struct {
	UtilApi
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
	repository := repositories.NewPageRepository(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	page, err := repository.FindPage(uint(id))
	if err != nil {
		p.HandleMessage(ctx, err.Error())
		return
	}
	if !page.IsPublished {
		p.HandleMessage(ctx, "Post is not exist")
		return
	}
	var sy sync.WaitGroup
	sy.Add(1)
	go func(view int) {
		attr := map[string]interface{}{
			"view": view,
		}
		_ = repository.UpdateAttr(&page, attr)
		sy.Done()
	}(page.View + 1)
	sy.Wait()
	user, _ := p.ClientUser(ctx)
	ctx.HTML(http.StatusOK, "page/display.html", gin.H{
		"page": page,
		"user": user,
	})
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

func (p *PageApi) commentList(ctx *gin.Context, pageId uint) ([]*models.Comment, error) {
	repository := repositories.NewCommentRepository()
	comments, err := repository.ListCommentByPostID(pageId)
	return comments, err
}