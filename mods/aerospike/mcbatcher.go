package aerospike

import (
	"context"
	"sync/atomic"

	as "github.com/aerospike/aerospike-client-go"
)

// MCBatcher implements multi-client Aerospike batcher.
// Experimental feature.
type MCBatcher struct {
	Namespace string
	SetName   string
	Bins      []string
	Policy    *as.BatchPolicy
	Clients   []*as.Client

	c uint64
}

func (b MCBatcher) Batch(dst []any, keys []any, ctx context.Context) ([]any, error) {
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
	if len(b.Clients) == 0 {
		return dst, ErrNoClients
	}
	cln := b.Clients[atomic.AddUint64(&b.c, 1)%uint64(len(b.Clients))]
	if cln == nil {
		return dst, ErrNoClient
	}
	var err error
	dst, err = fetch(cln, b.Policy, b.Namespace, b.SetName, b.Bins, dst, keys, ctx)
	return dst, err
}

func (b MCBatcher) MatchKey(key, val any) bool {
	return matchKey(key, val, b.Namespace, b.SetName)
}
