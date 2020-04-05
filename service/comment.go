package service

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
)

type ICommentService interface {
	CountComment() int
	ListUnreadComment() ([]*models.Comment, error)
	ListCommentByPostID(postId uint) ([]*models.Comment, error)
}

type CommentService struct {
	Model *models.Comment
}

func (c *CommentService) ListCommentByPostID(postId uint) ([]*models.Comment, error) {
	var comments []*models.Comment
	rows, err := database.DBCon.Raw("select c.*,u.github_login_id nick_name,u.avatar_url,u.github_url from comments c inner join users u on c.user_id = u.id where c.post_id = ? order by created_at desc", postId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		_ = database.DBCon.ScanRows(rows, &comment)
		comments = append(comments, &comment)
	}
	return comments, err
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
