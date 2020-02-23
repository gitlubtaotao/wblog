package main

import (
	"flag"
	"github.com/cihub/seelog"
	"github.com/claudiu/gocron"
	
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/controllers"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/gitlubtaotao/wblog/wrouter"
	
	"os"
	"path/filepath"
	"strings"
	
	"html/template"
)

func main()  {
	//go里面的flag包，主要是用于解析命令行参数
	configFilePath := flag.String("C", "conf/conf.yaml", "config file path")
	logConfigPath := flag.String("L", "conf/seelog.xml", "log config file path")
	flag.Parse()
	logger, err := seelog.LoggerFromConfigAsFile(*logConfigPath)
	
	if err != nil {
		seelog.Critical("err parsing seelog config file", err)
		return
	}
	seelog.ReplaceLogger(logger)
	
	if err := system.LoadConfiguration(*configFilePath); err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	defer seelog.Flush()
	
	//初始化数据库
	db, err := models.InitDB()
	if err != nil {
		seelog.Critical("err open databases", err)
		panic(err)
		return
	}
	defer db.Close()
	
	router := gin.Default()
	//设置设置gin模式。参数可以传递：gin.DebugMode、gin.ReleaseMode、gin.TestMode
	gin.SetMode(gin.DebugMode)
	
	setTemplate(router)
	setSessions(router)
	router.Use(SharedData())
	
	goCron() //Periodic tasks
	//获取静态资源文件
	router.Static("/static", "./static")
	//router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
	//路由不存在
	//注册路由
	wrouter.InitRouter(router)
	router.Run(system.GetConfiguration().Addr)
}


//定时任务
func goCron()  {
	gocron.Every(1).Day().Do(controllers.CreateXMLSitemap)
	gocron.Every(7).Days().Do(controllers.Backup)
	gocron.Start()
}
//定义模版
func setTemplate(engine *gin.Engine) {
	funcMap := template.FuncMap{
		"dateFormat": helpers.DateFormat,
		"substring":  helpers.Substring,
		"isOdd":      helpers.IsOdd,
		"isEven":     helpers.IsEven,
		"truncate":   helpers.Truncate,
		"add":        helpers.Add,
		"minus":      helpers.Minus,
		"listtag":    helpers.ListTag,
	}
	engine.SetFuncMap(funcMap)
	//engine.LoadHTMLGlob(filepath.Join(getCurrentDirectory(), "./views/**/*"))
	engine.LoadHTMLGlob("views/**/*")
}

//setSessions initializes sessions & csrf middlewares
func setSessions(router *gin.Engine) {
	config := system.GetConfiguration()
	//https://github.com/gin-gonic/contrib/tree/master/sessions
	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
	//https://github.com/utrack/gin-csrf
	/*router.Use(csrf.Middleware(csrf.Options{
		Secret: config.SessionSecret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))*/
}

//SharedData fills in common data, such as user info, etc...
func SharedData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if uID := session.Get(controllers.SESSION_KEY); uID != nil {
			user, err := models.GetUser(uID)
			if err == nil {
				c.Set(controllers.CONTEXT_USER_KEY, user)
			}
		}
		if system.GetConfiguration().SignupEnabled {
			c.Set("SignupEnabled", true)
		}
		c.Next()
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		seelog.Critical(err)
	}
	path := strings.Replace(dir, "\\", "/", -1)
	return path
}








