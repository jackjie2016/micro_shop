package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"orders_web/models"
)

func IsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")

		currentUser := claims.(*models.CustomClaims)
		zap.S().Infof("用户权限：%d", currentUser)
		if currentUser.AuthorityID != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
