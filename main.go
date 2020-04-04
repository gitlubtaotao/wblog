package main
//
//import (
//	"flag"
//	"github.com/cihub/seelog"
//	"github.com/gitlubtaotao/wblog/api"
//	"github.com/gitlubtaotao/wblog/database"
//	"github.com/gitlubtaotao/wblog/encrypt"
//	"github.com/gitlubtaotao/wblog/service"
//	"strconv"
//
//	"github.com/gitlubtaotao/wblog/migration"
//	"github.com/gitlubtaotao/wblog/schedule"
//
//	"github.com/gin-contrib/sessions"
//	"github.com/gin-contrib/sessions/cookie"
//	"github.com/gin-gonic/gin"
//	"github.com/gitlubtaotao/wblog/helpers"
//	"github.com/gitlubtaotao/wblog/system"
//
//	"github.com/gitlubtaotao/wblog/tools"
//
//	"os"
//	"path/filepath"
//	"strings"
//
//	"html/template"
//)
//
//func main() {
//	//go里面的flag包，主要是用于解析命令行参数
//	configFilePath := flag.String("C", "conf/conf.yaml", "config file path")
//	logConfigPath := flag.String("L", "conf/seelog.xml", "log config file path")
//	flag.Parse()
//	logger, err := seelog.LoggerFromConfigAsFile(*logConfigPath)
//	if err != nil {
//		_ = seelog.Critical("err parsing seelog config file", err)
//		return
//	}
//	_ = seelog.ReplaceLogger(logger)
//
//	//加载配置文件
//	if err := system.LoadConfiguration(*configFilePath); err != nil {
//		seelog.Critical("err parsing config log file", err)
//		return
//	}
//	defer seelog.Flush()
//	//初始化数据库
//	database.InitDB()
//	defer database.DBCon.Close()
//	migration.Migrate()
//
//	//设置设置gin模式。参数可以传递：gin.DebugMode、gin.ReleaseMode、gin.TestMode
//	//gin.SetMode(gin.DebugMode)
//	gin.SetMode(gin.DebugMode)
//	router := gin.Default()
//	setTemplate(router)
//	setSessions(router)
//	router.Use(SharedData())
//	schedule.GoCron() //Periodic tasks
//
//	//获取静态资源文件
//	router.Static("/static", "./static")
//	//router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
//
//	//注册路由
//	tools.NewRoutes(router).InitRouter()
//	err = router.Run(system.GetConfiguration().Addr)
//	if err != nil {
//		panic(err)
//	}
//}
//
////定义模版
//func setTemplate(engine *gin.Engine) {
//	funcMap := template.FuncMap{
//		"dateFormat": tools.DateFormat,
//		"substring":  tools.Substring,
//		"isOdd":      tools.IsOdd,
//		"isEven":     tools.IsEven,
//		"truncate":   helpers.Truncate,
//		"add":        tools.Add,
//		"minus":      tools.Minus,
//		"listtag":    tools.ListTag,
//	}
//	engine.SetFuncMap(funcMap)
//	//engine.LoadHTMLGlob(filepath.Join(getCurrentDirectory(), "./views/**/*"))
//	engine.LoadHTMLGlob("views/**/*")
//}
//
////setSessions initializes sessions & csrf middlewares
//func setSessions(router *gin.Engine) {
//	config := system.GetConfiguration()
//	//https://github.com/gin-gonic/contrib/tree/master/sessions
//	store := cookie.NewStore([]byte(config.SessionSecret))
//	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"}) //Also set Secure: true if using SSL, you should though
//	router.Use(sessions.Sessions("gin-session", store))
//	//go get github.com/utrack/gin-csrf
//	/*router.Use(csrf.Middleware(csrf.Options{
//		Secret: config.SessionSecret,
//		ErrorFunc: func(c *gin.Context) {
//			c.String(400, "CSRF token mismatch")
//			c.Abort()
//		},
//	}))*/
//}
//
////SharedData fills in common data, such as user info, etc...
//func SharedData() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//静态资源不做判断
//		if strings.Contains(c.Request.URL.String(), "/static") {
//			c.Next()
//			return
//		}
//		session := sessions.Default(c)
//		if uID := session.Get(api.SESSION_KEY); uID != nil {
//			userString, err := encrypt.DeCryptData(uID.(string), true)
//			intId, _ := strconv.ParseInt(userString, 10, 64)
//			user, err := service.NewUserService().GetUserByID(intId)
//			if err == nil {
//				c.Set(api.CONTEXT_USER_KEY, user)
//			} else {
//				_ = seelog.Error(err)
//			}
//		}
//		if system.GetConfiguration().SignupEnabled {
//			c.Set("SignupEnabled", true)
//		}
//		c.Next()
//	}
//}
//
//func getCurrentDirectory() string {
//	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
//	if err != nil {
//		seelog.Critical(err)
//	}
//	path := strings.Replace(dir, "\\", "/", -1)
//	return path
//}
