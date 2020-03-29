package admin

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	"strconv"
)

type UserApi struct {
	repository repositories.IUserRepository
	*api.BaseApi
}

func (u *UserApi) ProfileGet(ctx *gin.Context) {
	repository := repositories.NewUserRepository(ctx)
	u.repository = repository
	tempUser, err := u.CurrentUser(ctx)
	if err == nil {
		err := u.repository.ReloadGithub(tempUser)
		if err != nil {
			_ = seelog.Critical(err)
		}
	}
	var url = "/admin/user/" + strconv.Itoa(int(tempUser.ID)) + "/profile"
	u.RenderHtml(ctx, "user/show.html",
		u.RenderComments(gin.H{"user": tempUser, "url": url,}))
}

func (u *UserApi) ProfileUpdate(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer u.WriteJSON(c, res)
	user, err := u.CurrentUser(c)
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
		password, err = encrypt.EnCryptData(password)
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
	if c.PostForm("telephone") != ""{
		attr["telephone"] = c.PostForm("telephone")
	}
	if password != "" {
		attr["password"] = password
	}
	if c.PostForm("secret_key") != ""{
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

func UserIndex(c *gin.Context) {
	users, _ := models.ListUsers()
	user, _ := c.Get(api.CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "user/index.html", gin.H{
		"users":    users,
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func UserLock(c *gin.Context) {
	var (
		err  error
		_id  uint64
		res  = gin.H{}
		user *models.User
	)
	defer api.WriteJSON(c, res)
	id := c.Param("id")
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	user, err = models.GetUser(uint(_id))
	if err != nil {
		res["message"] = err.Error()
		return
	}
	user.LockState = !user.LockState
	err = user.Lock()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
