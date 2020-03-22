package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IController interface {
	Index(ctx *gin.Context)
	Get(ctx *gin.Context)
	New(ctx *gin.Context)
	Create(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type BaseController struct {

}

//return json 格式
//输出json格式
func (b *BaseController) WriteJSON(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	ctx.JSON(http.StatusOK, h)
}

