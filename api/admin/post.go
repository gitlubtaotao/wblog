package admin

import (
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/jinzhu/gorm"
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
	user, _ := p.CurrentUser(ctx)
	renderJson := p.RenderComments(gin.H{
		"posts":    posts,
		"Active":   "posts",
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
	ctx.HTML(http.StatusOK, "post/index.html", renderJson)
}

func (p *PostApi) New(c *gin.Context) {
	user, _ := p.CurrentUser(c)
	c.HTML(http.StatusOK, "post/edit.html", p.RenderComments(gin.H{
		"user": user,
	}))
}

func (p *PostApi) Create(c *gin.Context) {
	repository := p.getRepository(c)
	err := repository.Create()
	if err != nil {
		p.HandleMessage(c, err.Error())
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/posts")
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
	res["succeed"] = true
}

//编辑博文
func (p *PostApi) Edit(c *gin.Context) {
	id := p.stringToUnit(c.Param("id"))
	repository := p.getRepository(c)
	tagRepository := repositories.NewTagRepository(c)
	post, err := repository.GetPostById(id, false)
	if err != nil {
		p.HandleMessage(c, err.Error())
		return
	}
	post.Tags, _ = tagRepository.ListTagByPostId(id)
	c.HTML(http.StatusOK, "post/modify.html", gin.H{
		"post": post,
	})
}

//保存博文信息
func (p *PostApi) Update(ctx *gin.Context) {
	Id := p.stringToUnit(ctx.Param("id"))
	tags := ctx.PostForm("tags")
	repository := p.getRepository(ctx)
	tagRepository := repositories.NewTagRepository(ctx)
	p.post = repository
	p.tag = tagRepository
	post, err := repository.GetPostById(Id, false)
	if err != nil {
		p.HandleMessage(ctx, err.Error())
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
		if err = p.updateTag(tags, post.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		p.HandleMessage(ctx, err.Error())
		return
	}
	ctx.Redirect(http.StatusMovedPermanently, "/admin/posts")
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

func (p *PostApi) updateTag(tags string, id uint) error {
	if len(tags) < 0 {
		return nil
	}
	tagArr := strings.Split(tags, ",")
	for _, tag := range tagArr {
		tagId, err := strconv.ParseUint(tag, 10, 64)
		if err != nil {
			return err
		}
		pt := &models.PostTag{
			PostId: id,
			TagId:  uint(tagId),
		}
		if err = p.tag.PostTagCreate(pt); err != nil {
			return err
		}
	}
	return nil
}

func (p *PostApi) getRepository(ctx *gin.Context) repositories.IPostRepository {
	return repositories.NewPostRepository(ctx)
}

func (p *PostApi) stringToUnit(id string) uint {
	Id, _ := strconv.ParseUint(id, 10, 64)
	return uint(Id)
}
