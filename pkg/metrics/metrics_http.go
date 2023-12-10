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
	totalStatuses   *prometheus.CounterVec
	durationSummary *prometheus.SummaryVec
}

func NewMetricManagerHTTP(serviceName string) *MetricManagerHTTP {
	labelTotalStatuses := []string{"path", "method", "status"}
	totalStatuses := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
			Namespace: serviceName,
			Name:      "http_request_statuses_total",
			Help:      "count_of_all_request_with_status",
		}, labelTotalStatuses)
	prometheus.MustRegister(totalStatuses)

	labelDurationSummary := []string{"path", "method", "status"}
	durationSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{ //nolint:exhaustruct
			Namespace: serviceName,
			Name:      "http_duration",
			Help:      "duration_of_all_request_with_status",
			Objectives: map[float64]float64{
				0.5:  0.05,  //nolint:gomnd
				0.9:  0.01,  //nolint:gomnd
				0.99: 0.001, //nolint:gomnd
			},
		},
		labelDurationSummary,
	)
	prometheus.MustRegister(durationSummary)

	return &MetricManagerHTTP{totalStatuses: totalStatuses, durationSummary: durationSummary}
}

func (m *MetricManagerHTTP) IncreaseTotal(path, method, status string) {
	m.totalStatuses.WithLabelValues(path, method, status).Inc()
}

func (m *MetricManagerHTTP) AddDuration(path, method, status string, duration time.Duration) {
	m.durationSummary.WithLabelValues(path, method, status).Observe(duration.Seconds())
}
