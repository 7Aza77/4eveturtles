package handler

import (
	"goevent/internal/entity"
	"goevent/internal/usecase"
	"goevent/pkg/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func NewHandler(
	authUseCase usecase.AuthUseCase,
	eventUseCase usecase.EventUseCase,
	regUseCase usecase.RegistrationUseCase,
	tokenManager auth.TokenManager,
) *Handler {
	return &Handler{
		authHandler:         NewAuthHandler(authUseCase),
		eventHandler:        NewEventHandler(eventUseCase),
		registrationHandler: NewRegistrationHandler(regUseCase),
		tokenManager:        tokenManager,
	}
}

func (h *Handler) InitRouter(rdb *redis.Client) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Метрики и swagger — до rate limiting
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Глобальные middleware
	r.Use(MetricsMiddleware())
	r.Use(h.rateLimit(rdb, 100, time.Minute))

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "pong"})
	})

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/sign-up", h.authHandler.signUp)
		authGroup.POST("/sign-in", h.authHandler.signIn)
	}

	events := r.Group("/api/v1/events")
	{
		events.GET("", h.eventHandler.list)
		events.GET("/:id", h.eventHandler.getByID)

		authorized := events.Group("")
		authorized.Use(h.userIdentity(h.tokenManager))
		{
			authorized.POST(
				"",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.eventHandler.create,
			)
			authorized.PUT(
				"/:id",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.eventHandler.update,
			)
			authorized.DELETE(
				"/:id",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.eventHandler.delete,
			)
			authorized.POST(
				"/:id/register",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator), string(entity.RoleStudent)),
				h.registrationHandler.register,
			)
			authorized.DELETE(
				"/:id/unregister",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator), string(entity.RoleStudent)),
				h.registrationHandler.cancel,
			)
		}
	}

	return r
}
