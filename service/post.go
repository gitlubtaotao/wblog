package service

import (
	"database/sql"
	"fmt"
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/jinzhu/gorm"
	"strconv"
)

type IPostService interface {
	PublishPost(per, page uint, attr map[string]interface{}, columns []string, isTag bool) ([]*models.Post, error)
	NotPublishPost(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Post, error)
	TagsPost(per, page uint, attr map[string]interface{}, columns []string, tag string) ([]*models.Post, error)
	ListPost(per, page uint, attr map[string]interface{}, columns []string) ([]*models.Post, error)
	AllListPost(attr map[string]interface{}, columns []string) ([]*models.Post, error)
	Create(post *models.Post) error
	Delete(id uint) error
	GetPostById(id uint, isTags bool) (*models.Post, error)
	Update(post *models.Post) error
	UpdateAttr(post *models.Post, attr map[string]interface{}) error
	CountPostByTag(tag uint) (count int, err error)
	CountPost(attr map[string]interface{}) (count int, err error)
	ListMaxReadPost(column []string) ([]*models.Post, error)
	ListMaxCommentPost(columns []string)([]*models.Post,error)
}

type PostService struct {
	Model *models.Post
}

func (p *PostService) ListMaxCommentPost(columns []string) (posts []*models.Post, err error) {
	var rows *sql.Rows
	rows, err = database.DBCon.Raw("select p.title,p.id,c.total comment_total from posts p inner join (select post_id,count(*) total from comments group by post_id) c on p.id = c.post_id order by c.total desc limit 5").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		_ = database.DBCon.ScanRows(rows, &post)
		posts = append(posts, &post)
	}
	return
}
func (p *PostService) ListMaxReadPost(column []string) (posts []*models.Post, err error) {
	temp := database.DBCon
	if len(column) > 0 {
		temp = temp.Select(column)
	}
	err = temp.Where("is_published = ?", true).Order("view desc").Limit(5).Find(&posts).Error
	return
}

func (p *PostService) CountPost(attr map[string]interface{}) (count int, err error) {
	temp := database.DBCon.Model(&models.Post{})
	if len(attr) > 0 {
		temp = temp.Where(attr)
	}
	err = temp.Count(&count).Error
	return
}

func (p *PostService) CountPostByTag(tag uint) (count int, err error) {
	if tag > 0 {
		err = database.DBCon.Raw("select count(*) from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? and p.is_published = ?", tag, true).Row().Scan(&count)
	} else {
		err = database.DBCon.Raw("select count(*) from posts p where p.is_published = ?", true).Row().Scan(&count)
	}
	return
}

func (p *PostService) Create(post *models.Post) error {
	return database.DBCon.Create(&post).Error
}

/*
 查询已经发布过的文章
*/
func (p *PostService) PublishPost(per, page uint, attr map[string]interface{}, columns []string, isTag bool) (posts []*models.Post, err error) {
	if per == 0 {
		per = uint(system.GetConfiguration().PageSize)
	}
	temp := database.DBCon
	if page != 0 {
		temp = temp.Limit(per).Offset((page - 1) * per)
	}
	if len(attr) > 0 {
		temp = temp.Where(attr)
	}
	if isTag {
		temp = temp.Preload("Tags")
	}
	temp = temp.Where("is_published =?", true)
	if len(columns) > 0 {
		temp = temp.Select(columns)
	}
	err = temp.Find(&posts).Error
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
