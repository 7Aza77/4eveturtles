package handler

import (
	"goevent/internal/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricsMiddleware собирает HTTP метрики для каждого запроса.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method

		metrics.HttpRequestsTotal.With(prometheus.Labels{
			"method":      method,
			"path":        path,
			"status_code": status,
		}).Inc()

		metrics.HttpRequestDuration.With(prometheus.Labels{
			"method": method,
			"path":   path,
		}).Observe(duration)
	}
}
