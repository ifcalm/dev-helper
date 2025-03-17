package controllers

import (
	"net/http"

	"dev-helper/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookController 处理书籍相关的请求
type BookController struct {
	db *gorm.DB
}

// NewBookController 创建新的BookController
func NewBookController(db *gorm.DB) *BookController {
	return &BookController{db: db}
}

// ListBooks 显示书籍列表页面
func (bc *BookController) ListBooks(c *gin.Context) {
	var books []models.Book
	bc.db.Find(&books)
	c.HTML(http.StatusOK, "books.html", gin.H{
		"title": "书籍列表",
		"books": books,
	})
}

// NewBookForm 显示添加新书籍的表单页面
func (bc *BookController) NewBookForm(c *gin.Context) {
	c.HTML(http.StatusOK, "book_form.html", gin.H{
		"title": "添加新书籍",
	})
}

// EditBookForm 显示编辑书籍的表单页面
func (bc *BookController) EditBookForm(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := bc.db.First(&book, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "错误",
			"error": "找不到该书籍",
		})
		return
	}
	c.HTML(http.StatusOK, "book_form.html", gin.H{
		"title": "编辑书籍",
		"book":  book,
	})
}

// CreateBook 创建新书籍
func (bc *BookController) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBind(&book); err != nil {
		c.HTML(http.StatusBadRequest, "book_form.html", gin.H{
			"title": "添加新书籍",
			"error": "表单数据无效",
		})
		return
	}

	if err := bc.db.Create(&book).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "book_form.html", gin.H{
			"title": "添加新书籍",
			"error": "创建书籍失败",
		})
		return
	}

	c.Redirect(http.StatusFound, "/books")
}

// UpdateBook 更新书籍信息
func (bc *BookController) UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := bc.db.First(&book, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "错误",
			"error": "找不到该书籍",
		})
		return
	}

	if err := c.ShouldBind(&book); err != nil {
		c.HTML(http.StatusBadRequest, "book_form.html", gin.H{
			"title": "编辑书籍",
			"book":  book,
			"error": "表单数据无效",
		})
		return
	}

	if err := bc.db.Save(&book).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "book_form.html", gin.H{
			"title": "编辑书籍",
			"book":  book,
			"error": "更新书籍失败",
		})
		return
	}

	c.Redirect(http.StatusFound, "/books")
}

// DeleteBook 删除书籍
func (bc *BookController) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	if err := bc.db.Delete(&models.Book{}, id).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title": "错误",
			"error": "删除书籍失败",
		})
		return
	}
	c.Redirect(http.StatusFound, "/books")
}
