package system

import (
	"encoding/json"
	"github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/tools"
	"html/template"
	"io/ioutil"
	
	"github.com/go-yaml/yaml"
)

type Configuration struct {
	SignupEnabled      bool   `yaml:"signup_enabled" json:"signup_enabled"`   // signup enabled or not
	QiniuAccessKey     string `yaml:"qiniu_accesskey" json:"qiniu_accesskey"` // qiniu
	QiniuSecretKey     string `yaml:"qiniu_secretkey" json:"qiniu_secretkey"`
	QiniuFileServer    string `yaml:"qiniu_fileserver" json:"qiniu_fileserver"`
	QiniuBucket        string `yaml:"qiniu_bucket" json:"qiniu_bucket"`
	GithubClientId     string `yaml:"github_clientid" json:"github_clientid"` // github
	GithubClientSecret string `yaml:"github_clientsecret" json:"github_clientsecret"`
	GithubAuthUrl      string `yaml:"github_authurl" json:"github_authurl"`
	GithubRedirectURL  string `yaml:"github_redirecturl" json:"github_redirecturl"`
	GithubTokenUrl     string `yaml:"github_tokenurl" json:"github_tokenurl"`
	GithubScope        string `yaml:"github_scope" json:"github_scope"`
	SmtpUsername       string `yaml:"smtp_username" json:"smtp_username"`   // username
	SmtpPassword       string `yaml:"smtp_password" json:"smtp_password"`   //password
	SmtpHost           string `yaml:"smtp_host" json:"smtp_host"`           //host
	SessionSecret      string `yaml:"session_secret" json:"session_secret"` //session_secret
	Domain             string `yaml:"domain" json:"domain"`                 //domain
	Public             string `yaml:"public" json:"public"`                 //public
	Addr               string `yaml:"addr" json:"addr"`                     //addr
	BackupKey          string `yaml:"backup_key" json:"backup_key"`         //backup_key
	DSN                string `yaml:"dsn" json:"dsn"`                       //database dsn
	NotifyEmails       string `yaml:"notify_emails" json:"notify_emails"`   //notify_emails
	PageSize           int    `yaml:"page_size" json:"page_size"`           //page_size
	SmmsFileServer     string `yaml:"smms_fileserver" json:"smms_fileserver"`
	PasswordValid      int64  `yaml:"password_valid" json:"password_valid"`
	AdminAddr          string `json:"admin_addr" yaml:"admin_addr"`
	ClientAddr         string `json:"client_addr" yaml:"client_addr"`
	AdminSecret        string `json:"admin_secret" yaml:"admin_secret"`
	ClientSecret       string `json:"client_secret" yaml:"client_secret"`
	AdminSessionKey    string `json:"admin_session_key" yaml:"admin_session_key"`
	ClientSessionKey   string `json:"client_session_key" yaml:"client_session_key"`
	AdminUser          string `json:"admin_user" yaml:"admin_user"`
	ClientUser         string `json:"client_user" yaml:"client_user"`
}

const (
	DEFAULT_PAGESIZE = 10
)

var configuration *Configuration

func LoadConfiguration(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	if config.PageSize <= 0 {
		config.PageSize = DEFAULT_PAGESIZE
	}
	configuration = &config
	return err
}

func GetConfiguration() *Configuration {
	return configuration
}

/*
@title: 配置seelog
@description: 配置系统日志管理
@auth: taotao
@date: 2020.4.4
*/
func SetSeelogPath(logConfigPath string) {
	logger, err := seelog.LoggerFromConfigAsFile(logConfigPath)
	if err != nil {
		_ = seelog.Critical("err parsing seelog config file", err)
		return
	}
	_ = seelog.ReplaceLogger(logger)
}

/*
@title: 设置不同的环境变量
@auth: taotao
*/
func LoadEnvConfiguration(env string) error {
	data, err := ioutil.ReadFile("../conf/config.yaml")
	if err != nil {
		return err
	}
	config := make(map[string]interface{})
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	m := make(map[string]interface{})
	for k, v := range config[env].(map[interface{}]interface{}) {
		s := k.(string)
		m[s] = v
	}
	bytes, err := json.Marshal(m)
	err = json.Unmarshal(bytes, &configuration)
	return err
}

func GetGinMode(env string) string {
	switch env {
	case "production":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}

/*
@title: set common template

*/
func SetCommonTemplate(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"dateFormat": tools.DateFormat,
		"substring":  tools.Substring,
		"isOdd":      tools.IsOdd,
		"isEven":     tools.IsEven,
		"truncate":   helpers.Truncate,
		"add":        tools.Add,
		"minus":      tools.Minus,
		"listtag":    tools.ListTag,
	}
	engine.SetFuncMap(funcMap)
	//engine.LoadHTMLGlob(filepath.Join(getCurrentDirectory(), "./views/**/*"))
	engine.LoadHTMLGlob("views/**/*")
}

/*
 @title: 设置session
*/
func SetSessions(router *gin.Engine, env string, options map[string]interface{}) {
	config := GetConfiguration()
	var sessionSecret = ""
	if env == "admin" {
		sessionSecret = config.AdminSecret
	} else {
		sessionSecret = config.ClientSecret
	}
	//https://github.com/gin-gonic/contrib/tree/master/sessions
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		HttpOnly: options["http_only"].(bool),
		MaxAge:   options["MaxAge"].(int),
		Path:     options["path"].(string)}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
	//https://github.com/utrack/gin-csrf
	
	//router.Use(csrf.Middleware(csrf.Options{
	//	Secret: sessionSecret,
	//	ErrorFunc: func(c *gin.Context) {
	//		c.String(400, "CSRF token mismatch")
	//		c.Abort()
	//	},
	//}))
}
