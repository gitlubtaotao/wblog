package admin

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin/binding"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/jinzhu/gorm"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
	"strconv"
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
)

type PostApi struct {
	*api.BaseApi
	post repositories.IPostRepository
	tag  repositories.ITagRepository
}

func (p *PostApi) Index(ctx *gin.Context) {
	repository := repositories.NewPostRepository(ctx)
	posts, _ := repository.ListAll(map[string]interface{}{}, []string{})
	user, _ := p.AdminUser(ctx)
	renderJson := p.RenderComments(gin.H{
		"posts": posts,
		"user":  user,
		"token": csrf.GetToken(ctx),
	})
	ctx.HTML(http.StatusOK, "post/index.html", renderJson)
}

func (p *PostApi) New(c *gin.Context) {
	user, _ := p.AdminUser(c)
	c.HTML(http.StatusOK, "post/edit.html", p.RenderComments(gin.H{
		"user":   user,
		"post":   models.Post{},
		"token":  csrf.GetToken(c),
		"submit": "/admin/posts",
	}))
}

func (p *PostApi) Create(ctx *gin.Context) {
	var res = gin.H{}
	defer p.WriteJSON(ctx, res)
	var post models.Post
	err := ctx.ShouldBindWith(&post, binding.Form)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	isPublished := ctx.PostForm("isPublished")
	repository := p.getRepository(ctx)
	post.IsPublished = "on" == isPublished
	tags := ctx.PostForm("tags")
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		err := repository.Create(&post)
		if err != nil {
			return err
		}
		if len(tags) > 0 {
			if err := p.CreatePostTag(ctx, tags, post.ID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		_ = seelog.Critical(err)
		res["message"] = "create post is error"
		return
	}
	res["message"] = "Create post is successful"
	res["succeed"] = true
}

func (p *PostApi) Delete(ctx *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	repository := p.getRepository(ctx)
	defer p.WriteJSON(ctx, res)
	pid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	err = repository.Delete(uint(pid))
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["id"] = pid
	res["succeed"] = true
}

//编辑博文
func (p *PostApi) Edit(c *gin.Context) {
	id := p.stringToUnit(c.Param("id"))
	repository := p.getRepository(c)
	user, _ := p.AdminUser(c)
	post, err := repository.GetPostById(id, true)
	if err != nil {
		_ = seelog.Error(err)
		p.HandleMessage(c, err.Error())
		return
	}
	fmt.Println(post.Tags)
	c.HTML(http.StatusOK, "post/edit.html", p.RenderComments(gin.H{
		"post":   post,
		"user":   user,
		"token":  csrf.GetToken(c),
		"submit": "/admin/post/" + strconv.FormatInt(int64(post.ID), 10) + "/update",
	}))
}

//保存博文信息
func (p *PostApi) Update(ctx *gin.Context) {
	var res = gin.H{}
	defer p.WriteJSON(ctx, res)
	Id := p.stringToUnit(ctx.Param("id"))
	tags := ctx.PostForm("tags")
	repository := p.getRepository(ctx)
	tagRepository := repositories.NewTagRepository(ctx)
	post, err := repository.GetPostById(Id, false)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	err = ctx.ShouldBindWith(&post, binding.Form)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	isPublished := ctx.PostForm("isPublished")
	post.IsPublished = "on" == isPublished
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		if err = repository.Update(post); err != nil {
			return err
		}
		if err = tagRepository.DeletePostTagByPostId(post.ID); err != nil {
			return err
		}
		if err = p.CreatePostTag(ctx, tags, post.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["message"] = "Update post is successful"
	res["succeed"] = true
}

func PostGet(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil || !post.IsPublished {
		_ = api.HandlerError("post not published ", err)
		api.Handle404(c)
		return
	}
	post.View++
	_ = post.UpdateView()
	post.Tags, _ = models.ListTagByPostId(id)
	post.Comments, _ = models.ListCommentByPostID(id)
	user, _ := c.Get(api.CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "post/display.html", gin.H{
		"post": post,
		"user": user,
	})
}

func (p *PostApi) PostPublish(c *gin.Context) {
	var (
		err  error
		res  = gin.H{}
		post *models.Post
	)
	defer api.WriteJSON(c, res)
	Id := p.stringToUnit(c.Param("id"))
	repository := repositories.NewPostRepository(c)
	post, err = repository.GetPostById(Id, false)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	attr := map[string]interface{}{
		"is_published": !post.IsPublished,
	}
	err = repository.UpdateAttr(post, attr)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func (p *PostApi) getRepository(ctx *gin.Context) repositories.IPostRepository {
	return repositories.NewPostRepository(ctx)
}

func (p *PostApi) stringToUnit(id string) uint {
	Id, _ := strconv.ParseUint(id, 10, 64)
	return uint(Id)
}

func (p *PostApi) CreatePostTag(ctx *gin.Context, tags string, postId uint) error {
	tagArr := strings.Split(tags, ",")
	tagRepository := repositories.NewTagRepository(ctx)
	for _, tag := range tagArr {
		tagId, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			continue
		}
		pt := &models.PostTag{PostId: postId, TagId: uint(tagId),}
		err = tagRepository.PostTagCreate(pt)
		if err != nil {
			return err
		}
	}
	return nil
}
