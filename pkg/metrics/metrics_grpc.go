package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type IMetricManagerGrpc interface {
	IncTotal(method string)
	IncTotalErr(method string)
	AddDuration(method string, duration time.Duration)
}

var _ IMetricManagerGrpc = (*MetricManagerGrpc)(nil)

type MetricManagerGrpc struct {
	serviceName     string
	total           *prometheus.CounterVec
	totalErr        *prometheus.CounterVec
	durationSummary *prometheus.HistogramVec
}

func NewMetricManagerGrpc(serviceName string) *MetricManagerGrpc {
	labelTotal := []string{"service", "method"}
	total := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
			Name: "grpc_call_total",
			Help: "count_of_all_call",
		}, labelTotal)
	prometheus.MustRegister(total)

	labelErrTotal := []string{"service", "method"}
	totalErr := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
			Name: "grpc_call_error_total",
			Help: "count_of_all_call_error",
		}, labelErrTotal)
	prometheus.MustRegister(totalErr)

	labelDurationSummary := []string{"service", "method"}
	buckets := []float64{0.001, 0.1, 1, 2, 5, 10, 100}
	durationSummary := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{ //nolint:exhaustruct
			Name:    "grpc_duration",
			Help:    "duration_of_all_call",
			Buckets: buckets,
		},
		labelDurationSummary,
	)
	prometheus.MustRegister(durationSummary)

	return &MetricManagerGrpc{
		serviceName:     serviceName,
		total:           total,
		totalErr:        totalErr,
		durationSummary: durationSummary,
	}
}

func (m *MetricManagerGrpc) IncTotal(method string) {
	m.total.WithLabelValues(m.serviceName, method).Inc()
}

func (m *MetricManagerGrpc) IncTotalErr(method string) {
	m.totalErr.WithLabelValues(m.serviceName, method).Inc()
}

func (m *MetricManagerGrpc) AddDuration(method string, duration time.Duration) {
	m.durationSummary.WithLabelValues(m.serviceName, method).Observe(duration.Seconds())
}
