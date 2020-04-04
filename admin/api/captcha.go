package admin

import (
	"bytes"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/gitlubtaotao/wblog/system"
	"net/http"
	"path"
	"time"
)

type CaptchaApi struct {
	*api.BaseApi
	repository repositories.ICaptchaRepository
}

func (c *CaptchaApi) Get(ctx *gin.Context) {
	key := system.GetConfiguration().GinCaptcha
	repository := c.getRepository(ctx)
	session := sessions.Default(ctx)
	result := repository.GetCaptchaImage(4)
	session.Delete(key)
	reload := ctx.Param("reload")
	if reload != "" {
		captcha.Reload(result.CaptchaId)
	}
	session.Set(key, result.CaptchaId)
	_ = session.Save()
	_ = captcha.WriteImage(ctx.Writer, result.CaptchaId, 100, 40)
}

//校验验证码
func (c *CaptchaApi) Verify(ctx *gin.Context) {
	value := ctx.Query("value")
	session := sessions.Default(ctx)
	captchaId := session.Get(system.GetConfiguration().GinCaptcha)
	repository := c.getRepository(ctx)
	errs := repository.VerifyCaptcha(captchaId.(string), value)
	var json = gin.H{}
	json["contentType"] = ctx.ContentType()
	if errs != nil {
		json["message"] = errs.Error()
	} else {
		json["succeed"] = true
	}
	c.WriteJSON(ctx, json)
}

func (c *CaptchaApi) Image(ctx *gin.Context) {
	repository := c.getRepository(ctx)
	_, file := path.Split(ctx.Request.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	reload := false
	if ctx.Request.FormValue("reload") != "" {
		reload = true
	}
	content, err := repository.ServeHTTPCaptcha(id, reload)
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	ctx.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Writer.Header().Set("Pragma", "no-cache")
	ctx.Writer.Header().Set("Expires", "0")
	ctx.Writer.Header().Set("Content-Type", "image/png")
	http.ServeContent(ctx.Writer, ctx.Request, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
}

func (c *CaptchaApi) getRepository(ctx *gin.Context) repositories.ICaptchaRepository  {
	return repositories.NewCaptchaRepository(ctx)
}
