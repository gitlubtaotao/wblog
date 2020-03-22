package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/controllers"
)

type SessionController struct {
	*controllers.BaseController
}

func (s *SessionController) AuthGet(c *gin.Context) {

}
