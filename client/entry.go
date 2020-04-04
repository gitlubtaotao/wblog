package main

import (
	"flag"
	"github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/migration"
	"github.com/gitlubtaotao/wblog/service"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/gitlubtaotao/wblog/tools"
	"html/template"
	"strconv"
	"strings"
)

func main() {
	//配置环境变量
	configEnv := flag.String("env", "development", "set env development or production")
	flag.Parse()
	setSeelogPath("../conf/seelog.xml")
	if err := system.LoadEnvConfiguration(*configEnv); err != nil {
		_ = seelog.Critical("err parsing config log file", err)
		return
	}
	defer seelog.Flush()
	database.InitDB()
	defer database.DBCon.Close()
	migration.Migrate()
	gin.SetMode(system.GetGinMode(*configEnv))
	router := gin.Default()
	router.Static("../static", "./static")
	router.SetFuncMap(setCommonTemplate())
	router.LoadHTMLGlob("view/**/*")
	setSessions(router)
	//schedule.GoCron()
	//router.Use(SharedData())
	err := router.Run(system.GetConfiguration().ClientAddr)
	if err != nil {
		panic(err)
	}
}

func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		//静态资源不做判断
		if strings.Contains(c.Request.URL.String(), "/static") {
			c.Next()
			return
		}
		config := system.GetConfiguration()
		session := sessions.Default(c)
		if uID := session.Get(config.AdminSessionKey); uID != nil {
			userString, err := encrypt.DeCryptData(uID.(string), true)
			intId, _ := strconv.ParseInt(userString, 10, 64)
			user, err := service.NewUserService().GetUserByID(intId)
			if err == nil {
				c.Set(config.AdminUser, user)
			} else {
				_ = seelog.Error(err)
			}
		}
		if system.GetConfiguration().SignupEnabled {
			c.Set("SignupEnabled", true)
		}
		c.Next()
	}
}

/*
 @title: 设置session
*/
func setSessions(router *gin.Engine) {
	config := system.GetConfiguration()
	var sessionSecret = config.ClientSecret
	//https://github.com/gin-gonic/contrib/tree/master/sessions
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   7 * 86400,
		Path:     "/admin",
	}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
}

/*
@title: 配置seelog
@description: 配置系统日志管理
@auth: taotao
@date: 2020.4.4
*/
func setSeelogPath(logConfigPath string) {
	logger, err := seelog.LoggerFromConfigAsFile(logConfigPath)
	if err != nil {
		_ = seelog.Critical("err parsing seelog config file", err)
		return
	}
	_ = seelog.ReplaceLogger(logger)
}

/*
@title: set common template

*/
func setCommonTemplate() template.FuncMap {
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
	return funcMap
}
