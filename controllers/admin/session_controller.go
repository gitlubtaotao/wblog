package admin

import (
	"github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/controllers"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/services"
	"net/http"
)

type SessionController struct {
	*controllers.BaseController
}

func (s *SessionController) GetSignIn(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"title": "Wblog | Log in",
	})
}

//用户进行登录
func (s *SessionController) PostSignIn(ctx *gin.Context) {
	var (
		res      = gin.H{}
		remember bool
	)
	defer s.WriteJSON(ctx, res)
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	if account == "" || password == "" {
		res["message"] = "username or password cannot be null"
		return
	}
	if ctx.PostForm("checkbox") != "" {
		remember = true
	}
	service := services.NewUserService(ctx)
	user, err := service.SignIn(account, password)
	if err != nil {
		_ = seelog.Error(err)
		res["message"] = "Your account not exist"
		return
	}
	if user.LockState {
		res["message"] = "Your account have been locked"
		return
	}
	session := sessions.Default(ctx)
	session.Clear()
	key, err := encrypt.EnCryptData(string(user.ID))
	if err != nil {
		_ = seelog.Error(err)
		res["message"] = "Your account not exist"
		return
	}
	session.Set(controllers.SESSION_KEY, key)
	_ = session.Save()
	res["succeed"] = true
	res["remember"] = remember
	res["contentType"] = ctx.ContentType()
	//进行session id 加密
	
}

func (s *SessionController) AuthGet(c *gin.Context) {

}
