package admin

import (
	"github.com/gitlubtaotao/wblog/api"
	"strings"
	
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
)

type MailApi struct {
	*BaseApi
}

func (m *MailApi) Send(c *gin.Context) {
	var (
		err        error
		res        = gin.H{}
		uid        uint64
		subscriber *models.Subscriber
	)
	defer api.WriteJSON(c, res)
	subject := c.PostForm("subject")
	content := c.PostForm("content")
	userId := c.Query("userId")
	if subject == "" || content == "" || userId == "" {
		res["message"] = "error parameter"
		return
	}
	uid, err = strconv.ParseUint(userId, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	subscriber, err = models.GetSubscriberById(uint(uid))
	if err != nil {
		res["message"] = err.Error()
		return
	}
	err = m.SendMailHtml(subscriber.Email, subject, content)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (m *MailApi) SendBatch(c *gin.Context) {
	var (
		err         error
		res         = gin.H{}
		subscribers []*models.Subscriber
		emails      = make([]string, 0)
	)
	defer api.WriteJSON(c, res)
	subject := c.PostForm("subject")
	content := c.PostForm("content")
	if subject == "" || content == "" {
		res["message"] = "error parameter"
		return
	}
	subscribers, err = models.ListSubscriber(true)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	for _, subscriber := range subscribers {
		emails = append(emails, subscriber.Email)
	}
	err = m.SendMailHtml(strings.Join(emails, ";"), subject, content)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
