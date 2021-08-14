package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Handlers struct {
	DbConn    *gorm.DB
	RedisConn *redis.Client
}


func (e *Handlers) HomeHandler(c *gin.Context){
	c.JSON(200, gin.H{"hello": "world"})
}