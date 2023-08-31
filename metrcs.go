package batch_query

type MetricsWriter interface {
	Fetch()
	OK()
	NotFound()
	Timeout()
	Fail()
	Batch()
	BatchOK()
	BatchFail()
}
