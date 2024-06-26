package batch_query

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	single = "single"
	batch  = "batch"
	buffer = "buffer"

	ioIn   = "in"
	ioOK   = "success"
	ioTO   = "timeout"
	ioInt  = "interrupt"
	io404  = "not_found"
	ioFail = "fail"
)

// PrometheusMetrics is a Prometheus implementation of batch_query.MetricsWriter.
type PrometheusMetrics struct {
	name string
	prec time.Duration
}

var (
	promSize   *prometheus.GaugeVec
	promIO     *prometheus.CounterVec
	promBufIO  *prometheus.CounterVec
	promTiming *prometheus.HistogramVec

	_ = NewPrometheusMetrics
)

func NewPrometheusMetrics(name string) *PrometheusMetrics {
	return NewPrometheusMetricsWP(name, time.Nanosecond)
}

func NewPrometheusMetricsWP(name string, precision time.Duration) *PrometheusMetrics {
	if precision == 0 {
		precision = time.Nanosecond
	}
	m := &PrometheusMetrics{
		name: name,
		prec: precision,
	}
	return m
}

func init() {
	promSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "batch_query_size",
		Help: "Indicates entities distribution by types.",
	}, []string{"query", "entity"})
	promIO = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "batch_query_io",
		Help: "How many entities processed.",
	}, []string{"query", "entity", "type"})
	promBufIO = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "batch_query_bufio",
		Help: "Buffer operations.",
	}, []string{"query", "reason"})

	buckets := append(prometheus.DefBuckets, []float64{15, 20, 30, 40, 50, 100, 150, 200, 250, 500, 1000, 1500, 2000, 3000, 5000}...)
	promTiming = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "batch_query_timing",
		Help:    "How many worker waits due to delayed execution.",
		Buckets: buckets,
	}, []string{"query", "entity"})

	prometheus.MustRegister(promSize, promIO, promBufIO, promTiming)
}

func (m PrometheusMetrics) Fetch() {
	promSize.WithLabelValues(m.name, single).Inc()
	promIO.WithLabelValues(m.name, single, ioIn).Inc()
}

func (m PrometheusMetrics) OK(dur time.Duration) {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, ioOK).Inc()
	promTiming.WithLabelValues(m.name, single).Observe(float64(dur / m.prec))
}

func (m PrometheusMetrics) NotFound() {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, io404).Inc()
}

func (m PrometheusMetrics) Timeout() {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, ioTO).Inc()
}

func (m PrometheusMetrics) Interrupt() {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, ioInt).Inc()
}

func (m PrometheusMetrics) Fail() {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, ioFail).Inc()
}

func (m PrometheusMetrics) Batch() {
	promSize.WithLabelValues(m.name, batch).Inc()
	promIO.WithLabelValues(m.name, batch, ioIn).Inc()
}

func (m PrometheusMetrics) BatchOK(dur time.Duration) {
	promSize.WithLabelValues(m.name, batch).Dec()
	promIO.WithLabelValues(m.name, batch, ioOK).Inc()
	promTiming.WithLabelValues(m.name, batch).Observe(float64(dur / m.prec))
}

func (m PrometheusMetrics) BatchFail() {
	promSize.WithLabelValues(m.name, batch).Dec()
	promIO.WithLabelValues(m.name, batch, ioFail).Inc()
}

func (m PrometheusMetrics) BufferIn(reason string) {
	promSize.WithLabelValues(m.name, buffer).Inc()
	promBufIO.WithLabelValues(m.name, reason).Inc()
}

func (m PrometheusMetrics) BufferOut() {
	promSize.WithLabelValues(m.name, buffer).Dec()
}
