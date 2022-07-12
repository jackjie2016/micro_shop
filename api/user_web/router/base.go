package router

import (
	"github.com/gin-gonic/gin"
	"user_web/api"
)

func InitBaseRouter(Router *gin.RouterGroup) *gin.RouterGroup {
	BaseGroup := Router.Group("base")
	{
		BaseGroup.GET("captcha", api.GetCapthcha)
		BaseGroup.POST("send_sms", api.SendSms)
	}
	return BaseGroup

}
