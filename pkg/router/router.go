package router

import (
	"dev-helper/pkg/controller/book"
	"dev-helper/pkg/testcase"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() {
	// 创建一个默认的Gin路由
	r := gin.Default()

	// 定义根路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// 定义另一个路由
	r.GET("/book/list", book.GetBook)

	v1 := r.Group("/redis")
	v1.POST("/string/set/:key/:value", testcase.SetString)
	v1.GET("/string/get/:key", testcase.GetString)

	// 启动服务器，监听本地8080端口
	r.Run(":8080")
}
