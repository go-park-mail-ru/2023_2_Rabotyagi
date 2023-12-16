package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type IMetricManagerHTTP interface {
	IncreaseTotal(path, method, status string)
	AddDuration(path, method, status string, duration time.Duration)
}

var _ IMetricManagerHTTP = (*MetricManagerHTTP)(nil)

type MetricManagerHTTP struct {
	serviceName     string
	totalStatuses   *prometheus.CounterVec
	durationSummary *prometheus.HistogramVec
}

func NewMetricManagerHTTP(serviceName string) *MetricManagerHTTP {
	labelTotalStatuses := []string{"service", "path", "method", "status"}
	totalStatuses := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
			Name: "http_request_statuses_total",
			Help: "count_of_all_request_with_status",
		}, labelTotalStatuses)
	prometheus.MustRegister(totalStatuses)

	labelDurationSummary := []string{"service", "path", "method", "status"}
	buckets := []float64{0.001, 0.1, 1, 2, 5, 10, 100}
	durationSummary := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{ //nolint:exhaustruct
			Name:    "http_duration",
			Help:    "duration_of_all_request_with_status",
			Buckets: buckets,
		},
		labelDurationSummary,
	)
	prometheus.MustRegister(durationSummary)

	return &MetricManagerHTTP{
		serviceName:     serviceName,
		totalStatuses:   totalStatuses,
		durationSummary: durationSummary,
	}
}

func (m *MetricManagerHTTP) IncreaseTotal(path, method, status string) {
	m.totalStatuses.WithLabelValues(m.serviceName, path, method, status).Inc()
}

func (m *MetricManagerHTTP) AddDuration(path, method, status string, duration time.Duration) {
	m.durationSummary.WithLabelValues(m.serviceName, path, method, status).Observe(duration.Seconds())
}
