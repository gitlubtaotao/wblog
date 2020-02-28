package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/controllers"
	"github.com/gitlubtaotao/wblog/models"
	
	"net/http"
)

//
func Home(c *gin.Context)  {
	
	user,_ := c.Get(controllers.CONTEXT_USER_KEY)
	c.HTML(http.StatusOK,"admin/home.html",gin.H{
		"pageCount":    CountPage(),
		"postCount":    CountPost(),
		"tagCount":     CountTag(),
		"commentCount": CountComment(),
		"user":         user,
		"comments":     MustListUnreadComment(),
	})
}
