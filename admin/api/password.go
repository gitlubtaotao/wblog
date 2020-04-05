package admin

import (
	"errors"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/gitlubtaotao/wblog/system"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
	"sync"
	"time"
)

type PasswordApi struct {
	*api.BaseApi
}

func (p *PasswordApi) New(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "password/new.html",
		gin.H{
			"title": "Wblog | Modify Your Password",
			"token": csrf.GetToken(ctx),
		})
}

func (p *PasswordApi) Create(ctx *gin.Context) {
	repository := repositories.NewUserRepository(ctx)
	password := ctx.PostForm("password")
	confirmPassword := ctx.PostForm("confirm_password")
	email := ctx.PostForm("email")
	h := gin.H{
		"message": "Passwords entered twice are inconsistent",
		"email":   email,
	}
	if password != confirmPassword {
		p.errorHandler(ctx, errors.New("Passwords entered twice are inconsistent"),
			"passowrd/modify.html", h)
		return
	}
	_, err := repository.FirstUserByEmail(email)
	if err != nil {
		p.errorHandler(ctx, err, "passowrd/modify.html", h)
		return
	}
	password, _ = encrypt.HashAndSalt(password)
	var attr map[string]interface{}
	attr = make(map[string]interface{}, 1)
	attr["password"] = password
	err = repository.UpdateUserAttr(attr)
	if err != nil {
		p.errorHandler(ctx, err, "passowrd/modify.html", h)
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, "/admin/login")
}

func (p *PasswordApi) SendEmail(ctx *gin.Context) {
	repository := repositories.NewUserRepository(ctx)
	path := "password/new.html"
	email := ctx.PostForm("email")
	if email == "" {
		p.errorHandler(ctx,
			errors.New("Email not eq null"),
			path, gin.H{"message": "Email not eq null"},
		)
		return
	}
	message := gin.H{"message": "Your Account is not exist",}
	_, err := repository.FirstUserByEmail(email)
	if err != nil {
		p.errorHandler(ctx, err, path, message)
		return
	}
	//生成对于的修改密码链接
	modifyPasswordHash, err := encrypt.EnCryptData(email, "admin")
	if err != nil {
		p.errorHandler(ctx, err, path, message)
		return
	}
	var attr map[string]interface{}
	attr = make(map[string]interface{}, 2)
	attr["ModifyPasswordHash"] = modifyPasswordHash
	attr["ModifyPasswordTime"] = time.Now()
	err = repository.UpdateUserAttr(attr)
	if err != nil {
		p.errorHandler(ctx, err, path, message)
		return
	}
	//sync.WaitGroup 进行管理
	var wg sync.WaitGroup
	wg.Add(1)
	go func(email, modifyPasswordHash string) {
		err := p.SendMailHtml(email, "Reset password", p.sendPasswordContext(modifyPasswordHash))
		if err != nil {
			_ = seelog.Error(err)
		}
		wg.Done()
	}(email, modifyPasswordHash)
	wg.Wait()
	ctx.HTML(http.StatusOK, "session/new.html", gin.H{
		"title":   "Wblog | Log in",
		"message": "Reset password has been sent to your email",
	})
}

func (p *PasswordApi) Modify(ctx *gin.Context) {
	path := "password/modify.html"
	repository := repositories.NewUserRepository(ctx)
	hash := ctx.Param("hash")
	email, err := encrypt.DeCryptData(hash, false, "admin")
	if err != nil {
		p.errorHandler(ctx, err, path, gin.H{
			"message": "Your Account is not exist",
		})
		return
	}
	user, err := repository.FirstUserByEmail(email)
	if err != nil {
		p.errorHandler(ctx, err, path, gin.H{
			"message": "Your Account is not exist",
		})
		return
	}
	b := time.Now().Sub(user.ModifyPasswordTime).Minutes()
	if b > float64(system.GetConfiguration().PasswordValid) {
		p.errorHandler(ctx, err, path, gin.H{
			"message": "Verification code has expired",
		})
		return
	}
	ctx.HTML(http.StatusOK, path, gin.H{
		"title": "Wblog | Modify Your Password",
		"email": email,
		"token": csrf.GetToken(ctx),
	})
}

//错误处理
func (p *PasswordApi) errorHandler(ctx *gin.Context, err error, path string, h gin.H) {
	_ = seelog.Error(err)
	p.RenderHtml(ctx, path, h)
	ctx.Abort()
}

//生成重置密码对应的内容
func (p *PasswordApi) sendPasswordContext(hash string) (content string) {
	
	url := "http://localhost:8081/admin/password/modify/" + hash
	content = `<div>Please click the link to reset your password` + `<a href=` + url + `>` + url + `</a>` +
		`Effective time is 60 minutes` + `</div>`
	return content
}
