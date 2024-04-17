package redis

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/koykov/byteconv"
)

type Batcher struct {
	Client *redis.Client
}

func (b Batcher) Batch(dst []any, keys []any, _ context.Context) ([]any, error) {
	if b.Client == nil {
		return dst, ErrNoClient
	}

	skeys := make([]string, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		switch x := keys[i].(type) {
		case string:
			skeys = append(skeys, x)
		case []byte:
			skeys = append(skeys, byteconv.B2S(x))
		}
	}
	vals, err := b.Client.MGet(skeys...).Result()
	if err != nil {
		return dst, err
	}
	for i := 0; i < len(vals); i++ {
		dst = append(dst, Tuple{
			Key:   skeys[i],
			Value: vals[i],
		})
	}
	return dst, nil
}

func (b Batcher) MatchKey(key, val any) bool {
	var skey string
	switch x := key.(type) {
	case string:
		skey = x
	case []byte:
		skey = byteconv.B2S(x)
	}

	switch x := val.(type) {
	case Tuple:
		return skey == x.Key
	case *Tuple:
		return skey == x.Key
	}
	return false
}
