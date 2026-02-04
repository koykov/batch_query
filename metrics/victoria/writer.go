package victoria

import (
	"time"

	"github.com/koykov/vmchain"
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

// writer is a VictoriaMetrics implementation of batch_query.MetricsWriter.
type writer struct {
	name string
	prec time.Duration
}

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

func (m writer) Fetch() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Inc()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioIn).Inc()
}

func (m writer) OK(dur time.Duration) {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioOK).Inc()
	vmchain.Histogram("batch_query_timing").WithLabel("query", m.name).WithLabel("entity", single).Update(float64(dur / m.prec))
}

func (m writer) NotFound() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", io404).Inc()
}

func (m writer) Timeout() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioTO).Inc()
}

func (m writer) Interrupt() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioInt).Inc()
}

func (m writer) Fail() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioFail).Inc()
}

func (m writer) Batch() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", batch).Inc()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", batch).WithLabel("type", ioIn).Inc()
}

func (m writer) BatchOK(dur time.Duration) {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", batch).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", batch).WithLabel("type", ioOK).Inc()
	vmchain.Histogram("batch_query_timing").WithLabel("query", m.name).WithLabel("entity", batch).Update(float64(dur / m.prec))
}

func (m writer) BatchFail() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", batch).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", batch).WithLabel("type", ioFail).Inc()
}

func (m writer) BufferIn(reason string) {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", buffer).Inc()
	vmchain.Counter("batch_query_bufio").WithLabel("query", m.name).WithLabel("reason", reason).Inc()
}

func (m writer) BufferOut() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", buffer).Dec()
}

var _ = NewWriter
