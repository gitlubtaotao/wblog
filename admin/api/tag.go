package admin

import (
	"github.com/cihub/seelog"
	"github.com/gitlubtaotao/wblog/api"
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
	*api.BaseApi
	repository repositories.ITagRepository
}

func (t *TagApi) Index(ctx *gin.Context) {
	repository := repositories.NewTagRepository(ctx)
	t.repository = repository
	if ctx.Param("format") == "json" {
		t.WriteJSON(ctx, t.json(ctx))
	}
	
}

//todo-taotao 新增标签和删除对应的标签管理
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

func (t *TagApi) Delete(ctx *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer t.WriteJSON(ctx, res)
	repository := repositories.NewTagRepository(ctx)
	id := ctx.Param("id")
	Id, _ := strconv.ParseUint(id, 10, 64)
	err = repository.Delete(uint(Id))
	if err != nil {
		_ = seelog.Critical(err)
		res["message"] = "Delete tag is error "
		return
	}
	res["succeed"] = true
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

func (t *TagApi) json(ctx *gin.Context) gin.H {
	var (
		res  = gin.H{}
		data []map[string]interface{}
	)
	tags, err := t.repository.AllTag(map[string]interface{}{}, []string{"id", "name"})
	if err != nil {
		_ = seelog.Warn(err)
		res["data"] = data
		return res
	}
	for _, v := range tags {
		temp := map[string]interface{}{
			"text":  v.Name,
			"value": v.ID,
		}
		data = append(data, temp)
	}
	res["succeed"] = true
	res["data"] = data
	return res
}
