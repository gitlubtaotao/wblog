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

func (l *LinkApi) LinkCreate(c *gin.Context) {
	repository := repositories.NewLinkRepository(c)
	var (
		err error
		res = gin.H{}
	)
	_, err = repository.Create()
	defer WriteJSON(c, res)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

//显示link info
func (l *LinkApi) Show(ctx *gin.Context) {
	repository := repositories.NewLinkRepository(ctx)
	var res = gin.H{}
	defer l.WriteJSON(ctx, res)
	link, err := repository.Show()
	if err != nil {
		_ = seelog.Error(err)
		res["message"] = "link not record"
		return
	}
	res["link"] = link
	res["succeed"] = true
}

func (l *LinkApi) LinkUpdate(c *gin.Context) {
	repository := repositories.NewLinkRepository(c)
	var (
		err error
		res = gin.H{}
	)
	defer WriteJSON(c, res)
	_, err = repository.UpdateAttr()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

//

func (l *LinkApi) LinkDelete(c *gin.Context) {
	reposition := repositories.NewLinkRepository(c)
	var (
		err error
		res = gin.H{}
		id  uint
	)
	defer WriteJSON(c, res)
	id, err = reposition.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["id"] = id
	res["succeed"] = true
}

func (l *LinkApi) LinkGet(c *gin.Context) {
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
