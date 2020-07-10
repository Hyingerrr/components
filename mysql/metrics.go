package mysql

import "github.com/prometheus/client_golang/prometheus"

const (
	_PromNamespace = "technology_finance"
)

func init() {
	prometheus.MustRegister(DBMiss)
	prometheus.MustRegister(DBDuration)
	prometheus.MustRegister(DBError)
	prometheus.MustRegister(DBStats)
}

var (
	DBDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: _PromNamespace,
			Subsystem: "mysql",
			Name:      "db_duration_seconds",
			Help:      "db duration distribution",
			Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"table", "sql"},
	)

	DBMiss = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: _PromNamespace,
			Subsystem: "mysql",
			Name:      "db_miss",
			Help:      "Total number of db miss",
		},
		[]string{"schema", "table"},
	)

	DBError = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: _PromNamespace,
			Subsystem: "mysql",
			Name:      "db_error",
			Help:      "Total number of db error",
		},
		[]string{"schema", "table"},
	)

	DBStats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: _PromNamespace,
			Subsystem: "mysql",
			Name:      "db_stats",
			Help:      "db statistics",
		},
		[]string{"schema", "stats"},
	)
)
