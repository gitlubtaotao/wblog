package controllers

import (
	"github.com/gin-contrib/sessions"
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

//render html
func (b *BaseController) RenderHtml(ctx *gin.Context, path string, h gin.H) {
	ctx.HTML(http.StatusOK, path, h)
	ctx.Abort()
}

//操作对于的session
func (b *BaseController) OperationSession(ctx *gin.Context, Key string, value interface{}) error {
	session := sessions.Default(ctx)
	session.Delete(Key)
	session.Set(SESSION_GITHUB_STATE, value)
	return session.Save()
}

func (b *BaseController) GetSessionValue(ctx *gin.Context, key string, isDelete bool) (value interface{}, err error) {
	session := sessions.Default(ctx)
	value = session.Get(key)
	if isDelete {
		session.Delete(key)
		_ = session.Save()
	}
	return value, nil
}

func (b *BaseController)HandleMessage(c *gin.Context, message string) {
	c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
		"message": message,
	})
}