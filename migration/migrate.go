package migration

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
)

func Migrate() {
	database.DBCon.AutoMigrate(&models.Page{}, &models.Post{}, &models.Tag{},
		&models.PostTag{}, &models.User{}, &models.Comment{}, &models.Subscriber{}, &models.Link{}, &models.SmmsFile{})
	database.DBCon.Model(&models.PostTag{}).AddUniqueIndex("uk_post_tag", "post_id", "tag_id")
}
