package api

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/gitlubtaotao/wblog/system"
	"net/http"
)

type IBaseApi interface {
	Index(ctx *gin.Context)
	Get(ctx *gin.Context)
	New(ctx *gin.Context)
	Create(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type BaseApi struct {
}

//return json 格式
//输出json格式
func (b *BaseApi) WriteJSON(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	ctx.JSON(http.StatusOK, h)
}

//render html
func (b *BaseApi) RenderHtml(ctx *gin.Context, path string, h gin.H) {
	ctx.HTML(http.StatusOK, path, h)
	ctx.Abort()
}

//handler render html comments

func (b *BaseApi) RenderComments(h gin.H) gin.H {
	repository := repositories.NewCommentRepository()
	h["comments"], _ = repository.MustListUnreadComment()
	h["message"] = ""
	return h
}

//操作对于的session
func (b *BaseApi) OperationSession(ctx *gin.Context, Key string, value interface{}) error {
	session := sessions.Default(ctx)
	session.Delete(Key)
	session.Set(Key, value)
	return session.Save()
}

func (b *BaseApi) GetSessionValue(ctx *gin.Context, key string, isDelete bool) (value interface{}, err error) {
	session := sessions.Default(ctx)
	value = session.Get(key)
	if isDelete {
		session.Delete(key)
		_ = session.Save()
	}
	return value, nil
}

func (b *BaseApi) Handle404(c *gin.Context) {
	HandleMessage(c, "Sorry,I lost myself!")
}

func (b *BaseApi) HandleMessage(c *gin.Context, message string) {
	c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
		"message": message,
	})
	c.Abort()
}

func (b *BaseApi) CurrentUser(c *gin.Context) (*models.User, error) {
	
	sessionUser, exists := c.Get(CONTEXT_USER_KEY)
	if !exists {
		return nil, errors.New("current user is not exist")
	}
	
	user, ok := sessionUser.(*models.User)
	if !ok {
		return nil, errors.New("server interval error")
	}
	return user, nil
}

//发送邮件
func (b *BaseApi) SendMailHtml(to, subject, body string) error {
	repository := repositories.NewMailRepository(subject, body, "html")
	return repository.SendToMail(to)
}

//系统默认推送方式
func (b *BaseApi) DefaultNoticeMailHtml(subject, body string) error {
	repository := repositories.NewMailRepository(subject, body, "html")
	return repository.SystemDefaultNotify()
}

func (b *BaseApi) AdminUser(ctx *gin.Context) (*models.User, error) {
	sessionUser, exists := ctx.Get(system.GetConfiguration().AdminUser)
	if !exists {
		return nil, errors.New("current user is not exist")
	}
	user, ok := sessionUser.(*models.User)
	if !ok {
		return nil, errors.New("server interval error")
	}
	return user, nil
}
