package database

import (
"github.com/cihub/seelog"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
"github.com/jinzhu/gorm"
)

var DBCon *gorm.DB

/*初始化数据库*/
func InitDB() {
	//db, err := gorm.Open("sqlite3", system.GetConfiguration().DSN)
	db, err := gorm.Open("mysql", system.GetConfiguration().DSN)
	if err != nil {
		seelog.Critical("err open databases", err)
		panic(err)
		return
	}
	DBCon = db
	models.DB = db
	db.LogMode(true)
	db.DB().SetMaxIdleConns(1000)
	db.DB().SetMaxOpenConns(2000)
	if err = db.DB().Ping(); err != nil {
		seelog.Critical("err open databases", err)
		panic(err)
	}
}
