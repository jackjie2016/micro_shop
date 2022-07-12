package router

import (
	"github.com/gin-gonic/gin"
	"goods_web/middlewares"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	_ = Router.Group("base").Use(middlewares.Tracing())
	{

	}

}
