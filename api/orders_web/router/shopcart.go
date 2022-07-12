package router

import (
	"github.com/gin-gonic/gin"

	"orders_web/api/shop_cart"
	"orders_web/middlewares"
)

func InitCartRouter(Router *gin.RouterGroup) {
	CartGroup := Router.Group("shop_cart").Use(middlewares.JWTAuth())
	{
		CartGroup.GET("", shop_cart.List)          // 订单列表
		CartGroup.POST("", shop_cart.New)          //新建订单
		CartGroup.PATCH("/:id", shop_cart.Update)  //修改条目
		CartGroup.DELETE("/:id", shop_cart.Delete) // 订单删除
	}
}
