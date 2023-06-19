package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/models"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, _ := context.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			context.JSON(403, gin.H{
				"msg": "无权限",
			})
			context.Abort()
			return
		}
		context.Next()
	}
}
