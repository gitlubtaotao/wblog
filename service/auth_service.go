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

func NewAuthService() IAuthService {
	return &AuthService{}
}

func (g *AuthService)FirstOrCreate(github *models.GithubUserInfo) (*models.GithubUserInfo, error){
	err := database.DBCon.Where(models.GithubUserInfo{Login: github.Login}).Attrs(github).FirstOrCreate(&github).Error
	return github, err
}
