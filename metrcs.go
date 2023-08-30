package batch_query

type MetricsWriter interface {
	Fetch()
	OK()
	FetchTimeout()
	Fail()
	Batch()
	BatchOK()
	BatchFail()
}
