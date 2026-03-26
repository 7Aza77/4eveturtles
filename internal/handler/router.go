package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter настраивает все пути (endpoints) нашего API
func InitRouter() *gin.Engine {
	r := gin.Default() // Создаем стандартный сервер Gin

	// Простая проверка, что сервер жив
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Группа для API мероприятий (пока пустая)
	events := r.Group("/api/v1/events")
	{
		events.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "Список мероприятий скоро будет здесь"})
		})
	}

	return r
}
