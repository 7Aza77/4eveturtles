package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func (h *Handler) rateLimit(rdb *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "rate_limit:" + c.ClientIP()
		ctx := context.Background()

		count, err := rdb.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.Next() // В случае ошибки редиса пропускаем (fail-open)
			return
		}

		if count >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}

		pipe := rdb.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, _ = pipe.Exec(ctx)

		c.Next()
	}
}
