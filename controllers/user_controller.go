package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"net/http"
	"strconv"
)
func ProfileGet(c *gin.Context) {
	sessionUser, exists := c.Get(CONTEXT_USER_KEY)
	if exists {
		c.HTML(http.StatusOK, "admin/profile.html", gin.H{
			"user":     sessionUser,
			"comments": models.MustListUnreadComment(),
		})
	}
}

func ProfileUpdate(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer WriteJSON(c, res)
	avatarUrl := c.PostForm("avatarUrl")
	nickName := c.PostForm("nickName")
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
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
	defer WriteJSON(c, res)
	email := c.PostForm("email")
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
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

func UnbindEmail(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer WriteJSON(c, res)
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
	user, ok := sessionUser.(*models.User)
	if !ok {
		res["message"] = "server interval error"
		return
	}
	if user.Email == "" {
		res["message"] = "email haven't bound"
		return
	}
	err = user.UpdateEmail("")
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func UnbindGithub(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer WriteJSON(c, res)
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
	user, ok := sessionUser.(*models.User)
	if !ok {
		res["message"] = "server interval error"
		return
	}
	if user.GithubLoginId == "" {
		res["message"] = "github haven't bound"
		return
	}
	user.GithubLoginId = ""
	err = user.UpdateGithubUserInfo()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func UserIndex(c *gin.Context) {
	users, _ := models.ListUsers()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/user.html", gin.H{
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
	defer WriteJSON(c, res)
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
