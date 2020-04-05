package client

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/models"
	"github.com/gitlubtaotao/wblog/repositories"
	"net/http"
	"strconv"
	"sync"
)

type PostApi struct {
	*UtilApi
}

func (p *PostApi) Show(ctx *gin.Context) {
	repository := repositories.NewPostRepository(ctx)
	id, _ := strconv.Atoi(ctx.Param("id"))
	post, err := repository.GetPostById(uint(id), true)
	if err != nil {
		p.HandleMessage(ctx, err.Error())
		return
	}
	if !post.IsPublished {
		p.HandleMessage(ctx, "Post is not exist")
		return
	}
	var sy sync.WaitGroup
	sy.Add(1)
	go func(view int) {
		attr := map[string]interface{}{
			"view": view,
		}
		_ = repository.UpdateAttr(post, attr)
		sy.Done()
	}(post.View + 1)
	sy.Wait()
	post.Comments, _ = p.commentList(ctx, post.ID)
	user, _ := p.ClientUser(ctx)
	ctx.HTML(http.StatusOK, "post/display.html", gin.H{
		"post": post,
		"user": user,
	})
}

func (p *PostApi) Index(ctx *gin.Context) {

}

func (p *PostApi) commentList(ctx *gin.Context, postId uint) ([]*models.Comment, error) {
	repository := repositories.NewCommentRepository()
	comments, err := repository.ListCommentByPostID(postId)
	return comments, err
}
