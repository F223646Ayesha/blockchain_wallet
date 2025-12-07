package utils

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

func Error(c *gin.Context, msg string) {
	c.JSON(400, gin.H{
		"success": false,
		"message": msg,
	})
}
