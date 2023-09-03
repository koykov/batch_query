package batch_query

import "time"

// DummyMetrics is a stub metrics writer handler that uses by default and does nothing.
// Need just to reduce checks in code.
type DummyMetrics struct{}

func (DummyMetrics) Fetch()                  {}
func (DummyMetrics) OK(_ time.Duration)      {}
func (DummyMetrics) Timeout()                {}
func (DummyMetrics) NotFound()               {}
func (DummyMetrics) Fail()                   {}
func (DummyMetrics) Batch()                  {}
func (DummyMetrics) BatchOK(_ time.Duration) {}
func (DummyMetrics) BatchFail()              {}
func (DummyMetrics) BufferIn(_ string)       {}
func (DummyMetrics) BufferOut()              {}
