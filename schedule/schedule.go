package schedule

import (
	"github.com/claudiu/gocron"
	"github.com/gitlubtaotao/wblog/admin/api"
	"github.com/gitlubtaotao/wblog/api"
)

func GoCron()  {
	gocron.Every(1).Day().Do(api.CreateXMLSitemap)
	backup := admin.BackUpApi{}
	gocron.Every(7).Days().Do(backup.Backup)
	gocron.Start()
}
