package handler

import (
	"goevent/internal/entity"
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
	r := gin.Default()

	r.Use(h.rateLimit(rdb, 100, time.Minute))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/sign-up", h.authHandler.signUp)
		authGroup.POST("/sign-in", h.authHandler.signIn)
	}

	events := r.Group("/api/v1/events")
	{
		events.GET("/", h.eventHandler.list)
		events.GET("/:id", h.eventHandler.getByID)

		authorized := events.Group("/")
		authorized.Use(h.userIdentity(h.tokenManager))
		{
			authorized.POST(
				"/",
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
