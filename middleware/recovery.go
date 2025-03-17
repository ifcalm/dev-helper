package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				gin.DefaultErrorWriter.Write([]byte(
					"[GIN] " + c.Request.URL.Path + " | " + c.Request.Method + " | " + err.(string) + "\n",
				))

				// 返回500错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
		}()

		c.Next()
	}
}
