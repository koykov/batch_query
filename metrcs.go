package batch_query

import "time"

type MetricsWriter interface {
	Fetch()
	OK(duration time.Duration)
	NotFound()
	Timeout()
	Fail()
	Batch()
	BatchOK(duration time.Duration)
	BatchFail()
}
