package observability

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	MessagesTotal      *prometheus.CounterVec
	ProcessingDuration *prometheus.HistogramVec
	ErrorsTotal        *prometheus.CounterVec
	MessagesInProgress *prometheus.GaugeVec
}

func NewMetrics(namespace string) *Metrics {
	m := &Metrics{
		MessagesTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "messages_total",
			Help:      "Total number of messages processed",
		},
			[]string{"stage"},
		),

		ProcessingDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "processing_duration_seconds",
			Help:      "Processing duration of messages",
		}, []string{"stage"}),

		ErrorsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "errors_total",
			Help:      "Total number of errors",
		}, []string{"stage", "type"}),

		MessagesInProgress: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "messages_in_progress",
			Help:      "Number of messages in progress",
		}, []string{"stage"}),
	}

	prometheus.MustRegister(
		m.MessagesTotal,
		m.ProcessingDuration,
		m.ErrorsTotal,
		m.MessagesInProgress)

	return m
}
