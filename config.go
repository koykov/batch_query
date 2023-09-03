package batch_query

import (
	"math"
	"time"
)

const (
	defaultBatchSize       = 64
	defaultCollectInterval = time.Second
	defaultTimeoutInterval = math.MaxInt64
	defaultBuffer          = 16
)

// Config describes query properties and behavior.
type Config struct {
	// Max capacity of one batch.
	// If query collects that amount of requests before reach CollectInterval then batch will process with reason
	// 'size reach'.
	// If this param omit defaultBatchSize (64) will use instead.
	BatchSize uint64
	// How long collect requests before process the batch.
	// Timer starts by first request incoming and stops after process the batch. After reach that interval betch will
	// process even contains only one request.
	// If this param omit defaultCollectInterval (1 second) will use instead.
	CollectInterval time.Duration
	// How long request may wait collecting and processing. Must be greater that CollectInterval.
	TimeoutInterval time.Duration
	// Internal workers count to process batches.
	Workers uint
	// Internal buffer size to collect batches.
	// If this param omit defaultBuffer (16) will use instead.
	Buffer uint64
	// Batch processor.
	// Mandatory param.
	Batcher Batcher

	// Metrics writer handler.
	MetricsWriter MetricsWriter

	// Logger handler.
	Logger Logger
}

func (c *Config) Copy() *Config {
	cpy := *c
	return &cpy
}
