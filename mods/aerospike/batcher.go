package aerospike

import (
	"context"

	as "github.com/aerospike/aerospike-client-go"
)

// Batcher implements Aerospike batcher.
type Batcher struct {
	Namespace string
	SetName   string
	Bins      []string
	Policy    *as.BatchPolicy
	Client    *as.Client
}

func (b Batcher) Batch(dst []any, keys []any, ctx context.Context) ([]any, error) {
	if len(b.Namespace) == 0 {
		return dst, ErrNoNS
	}
	if len(b.SetName) == 0 {
		return dst, ErrNoSet
	}
	if len(b.Bins) == 0 {
		return dst, ErrNoBins
	}
	if b.Policy == nil {
		return dst, ErrNoPolicy
	}
	if b.Client == nil {
		return dst, ErrNoClient
	}
	var err error
	dst, err = fetch(b.Client, b.Policy, b.Namespace, b.SetName, b.Bins, dst, keys, ctx)
	return dst, err
}

func (b Batcher) MatchKey(key, val any) bool {
	return matchKey(key, val, b.Namespace, b.SetName)
}
