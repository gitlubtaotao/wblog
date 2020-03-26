package repositories

import (
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/service"
)

type ICommentRepository interface {
	MustListUnreadComment() ([]*models.Comment, error)
	CountComment() int
}

type CommentRepository struct {
	service service.ICommentService
}

func NewCommentRepository() ICommentRepository {
	return &CommentRepository{service: service.NewCommentService()}
}

func (c *CommentRepository) MustListUnreadComment() ([]*models.Comment, error) {
	return c.service.ListUnreadComment()
}

func (c *CommentRepository) CountComment() int {
	return c.service.CountComment()
}
