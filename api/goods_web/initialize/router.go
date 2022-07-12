package initialize

import (
	"github.com/gin-gonic/gin"
	"goods_web/middlewares"
	"net/http"

	router2 "goods_web/router"
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
	router2.InitGoodsRouter(ApiGroup)
	router2.InitCategoryRouter(ApiGroup)
	router2.InitBrandRouter(ApiGroup)
	router2.InitBannerRouter(ApiGroup)
	return Router
}
