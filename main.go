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
	"net/http"
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
		return
	}
	defer db.Close()
	
	
	router := gin.Default()
	
	setTemplate(router)
	setSessions(router)
	
	//设置设置gin模式。参数可以传递：gin.DebugMode、gin.ReleaseMode、gin.TestMode
	gin.SetMode(gin.DebugMode)
	
	router.Use(SharedData())
	//Periodic tasks
	
	//定时任务处理器
	gocron.Every(1).Day().Do(controllers.CreateXMLSitemap)
	gocron.Every(7).Days().Do(controllers.Backup)
	gocron.Start()
	//获取静态资源文件
	router.Static("/static", "./static")
	//router.Static("/static", filepath.Join(getCurrentDirectory(), "./static"))
	//路由不存在
	router.NoRoute(controllers.Handle404)
	
	router.GET("/", controllers.IndexGet)
	router.GET("/index", controllers.IndexGet)
	router.GET("/rss", controllers.RssGet)
	
	//注册路由
	if system.GetConfiguration().SignupEnabled {
		signUp(router)
	}
	//登录和退出
	signInAndOut(router)
	// captcha
	router.GET("/captcha", controllers.CaptchaGet)
	//访问者路由
	visitorRouter(router)
	//订阅者路由
	subscriberRouter(router)
	//other router
	otherRouter(router)
	//后台路由
	adminRouter(router)
	router.Run(system.GetConfiguration().Addr)
}

func signUp(engine *gin.Engine)  {
	engine.GET("/signup", controllers.SignupGet)
	engine.POST("/signup", controllers.SignupPost)
}

//登录和退出
func signInAndOut(engine *gin.Engine)  {
	engine.GET("/signin",controllers.SigninGet)
	engine.POST("/signin", controllers.SigninPost)
	engine.GET("/logout", controllers.LogoutGet)
	engine.GET("/oauth2callback", controllers.Oauth2Callback)
	engine.GET("/auth/:authType", controllers.AuthGet)
}

func visitorRouter(engine *gin.Engine)  {
	visitor := engine.Group("/visitor")
	visitor.Use(AuthRequired())
	{
		visitor.POST("/new_comment", controllers.CommentPost)
		visitor.POST("/comment/:id/delete", controllers.CommentDelete)
	}
}

//订阅者访问
func subscriberRouter(engine *gin.Engine){
	engine.GET("/subscribe", controllers.SubscribeGet)
	engine.POST("/subscribe", controllers.Subscribe)
	engine.GET("/active", controllers.ActiveSubscriber)
	engine.GET("/unsubscribe", controllers.UnSubscribe)
}

func otherRouter(engine *gin.Engine)  {
	engine.GET("/page/:id", controllers.PageGet)
	engine.GET("/post/:id", controllers.PostGet)
	engine.GET("/tag/:tag", controllers.TagGet)
	engine.GET("/archives/:year/:month", controllers.ArchiveGet)
	engine.GET("/link/:id", controllers.LinkGet)
}
//后台路由
func adminRouter(engine *gin.Engine)  {
	authorized := engine.Group("/admin")
	authorized.Use(AdminScopeRequired())
	{
		authorized.GET("/index", controllers.AdminIndex)
		authorized.POST("/upload", controllers.Upload)
		authorized.GET("/page", controllers.PageIndex)
		authorized.GET("/new_page", controllers.PageNew)
		authorized.POST("/new_page", controllers.PageCreate)
		authorized.GET("/page/:id/edit", controllers.PageEdit)
		authorized.POST("/page/:id/edit", controllers.PageUpdate)
		authorized.POST("/page/:id/publish", controllers.PagePublish)
		authorized.POST("/page/:id/delete", controllers.PageDelete)
		
		// post
		authorized.GET("/post", controllers.PostIndex)
		authorized.GET("/new_post", controllers.PostNew)
		authorized.POST("/new_post", controllers.PostCreate)
		authorized.GET("/post/:id/edit", controllers.PostEdit)
		authorized.POST("/post/:id/edit", controllers.PostUpdate)
		authorized.POST("/post/:id/publish", controllers.PostPublish)
		authorized.POST("/post/:id/delete", controllers.PostDelete)
		// tag
		authorized.POST("/new_tag", controllers.TagCreate)
		authorized.GET("/user", controllers.UserIndex)
		authorized.POST("/user/:id/lock", controllers.UserLock)
		// profile
		authorized.GET("/profile", controllers.ProfileGet)
		authorized.POST("/profile", controllers.ProfileUpdate)
		authorized.POST("/profile/email/bind", controllers.BindEmail)
		authorized.POST("/profile/email/unbind", controllers.UnbindEmail)
		authorized.POST("/profile/github/unbind", controllers.UnbindGithub)
		
		// subscriber
		authorized.GET("/subscriber", controllers.SubscriberIndex)
		authorized.POST("/subscriber", controllers.SubscriberPost)
		
		// link
		authorized.GET("/link", controllers.LinkIndex)
		authorized.POST("/new_link", controllers.LinkCreate)
		authorized.POST("/link/:id/edit", controllers.LinkUpdate)
		authorized.POST("/link/:id/delete", controllers.LinkDelete)
		// comment
		authorized.POST("/comment/:id", controllers.CommentRead)
		authorized.POST("/read_all", controllers.CommentReadAll)
		
		// backup
		authorized.POST("/backup", controllers.BackupPost)
		authorized.POST("/restore", controllers.RestorePost)
		
		// mail
		authorized.POST("/new_mail", controllers.SendMail)
		authorized.POST("/new_batchmail", controllers.SendBatchMail)
	}
	
}

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


func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(controllers.CONTEXT_USER_KEY); user != nil {
			if _, ok := user.(*models.User); ok {
				c.Next()
				return
			}
		}
		seelog.Warnf("User not authorized to visit %s", c.Request.RequestURI)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "Forbidden!",
		})
		c.Abort()
	}
}

//AuthRequired grants access to authenticated users, requires SharedData middleware
func AdminScopeRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if user, _ := c.Get(controllers.CONTEXT_USER_KEY); user != nil {
			if u, ok := user.(*models.User); ok && u.IsAdmin {
				c.Next()
				return
			}
		}
		seelog.Warnf("User not authorized to visit %s", c.Request.RequestURI)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "Forbidden!",
		})
		c.Abort()
	}
}



