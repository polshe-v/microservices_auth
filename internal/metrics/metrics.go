package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "auth"
	appName   = "auth_service"
	subsystem = "grpc"
)

// Metrics contains application metrics.
type Metrics struct {
	requestCounter        prometheus.Counter
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

var metrics *Metrics

// Init creates metrics object for metrics operations.
func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_requests_total",
				Help:      "Number of requests to server",
			},
		),
		responseCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_responses_total",
				Help:      "Number of responses from server",
			},
			[]string{"status", "method"},
		),
		histogramResponseTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      appName + "_histogram_response_time_seconds",
				Help:      "Server response time",
				Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
			},
			[]string{"status"},
		),
	}

	return nil
}

// IncRequestCounter increases number of requests to application.
func IncRequestCounter() {
	metrics.requestCounter.Inc()
}

// IncResponseCounter increases number of responses from application labelling with method name and response status.
func IncResponseCounter(status string, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}

// HistogramResponseTimeObserve writes application response time.
func HistogramResponseTimeObserve(status string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status).Observe(time)
}
