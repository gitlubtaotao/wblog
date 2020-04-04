package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/models"
	"net/http"
)

type HomeApi struct {
	*api.BaseApi
}

func (h *HomeApi) Index(ctx *gin.Context) {
	user, _ := h.AdminUser(ctx)
	fmt.Println(user)
	ctx.HTML(http.StatusOK, "home/index.html", h.RenderComments(gin.H{
		"pageCount":    models.CountPage(),
		"postCount":    models.CountPost(),
		"tagCount":     models.CountTag(),
		"commentCount": models.CountComment(),
		"user":         user,
	}))
}
