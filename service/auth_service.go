package service

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
)

type IAuthService interface {
	FirstOrCreate(github *models.GithubUserInfo) (*models.GithubUserInfo, error)
}
type AuthService struct {
	Model *models.GithubUserInfo
}

func NewGitHubService() IAuthService {
	return &AuthService{}
}

func (g *AuthService)FirstOrCreate(github *models.GithubUserInfo) (*models.GithubUserInfo, error){
	err := database.DBCon.FirstOrCreate(github, "login = ?", github.Login).Error
	return github, err
}
