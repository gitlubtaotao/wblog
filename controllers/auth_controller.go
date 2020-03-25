package controllers

import (
	"errors"
	"github.com/cihub/seelog"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/helpers"
)

//auth 其他登录开发
type AuthController struct {
	*BaseController
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
		ctx.Abort()
		return
	}
	//通过code换取对于的token
	token, err := a.Auth.GitHubExchangeTokenByCode(code)
	if err != nil {
		_ = seelog.Critical(err)
		ctx.Redirect(http.StatusMovedPermanently, "/signin")
		return
	}
	githubUser, err := a.Auth.GithubUserInfoByAccessToken(token)
	if err != nil {
		_ = seelog.Error(err)
		c.Redirect(http.StatusMovedPermanently, "/signin")
		return
	}
	//	联合创建
	sessionUser, exists := ctx.Get(CONTEXT_USER_KEY)
	if exists { // 已登录
		user, err = a.Auth.GithubUserBing(sessionUser, githubUser)
		if err != nil {
			_ = seelog.Error(err)
			ctx.Redirect(http.StatusMovedPermanently, "/signin")
			return
		}
	} else {
		user, err = a.Auth.GithubUserCreate(githubUser)
		if err != nil {
			_ = seelog.Error(err)
			ctx.Redirect(http.StatusMovedPermanently, "/signin")
			return
		}
	}
	if user.LockState {
		err = errors.New("Your account have been locked.")
		a.HandleMessage(c, "Your account have been locked.")
		return
	}
	_ = a.OperationSession(ctx, SESSION_KEY, user.ID)
	ctx.Redirect(http.StatusMovedPermanently, "/admin/index")
}
