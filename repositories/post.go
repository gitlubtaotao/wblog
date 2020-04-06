package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/service"
)

type IPostRepository interface {
	ListAll(attr map[string]interface{}, column []string) (posts []*models.Post, err error)
	Create(post *models.Post) error
	Delete(id uint) error
	UpdateAttr(post *models.Post, attr map[string]interface{}) error
	Update(post *models.Post) error
	GetPostById(id uint, isTags bool) (*models.Post, error)
	PublishPost(per, page uint, attr map[string]interface{}, columns []string, isTags bool) ([]*models.Post, error)
	CountPostByTag(tagId uint) (count int, err error)
	CountPost(attr map[string]interface{}) (count int, err error)
	ListMaxReadPost(column []string) ([]*models.Post, error)
	ListMaxCommentPost(columns []string) ([]*models.Post, error)
	TagsPost(per, page uint, attr map[string]interface{}, columns []string, tag string) (posts []*models.Post, err error)
}

type PostRepository struct {
	ctx     *gin.Context
	service service.IPostService
}

func (p *PostRepository) TagsPost(per, page uint, attr map[string]interface{}, columns []string, tag string) (posts []*models.Post, err error) {
	posts, err = p.service.TagsPost(per, page, attr, columns, tag)
	if err != nil {
		return
	}
	var postId []uint
	for _, post := range posts {
		postId = append(postId, post.ID)
	}
	//获取对应的tags
	
	return
}

func (p *PostRepository) ListMaxCommentPost(columns []string) ([]*models.Post, error) {
	return p.service.ListMaxCommentPost(columns)
}

func (p *PostRepository) ListMaxReadPost(column []string) ([]*models.Post, error) {
	return p.service.ListMaxReadPost(column)
}

func (p *PostRepository) CountPost(attr map[string]interface{}) (count int, err error) {
	return p.service.CountPost(attr)
}

func (p *PostRepository) CountPostByTag(tagId uint) (count int, err error) {
	return p.service.CountPostByTag(tagId)
}

func (p *PostRepository) PublishPost(per, page uint, attr map[string]interface{}, columns []string, isTag bool) ([]*models.Post, error) {
	return p.service.PublishPost(per, page, attr, columns, isTag)
}

func (p *PostRepository) UpdateAttr(post *models.Post, attr map[string]interface{}) error {
	return p.service.UpdateAttr(post, attr)
}

func (p *PostRepository) Update(post *models.Post) (err error) {
	return p.service.Update(post)
}
func (p *PostRepository) GetPostById(id uint, isTags bool) (*models.Post, error) {
	return p.service.GetPostById(id, isTags)
}
func (p *PostRepository) Delete(id uint) error {
	return p.service.Delete(id)
}

func (p *PostRepository) Create(post *models.Post) error {
	err := NewValidatorRepository(post).HandlerError()
	if err != nil {
		return err
	}
	return p.service.Create(post)
}

/*
@title: 获取所有的博客
*/
func (p PostRepository) ListAll(attr map[string]interface{}, column []string) (posts []*models.Post, err error) {
	return p.service.AllListPost(attr, column)
}

func NewPostRepository(ctx *gin.Context) IPostRepository {
	return &PostRepository{ctx: ctx, service: service.NewPostService()}
}
