package controllers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/helpers"
	"net/http"
	"path"
	"time"
)

type CaptchaController struct {
}

func (cap *CaptchaController) GetCaptcha(c *gin.Context) {
	result := helpers.GetCaptchaImage(0)
	c.JSON(http.StatusOK, gin.H{
		"captchaId": result.CaptchaId,
		"imageUrl":  result.ImageUrl,
	})
}

func (cap *CaptchaController) VerifyCaptcha(c *gin.Context) {

}
func (cap *CaptchaController) GetCaptchaPng(c *gin.Context) {
	source := c.Param("source")
	fmt.Println("GetCaptchaPng : " + source)
	_, file := path.Split(c.Request.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	reload := false
	if c.Request.FormValue("reload") != "" {
		reload = true
	}
	content, err := helpers.ServeHTTPCaptcha(id, reload)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Writer.Header().Set("Pragma", "no-cache")
	c.Writer.Header().Set("Expires", "0")
	c.Writer.Header().Set("Content-Type", "image/png")
	http.ServeContent(c.Writer, c.Request, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
}
