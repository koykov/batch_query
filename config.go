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

type Config struct {
	BatchSize       uint64
	CollectInterval time.Duration
	TimeoutInterval time.Duration
	Workers         uint
	Buffer          uint64
	Batcher         Batcher

	MetricsWriter MetricsWriter

	Logger Logger
}

func (c *Config) Copy() *Config {
	cpy := *c
	return &cpy
}
