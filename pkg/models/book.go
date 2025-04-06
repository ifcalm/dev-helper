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
