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
	organizationHandler *OrganizationHandler
	tagHandler          *TagHandler
	tokenManager        auth.TokenManager
}

func NewHandler(
	authUseCase usecase.AuthUseCase,
	eventUseCase usecase.EventUseCase,
	regUseCase usecase.RegistrationUseCase,
	tokenManager auth.TokenManager,
	orgUseCase usecase.OrganizationUseCase,
	tagUseCase usecase.TagUseCase,
) *Handler {
	return &Handler{
		authHandler:         NewAuthHandler(authUseCase),
		eventHandler:        NewEventHandler(eventUseCase),
		registrationHandler: NewRegistrationHandler(regUseCase),
		organizationHandler: NewOrganizationHandler(orgUseCase),
		tagHandler:          NewTagHandler(tagUseCase),
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

	// --- Events ---
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

			// Tags on events
			authorized.POST(
				"/:id/tags",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.tagHandler.addTagsToEvent,
			)

			authorized.DELETE(
				"/:id/tags",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.tagHandler.removeTagsFromEvent,
			)
		}
	}

	// --- Organizations ---
	orgs := r.Group("/api/v1/organizations")
	{
		orgs.GET("/", h.organizationHandler.list)
		orgs.GET("/:id", h.organizationHandler.getByID)

		authorizedOrg := orgs.Group("/")
		authorizedOrg.Use(h.userIdentity(h.tokenManager))
		{
			authorizedOrg.POST(
				"/",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.organizationHandler.create,
			)

			authorizedOrg.PUT(
				"/:id",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.organizationHandler.update,
			)

			authorizedOrg.DELETE(
				"/:id",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.organizationHandler.delete,
			)
		}
	}

	// --- Tags ---
	tags := r.Group("/api/v1/tags")
	{
		tags.GET("/", h.tagHandler.list)

		authorizedTag := tags.Group("/")
		authorizedTag.Use(h.userIdentity(h.tokenManager))
		{
			authorizedTag.POST(
				"/",
				h.roleRestriction(string(entity.RoleAdmin), string(entity.RoleModerator)),
				h.tagHandler.create,
			)
		}
	}

	return r
}


