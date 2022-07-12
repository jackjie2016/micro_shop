package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"orders_web/middlewares"

	router2 "orders_web/router"
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

	router2.InitCartRouter(ApiGroup)
	router2.InitOrderRouter(ApiGroup)

	return Router
}
