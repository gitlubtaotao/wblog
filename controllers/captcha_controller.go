package controllers

import (
	"bytes"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/helpers"
	"net/http"
	"path"
	"time"
)

type CaptchaController struct {
	*BaseController
}

//获取验证码图片
func (cap *CaptchaController) Captcha(context *gin.Context) {
	session := sessions.Default(context)
	result := helpers.GetCaptchaImage(4)
	session.Delete(SESSION_CAPTCHA)
	reload := context.Param("reload")
	if reload != "" {
		captcha.Reload(result.CaptchaId)
	}
	session.Set(SESSION_CAPTCHA, result.CaptchaId)
	_ = session.Save()
	_ = captcha.WriteImage(context.Writer, result.CaptchaId, 100, 40)
}

func (cap *CaptchaController) GetCaptcha(c *gin.Context) {
	result := helpers.GetCaptchaImage(0)
	c.JSON(http.StatusOK, gin.H{
		"captchaId": result.CaptchaId,
		"imageUrl":  result.ImageUrl,
	})
}

//校验验证码
func (cap *CaptchaController) VerifyCaptcha(c *gin.Context) {
	value := c.PostForm("value")
	session := sessions.Default(c)
	captchaId := session.Get(SESSION_CAPTCHA)
	errs := helpers.VerifyCaptcha(captchaId.(string), value)
	var json = gin.H{}
	if errs != nil {
		json["status"] = http.StatusInternalServerError
		json["message"] = errs
	} else {
		json["status"] = http.StatusOK
	}
	cap.WriteJSON(c, json)
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
