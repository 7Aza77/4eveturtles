package handler

import (
	"goevent/internal/usecase"
	"goevent/pkg/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	_ "goevent/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	authHandler         *AuthHandler
	eventHandler        *EventHandler
	registrationHandler *RegistrationHandler
	tokenManager        auth.TokenManager
}

func NewHandler(authUseCase usecase.AuthUseCase, eventUseCase usecase.EventUseCase, regUseCase usecase.RegistrationUseCase, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		authHandler:         NewAuthHandler(authUseCase),
		eventHandler:        NewEventHandler(eventUseCase),
		registrationHandler: NewRegistrationHandler(regUseCase),
		tokenManager:        tokenManager,
	}
}

// InitRouter настраивает все пути (endpoints) нашего API
func (h *Handler) InitRouter(rdb *redis.Client) *gin.Engine {
	r := gin.Default() // Создаем стандартный сервер Gin

	// Middleware
	r.Use(h.rateLimit(rdb, 100, time.Minute)) // 100 запросов в минуту

	// Простая проверка, что сервер жив
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", h.authHandler.signUp)
		auth.POST("/sign-in", h.authHandler.signIn)
	}

	// Группа для API мероприятий
	events := r.Group("/api/v1/events")
	{
		events.GET("/", h.eventHandler.list)
		events.GET("/:id", h.eventHandler.getByID)

		// Защищенные маршруты
		protected := events.Group("/", h.userIdentity(h.tokenManager))
		{
			protected.POST("/", h.eventHandler.create)
			protected.PUT("/:id", h.eventHandler.update)
			protected.DELETE("/:id", h.eventHandler.delete)

			// Регистрация
			protected.POST("/:id/register", h.registrationHandler.register)
			protected.DELETE("/:id/unregister", h.registrationHandler.cancel)
		}
	}

	return r
}
