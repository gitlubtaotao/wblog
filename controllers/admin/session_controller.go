package admin

import (
	"errors"
	"github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/controllers"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/gitlubtaotao/wblog/system"
	"net/http"
	"time"
)

type SessionController struct {
	*controllers.BaseController
	repository repositories.IUserRepository
}

func (s *SessionController) GetSignIn(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"title": "Wblog | Log in",
	})
}

//用户进行登录
func (s *SessionController) PostSignIn(ctx *gin.Context) {
	var (
		res      = gin.H{}
		remember bool
	)
	s.repository = repositories.NewUserRepository(ctx)
	defer s.WriteJSON(ctx, res)
	account := ctx.PostForm("account")
	password := ctx.PostForm("password")
	if account == "" || password == "" {
		res["message"] = "username or password cannot be null"
		return
	}
	if ctx.PostForm("checkbox") != "" {
		remember = true
	}
	user, err := s.repository.SignIn(account, password)
	if err != nil {
		_ = seelog.Error(err)
		res["message"] = "Your account not exist"
		return
	}
	if user.LockState {
		res["message"] = "Your account have been locked"
		return
	}
	
	key, err := encrypt.EnCryptData(string(user.ID))
	if err != nil {
		_ = seelog.Error(err)
		res["message"] = "Your account not exist"
		return
	}
	_ = s.OperationSession(ctx, controllers.SESSION_KEY, key)
	res["succeed"] = true
	res["remember"] = remember
	res["contentType"] = ctx.ContentType()
	//进行session id 加密
}

func (s *SessionController) LogoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(controllers.SESSION_KEY)
	_ = session.Save()
	c.Redirect(http.StatusSeeOther, "/admin/signin")
}

//修改密码功能
func (s *SessionController) GetPassword(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "auth/password.html", gin.H{
		"title": "Wblog | Modify Your Password",
	})
}

func (s *SessionController) ModifyPassword(ctx *gin.Context) {
	path := "auth/modify_password.html"
	s.repository = repositories.NewUserRepository(ctx)
	hash := ctx.Param("hash")
	email, err := encrypt.DeCryptData(hash, false)
	if err != nil {
		s.errorHandler(ctx, err, path, gin.H{
			"message": "Your Account is not exist",
		})
	}
	user, err := s.repository.FirstUserByEmail(email)
	if err != nil {
		s.errorHandler(ctx, err, path, gin.H{
			"message": "Your Account is not exist",
		})
	}
	b := time.Now().Sub(user.ModifyPasswordTime).Minutes()
	if b > float64(system.GetConfiguration().PasswordValid) {
		s.errorHandler(ctx, err, path, gin.H{
			"message": "Verification code has expired",
		})
		return
	}
	ctx.HTML(http.StatusOK, path, gin.H{
		"title": "Wblog | Modify Your Password",
		"email": email,
	})
}

func (s *SessionController) UpdatePassword(ctx *gin.Context) {
	s.repository = repositories.NewUserRepository(ctx)
	password := ctx.PostForm("password")
	confirmPassword := ctx.PostForm("confirm_password")
	email := ctx.PostForm("email")
	h := gin.H{
		"message": "Passwords entered twice are inconsistent",
		"email":   email,
	}
	if password != confirmPassword {
		s.errorHandler(ctx, errors.New("Passwords entered twice are inconsistent"),
			"auth/modify_password.html", h)
		return
	}
	_, err := s.repository.FirstUserByEmail(email)
	if err != nil {
		s.errorHandler(ctx, err, "auth/modify_password.html", h)
		return
	}
	password, _ = encrypt.HashAndSalt(password)
	var attr map[string]interface{}
	attr = make(map[string]interface{}, 1)
	attr["password"] = password
	err = s.repository.UpdateUserAttr(attr)
	if err != nil {
		s.errorHandler(ctx, err, "auth/modify_password.html", h)
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, "/admin/signin")
}

//发送邮件or验证码
func (s *SessionController) SendNotice(ctx *gin.Context) {
	s.repository = repositories.NewUserRepository(ctx)
	path := "auth/password.html"
	email := ctx.PostForm("email")
	if email == "" {
		s.errorHandler(ctx, errors.New("Email not eq null"),
			path, gin.H{"message": "Email not eq null"},
		)
		return
	}
	message := gin.H{"message": "Your Account is not exist",}
	_, err := s.repository.FirstUserByEmail(email)
	if err != nil {
		s.errorHandler(ctx, err, path, message)
		return
	}
	//生成对于的修改密码链接
	modifyPasswordHash, err := encrypt.EnCryptData(email)
	if err != nil {
		s.errorHandler(ctx, err, path, message)
		return
	}
	var attr map[string]interface{}
	attr = make(map[string]interface{}, 2)
	attr["ModifyPasswordHash"] = modifyPasswordHash
	attr["ModifyPasswordTime"] = time.Now()
	err = s.repository.UpdateUserAttr(attr)
	if err != nil {
		s.errorHandler(ctx, err, path, message)
		return
	}
	err = helpers.SendToMail(email, "Reset password", s.sendPasswordContext(modifyPasswordHash), "html")
	if err != nil {
		s.errorHandler(ctx, err, path, message)
		return
	}
	ctx.HTML(http.StatusOK, "auth/signin.html", gin.H{
		"title":   "Wblog | Log in",
		"message": "Reset password has been sent to your email",
	})
}

//错误处理
func (s *SessionController) errorHandler(ctx *gin.Context, err error, path string, h gin.H) {
	_ = seelog.Error(err)
	s.RenderHtml(ctx, path, h)
	ctx.Abort()
}

//生成重置密码对应的内容
func (s *SessionController) sendPasswordContext(hash string) (content string) {
	url := "http://localhost:8081/password/modifyPassword/" + hash
	content = `<div>Please click the link to reset your password` + `<a href=` + url + `>` + url + `</a>` +
		`Effective time is 60 minutes` + `</div>`
	return content
}
