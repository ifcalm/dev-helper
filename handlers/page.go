package handlers

import (
	"github.com/gin-gonic/gin"
)

// HomePage 首页
func HomePage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "书籍管理系统",
	})
}

// BooksPage 书籍列表页
func BooksPage(c *gin.Context) {
	c.HTML(200, "books.html", gin.H{
		"title": "书籍列表",
	})
}
