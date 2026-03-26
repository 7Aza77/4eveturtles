package handler

import (
	"errors"
	"goevent/pkg/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	roleCtx             = "userRole"
)

func (h *Handler) userIdentity(tokenManager auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		if len(headerParts[1]) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			return
		}

		userId, role, err := tokenManager.Parse(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set(userCtx, userId)
		c.Set(roleCtx, role)
	}
}

func (h *Handler) roleRestriction(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(roleCtx)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
			return
		}

		roleStr := role.(string)
		for _, r := range roles {
			if r == roleStr {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
	}
}
