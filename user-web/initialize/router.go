package initialize

import (
	"github.com/gin-gonic/gin"
	userRouter "mxshop-api/user-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	ApiGroup := Router.Group("/u/v1")
	userRouter.InitUserRouter(ApiGroup)
	return Router
}
