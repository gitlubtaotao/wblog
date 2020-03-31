package admin

import (
	"github.com/cihub/seelog"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type LinkApi struct {
	*api.BaseApi
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

func (l *LinkApi) Create(c *gin.Context) {
	repository := repositories.NewLinkRepository(c)
	var (
		err error
		res = gin.H{}
	)
	_, err = repository.Create()
	defer api.WriteJSON(c, res)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

//显示link info
func (l *LinkApi) Get(ctx *gin.Context) {
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

func (l *LinkApi) Update(c *gin.Context) {
	repository := repositories.NewLinkRepository(c)
	var (
		err error
		res = gin.H{}
	)
	defer api.WriteJSON(c, res)
	_, err = repository.UpdateAttr()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

//

func (l *LinkApi) Delete(c *gin.Context) {
	reposition := repositories.NewLinkRepository(c)
	var (
		err error
		res = gin.H{}
		id  uint
	)
	defer api.WriteJSON(c, res)
	id, err = reposition.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["id"] = id
	res["succeed"] = true
}

func (l *LinkApi) Edit(ctx *gin.Context) {

}

func (l *LinkApi) New(ctx *gin.Context) {

}

//func (l *LinkApi) LinkGet(c *gin.Context) {
////	id := c.Param("id")
////	_id, _ := strconv.ParseInt(id, 10, 64)
////	link, err := models.GetLinkById(uint(_id))
////	if err != nil {
////		c.AbortWithStatus(http.StatusInternalServerError)
////		return
////	}
////	link.View++
////	link.Update()
////	c.Redirect(http.StatusFound, link.Url)
////}
