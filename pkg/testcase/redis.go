package testcase

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"dev-helper/pkg/database/rds"
)

func SetString(c *gin.Context) {

	key := c.Param("key")
	value := c.Param("value")

	err := rds.SetKey(key, value)
	if err != nil {
		log.Fatalf("could not get key: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": value,
	})
}

func GetString(c *gin.Context) {
	key := c.Param("key")

	val, err := rds.GetKey(key)
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "键不存在"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": val,
	})
}
