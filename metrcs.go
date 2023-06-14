package batch_query

type MetricsWriter interface {
	FindIn()
	FindOut()
	FindFail()
	BatchIn()
	BatchOut()
	BatchFail()
}
