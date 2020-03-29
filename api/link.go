package api

import (
	"github.com/cihub/seelog"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
)

type LinkApi struct {
	*BaseApi
}

func (l *LinkApi) Index(ctx *gin.Context) {
	repository := repositories.NewLinkRepository(ctx)
	var columns []string
	links, err := repository.ListAllLink(columns)
	if err != nil {
		_ = seelog.Critical(err)
		l.HandleMessage(ctx, "service is inter error")
	}
	user, err := l.CurrentUser(ctx)
	if err != nil {
		_ = seelog.Critical(err)
		l.HandleMessage(ctx, err.Error())
		return
	}
	ctx.HTML(http.StatusOK, "admin/link.html", l.RenderComments(gin.H{
		"links": links,
		"user":  user,
	}))
}



func (l *LinkApi)LinkCreate(c *gin.Context) {
	var (
		err   error
		res   = gin.H{}
		_sort int64
	)
	defer WriteJSON(c, res)
	name := c.PostForm("name")
	url := c.PostForm("url")
	sort := c.PostForm("sort")
	if len(name) == 0 || len(url) == 0 {
		res["message"] = "error parameter"
		return
	}
	_sort, err = strconv.ParseInt(sort, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	link := &models.Link{
		Name: name,
		Url:  url,
		Sort: int(_sort),
	}
	err = link.Insert()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (l *LinkApi)LinkUpdate(c *gin.Context) {
	var (
		_id   uint64
		_sort int64
		err   error
		res   = gin.H{}
	)
	defer WriteJSON(c, res)
	id := c.Param("id")
	name := c.PostForm("name")
	url := c.PostForm("url")
	sort := c.PostForm("sort")
	if len(id) == 0 || len(name) == 0 || len(url) == 0 {
		res["message"] = "error parameter"
		return
	}
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	_sort, err = strconv.ParseInt(sort, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	link := &models.Link{
		Name: name,
		Url:  url,
		Sort: int(_sort),
	}
	link.ID = uint(_id)
	err = link.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (l *LinkApi)LinkGet(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.ParseInt(id, 10, 64)
	link, err := models.GetLinkById(uint(_id))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	link.View++
	link.Update()
	c.Redirect(http.StatusFound, link.Url)
}

func (l *LinkApi)LinkDelete(c *gin.Context) {
	var (
		err error
		_id uint64
		res = gin.H{}
	)
	defer WriteJSON(c, res)
	id := c.Param("id")
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	
	link := new(models.Link)
	link.ID = uint(_id)
	err = link.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
