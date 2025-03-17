package models

import (
	"time"

	"gorm.io/gorm"
)

// Book 书籍模型
type Book struct {
	gorm.Model
	Title       string `gorm:"size:255;not null" form:"title" binding:"required"`
	Author      string `gorm:"size:255;not null" form:"author" binding:"required"`
	Publisher   string `gorm:"size:255" form:"publisher"`
	PublishYear int    `gorm:"type:int" form:"publish_year"`
	ISBN        string `gorm:"size:13;unique" form:"isbn" binding:"required"`
	Description string `gorm:"type:text" form:"description"`
	Quantity    int    `gorm:"type:int;default:0" form:"quantity"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

// TableName 指定表名
func (Book) TableName() string {
	return "books"
}

// CreateBook 创建新书籍
func CreateBook(db *gorm.DB, book *Book) error {
	return db.Create(book).Error
}

// GetAllBooks 获取所有书籍
func GetAllBooks(db *gorm.DB) ([]Book, error) {
	var books []Book
	err := db.Find(&books).Error
	return books, err
}

// GetBookByID 根据ID获取书籍
func GetBookByID(db *gorm.DB, id uint) (Book, error) {
	var book Book
	err := db.Where("id = ?", id).First(&book).Error
	return book, err
}

// UpdateBook 更新书籍信息
func UpdateBook(db *gorm.DB, book *Book) error {
	return db.Save(book).Error
}

// DeleteBook 删除书籍
func DeleteBook(db *gorm.DB, id uint) error {
	return db.Where("id = ?", id).Delete(&Book{}).Error
}
