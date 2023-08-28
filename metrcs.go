package batch_query

type MetricsWriter interface {
	FetchIn()
	FetchOut()
	FetchTimeout()
	FetchFail()
	BatchIn()
	BatchOut()
	BatchFail()
}
