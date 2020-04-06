package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	"strconv"
)

type LinkApi struct {
	UtilApi
}

func (l *LinkApi) Show(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	repository :=  repositories.NewLinkRepository(ctx)
	link, err := repository.GetLinkById(uint(id))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	link.View++
	_ = repository.Update(&link)
	ctx.Redirect(http.StatusFound, link.Url)
}
