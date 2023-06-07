package batch_query

import "time"

type Config struct {
	ChunkSize       uint64
	CollectInterval time.Duration
	Workers         uint
}

func (c *Config) Copy() *Config {
	cpy := *c
	return &cpy
}
