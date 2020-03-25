package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	oauth "github.com/alimoeeny/gooauth2"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/services"
	"github.com/gitlubtaotao/wblog/system"
	"io/ioutil"
	"net/http"
)

type IAuthRepository interface {
	GitHubAccessURL(uuid string) (url string)
	GitHubExchangeTokenByCode(code string) (accessToken string, err error)
	GithubUserInfoByAccessToken(token string) (*models.GithubUserInfo, error)
	GithubUserCreate(github *models.GithubUserInfo) (*models.User, error)
	GithubUserBing(sessionUser interface{}, githubUser *models.GithubUserInfo)(user *models.User, err error)
}

type AuthRepository struct {
	gitHubService services.IGitHubService
	Ctx           *gin.Context
}

func NewAuthRepository() IAuthRepository {
	return &AuthRepository{}
}

func (a *AuthRepository) GitHubAccessURL(uuid string) (url string) {
	return fmt.Sprintf(system.GetConfiguration().GithubAuthUrl, system.GetConfiguration().GithubClientId, uuid)
}

func (a *AuthRepository) GitHubExchangeTokenByCode(code string) (accessToken string, err error) {
	config := system.GetConfiguration()
	transport := &oauth.Transport{Config: &oauth.Config{
		ClientId:     config.GithubClientId,
		ClientSecret: config.GithubClientSecret,
		RedirectURL:  config.GithubRedirectURL,
		TokenURL:     config.GithubTokenUrl,
		Scope:        config.GithubScope,
	}}
	token, err := transport.Exchange(code)
	if err != nil {
		return
	}
	accessToken = token.AccessToken
	tokenCache := oauth.CacheFile("./request.token")
	if err := tokenCache.PutToken(token); err != nil {
		return accessToken, err
	}
	return accessToken, nil
}

//通过github获取用户信息
func (a *AuthRepository) GithubUserInfoByAccessToken(token string) (*models.GithubUserInfo, error) {
	var (
		resp *http.Response
		body []byte
		err  error
	)
	resp, err = http.Get(fmt.Sprintf("https://api.github.com/user?access_token=%s", token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var userInfo models.GithubUserInfo
	err = json.Unmarshal(body, &userInfo)
	return &userInfo, err
}

func (a *AuthRepository) GithubUserBing(sessionUser interface{}, githubUser *models.GithubUserInfo) (user *models.User, err error) {
	service := services.NewUserService()
	user, _ = sessionUser.(*models.User)
	_ = service.SetModel(user)
	var attr map[string]interface{}
	attr = make(map[string]interface{}, 1)
	attr["github_login_id"] = githubUser.Login
	attr["id"] = user.ID
	//没有进行绑定，可以进行绑定操作
	_, err = service.FindUserAll(attr)
	if err != nil {
		user.GithubLoginId = githubUser.Login
		user.AvatarUrl = githubUser.AvatarURL
		user.GithubUrl = githubUser.HTMLURL
		_ = service.SetModel(user)
		return user, service.UpdateUser()
	} else {
		return nil, errors.New("this github loginId has bound another account.")
	}
}

//通过github auth进行用户的创建
func (a *AuthRepository) GithubUserCreate(github *models.GithubUserInfo) (user *models.User, err error) {
	service := services.NewUserService()
	user = &models.User{
		GithubLoginId: github.Login,
		AvatarUrl:     github.AvatarURL,
		GithubUrl:     github.HTMLURL,
	}
	_ = service.SetModel(user)
	user, err = service.FirstOrCreate(user)
	if err != nil {
		return nil, err
	}
	github, err = a.gitHubService.FirstOrCreate(github)
	return
}
