package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
	"math"
	"net/http"
)

type HomeApi struct {
	*UtilApi
}

func (h *HomeApi) Index(ctx *gin.Context) {
	var (
		pageSize = system.GetConfiguration().PageSize
		total    int
		err      error
		posts    []*models.Post
		policy   *bluemonday.Policy
	)
	pageIndex, _ := h.PageIndex(ctx)
	posts, err = h.listPost(ctx)
	if err != nil {
		h.HandlerError("", err)
		return
	}
	total, err = h.CountPostByTag(ctx)
	if err != nil {
		h.HandlerError("", err)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Body = policy.Sanitize(string(blackfriday.Run([]byte(post.Body), blackfriday.WithNoExtensions())))
	}
	ctx.HTML(http.StatusOK, "home/index.html", gin.H{
		"posts":           posts,
		"tags":            h.PublishTags(ctx),
		"archives":        models.MustListPostArchives(),
		"links":           h.ListLinks(ctx),
		"pageIndex":       pageIndex,
		"totalPage":       int(math.Ceil(float64(total) / float64(pageSize))),
		"path":            ctx.Request.URL.Path,
		"maxReadPosts":    h.MaxReadPost(ctx),
		"maxCommentPosts": h.MaxCommentPost(ctx),
	})
}

/*
 
 */
func (h *HomeApi) listPost(ctx *gin.Context) (posts []*models.Post, err error) {
	repository := repositories.NewPostRepository(ctx)
	page, _ := h.PageIndex(ctx)
	posts, err = repository.PublishPost(
		uint(system.GetConfiguration().PageSize),
		uint(page),
		map[string]interface{}{},
		[]string{},
		true,
	)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	return
}

func (h *HomeApi) CountPostByTag(ctx *gin.Context) (int, error) {
	repository := repositories.NewPostRepository(ctx)
	attr := map[string]interface{}{
		"is_published": true,
	}
	total, err := repository.CountPost(attr)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return 0, err
	}
	return total, err
}

func (h *HomeApi) PublishTags(ctx *gin.Context) []*models.Tag {
	repository := repositories.NewTagRepository(ctx)
	tags, _ := repository.PublishTagsList()
	return tags
}

func (h *HomeApi) ListLinks(ctx *gin.Context) []*models.Link {
	repository := repositories.NewLinkRepository(ctx)
	links, _ := repository.ListAllLink([]string{})
	return links
}

func (h *HomeApi) MaxReadPost(ctx *gin.Context) []*models.Post {
	columns := []string{"title", "id"}
	repository := repositories.NewPostRepository(ctx)
	posts, _ := repository.ListMaxReadPost(columns)
	return posts
}

func (h *HomeApi) MaxCommentPost(ctx *gin.Context) []*models.Post {
	columns := []string{"title", "id"}
	repository := repositories.NewPostRepository(ctx)
	posts, _ := repository.ListMaxCommentPost(columns)
	return posts
}
