package router

import (
	"github.com/gin-gonic/gin"
	"orders_web/api/orders"
	"orders_web/api/pay"
	"orders_web/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("order").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", orders.List)          // 订单列表
		OrderRouter.GET("/:id", orders.Detail)    // 订单详情
		OrderRouter.POST("", orders.New)          //新建订单
		OrderRouter.PUT("/:id", orders.Update)    //修改订单
		OrderRouter.DELETE("/:id", orders.Delete) // 订单删除

	}
	PayRouter := Router.Group("pay")
	{
		PayRouter.POST("alipay/notify", pay.Notify)
	}
}
