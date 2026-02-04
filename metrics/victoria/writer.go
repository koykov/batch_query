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

// Writer is a VictoriaMetrics implementation of batch_query.MetricsWriter.
type Writer struct {
	name string
	prec time.Duration
}

func NewWriter(name string, options ...Option) *Writer {
	w := &Writer{
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

func (m Writer) Fetch() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Inc()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioIn).Inc()
}

func (m Writer) OK(dur time.Duration) {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioOK).Inc()
	vmchain.Histogram("batch_query_timing").WithLabel("query", m.name).WithLabel("entity", single).Update(float64(dur / m.prec))
}

func (m Writer) NotFound() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", io404).Inc()
}

func (m Writer) Timeout() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioTO).Inc()
}

func (m Writer) Interrupt() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioInt).Inc()
}

func (m Writer) Fail() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", single).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", single).WithLabel("type", ioFail).Inc()
}

func (m Writer) Batch() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", batch).Inc()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", batch).WithLabel("type", ioIn).Inc()
}

func (m Writer) BatchOK(dur time.Duration) {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", batch).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", batch).WithLabel("type", ioOK).Inc()
	vmchain.Histogram("batch_query_timing").WithLabel("query", m.name).WithLabel("entity", batch).Update(float64(dur / m.prec))
}

func (m Writer) BatchFail() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", batch).Dec()
	vmchain.Counter("batch_query_io").WithLabel("query", m.name).WithLabel("entity", batch).WithLabel("type", ioFail).Inc()
}

func (m Writer) BufferIn(reason string) {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", buffer).Inc()
	vmchain.Counter("batch_query_bufio").WithLabel("query", m.name).WithLabel("reason", reason).Inc()
}

func (m Writer) BufferOut() {
	vmchain.Gauge("batch_query_size", nil).WithLabel("query", m.name).WithLabel("entity", buffer).Dec()
}

var _ = NewWriter
