package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"strconv"
)

/*
所有的api服务必须要继承UtilApi,所有的公共方法
*/
type IUtilApi interface {
	PageIndex(ctx *gin.Context) (page int, err error)
}

type UtilApi struct {
	*api.BaseApi
}

/*
@title: 获取当前分页数
*/
func (u *UtilApi) PageIndex(ctx *gin.Context) (page int, err error) {
	pageIndex, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		return 0, err
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	return pageIndex, nil
}
