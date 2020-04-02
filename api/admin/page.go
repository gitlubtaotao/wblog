package admin

import (
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

type PageApi struct {
	*BaseApi
}

func (p *PageApi) PageGet(c *gin.Context) {
	repository := p.repository(c)
	id := c.Param("id")
	uints, _ := strconv.ParseUint(id, 10, 64)
	page, err := repository.FindPage(uint(uints))
	if err != nil || !page.IsPublished {
		p.Handle404(c)
		return
	}
	page.View++
	_ = page.UpdateView()
	c.HTML(http.StatusOK, "page/display.html", gin.H{
		"page": page,
	})
}

func (p *PageApi) New(c *gin.Context) {
	user, _ := p.CurrentUser(c)
	repository := p.repository(c)
	page, _ := repository.New()
	c.HTML(http.StatusOK, "page/edit.html",
		p.RenderComments(gin.H{"user": user,
			"page":   page,
			"action": "/admin/page"}))
}

func (p *PageApi) Create(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer p.WriteJSON(c, res)
	repository := p.repository(c)
	err = repository.GinCreate()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (p *PageApi) Edit(c *gin.Context) {
	repository := p.repository(c)
	user, _ := p.CurrentUser(c)
	page, err := repository.FindPage(p.stringToUnit(c))
	if err != nil {
		p.Handle404(c)
		return
	}
	id := strconv.FormatInt(int64(page.ID), 10)
	c.HTML(http.StatusOK, "page/edit.html",
		p.RenderComments(gin.H{"user": user, "page": page,
			"action": "/admin/page/update/" + id}))
}

func (p *PageApi) Update(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer p.WriteJSON(c, res)
	id := p.stringToUnit(c)
	repository := p.repository(c)
	err = repository.GinUpdate(id)
	if err != nil {
		res["message"] = err.Error()
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	res["succeed"] = true
}

func (p *PageApi) Publish(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer api.WriteJSON(c, res)
	repository := p.repository(c)
	err = repository.Publish(p.stringToUnit(c))
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (p *PageApi) Delete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer api.WriteJSON(c, res)
	id := p.stringToUnit(c)
	repository := p.repository(c)
	err = repository.Delete(id)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (p *PageApi) Index(c *gin.Context) {
	repository := p.repository(c)
	pages, _ := repository.ListAllPage(map[string]interface{}{})
	user, _ := p.CurrentUser(c)
	
	c.HTML(http.StatusOK, "page/index.html",
		p.RenderComments(gin.H{"pages": pages, "user": user}))
}

func (p *PageApi) Get(c *gin.Context) {

}

func (p *PageApi) repository(c *gin.Context) repositories.IPageRepository {
	return repositories.NewPageRepository(c)
}

func (p *PageApi) stringToUnit(c *gin.Context) uint {
	id := c.Param("id")
	units, _ := strconv.ParseUint(id, 10, 64)
	return uint(units)
}
