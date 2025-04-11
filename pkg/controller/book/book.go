package book

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"dev-helper/pkg/crontab/btc"
)

func GetBook(c *gin.Context) {
	for i := 0; i <= 100; i++ {
		address := btc.Btc()
		btc.GetBitcoinBalance(address)
		time.Sleep(1 * time.Second)
	}

	c.JSON(http.StatusOK, gin.H{
		"title":  "金刚经",
		"author": "侠名",
		"msg":    "是法平等, 无有高下",
	})
}
