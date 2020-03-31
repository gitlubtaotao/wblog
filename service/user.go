package service

import (
	"fmt"
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/jinzhu/gorm"
)

type IUserService interface {
	FindUserByEmail(email string) (models.User, error)
	FirstUser() error
	Insert() error
	FirstOrCreate(user *models.User) (*models.User, error)
	GetUserByID(id int64) (*models.User, error)
	UpdateUser() (err error)
	Update(user *models.User, attr map[string]interface{}) error
	UpdateUserAttr(attr map[string]interface{}) error
	GetModel() (*models.User, error)
	SetModel(user *models.User) error
	FindUserAll(attr map[string]interface{}) ([]*models.User, error)
	ReloadGithub(user *models.User) error
	ListAdminUsers(per, page int, columns []string) ([]*models.User, error)
	Lock() error
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

func (u *UserService) GetUserByID(id int64) (*models.User, error) {
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
	err := database.DBCon.Where(models.User{GithubLoginId: user.GithubLoginId}).Attrs(user).FirstOrCreate(&user).Error
	return user, err
}

func (u *UserService) Update(user *models.User, attr map[string]interface{}) error {
	return database.DBCon.Model(&user).Update(attr).Error
}

//get model value
func (u *UserService) GetModel() (*models.User, error) {
	return u.Model, nil
}

func (u *UserService) SetModel(user *models.User) error {
	fmt.Println("sdssdsdsdsd")
	u.Model = user
	return nil
}

//reloadGithub: 加载github
func (u *UserService) ReloadGithub(user *models.User) error {
	return database.DBCon.Model(&user).Related(&user.GithubUserInfo).Error
}

func (u *UserService) ListAdminUsers(per, page int, columns []string) (users []*models.User, err error) {
	//取默认的分页数
	if per == 0 {
		per = system.GetConfiguration().PageSize
	}
	var temp *gorm.DB
	temp = database.DBCon.Find(&users, "is_admin = ?", true)
	if page != 0 {
		temp = temp.Limit(per).Offset((page - 1) * per)
	}
	if len(columns) == 0 {
		err = temp.Error
	} else {
		err = temp.Select(columns).Error
	}
	return
}

func (u *UserService) Lock() error {
	return database.DBCon.Model(&u.Model).Update(map[string]interface{}{
		"lock_state": u.Model.LockState,
	}).Error
}
