package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"strconv"
)

// table posts
type Post struct {
	BaseModel
	Title        string     `validate:"required"` // title
	Body         string     `validate:"required"` // body
	View         int        // view count
	IsPublished  bool       // published or not
	Tags         []*Tag     `gorm:"-"` // tags of post
	Comments     []*Comment `gorm:"-"` // comments of post
	CommentTotal int        `gorm:"-"` // count of comment
	Keyword      string     `gorm:"size:255;not null" validate:"required"`
}

//获取所有的博文
func ListAllPost(tag string) ([]*Post, error) {
	return listPost(tag, false, 0, 0)
}

//获取已发布的博文
func ListPublishedPost(tag string, pageIndex, pageSize int) ([]*Post, error) {
	return listPost(tag, true, pageIndex, pageSize)
}

//博文查询method
func listPost(tag string, published bool, pageIndex, pageSize int) ([]*Post, error) {
	var posts []*Post
	var err error
	//是否查询指定的标签
	if len(tag) <= 0 {
		posts, err = notTagsAllPosts(published, pageIndex, pageSize)
		return posts, err
	}
	tagId, _err := strconv.ParseUint(tag, 10, 64)
	if _err != nil {
		return nil, _err
	}
	var rows *sql.Rows
	if published {
		temp := selectTagsAndPost(tagId)
		if pageIndex > 0 {
			temp = temp.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
		}
		rows, err = temp.Where("is_published = ?", true).Rows()
	} else {
		rows, err = selectTagsAndPost(tagId).Rows()
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		_ = DB.ScanRows(rows, &post)
		posts = append(posts, &post)
	}
	return posts, err
}

func selectTagsAndPost(tagId uint64) *gorm.DB {
	db := DB.Raw("select p.* from posts p inner join post_tags pt on p.id = pt.post_id where pt.tag_id = ? order by created_at desc", tagId)
	return db
}

//没有指定标签查询所有的博文
func notTagsAllPosts(published bool, pageIndex, pageSize int) ([]*Post, error) {
	var (
		posts []*Post
		err   error
	)
	if published {
		temp := DB.Where("is_published = ?", true).Order("id desc")
		if pageIndex > 0 {
			temp = temp.Limit(pageSize).Offset((pageIndex - 1) * pageSize)
		}
		err = temp.Find(&posts).Error
	} else {
		err = DB.Order("id desc").Find(&posts).Error
	}
	return posts, err
}
