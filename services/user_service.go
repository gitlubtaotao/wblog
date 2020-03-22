package services

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/models"
	"time"
)

type IUserService interface {
	Register() (err error)
	Insert() error
	SignIn(account string, password string) (user *models.User, err error)
}
type UserService struct {
	Model   *models.User
	Context *gin.Context
}

func NewUserService(context *gin.Context) IUserService {
	return &UserService{Context: context}
}

//用户进行注册
func (r *UserService) Register() (err error) {
	err = r.Context.Bind(&r.Model)
	if err != nil {
		return err
	}
	r.Model.IsAdmin = true
	r.Model.OutTime = time.Now().AddDate(0, 0, 4)
	r.Model.Password = helpers.Md5(r.Model.Password)
	valid := ValidatorService{model: r.Model}
	err = valid.HandlerError()
	if err != nil {
		return err
	}
	return r.Insert()
}

//注册用户
func (r *UserService) Insert() error {
	return database.DBCon.Create(r.Model).Error
}

//用户进行登录
func (r *UserService) SignIn(account string, password string) (*models.User, error) {
	md5Password := helpers.Md5(password)
	r.Model = &models.User{Password: md5Password}
	//邮箱登录
	if helpers.MatchEmail(account) {
		r.Model.Email = account
	} else if helpers.MatchTelephone(account) {
		r.Model.Telephone = account
	} else {
		return nil, errors.New("account is error ")
	}
	var user models.User
	err := database.DBCon.Where(r.Model).First(&user).Error
	return &user, err
}
