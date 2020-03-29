package api

import (
	"fmt"
	"github.com/cihub/seelog"
	"log"
	"net/http"
	"os"
	"path"
	
	"strings"
	
	"github.com/denisbakhtin/sitemap"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
)

const (
	SESSION_KEY          = "_session_UserID" // session key
	CONTEXT_USER_KEY     = "User"            // context user key
	SESSION_GITHUB_STATE = "GITHUB_STATE"    // github state session key
	SESSION_CAPTCHA      = "GIN_CAPTCHA"     // captcha session key
)

func Handle404(c *gin.Context) {
	HandleMessage(c, "Sorry,I lost myself!")
}

func HandleMessage(c *gin.Context, message string) {
	c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
		"message": message,
	})
}


//发送邮件
func sendMail(to, subject, body string) error {
	return helpers.SendToMail(to, subject, body, "html")
}
// 通知邮件
func NotifyEmail(subject, body string) error {
	notifyEmailsStr := system.GetConfiguration().NotifyEmails
	if notifyEmailsStr != "" {
		notifyEmails := strings.Split(notifyEmailsStr, ";")
		emails := make([]string, 0)
		for _, email := range notifyEmails {
			if email != "" {
				emails = append(emails, email)
			}
		}
		if len(emails) > 0 {
			return sendMail(strings.Join(emails, ";"), subject, body)
		}
	}
	return nil
}

func CreateXMLSitemap() {
	configuration := system.GetConfiguration()
	folder := path.Join(configuration.Public, "sitemap")
	os.MkdirAll(folder, os.ModePerm)
	domain := configuration.Domain
	now := helpers.GetCurrentTime()
	items := make([]sitemap.Item, 0)
	
	items = append(items, sitemap.Item{
		Loc:        domain,
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1,
	})
	
	posts, err := models.ListPublishedPost("", 0, 0)
	if err == nil {
		for _, post := range posts {
			items = append(items, sitemap.Item{
				Loc:        fmt.Sprintf("%s/post/%d", domain, post.ID),
				LastMod:    post.UpdatedAt,
				Changefreq: "weekly",
				Priority:   0.9,
			})
		}
	}
	
	pages, err := models.ListPublishedPage()
	if err == nil {
		for _, page := range pages {
			items = append(items, sitemap.Item{
				Loc:        fmt.Sprintf("%s/page/%d", domain, page.ID),
				LastMod:    page.UpdatedAt,
				Changefreq: "monthly",
				Priority:   0.8,
			})
		}
	}
	
	if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items); err != nil {
		return
	}
	if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/static/sitemap/"); err != nil {
		return
	}
}

func WriteJSON(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	ctx.JSON(http.StatusOK, h)
}

func GetUser(c *gin.Context) *models.User {
	user, _ := c.Get(CONTEXT_USER_KEY)
	return user.(*models.User)
}

//处理共同错误信息
func HandlerError(message string, err error) bool {
	if err != nil {
		_ = seelog.Critical(message, err)
		log.Fatalf("%s:%s", message, err)
		return false
	}
	return true
}