package batch_query

import "time"

const (
	defaultChunkSize       = 64
	defaultCollectInterval = time.Second
)

type Config struct {
	ChunkSize       uint64
	CollectInterval time.Duration
	Workers         uint
}

func (c *Config) Copy() *Config {
	cpy := *c
	return &cpy
}
