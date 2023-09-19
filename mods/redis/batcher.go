package redis

import (
	"context"

	"github.com/go-redis/redis"
)

type Batcher struct {
	Client *redis.Client
}

func (b Batcher) Batch(dst []any, keys []any, ctx context.Context) ([]any, error) {
	_, _ = keys, ctx
	// todo implement me
	return dst, nil
}

func (b Batcher) MatchKey(key, val any) bool {
	_, _ = key, val
	// todo implement me
	return false
}
