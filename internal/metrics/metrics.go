package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP метрики
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goevent",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goevent",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"method", "path"},
	)

	// Бизнес метрики
	EventsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "goevent",
		Name:      "events_created_total",
		Help:      "Total number of created events",
	})

	EventsDeletedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "goevent",
		Name:      "events_deleted_total",
		Help:      "Total number of deleted events",
	})

	RegistrationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goevent",
			Name:      "registrations_total",
			Help:      "Total number of event registrations",
		},
		[]string{"action"},
	)

	ActiveUsersTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "goevent",
		Name:      "active_users_total",
		Help:      "Total number of active (logged-in) users",
	})

	// Инфраструктурные метрики
	CacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goevent",
			Name:      "cache_hits_total",
			Help:      "Total number of cache hits",
		},
		[]string{"key_type"},
	)

	CacheMissesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goevent",
			Name:      "cache_misses_total",
			Help:      "Total number of cache misses",
		},
		[]string{"key_type"},
	)

	DatabaseQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "goevent",
			Name:      "database_queries_total",
			Help:      "Total number of database queries",
		},
		[]string{"operation", "entity"},
	)
)
