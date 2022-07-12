package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user_web/middlewares"

	router2 "user_web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	ApiGroup := Router.Group("/v1")

	router2.InitBaseRouter(ApiGroup)
	router2.InitUserRouter(ApiGroup)

	return Router
}
