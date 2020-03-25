package api

import (
	"errors"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
)

//auth 其他登录开发
type AuthController struct {
	*BaseApi
	Auth repositories.IAuthRepository
}

//绑定不同的登录方式
func (a *AuthController) AuthGet(c *gin.Context) {
	a.Auth = repositories.NewAuthRepository()
	authType := c.Param("authType")
	uuid := helpers.UUID()
	_ = a.OperationSession(c, SESSION_GITHUB_STATE, uuid)
	authUrl := "/signin"
	switch authType {
	case "github":
		authUrl = a.Auth.GitHubAccessURL(uuid)
	case "weibo":
	case "qq":
	case "wechat":
	case "oschina":
	default:
	}
	c.Redirect(http.StatusFound, authUrl)
}

//github callback
func (a *AuthController) GithubCallback(ctx *gin.Context) {
	var user *models.User
	a.Auth = repositories.NewAuthRepository()
	code := ctx.Query("code")
	state := ctx.Query("state")
	systemState, _ := a.GetSessionValue(ctx, SESSION_GITHUB_STATE, true)
	//验证失败
	if len(state) == 0 || state != systemState {
		a.handlerError(ctx, errors.New("state is error "))
		return
	}
	//通过code换取对于的token
	token, err := a.Auth.GitHubExchangeTokenByCode(code)
	if err != nil {
		a.handlerError(ctx, err)
		return
	}
	githubUser, err := a.Auth.GithubUserInfoByAccessToken(token)
	if err != nil {
		a.handlerError(ctx, err)
		return
	}
	//	联合创建
	sessionUser, exists := ctx.Get(CONTEXT_USER_KEY)
	fmt.Println("ssss22222")
	if exists { // 已登录
		user, err = a.Auth.GithubUserBing(sessionUser, githubUser)
		if err != nil {
			a.handlerError(ctx, err)
			return
		}
	} else {
		user, err = a.Auth.GithubUserCreate(githubUser)
		fmt.Println(user,"sssss")
		if err != nil {
			a.handlerError(ctx, err)
			return
		}
	}
	if user.LockState {
		err = errors.New("Your account have been locked.")
		a.HandleMessage(ctx, "Your account have been locked.")
		return
	}
	key, err := encrypt.EnCryptData(string(user.ID))
	_ = a.OperationSession(ctx, SESSION_KEY, key)
	if err != nil {
		a.handlerError(ctx, err)
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, "/admin/index")
}

func (a *AuthController) handlerError(ctx *gin.Context, err error) {
	_ = seelog.Error(err)
	ctx.Redirect(http.StatusMovedPermanently, "/admin/signin")
	ctx.Abort()
}
