package route

import (
	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.Default()
	r.GET("/select/domain/beian")
}
