package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	service2 "github.com/gitlubtaotao/wblog/service"
)

type ITagRepository interface {
	ListTagByPostId(id uint) (tags []*models.Tag, err error)
	DeletePostTagByPostId(postId uint) error
	Create(tag models.Tag) (models.Tag, error)
	PostTagCreate(tag *models.PostTag) error
	Delete(id uint) error
	AllTag(attr map[string]interface{},columns []string)([]models.Tag,error)
}

type TagRepository struct {
	ctx     *gin.Context
	service service2.ITagService
}

func (t *TagRepository) AllTag(attr map[string]interface{}, columns []string)([]models.Tag,error) {
	return t.service.ListTag(0,0,attr,columns)
}

func (t *TagRepository) Delete(id uint) error {
	return t.service.Delete(id)
}

func (t *TagRepository) PostTagCreate(tag *models.PostTag) error {
	return t.service.PostTagCreate(tag)
}

func (t *TagRepository) Create(tag models.Tag) (models.Tag, error) {
	err := NewValidatorRepository(tag).HandlerError()
	if err != nil {
		return models.Tag{}, err
	}
	newTag, err := t.service.Create(tag)
	return newTag, err
}

func (t *TagRepository) DeletePostTagByPostId(postId uint) error {
	return t.service.DeletePostTagByPostId(postId)
}

func (t TagRepository) ListTagByPostId(id uint) (tags []*models.Tag, err error) {
	return t.service.ListTagByPostId(id)
}

func NewTagRepository(ctx *gin.Context) ITagRepository {
	return &TagRepository{ctx: ctx, service: service2.NewTagService()}
}
