package services

import "github.com/gitlubtaotao/wblog/repositories"

type IPostService interface {

}

type PostService struct {
	Repository *repositories.PostRepository
}
