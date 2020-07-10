package redis

import "github.com/prometheus/client_golang/prometheus"

const (
	_PromNamespace = "technology_finance"

	GORedisStatus      = "200"
	GoRedisMissStatus  = "404"
	GoRedisErrorStatus = "400"
)

func init() {
	prometheus.MustRegister(GoRedisError)
	prometheus.MustRegister(GoRedisDuration)
	prometheus.MustRegister(GoRedisStats)
	prometheus.MustRegister(GoRedisTotal)
}

var (
	GoRedisError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: _PromNamespace,
			Subsystem: "redis",
			Name:      "redis_error",
			Help:      "Total number of redis error",
		},
		[]string{"schema", "cmd", "status"},
	)

	GoRedisDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: _PromNamespace,
			Subsystem: "redis",
			Name:      "redis_duration",
			Help:      "redis request time",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1},
		}, []string{"schema", "cmd", "cost"},
	)

	GoRedisStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: _PromNamespace,
			Subsystem: "redis",
			Name:      "redis_pool",
			Help:      "pool's statistics",
		},
		[]string{"name", "pool"},
	)

	GoRedisTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: _PromNamespace,
			Subsystem: "redis",
			Name:      "redis_stats_total",
			Help:      "Number of hello requests in total",
		},
		[]string{"schema", "cmd"},
	)
)
