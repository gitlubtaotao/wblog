package repositories

import (
	"bytes"
	"errors"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type ICaptchaRepository interface {
	GetCaptchaImage(length int) CaptchaImageResponse
	ServeHTTPCaptcha(id string, reload bool) (bytes.Buffer, error)
	ServeImage(id string, width, height int, reload bool) (content bytes.Buffer, err error)
	VerifyCaptcha(captchaId string, value string) (err error)
}
type CaptchaRepository struct {
	ctx       *gin.Context
}
type CaptchaImageResponse struct {
	CaptchaId string `json:"captchaId"`
	ImageUrl  string `json:"imageUrl"`
}

func NewCaptchaRepository(ctx *gin.Context) ICaptchaRepository {
	return &CaptchaRepository{ctx: ctx}
}

func (c *CaptchaRepository) GetCaptchaImage(length int) CaptchaImageResponse {
	var captchaId string
	if length == 0 {
		captchaId = captcha.New()
	} else {
		captchaId = captcha.NewLen(length)
	}
	return CaptchaImageResponse{CaptchaId: captchaId, ImageUrl: "/captcha/" + captchaId + ".png"}
}

func (c *CaptchaRepository) ServeHTTPCaptcha(id string, reload bool) (content bytes.Buffer, err error) {
	if id == "" {
		err = captcha.ErrNotFound
		return
	}
	content, err = c.ServeImage(id, 0, 0, reload)
	return
}

func (c *CaptchaRepository) ServeImage(id string, width, height int, reload bool) (content bytes.Buffer, err error) {
	if reload {
		captcha.Reload(id)
	}
	if width == 0 {
		width = captcha.StdWidth
	}
	if height == 0 {
		height = captcha.StdHeight
	}
	err = captcha.WriteImage(&content, id, width, height)
	return
}

func (c *CaptchaRepository) VerifyCaptcha(captchaId string, value string) (err error) {
	if captchaId == "" || value == "" {
		return errors.New("参数错误")
	}
	if captcha.VerifyString(captchaId, value) {
		return
	} else {
		return errors.New("验证失败")
	}
}
