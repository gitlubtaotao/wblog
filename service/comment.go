package service

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
)

type ICommentService interface {
	CountComment() int
	ListUnreadComment() ([]*models.Comment, error)
}

type CommentService struct {
	Model *models.Comment
}

func NewCommentService() ICommentService {
	return &CommentService{}
}

//统计评论的数量
func (c *CommentService) CountComment() int {
	var count int
	database.DBCon.Model(&models.Comment{}).Count(&count)
	return count
}

/*
	@method: 查询未读的评论
 */
func (c *CommentService) ListUnreadComment() ([]*models.Comment, error) {
	var comments []*models.Comment
	err := database.DBCon.Where("read_state = ?", false).Order("created_at desc").Find(&comments).Error
	return comments, err
}
