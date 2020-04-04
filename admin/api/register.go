package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/repositories"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

type RegisterApi struct {
	*api.BaseApi
	repository repositories.IUserRepository
}

func (r *RegisterApi) New(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register/new.html", gin.H{
		"title": "Wblog | Registeration Page",
		"token": csrf.GetToken(ctx),
	})
}
func (r *RegisterApi) Create(ctx *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer r.WriteJSON(ctx, res)
	if ctx.PostForm("password") != ctx.PostForm("confirm_password") {
		res["message"] = "Inconsistent password entry"
		return
	}
	repository := repositories.NewUserRepository(ctx)
	err = repository.Register()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
