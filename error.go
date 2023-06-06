package batch_query

import "errors"

var (
	ErrUselessChunkSize = errors.New("useless chunk size")
	ErrNegativeDuration = errors.New("negative duration")
	ErrNotFound         = errors.New("record not found")
)
