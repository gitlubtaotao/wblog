package repositories

import (
	"database/sql"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/jinzhu/gorm"
	"strconv"
)

type IPostRepository interface {
}

type PostRepository struct {
	Object *models.Post
	DB     *gorm.DB
}

func NewPostRepository() *PostRepository {
	return &PostRepository{Object: &models.Post{}}
}

func (p *PostRepository) ListAll(tag string) (post []*models.Post, err error) {
	return p.listPost(tag, false, 0, 0)
}

func (p *PostRepository) ListPublishedPost(tag string, pageIndex, pageSize int) ([]*models.Post, error) {
	return p.listPost(tag, true, pageIndex, pageSize)
}

func (p *PostRepository) Delete() error {
	return models.DB.Delete(p.Object).Error
}

func (p *PostRepository) GetPostById(id string) (*models.Post, error) {
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	err = models.DB.First(&p.Object, "id=?", pid).Error
	return p.Object, err
}

//更新博文信息
func (p *PostRepository) Update() error {
	
	return models.DB.Model(&p.Object).Updates(map[string]interface{}{
		"title":        p.Object.Title,
		"body":         p.Object.Body,
		"is_published": p.Object.IsPublished,
	}).Error
}

func (p *PostRepository) listPost(tag string, published bool, pageIndex, pageSize int) ([]*models.Post, error) {
	var posts []*models.Post
	var err error
	if len(tag) <= 0 {
		posts, err = p.notTagsAllPosts(published, pageIndex, pageSize)
		return posts, err
	}
	tagId, _err := strconv.ParseUint(tag, 10, 64)
	if _err != nil {
		return nil, _err
	}
	var rows *sql.Rows
	if published {
		temp := p.selectTagsAndPost(tagId)
		if pageIndex > 0 {
			temp = temp.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
		}
		rows, err = temp.Where("is_published = ?", true).Rows()
	} else {
		rows, err = p.selectTagsAndPost(tagId).Rows()
	}
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

func (p *PostRepository) notTagsAllPosts(published bool, pageIndex, pageSize int) ([]*models.Post, error) {
	var (
		posts []*models.Post
		err   error
	)
	if published {
		temp := models.DB.Where("is_published = ?", true).Order("id desc")
		if pageIndex > 0 {
			temp = temp.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
		}
		err = temp.Find(&posts).Error
	} else {
		err = models.DB.Order("id desc").Find(&posts).Error
	}
	return posts, err
}

func (p *PostRepository) selectTagsAndPost(tagId uint64) *gorm.DB {
	db := models.DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? order by created_at desc", tagId)
	return db
}