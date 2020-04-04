package api

import (
	"errors"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
)

//auth 其他登录开发
type AuthApi struct {
	*BaseApi
	Auth repositories.IAuthRepository
}

//绑定不同的登录方式
func (a *AuthApi) AuthGet(c *gin.Context) {
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
func (a *AuthApi) GithubCallback(ctx *gin.Context) {
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
	sessionUser, exists := a.CurrentUser(ctx)
	if exists == nil { // 已登录
		a.bindUser(ctx, sessionUser, githubUser)
	} else {
		a.createUser(ctx, githubUser)
	}
}

//对github 进行解绑
func (a *AuthApi) UnbindGithub(ctx *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	repository := repositories.NewUserRepository(ctx)
	defer a.WriteJSON(ctx, res)
	currentUser, err := a.CurrentUser(ctx)
	if err != nil {
		res["message"] = "server interval error"
		return
	}
	if currentUser.GithubLoginId == "" {
		res["message"] = "github haven't bound"
		return
	}
	attr := map[string]interface{}{
		"GithubLoginId": "",
	}
	err = repository.Update(currentUser, attr)
	if err != nil {
		res["message"] = "Update User Info is Error "
		return
	}
	res["message"] = "UnBind user is successful"
	res["succeed"] = true
}

//对邮件进行解绑
func (a *AuthApi) UnbindEmail(ctx *gin.Context) {
	var res = gin.H{}
	repository := repositories.NewUserRepository(ctx)
	defer a.WriteJSON(ctx, res)
	currentUser, err := a.CurrentUser(ctx)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	if currentUser.Email == "" {
		res["message"] = "email haven't bound"
		return
	}
	_ = repository.SetUser(currentUser)
	err = repository.UpdateUserAttr(map[string]interface{}{"email": ""})
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

//绑定邮箱
func (a *AuthApi) BindEmail(ctx *gin.Context) {
	var res = gin.H{}
	defer a.WriteJSON(ctx, res)
	repository := repositories.NewUserRepository(ctx)
	email := ctx.PostForm("email")
	if email == "" {
		res["message"] = "email have not bound"
		return
	}
	user, err := a.CurrentUser(ctx)
	if err != nil {
		_ = seelog.Error(err)
		res["message"] = err.Error()
		return
	}
	if user.Email != "" {
		res["message"] = "email have bound"
		return
	}
	_, err = repository.FirstUserByEmail(email)
	//邮箱已经被注册过
	if err == nil {
		res["message"] = "email have be registered"
		return
	}
	_ = repository.SetUser(user)
	err = repository.UpdateUserAttr(map[string]interface{}{"email": email})
	if err != nil {
		res["message"] = "Bind email is error"
		return
	}
	res["message"] = "Bind email is successful"
	res["succeed"] = true
}

func (a *AuthApi) handlerError(ctx *gin.Context, err error) {
	_ = seelog.Error(err)
	a.HandleMessage(ctx, err.Error())
	return
}

//method: bind user
func (a *AuthApi) bindUser(ctx *gin.Context, sessionUser *models.User, githubUser *models.GithubUserInfo) {
	if _, err := a.Auth.GithubUserBing(sessionUser, githubUser); err != nil {
		a.handlerError(ctx, err)
		return
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/admin/user/profile")
		return
	}
}

func (a *AuthApi) createUser(ctx *gin.Context, githubUser *models.GithubUserInfo) {
	user, err := a.Auth.GithubUserCreate(githubUser)
	if err != nil {
		a.handlerError(ctx, err)
		return
	}
	if user.LockState {
		a.HandleMessage(ctx, "Your account have been locked.")
		return
	}
	key, err := encrypt.EnCryptData(string(user.ID),"admin")
	_ = a.OperationSession(ctx, SESSION_KEY, key)
	if err != nil {
		a.handlerError(ctx, err)
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, "/admin/index")
}
