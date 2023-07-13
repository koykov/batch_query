package aerospike

import (
	"context"

	"github.com/aerospike/aerospike-client-go"
)

type Batcher struct {
	Namespace string
	SetName   string
	Bins      []string
	Policy    *aerospike.BatchPolicy
	Client    *aerospike.Client
}

func (b Batcher) Batch(dst []any, keys []any, ctx context.Context) ([]any, error) {
	// todo implement me
	return dst, nil
}

func (b Batcher) CheckKey(key, val any) bool {
	var (
		ask, asv *aerospike.Key
		ok       bool
	)
	if ask, ok = key.(*aerospike.Key); !ok || ask == nil {
		return false
	}
	if asv, ok = val.(*aerospike.Key); !ok || asv == nil {
		return false
	}
	return ask.Equals(asv)
}
