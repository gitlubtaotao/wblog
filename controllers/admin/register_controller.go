package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/controllers"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
)

type RegisterController struct {
	*controllers.BaseController
	repository repositories.IUserRepository
}

//注册页面
func (r *RegisterController) SignUpGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", gin.H{
		"title": "Wblog | Registeration Page",
	})
}

func (r *RegisterController) SignUpPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer r.WriteJSON(c, res)
	if c.PostForm("password") != c.PostForm("confirm_password") {
		res["message"] = "Inconsistent password entry"
		return
	}
	err = r.repository.Register()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["contentType"] = c.ContentType()
	res["succeed"] = true
}
