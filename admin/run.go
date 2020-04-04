package admin

import (
	"flag"
	"github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/database"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/migration"
	"github.com/gitlubtaotao/wblog/schedule"
	"github.com/gitlubtaotao/wblog/service"
	"github.com/gitlubtaotao/wblog/system"
	"strconv"
	"strings"
)

func main() {
	//配置环境变量
	configEnv := flag.String("env", "development", "set env development or production")
	flag.Parse()
	system.SetSeelogPath("../conf/seelog.xml")
	
	if err := system.LoadEnvConfiguration(*configEnv); err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	defer seelog.Flush()
	database.InitDB()
	defer database.DBCon.Close()
	migration.Migrate()
	gin.SetMode(system.GetGinMode(*configEnv))
	router := gin.Default()
	router.Static("../static", "./static")
	system.SetCommonTemplate(router)
	system.SetSessions(router, "admin", map[string]interface{}{
		"HttpOnly": true, "MaxAge": 7 * 86400, "Path": "/admin",
	})
	schedule.GoCron()
	router.Use(SharedData())
	err := router.Run(system.GetConfiguration().AdminAddr)
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
