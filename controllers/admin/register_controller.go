package admin

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/controllers"
	"github.com/gitlubtaotao/wblog/services"
	
	"net/http"
)

type RegisterController struct {
	*controllers.BaseController
}

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
	service := services.NewUserService(c)
	err = service.Register()
	if err != nil {
		_ = seelog.Critical(err)
		res["message"] = err
		return
	}
	res["succeed"] = true
}
