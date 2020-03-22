package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IController interface {
	Index(c *gin.Context)
	Get(c *gin.Context)
	New(c *gin.Context)
	Create(c *gin.Context)
	Edit(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type BaseController struct {
	Ctx *gin.Context
}

//初始化newBase
func NewBase(ctx *gin.Context) *BaseController {
	return &BaseController{Ctx: ctx}
}

//return json 格式
//输出json格式
func (b *BaseController) WriteJSON(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	ctx.JSON(http.StatusOK, h)
}
