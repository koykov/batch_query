package batch_query

import "errors"

var (
	ErrNoConfig     = errors.New("no config provided")
	ErrNoWorkers    = errors.New("no workers available")
	ErrNoBatcher    = errors.New("no batcher provided")
	ErrBadIntervals = errors.New("bad intervals: timeout less that collect")
	ErrQueryClosed  = errors.New("query closed")
	ErrNotFound     = errors.New("record not found")
	ErrInterrupt    = errors.New("interrupt")
	ErrTimeout      = errors.New("timeout")
)
