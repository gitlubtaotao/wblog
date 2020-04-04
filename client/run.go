package main
//
//import (
//	"flag"
//	"github.com/cihub/seelog"
//	"github.com/gin-contrib/sessions"
//	"github.com/gin-contrib/sessions/cookie"
//	"github.com/gin-gonic/gin"
//	"github.com/gitlubtaotao/wblog/database"
//	"github.com/gitlubtaotao/wblog/migration"
//	"github.com/gitlubtaotao/wblog/system"
//)
//
//func main() {
//	configEnv := flag.String("env", "development", "set env development or production")
//	flag.Parse()
//	system.SetSeelogPath("../conf/seelog.xml")
//	if err := system.LoadEnvConfiguration(*configEnv); err != nil {
//		seelog.Critical("err parsing config log file", err)
//		return
//	}
//	defer seelog.Flush()
//	database.InitDB()
//	defer database.DBCon.Close()
//	migration.Migrate()
//	gin.SetMode(system.GetGinMode(*configEnv))
//	router := gin.Default()
//	//获取静态资源文件
//	router.Static("../static", "./static")
//	router.SetFuncMap(system.SetCommonTemplate())
//	router.LoadHTMLGlob("views/**/*")
//	setSessions(router)
//	err := router.Run(system.GetConfiguration().ClientAddr)
//	if err != nil {
//		panic(err)
//	}
//}
//
//func setSessions(router *gin.Engine) {
//	config := system.GetConfiguration()
//	var sessionSecret = config.ClientSecret
//	//https://github.com/gin-gonic/contrib/tree/master/sessions
//	store := cookie.NewStore([]byte(sessionSecret))
//	store.Options(sessions.Options{
//		HttpOnly: true,
//		MaxAge:   7 * 86400,
//		Path:     "/client",
//	}) //Also set Secure: true if using SSL, you should though
//	router.Use(sessions.Sessions("gin-session", store))
//}
