package batch_query

type MetricsWriter interface {
	FindIn()
	FindOut()
	FindTimeout()
	FindFail()
	BatchIn()
	BatchOut()
	BatchFail()
}
