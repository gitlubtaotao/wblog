package schedule

import (
	"github.com/claudiu/gocron"
	"github.com/gitlubtaotao/wblog/api"
)

func GoCron()  {
	gocron.Every(1).Day().Do(api.CreateXMLSitemap)
	gocron.Every(7).Days().Do(api.Backup)
	gocron.Start()
}
