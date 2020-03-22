package services

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/models"
	"time"
)

type IUserService interface {
	Register() (err error)
	Insert() error
}
type UserService struct {
	model   *models.User
	context *gin.Context
}

func NewUserService(context *gin.Context) IUserService {
	return &UserService{context: context}
}

//用户进行注册
func (r *UserService) Register() (err error) {
	err = r.context.Bind(&r.model)
	if err != nil {
		return err
	}
	r.model.IsAdmin = true
	r.model.OutTime = time.Now().AddDate(0, 0, 4)
	r.model.Password = helpers.Md5(r.model.Email + r.model.Password)
	valid := ValidatorService{model: r.model}
	err = valid.HandlerError()
	if err != nil {
		return err
	}
	return r.Insert()
}

//注册用户
func (r *UserService) Insert() error {
	return database.DBCon.Create(r.model).Error
}
