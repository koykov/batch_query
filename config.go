package batch_query

import "time"

const (
	defaultChunkSize       = 64
	defaultCollectInterval = time.Second
	defaultBuffer          = 16
)

type Config struct {
	ChunkSize       uint64
	CollectInterval time.Duration
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
