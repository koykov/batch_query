package batch_query

type MetricsWriter interface {
	Fetch()
	OK()
	Timeout()
	Fail()
	Batch()
	BatchOK()
	BatchFail()
}
