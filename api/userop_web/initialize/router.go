package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userop_web/middlewares"
	"userop_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())

	// “/up/v1”
	ApiGroup := Router.Group("/v1")
	router.InitUserFavRouter(ApiGroup)
	router.InitMessageRouter(ApiGroup)
	router.InitAddressRouter(ApiGroup)

	return Router
}
