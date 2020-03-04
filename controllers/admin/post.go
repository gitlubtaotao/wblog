package admin

import (
	"github.com/gitlubtaotao/wblog/controllers"
	"github.com/gitlubtaotao/wblog/repositories"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
	
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
)

type PostController struct {
}

func (p *PostController) Index(c *gin.Context) {
	res := repositories.NewPostRepository()
	posts, _ := res.ListAll("")
	user := controllers.GetUser(c)
	c.HTML(http.StatusOK, "admin/post.html", gin.H{
		"posts":    posts,
		"Active":   "posts",
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func (p *PostController) New(c *gin.Context) {
	c.HTML(http.StatusOK, "post/new.html", nil)
}

func (p *PostController) Create(c *gin.Context) {
	tags := c.PostForm("tags")
	array := c.PostFormMap("post")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	post := &models.Post{
		Title:       array["title"],
		Body:        array["body"],
		IsPublished: published,
	}
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		err := post.Insert()
		if err != nil {
			return err
		}
		if len(tags) > 0 {
			tagArr := strings.Split(tags, ",")
			for _, tag := range tagArr {
				tagId, err := strconv.ParseUint(tag, 10, 64)
				if err != nil {
					continue
				}
				pt := &models.PostTag{
					PostId: post.ID,
					TagId:  uint(tagId),
				}
				err = pt.Insert()
				if err != nil {
					return nil
				}
			}
		}
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		// return nil will commit
		return nil
	})
	if err != nil {
		c.HTML(http.StatusOK, "post/new.html", gin.H{
			"post":    post,
			"message": err.Error(),
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/posts")
}

func PostGet(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil || !post.IsPublished {
		_ = controllers.HandlerError("post not published ", err)
		controllers.Handle404(c)
		return
	}
	post.View++
	_ = post.UpdateView()
	post.Tags, _ = models.ListTagByPostId(id)
	post.Comments, _ = models.ListCommentByPostID(id)
	user, _ := c.Get(controllers.CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "post/display.html", gin.H{
		"post": post,
		"user": user,
	})
}

func PostNew(c *gin.Context) {
	c.HTML(http.StatusOK, "post/new.html", nil)
}

func PostEdit(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil {
		controllers.Handle404(c)
		return
	}
	post.Tags, _ = models.ListTagByPostId(id)
	c.HTML(http.StatusOK, "post/modify.html", gin.H{
		"post": post,
	})
}

func PostUpdate(c *gin.Context) {
	id := c.Param("id")
	tags := c.PostForm("tags")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished")
	published := "on" == isPublished
	
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		controllers.Handle404(c)
		return
	}
	
	post := &models.Post{
		Title:       title,
		Body:        body,
		IsPublished: published,
	}
	post.ID = uint(pid)
	err = post.Update()
	if err != nil {
		c.HTML(http.StatusOK, "post/modify.html", gin.H{
			"post":    post,
			"message": err.Error(),
		})
		return
	}
	// 删除tag
	_ = models.DeletePostTagByPostId(post.ID)
	// 添加tag
	if len(tags) > 0 {
		tagArr := strings.Split(tags, ",")
		for _, tag := range tagArr {
			tagId, err := strconv.ParseUint(tag, 10, 64)
			if err != nil {
				continue
			}
			pt := &models.PostTag{
				PostId: post.ID,
				TagId:  uint(tagId),
			}
			pt.Insert()
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/admin/post")
}

func PostPublish(c *gin.Context) {
	var (
		err  error
		res  = gin.H{}
		post *models.Post
	)
	defer controllers.WriteJSON(c, res)
	id := c.Param("id")
	post, err = models.GetPostById(id)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	post.IsPublished = !post.IsPublished
	err = post.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func PostDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer controllers.WriteJSON(c, res)
	id := c.Param("id")
	pid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	post := &models.Post{}
	post.ID = uint(pid)
	err = post.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	_ = models.DeletePostTagByPostId(uint(pid))
	res["succeed"] = true
}
