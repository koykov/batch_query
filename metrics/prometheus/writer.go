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

type Writer interface {
	Fetch()
	OK(duration time.Duration)
	NotFound()
	Timeout()
	Interrupt()
	Fail()
	Batch()
	BatchOK(duration time.Duration)
	BatchFail()
	BufferIn(reason string)
	BufferOut()
}

// writer is a Prometheus implementation of batch_query.MetricsWriter.
type writer struct {
	name string
	prec time.Duration
}

var (
	promSize   *prometheus.GaugeVec
	promFlush  *prometheus.GaugeVec
	promIO     *prometheus.CounterVec
	promBufIO  *prometheus.CounterVec
	promTiming *prometheus.HistogramVec
)

func NewWriter(name string, options ...Option) Writer {
	w := &writer{
		name: name,
		prec: time.Nanosecond,
	}
	for _, fn := range options {
		fn(w)
	}
	if w.prec <= 0 {
		w.prec = time.Nanosecond
	}
	return w
}

// NewPrometheusMetrics makes new instance of Writer.
// Deprecated: use NewWriter instead.
func NewPrometheusMetrics(name string) Writer {
	return NewWriter(name, WithPrecision(time.Nanosecond))
}

// NewPrometheusMetricsWP makes new instance of Writer with given precision.
// Deprecated: use NewWriter instead.
func NewPrometheusMetricsWP(name string, precision time.Duration) Writer {
	return NewWriter(name, WithPrecision(precision))
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
	}, []string{"query"})
	promFlush = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "batch_query_flush",
		Help: "Indicates flush events distribution by reason.",
	}, []string{"query", "reason"})

	buckets := append(prometheus.DefBuckets, []float64{15, 20, 30, 40, 50, 100, 150, 200, 250, 500, 1000, 1500, 2000, 3000, 5000}...)
	promTiming = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "batch_query_timing",
		Help:    "How many worker waits due to delayed execution.",
		Buckets: buckets,
	}, []string{"query", "entity"})

	prometheus.MustRegister(promSize, promFlush, promIO, promBufIO, promTiming)
}

func (m writer) Fetch() {
	promSize.WithLabelValues(m.name, single).Inc()
	promIO.WithLabelValues(m.name, single, ioIn).Inc()
}

func (m writer) OK(dur time.Duration) {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, ioOK).Inc()
	promTiming.WithLabelValues(m.name, single).Observe(float64(dur / m.prec))
}

func (m writer) NotFound() {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, io404).Inc()
}

func (m writer) Timeout() {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, ioTO).Inc()
}

func (m writer) Interrupt() {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, ioInt).Inc()
}

func (m writer) Fail() {
	promSize.WithLabelValues(m.name, single).Dec()
	promIO.WithLabelValues(m.name, single, ioFail).Inc()
}

func (m writer) Batch() {
	promSize.WithLabelValues(m.name, batch).Inc()
	promIO.WithLabelValues(m.name, batch, ioIn).Inc()
}

func (m writer) BatchOK(dur time.Duration) {
	promSize.WithLabelValues(m.name, batch).Dec()
	promIO.WithLabelValues(m.name, batch, ioOK).Inc()
	promTiming.WithLabelValues(m.name, batch).Observe(float64(dur / m.prec))
}

func (m writer) BatchFail() {
	promSize.WithLabelValues(m.name, batch).Dec()
	promIO.WithLabelValues(m.name, batch, ioFail).Inc()
}

func (m writer) BufferIn(reason string) {
	promSize.WithLabelValues(m.name, buffer).Inc()
	promFlush.WithLabelValues(m.name, reason)
	promBufIO.WithLabelValues(m.name).Inc()
}

func (m writer) BufferOut() {
	promSize.WithLabelValues(m.name, buffer).Dec()
	promBufIO.WithLabelValues(m.name).Add(-1)
}
