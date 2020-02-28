package schedule

import (
	"github.com/claudiu/gocron"
	"github.com/gitlubtaotao/wblog/controllers"
)

func GoCron()  {
	gocron.Every(1).Day().Do(controllers.CreateXMLSitemap)
	gocron.Every(7).Days().Do(controllers.Backup)
	gocron.Start()
}
