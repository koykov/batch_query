package batch_query

import "errors"

var (
	ErrNoConfig  = errors.New("no config provided")
	ErrNoWorkers = errors.New("no workers available")
	ErrNotFound  = errors.New("record not found")
)
