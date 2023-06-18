package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("初始化user路由")
	{
		UserRouter.GET("list", api.GetUserList)
		UserRouter.GET("pwd_login", api.PassWordLogin)
	}
}
