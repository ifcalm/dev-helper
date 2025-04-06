package main

import (
	"dev-helper/pkg/database/rds"
	"dev-helper/pkg/router"
)

func main() {
	// redis初始化
	rds.InitRedis()
	// 启动路由
	router.SetupRoutes()
}
