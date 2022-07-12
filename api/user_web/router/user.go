package router

import (
	"github.com/gin-gonic/gin"
	"user_web/api"
	"user_web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) *gin.RouterGroup {
	UserGroup := Router.Group("user")
	{
		//,middlewares.IsAdmin()
		UserGroup.GET("list", middlewares.JWTAuth(), api.GetUserlist)
		UserGroup.POST("pwd_login", api.PassWordLogin)
		UserGroup.POST("register", api.RegisterForm)
	}
	return UserGroup

}
