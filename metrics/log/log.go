package batch_query

import "log"

// LogMetrics is Log implementation of batch_query.MetricsWriter.
//
// Don't use in production. Only for debug purposes.
type LogMetrics struct {
	name string
}

var _ = NewLogMetrics

func NewLogMetrics(name string) *LogMetrics {
	m := &LogMetrics{name}
	return m
}

func (m LogMetrics) FindIn() {
	log.Printf("batch_query %s: new item income\n", m.name)
}

func (m LogMetrics) FindOut() {
	log.Printf("batch_query %s: item outcome\n", m.name)
}

func (m LogMetrics) FindFail() {
	log.Printf("batch_query %s: item processing fail\n", m.name)
}

func (m LogMetrics) BatchIn() {
	log.Printf("batch_query %s: new batch completed\n", m.name)
}

func (m LogMetrics) BatchOut() {
	log.Printf("batch_query %s: batch processed\n", m.name)
}

func (m LogMetrics) BatchFail() {
	log.Printf("batch_query %s: batch failed\n", m.name)
}
