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
	total           *prometheus.CounterVec
	totalErr        *prometheus.CounterVec
	durationSummary *prometheus.SummaryVec
}

func NewMetricManagerGrpc(serviceName string) *MetricManagerGrpc {
	labelTotal := []string{"method"}
	total := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
			Namespace: serviceName,
			Name:      "grpc_call_total",
			Help:      "count_of_all_call",
		}, labelTotal)
	prometheus.MustRegister(total)

	labelErrTotal := []string{"method"}
	totalErr := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct
			Namespace: serviceName,
			Name:      "grpc_call_error_total",
			Help:      "count_of_all_call_error",
		}, labelErrTotal)
	prometheus.MustRegister(totalErr)

	labelDurationSummary := []string{"method"}
	durationSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{ //nolint:exhaustruct
			Namespace: serviceName,
			Name:      "grpc_duration",
			Help:      "duration_of_all_call",
			Objectives: map[float64]float64{
				0.5:  0.05,  //nolint:gomnd
				0.9:  0.01,  //nolint:gomnd
				0.99: 0.001, //nolint:gomnd
			},
		},
		labelDurationSummary,
	)
	prometheus.MustRegister(durationSummary)

	return &MetricManagerGrpc{total: total, totalErr: totalErr, durationSummary: durationSummary}
}

func (m *MetricManagerGrpc) IncTotal(method string) {
	m.total.WithLabelValues(method).Inc()
}

func (m *MetricManagerGrpc) IncTotalErr(method string) {
	m.totalErr.WithLabelValues(method).Inc()
}

func (m *MetricManagerGrpc) AddDuration(method string, duration time.Duration) {
	m.durationSummary.WithLabelValues(method).Observe(duration.Seconds())
}
