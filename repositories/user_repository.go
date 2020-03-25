package repositories

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/service"
	"time"
)

type IUserRepository interface {
	Register() (err error)
	SignIn(account string, password string) (user *models.User, err error)
	FirstUserByEmail(email string) (models.User, error)
	UpdateUserAttr(attr map[string]interface{}) error
}

type UserRepository struct {
	userService service.IUserService
	Ctx         *gin.Context
}

func NewUserRepository(ctx *gin.Context) IUserRepository {
	return &UserRepository{Ctx: ctx, userService: service.NewUserService()}
}

//系统后台进行注册
func (u *UserRepository) Register() (err error) {
	var user *models.User
	err = u.Ctx.Bind(&user)
	if err != nil {
		return err
	}
	user.IsAdmin = true
	user.OutTime = time.Now().AddDate(0, 0, 4)
	password, err := encrypt.HashAndSalt(user.Password)
	if err != nil {
		return err
	}
	user.Password = password
	valid := ValidatorRepository{model: user}
	err = valid.HandlerError()
	if err != nil {
		return err
	}
	_ = u.userService.SetModel(user)
	return u.userService.Insert()
}

//用户进行登录
func (u *UserRepository) SignIn(account string, password string) (*models.User, error) {
	user := &models.User{}
	//邮箱登录
	if helpers.MatchEmail(account) {
		user.Email = account
	} else if helpers.MatchTelephone(account) {
		user.Telephone = account
	} else {
		return nil, errors.New("account is error ")
	}
	_ = u.userService.SetModel(user)
	err := u.userService.FirstUser()
	if err != nil {
		return nil, err
	}
	if err := encrypt.CompareHashSalt(user.Password, password); err != nil {
		return nil, err
	}
	return user, err
}

func (u *UserRepository) FirstUserByEmail(email string) (models.User, error) {
	return u.userService.FindUserByEmail(email)
}

func (u *UserRepository)UpdateUserAttr(attr map[string]interface{}) error  {
	return u.userService.UpdateUserAttr(attr)
}