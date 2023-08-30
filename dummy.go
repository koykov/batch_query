package batch_query

type DummyMetrics struct{}

func (DummyMetrics) Fetch()        {}
func (DummyMetrics) OK()           {}
func (DummyMetrics) FetchTimeout() {}
func (DummyMetrics) Fail()         {}
func (DummyMetrics) Batch()        {}
func (DummyMetrics) BatchOK()      {}
func (DummyMetrics) BatchFail()    {}
