package services

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
)

//
type IGitHubService interface {
	FirstOrCreate(github *models.GithubUserInfo) (*models.GithubUserInfo, error)
}
type GitHubService struct {
	Model *models.GithubUserInfo
}

func NewGitHubService() IGitHubService {
	return &GitHubService{}
}

func (g *GitHubService)FirstOrCreate(github *models.GithubUserInfo) (*models.GithubUserInfo, error){
	err := database.DBCon.FirstOrCreate(github, "login = ?", github.Login).Error
	return github, err
}
