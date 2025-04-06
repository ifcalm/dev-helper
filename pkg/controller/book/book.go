package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title":  "金刚经",
		"author": "侠名",
		"msg":    "是法平等, 无有高下",
	})
}
