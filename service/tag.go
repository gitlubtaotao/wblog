package service

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
)

type ITagService interface {
	Create(tag models.Tag) (models.Tag, error)
	Delete(id uint) error
	PostTagCreate(tag *models.PostTag) error
	ListTagByPostId(postId uint) ([]*models.Tag, error)
	DeletePostTagByPostId(postId uint) error
	ListTag(per, page uint, attr map[string]interface{}, columns []string) ([]models.Tag, error)
}

//
type TagService struct {
	Model *models.Tag
}

func (t *TagService) ListTag(per, page uint, attr map[string]interface{}, columns []string) (tags []models.Tag, err error) {
	var temp = database.DBCon
	if len(columns) > 0 {
		temp = database.DBCon.Select(columns)
	}
	temp = temp.Find(&tags)
	if per == 0 {
		per = uint(system.GetConfiguration().PageSize)
	}
	if page > 0 {
		temp = temp.Limit(per).Offset((page - 1) * per)
	}
	if len(attr) > 0 {
		temp = temp.Where(attr)
	}
	err = temp.Error
	return tags, err
}

func (t *TagService) Delete(id uint) error {
	return database.DBCon.Delete(&t.Model, "id=?", id).Error
}

func (t *TagService) DeletePostTagByPostId(postId uint) error {
	return database.DBCon.Delete(&models.PostTag{}, "post_id = ?", postId).Error
}

/*
 @title: list tags by post id
*/
func (t *TagService) ListTagByPostId(postId uint) (tags []*models.Tag, err error) {
	rows, err := database.DBCon.Table("tags").Joins("inner joins post_tags on post_tags.tag_id=tags.id and post_tags.post_id= ?",
		postId).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tag models.Tag
		_ = database.DBCon.ScanRows(rows, &tag)
		tags = append(tags, &tag)
	}
	return tags, nil
}

func (t *TagService) PostTagCreate(tag *models.PostTag) error {
	return database.DBCon.Create(&tag).Error
}

//
func (t *TagService) Create(tag models.Tag) (models.Tag, error) {
	err := database.DBCon.Create(&tag).Error
	return tag, err
}

func NewTagService() ITagService {
	return &TagService{Model: &models.Tag{}}
}
