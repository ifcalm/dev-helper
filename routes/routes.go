package routes

import (
	"dev-helper/handlers"
	"dev-helper/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 设置静态文件目录
	r.Static("/static", "./static")

	// 设置中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// API路由组
	api := r.Group("/api")
	{
		// 书籍相关路由
		books := api.Group("/books")
		{
			books.GET("", handlers.GetBooks)
			books.GET("/:id", handlers.GetBook)
			books.POST("", handlers.CreateBook)
			books.PUT("/:id", handlers.UpdateBook)
			books.DELETE("/:id", handlers.DeleteBook)
		}
	}

	// 页面路由
	r.GET("/", handlers.HomePage)
	r.GET("/books", handlers.BooksPage)
}
