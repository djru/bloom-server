package handlers

import "github.com/gin-gonic/gin"

func MalformedErr(c *gin.Context) {
	c.JSON(400, gin.H{"status": "failed", "message": "bad payload"})
	return
}