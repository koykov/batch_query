package batch_query

type DummyMetrics struct{}

func (DummyMetrics) FindIn()    {}
func (DummyMetrics) FindOut()   {}
func (DummyMetrics) FindFail()  {}
func (DummyMetrics) BatchIn()   {}
func (DummyMetrics) BatchOut()  {}
func (DummyMetrics) BatchFail() {}
