package admin

import (
	"github.com/gitlubtaotao/wblog/api"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
)

type PageApi struct {
	*api.BaseApi
}

func PageGet(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageById(id)
	if err != nil || !page.IsPublished {
		api.Handle404(c)
		return
	}
	page.View++
	page.UpdateView()
	c.HTML(http.StatusOK, "page/display.html", gin.H{
		"page": page,
	})
}

func (p *PageApi) New(c *gin.Context) {
	c.HTML(http.StatusOK, "page/new.html", nil)
}

func (p *PageApi) Create(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	page := &models.Page{
		Title:       title,
		Body:        body,
		IsPublished: published,
	}
	err := page.Insert()
	if err != nil {
		c.HTML(http.StatusOK, "page/new.html", gin.H{
			"message": err.Error(),
			"page":    page,
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/page")
}

func (p *PageApi) Edit(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageById(id)
	if err != nil {
		api.Handle404(c)
	}
	c.HTML(http.StatusOK, "page/modify.html", gin.H{
		"page": page,
	})
}

func (p *PageApi) Update(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	page := &models.Page{Title: title, Body: body, IsPublished: published}
	page.ID = uint(pid)
	err = page.Update()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/page")
}

func (p *PageApi) PagePublish(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer api.WriteJSON(c, res)
	id := c.Param("id")
	page, err := models.GetPageById(id)
	if err == nil {
		res["message"] = err.Error()
		return
	}
	page.IsPublished = !page.IsPublished
	err = page.Update()
	if err == nil {
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
	id := c.Param("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	page := &models.Page{}
	page.ID = uint(pid)
	err = page.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (p *PageApi) Index(c *gin.Context) {
	pages, _ := models.ListAllPage()
	user, _ := c.Get(api.CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/page.html", gin.H{
		"pages":    pages,
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func (p *PageApi) Get(c *gin.Context) {

}
