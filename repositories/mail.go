package repositories

import (
	"github.com/gitlubtaotao/wblog/system"
	"net/smtp"
	"strings"
)

type IMailRepository interface {
	SendToMail(to string) error
	SystemDefaultNotify() error
	GetSystemNotify() string
}

/*
 邮箱发送功能
 @params to: 接受者邮箱
 @params subject: 发送主体
 @params body: 主体内容
 @params mailType: 发送的内容
*/
type MailRepository struct {
	subject  string
	body     string
	mailType string
}

func (m *MailRepository) SystemDefaultNotify() error {
	to := m.GetSystemNotify()
	return m.SendToMail(to)
}


func (m *MailRepository) GetSystemNotify() string {
	notifyEmailsStr := system.GetConfiguration().NotifyEmails
	emails := make([]string, 0)
	if notifyEmailsStr != "" {
		notifyEmails := strings.Split(notifyEmailsStr, ";")
		for _, email := range notifyEmails {
			if email != "" {
				emails = append(emails, email)
			}
		}
	}
	return strings.Join(emails, ";")
}


func (m *MailRepository) SendToMail(to string) error {
	config := system.GetConfiguration()
	user := config.SmtpUsername
	password := config.SmtpPassword
	host := config.SmtpHost
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if m.mailType == "html" {
		contentType = "Content-Type: text/" + m.mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + m.subject + "\r\n" + contentType + "\r\n\r\n" + m.body)
	sendTo := strings.Split(to, ";")
	return smtp.SendMail(host, auth, user, sendTo, msg)
}

func NewMailRepository(subject, body, mailType string) IMailRepository {
	return &MailRepository{subject: subject, body: body, mailType: mailType}
}
