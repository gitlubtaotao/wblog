package service

import (
	"fmt"
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/jinzhu/gorm"
	"strconv"
)

type IPostService interface {
	PublishPost(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Post, error)
	NotPublishPost(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Post, error)
	TagsPost(per, page uint, attr map[string]interface{}, columns []string, tag string) ([]*models.Post, error)
	ListPost(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Post, error)
	AllListPost(attr map[string]interface{}, columns []string) ([]*models.Post, error)
	Create(post *models.Post) error
	Delete(id uint) error
	GetPostById(id uint, isTags bool) (*models.Post, error)
	Update(post *models.Post) error
	UpdateAttr(post *models.Post, attr map[string]interface{}) error
}

type PostService struct {
	Model *models.Post
}

func (p *PostService) Create(post *models.Post) error {
	return database.DBCon.Create(&post).Error
}

/*
 查询已经发布过的文章
*/
func (p *PostService) PublishPost(per, page uint, attr map[string]interface{}, columns []string) (posts []*models.Post, err error) {
	temp := p.tempListPost(per, page, attr)
	temp = temp.Where("is_published =?", true)
	if len(columns) > 0 {
		temp = temp.Select(columns)
	}
	err = temp.Scan(&posts).Error
	return
}

func (p *PostService) NotPublishPost(per, page uint, attr map[string]interface{}, columns []string) (posts []*models.Post, err error) {
	temp := p.tempListPost(per, page, attr)
	temp = temp.Where("is_published =?", false)
	if len(columns) > 0 {
		temp = temp.Select(columns)
	}
	err = temp.Scan(&posts).Error
	return
}

/*
@title: 查询包涵标签的博文
*/
func (p *PostService) TagsPost(per, page uint, attr map[string]interface{}, columns []string, tag string) (posts []*models.Post, err error) {
	var tagId uint64
	tagId, err = strconv.ParseUint(tag, 10, 64)
	if err != nil {
		return nil, err
	}
	temp := p.tempListPost(per, page, attr)
	temp = temp.Joins("inner join post_tags pt on post.id = pt.post_id where pt.tag_id = ?", tagId)
	if len(columns) > 0 {
		temp = temp.Select(columns)
	}
	rows, err := temp.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		_ = models.DB.ScanRows(rows, &post)
		posts = append(posts, &post)
	}
	return posts, err
}

/*
@title: 查询的博文
@description:
@auth   Xutaotao   2020.4.1
@param attr  需要过滤的条件
@param columns 需要查询的字段
@param per limit
@param page offset
@return *gorm.DB
*/
func (p *PostService) ListPost(per, page uint, attr map[string]interface{}, columns []string) (posts []*models.Post, err error) {
	temp := p.tempListPost(per, page, attr)
	if len(columns) > 0 {
		temp = temp.Select(columns)
	}
	err = temp.Scan(&posts).Error
	return
}

/*
@title: 查询的博文
@description:
@auth   Xutaotao   2020.4.1
@param attr  需要过滤的条件
@param columns 需要查询的字段
@return
*/
func (p *PostService) AllListPost(attr map[string]interface{}, columns []string) ([]*models.Post, error) {
	return p.ListPost(0, 0, attr, columns)
}

func (p *PostService) tempListPost(per, page uint, attr map[string]interface{}) (temp *gorm.DB) {
	if per == 0 {
		per = uint(system.GetConfiguration().PageSize)
	}
	temp = database.DBCon.Table("posts")
	if page != 0 {
		temp = temp.Limit(per).Offset((page - 1) * per)
	}
	if len(attr) > 0 {
		temp = temp.Where(attr)
	}
	return temp
}

func (p *PostService) Delete(id uint) error {
	return database.DBCon.Where("id = ?", id).Delete(&models.Post{}).Error
}

func (p *PostService) GetPostById(id uint, isTags bool) (*models.Post, error) {
	var temp *gorm.DB
	var post models.Post
	temp = database.DBCon
	fmt.Println(isTags)
	if isTags {
		temp = temp.Preload("Tags")
	}
	err := temp.Where("id=?", id).First(&post).Order("id desc").Error
	return &post, err
}

//更新博文信息
func (p *PostService) UpdatePost() error {
	return database.DBCon.Save(&p.Model).Error
}
func (p *PostService) Update(post *models.Post) error {
	return database.DBCon.Save(&post).Error
}

func (p *PostService) UpdateAttr(post *models.Post, attr map[string]interface{}) error {
	return database.DBCon.Model(&post).Update(attr).Error
}

func (p *PostService) selectTagsAndPost(tagId uint64) *gorm.DB {
	db := models.DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? order by created_at desc", tagId)
	return db
}

func NewPostService() IPostService {
	return &PostService{Model: &models.Post{}}
}
