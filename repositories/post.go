package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/service"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type IPostRepository interface {
	ListAll(attr map[string]interface{}, column []string) (posts []*models.Post, err error)
	Create() error
	Delete(id uint) error
	UpdateAttr(post *models.Post, attr map[string]interface{}) error
	Update(post *models.Post) error
	GetPostById(id uint, isTags bool) (*models.Post, error)
}

type PostRepository struct {
	ctx     *gin.Context
	service service.IPostService
}

func (p *PostRepository) UpdateAttr(post *models.Post, attr map[string]interface{}) error {
	return p.service.UpdateAttr(post, attr)
}

func (p *PostRepository) Update(post *models.Post) (err error) {
	return p.service.Update(post)
}
func (p *PostRepository) GetPostById(id uint, isTags bool) (*models.Post, error) {
	return p.service.GetPostById(id,isTags)
}
func (p *PostRepository) Delete(id uint) error {
	return p.service.Delete(id)
}

func (p *PostRepository) Create() error {
	tagService := service.NewTagService()
	var post models.Post
	err := p.ctx.ShouldBindWith(&post, binding.Form)
	if err != nil {
		return nil
	}
	isPublished := p.ctx.PostForm("isPublished")
	post.IsPublished = "on" == isPublished
	tags := p.ctx.PostForm("tags")
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		err := p.service.Create(&post)
		if err != nil {
			return err
		}
		if len(tags) > 0 {
			tagArr := strings.Split(tags, ",")
			for _, tag := range tagArr {
				tagId, err := strconv.ParseUint(tag, 10, 64)
				if err != nil {
					continue
				}
				pt := &models.PostTag{PostId: post.ID, TagId: uint(tagId),}
				err = tagService.PostTagCreate(pt)
				if err != nil {
					return nil
				}
			}
		}
		return nil
	})
	return err
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
