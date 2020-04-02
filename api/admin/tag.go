package admin

import (
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	"strconv"
	
	"math"
	
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

type TagApi struct {
	*BaseApi
	repository repositories.ITagRepository
}

func (t *TagApi) Create(ctx *gin.Context) {
	repository := repositories.NewTagRepository(ctx)
	
	var (
		err error
		res = gin.H{}
	)
	defer t.WriteJSON(ctx, res)
	name := ctx.PostForm("value")
	tag := models.Tag{Name: name}
	tag, err = repository.Create(tag)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
	res["data"] = tag
}

func TagGet(c *gin.Context) {
	var (
		tagName   string
		page      string
		pageIndex int
		pageSize  = system.GetConfiguration().PageSize
		total     int
		err       error
		policy    *bluemonday.Policy
		posts     []*models.Post
	)
	tagName = c.Param("tag")
	page = c.Query("page")
	pageIndex, _ = strconv.Atoi(page)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	posts, err = models.ListPublishedPost(tagName, pageIndex, pageSize)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	total, err = models.CountPostByTag(tagName)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostId(strconv.FormatUint(uint64(post.ID), 10))
		post.Body = policy.Sanitize(string(blackfriday.Run([]byte(post.Body), blackfriday.WithNoExtensions())))
	}
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"posts":           posts,
		"tags":            models.MustListTag(),
		"archives":        models.MustListPostArchives(),
		"links":           models.MustListLinks(),
		"pageIndex":       pageIndex,
		"totalPage":       int(math.Ceil(float64(total) / float64(pageSize))),
		"maxReadPosts":    models.MustListMaxReadPost(),
		"maxCommentPosts": models.MustListMaxCommentPost(),
	})
}
