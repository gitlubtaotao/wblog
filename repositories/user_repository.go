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
	FirstUser() (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
	FirstUserByEmail(email string) (models.User, error)
	Update(user *models.User, attr map[string]interface{}) error
	UpdateUserAttr(attr map[string]interface{}) error
	ListAllAdminUsers(columns []string) ([]*models.User, error)
	ListAdminUsers(per, page int, columns []string) ([]*models.User, error)
	ReloadGithub(user *models.User) (err error)
	Lock(user *models.User) (err error)
	GetUser() (*models.User, error)
	SetUser(user *models.User) error
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

func (u *UserRepository) FirstUser() (user *models.User, err error) {
	err = u.userService.FirstUser()
	user, _ = u.userService.GetModel()
	return
}

/*
 
 */
func (u *UserRepository) GetUserByID(id int64) (user *models.User, err error) {
	return u.userService.GetUserByID(id)
}

func (u *UserRepository) Update(user *models.User, attr map[string]interface{}) error {
	return u.userService.Update(user, attr)
}
func (u *UserRepository) UpdateUserAttr(attr map[string]interface{}) error {
	return u.userService.UpdateUserAttr(attr)
}

func (u *UserRepository) Lock(user *models.User) error {
	_ = u.userService.SetModel(user)
	return u.userService.Lock()
}
func (u *UserRepository) GetUser() (user *models.User, err error) {
	return u.userService.GetModel()
}
func (u *UserRepository) SetUser(user *models.User) error {
	return u.userService.SetModel(user)
}

func (u *UserRepository) ReloadGithub(user *models.User) (err error) {
	return u.userService.ReloadGithub(user)
}

func (u *UserRepository) ListAllAdminUsers(columns []string) ([]*models.User, error) {
	return u.ListAdminUsers(0, 0, columns)
}

func (u *UserRepository) ListAdminUsers(per, page int, columns []string) ([]*models.User, error) {
	return u.userService.ListAdminUsers(per, page, columns)
}
