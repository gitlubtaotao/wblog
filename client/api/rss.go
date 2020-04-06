package client

import (
	"fmt"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/gitlubtaotao/wblog/tools/helpers"
	"github.com/gorilla/feeds"
)

type RssApi struct {
	UtilApi
}

func (r *RssApi) RssGet(c *gin.Context) {
	now := helpers.GetCurrentTime()
	domain := system.GetConfiguration().Domain
	feed := &feeds.Feed{
		Title:       "Wblog",
		Link:        &feeds.Link{Href: domain},
		Description: "Wblog,talk about golang,java and so on.",
		Author:      &feeds.Author{Name: "Xutaotao", Email: "xtt691373656@iCloud.com"},
		Created:     now,
	}
	
	feed.Items = make([]*feeds.Item, 0)
	posts, err := r.listPost(c)
	if err != nil {
		_ = seelog.Error(err)
		return
	}
	for _, post := range posts {
		item := &feeds.Item{
			Id:          fmt.Sprintf("%s/post/%d", domain, post.ID),
			Title:       post.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/post/%d", domain, post.ID)},
			Description: string(post.Excerpt()),
			Created:     now,
		}
		feed.Items = append(feed.Items, item)
	}
	rss, err := feed.ToRss()
	if err != nil {
		_ = seelog.Error(err)
		return
	}
	_, _ = c.Writer.WriteString(rss)
}

func (r *RssApi) listPost(ctx *gin.Context) (posts []*models.Post, err error) {
	repository := repositories.NewPostRepository(ctx)
	posts, err = repository.PublishPost(0,0, map[string]interface{}{}, []string{}, false, )
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	return
}