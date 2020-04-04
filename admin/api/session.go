package admin

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/gitlubtaotao/wblog/system"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

type SessionApi struct {
	*api.BaseApi
}

func (s *SessionApi) New(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "session/new.html", gin.H{
		"title": "Wblog | Log in",
		"token": csrf.GetToken(ctx),
	})
}

func (s *SessionApi) Create(ctx *gin.Context) {
	var (
		res      = gin.H{}
	)
	repository := repositories.NewUserRepository(ctx)
	defer s.WriteJSON(ctx, res)
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	if account == "" || password == "" {
		res["message"] = "username or password cannot be null"
		return
	}
	if ctx.PostForm("checkbox") != "" {
	}
	user, err := repository.SignIn(account, password)
	fmt.Println(err)
	if err != nil {
		_ = seelog.Error(err)
		res["message"] = "Your account not exist"
		return
	}
	if user.LockState {
		res["message"] = "Your account have been locked"
		return
	}
	key, err := encrypt.EnCryptData(string(user.ID),"admin")
	fmt.Println(err)
	if err != nil {
		_ = seelog.Error(err)
		res["message"] = "Your account not exist"
		return
	}
	
	_ = s.OperationSession(ctx, system.GetConfiguration().AdminSessionKey, key)
	res["succeed"] = true
}

func (s *SessionApi) Destroy(ctx *gin.Context) {

}
