package batch_query

type DummyMetrics struct{}

func (DummyMetrics) FetchIn()      {}
func (DummyMetrics) FetchOut()     {}
func (DummyMetrics) FetchTimeout() {}
func (DummyMetrics) FetchFail()    {}
func (DummyMetrics) BatchIn()      {}
func (DummyMetrics) BatchOut()     {}
func (DummyMetrics) BatchFail()    {}
