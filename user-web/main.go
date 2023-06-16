package main

import (
	"fmt"
	"go.uber.org/zap"
	"mxshop-api/user-web/initialize"
)

func main() {
	port := 8021
	// 1. 初始化logger
	initialize.InitLogger()
	// 2. 初始化routers
	router := initialize.Routers()
	// 3. 启动服务器
	zap.S().Infof("启动服务器，端口号为:%d", port)

	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Panicf("启动失败:%s", err.Error())
	}
}
