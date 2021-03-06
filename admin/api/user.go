package admin

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
	"strconv"
)

type UserApi struct {
	repository repositories.IUserRepository
	*api.BaseApi
}

func (u *UserApi) Get(ctx *gin.Context) {
	var (
		tempUser *models.User
		err      error
	)
	repository := repositories.NewUserRepository(ctx)
	u.repository = repository
	userId := ctx.Query("id")
	if userId == "" {
		tempUser, err = u.AdminUser(ctx)
	} else {
		id, _ := strconv.ParseInt(userId, 10, 64)
		tempUser, err = repository.GetUserByID(id)
	}
	if err != nil {
		_ = seelog.Critical(err)
		u.HandleMessage(ctx, "service inter is error")
		return
	}
	err = u.repository.ReloadGithub(tempUser)
	if err != nil {
		_ = seelog.Critical(err)
	}
	var url = "/admin/user/" + strconv.Itoa(int(tempUser.ID)) + "/profile"
	u.RenderHtml(
		ctx,
		"user/show.html",
		u.RenderComments(gin.H{
			"user":  tempUser,
			"url":   url,
			"token": csrf.GetToken(ctx),
		}))
}

func (u *UserApi) Update(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer u.WriteJSON(c, res)
	user, err := u.AdminUser(c)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")
	if password != "" {
		if password != confirmPassword {
			res["message"] = "password is error"
			return
		}
		password, err = encrypt.EnCryptData(password, "admin")
	}
	if err != nil {
		res["message"] = "pssword is error"
		return
	}
	var attr map[string]interface{}
	attr = make(map[string]interface{}, 1)
	if c.PostForm("avatar_url") != "" {
		attr["avatar_url"] = c.PostForm("avatar_url")
	}
	if c.PostForm("nick_name") != "" {
		attr["nick_name"] = c.PostForm("nick_name")
	}
	if c.PostForm("telephone") != "" {
		attr["telephone"] = c.PostForm("telephone")
	}
	if password != "" {
		attr["password"] = password
	}
	if c.PostForm("secret_key") != "" {
		attr["secret_key"] = c.PostForm("secret_key")
	}
	repository := repositories.NewUserRepository(c)
	err = repository.Update(user, attr)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (u *UserApi) Index(c *gin.Context) {
	repository := repositories.NewUserRepository(c)
	columns := []string{"telephone", "email", "nick_name", "github_login_id",
		"created_at", "id", "is_admin", "avatar_url", "secret_key"}
	users, _ := repository.ListAllAdminUsers(columns)
	user, _ := u.AdminUser(c)
	c.HTML(http.StatusOK, "user/index.html",
		u.RenderComments(gin.H{"user": user, "users": users,}))
}

func (u *UserApi) Lock(c *gin.Context) {
	var (
		err  error
		_id  uint64
		res  = gin.H{}
		user *models.User
	)
	defer u.WriteJSON(c, res)
	id := c.Param("id")
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	repository := repositories.NewUserRepository(c)
	user, err = repository.GetUserByID(int64(_id))
	if err != nil {
		res["message"] = err.Error()
		return
	}
	user.LockState = !user.LockState
	err = repository.Lock(user)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
