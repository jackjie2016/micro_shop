package router

import (
	"github.com/gin-gonic/gin"
	"goods_web/api/goods"
	"goods_web/middlewares"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsGroup := Router.Group("goods").Use(middlewares.Tracing())
	{
		GoodsGroup.GET("", goods.List)
		GoodsGroup.POST("", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.New) //改接口需要管理员权限
		GoodsGroup.GET("/:id", goods.Detail)
		GoodsGroup.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.Delete)
		GoodsGroup.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.Update)
		GoodsGroup.PATCH("/:id/UpdateStatus", middlewares.JWTAuth(), middlewares.IsAdmin(), goods.UpdateStatus)
		GoodsGroup.GET("/:id/stocks", goods.Stocks)
	}

}
