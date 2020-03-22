package helpers

import (
	"bytes"
	"errors"
	"github.com/dchest/captcha"
)

type CaptchaImageResponse struct {
	CaptchaId string `json:"captchaId"`
	ImageUrl  string `json:"imageUrl"`
}

//new captcha
func NewCaptchaResponse(CaptchaId string, ImageUrl string) CaptchaImageResponse {
	return CaptchaImageResponse{CaptchaId: CaptchaId, ImageUrl: ImageUrl}
}

//get Captcha
func GetCaptchaImage(length int) (response CaptchaImageResponse) {
	var captchaId string
	if length == 0 {
		captchaId = captcha.New()
	} else {
		captchaId = captcha.NewLen(length)
	}
	response = NewCaptchaResponse(captchaId, "/captcha/"+captchaId+".png")
	return
}

func ServeHTTPCaptcha(id string, reload bool) (content bytes.Buffer, err error) {
	if id == "" {
		err = captcha.ErrNotFound
		return
	}
	content, err = ServeImage(id, 0, 0, reload)
	return
}

func ServeImage(id string, width, height int, reload bool) (content bytes.Buffer, err error) {
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

func VerifyCaptcha(captchaId string, value string) (err error) {
	if captchaId == "" || value == "" {
		return errors.New("参数错误")
	}
	if captcha.VerifyString(captchaId, value) {
		return
	} else {
		return errors.New("验证失败")
	}
}


