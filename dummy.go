package batch_query

import "time"

type DummyMetrics struct{}

func (DummyMetrics) Fetch()                  {}
func (DummyMetrics) OK(_ time.Duration)      {}
func (DummyMetrics) Timeout()                {}
func (DummyMetrics) NotFound()               {}
func (DummyMetrics) Fail()                   {}
func (DummyMetrics) Batch()                  {}
func (DummyMetrics) BatchOK(_ time.Duration) {}
func (DummyMetrics) BatchFail()              {}
