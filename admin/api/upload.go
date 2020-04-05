package admin

import (
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/tools/upload/qiniu"
	"mime/multipart"
	
	"github.com/gin-gonic/gin"
)

type UploadApi struct {
	*api.BaseApi
}

func (u *UploadApi) Upload(c *gin.Context) {
	var (
		err  error
		res  = gin.H{}
		url  string
		key  string
		file multipart.File
		fh   *multipart.FileHeader
	)
	defer u.WriteJSON(c, res)
	file, fh, err = c.Request.FormFile("file")
	if err != nil {
		res["message"] = err.Error()
		return
	}
	uploader := qiniu.NewUploaderDefault()
	url, key, err = uploader.Upload(file, fh)
	url = uploader.PrivateReadURL(key)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
	res["url"] = url
	res["key"] = key
	return
}
