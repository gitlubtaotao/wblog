package admin

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
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
	user, _ := ctx.Get(api.CONTEXT_USER_KEY)
	repository := repositories.NewUserRepository(ctx)
	u.repository = repository
	tempUser, ok := user.(*models.User)
	if ok {
		err := u.repository.ReloadGithub(tempUser)
		if err != nil {
			_ = seelog.Critical(err)
		}
	}
	u.RenderHtml(ctx, "user/show.html", u.RenderComments(gin.H{"user": tempUser,}))
}

func ProfileUpdate(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer api.WriteJSON(c, res)
	avatarUrl := c.PostForm("avatarUrl")
	nickName := c.PostForm("nickName")
	sessionUser, _ := c.Get(api.CONTEXT_USER_KEY)
	user, ok := sessionUser.(*models.User)
	if !ok {
		res["message"] = "server interval error"
		return
	}
	err = user.UpdateProfile(avatarUrl, nickName)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
	res["user"] = models.User{AvatarUrl: avatarUrl, NickName: nickName}
}

func BindEmail(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer api.WriteJSON(c, res)
	email := c.PostForm("email")
	sessionUser, _ := c.Get(api.CONTEXT_USER_KEY)
	user, ok := sessionUser.(*models.User)
	if !ok {
		res["message"] = "server interval error"
		return
	}
	if len(user.Email) > 0 {
		res["message"] = "email have bound"
		return
	}
	_, err = models.GetUserByUsername(email)
	if err == nil {
		res["message"] = "email have be registered"
		return
	}
	err = user.UpdateEmail(email)
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
