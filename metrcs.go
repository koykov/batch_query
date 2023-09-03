package batch_query

import "time"

// MetricsWriter is an interface of query metrics handler.
// See example of implementations https://github.com/koykov/metrics_writers/tree/master/batch_query.
type MetricsWriter interface {
	// Fetch registers income single request.
	Fetch()
	// OK register successful processing of single request.
	OK(duration time.Duration)
	// NotFound registers single request failed due to empty response.
	NotFound()
	// Timeout registers single request filed due to timeout error.
	Timeout()
	// Fail registers any other error encountered.
	Fail()
	// Batch register processing start of batch.
	Batch()
	// BatchOK register successful processing of batch.
	BatchOK(duration time.Duration)
	// BatchFail registers failed batch processing.
	BatchFail()
}
