package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MessageProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "messages_processed_total",
		Help: "The total number of processed messages",
	}, []string{"tenant_id", "status"})

	WorkerCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "worker_count",
		Help: "Number of active workers per tenant",
	}, []string{"tenant_id"})

	ProcessingTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "message_processing_duration_seconds",
		Help:    "Time spent processing messages",
		Buckets: prometheus.DefBuckets,
	}, []string{"tenant_id"})
)
