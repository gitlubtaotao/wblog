package services

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/models"
	"time"
)

type IUserService interface {
	Register() (err error)
	Insert() error
	SignIn(account string, password string) (user *models.User, err error)
	FindUserByEmail(email string) (user *models.User, err error)
	UpdateUser(user *models.User, fields string) (err error)
}
type UserService struct {
	Model   *models.User
	Context *gin.Context
}

func NewUserService(context *gin.Context) IUserService {
	return &UserService{Context: context}
}

//用户进行注册
func (u *UserService) Register() (err error) {
	err = u.Context.Bind(&u.Model)
	if err != nil {
		return err
	}
	u.Model.IsAdmin = true
	u.Model.OutTime = time.Now().AddDate(0, 0, 4)
	password, err := encrypt.HashAndSalt(u.Model.Password)
	if err != nil {
		return err
	}
	u.Model.Password = password
	valid := ValidatorService{model: u.Model}
	err = valid.HandlerError()
	if err != nil {
		return err
	}
	return u.Insert()
}

//注册用户
func (u *UserService) Insert() error {
	return database.DBCon.Create(u.Model).Error
}

//用户进行登录
func (u *UserService) SignIn(account string, password string) (*models.User, error) {
	u.Model = &models.User{}
	//邮箱登录
	if helpers.MatchEmail(account) {
		u.Model.Email = account
	} else if helpers.MatchTelephone(account) {
		u.Model.Telephone = account
	} else {
		return nil, errors.New("account is error ")
	}
	var user models.User
	err := database.DBCon.Where(u.Model).First(&user).Error
	if err != nil {
		return nil, err
	}
	if err := encrypt.CompareHashSalt(user.Password, password); err != nil {
		return nil, err
	}
	return &user, err
}

//查询user
func (u *UserService) FindUserByEmail(email string) (user *models.User, err error) {
	err = database.DBCon.Where("email = ?", email).First(&user).Error
	return user, err
}

//根据不同的字段进行根据
//TODO-taotao 根据用户
func (u *UserService) UpdateUser(user *models.User, fields string) (err error) {
	//更新所有的字段
	if fields == "" {
	
	}
	return
	
}
