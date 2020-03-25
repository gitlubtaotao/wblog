package service

import (
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
)

type IUserService interface {
	FindUserByEmail(email string) (models.User, error)
	FirstUser() error
	Insert() error
	FirstOrCreate(user *models.User) (*models.User, error)
	GetUserByID(id interface{}) (*models.User, error)
	UpdateUser() (err error)
	UpdateUserAttr(attr map[string]interface{}) error
	GetModel() (*models.User, error)
	SetModel(user *models.User) error
	FindUserAll(attr map[string]interface{}) ([]*models.User, error)
}
type UserService struct {
	Model *models.User
}

func NewUserService() IUserService {
	return &UserService{}
}

//注册用户
func (u *UserService) Insert() error {
	return database.DBCon.Create(u.Model).Error
}

//查询user
func (u *UserService) FindUserByEmail(email string) (models.User, error) {
	var user models.User
	err := database.DBCon.Where("email = ?", email).First(&user).Error
	if err == nil {
		u.Model = &user
	}
	return user, err
}

func (u *UserService) FirstUser() error {
	return database.DBCon.First(&u.Model).Error
}

func (u *UserService) GetUserByID(id interface{}) (*models.User, error) {
	var user models.User
	err := database.DBCon.First(&user, id).Error
	return &user, err
}

func (u *UserService) FindUserAll(attr map[string]interface{}) (users []*models.User, err error) {
	err = database.DBCon.Where(attr).Find(&users).Error
	return users, err
}

//根据不同的字段进行根据
//TODO-taotao 根据用户
func (u *UserService) UpdateUser() (err error) {
	//更新所有的字段
	return database.DBCon.Save(&u.Model).Error
}

func (u *UserService) UpdateUserAttr(attr map[string]interface{}) error {
	return database.DBCon.Model(&u.Model).Update(attr).Error
}

func (u *UserService) FirstOrCreate(user *models.User) (*models.User, error) {
	err := database.DBCon.FirstOrCreate(user, "github_login_id = ?", user.GithubLoginId).Error
	return user, err
}

//get model value
func (u *UserService) GetModel() (*models.User, error) {
	return u.Model, nil
}

func (u *UserService) SetModel(user *models.User) error {
	u.Model = user
	return nil
}
