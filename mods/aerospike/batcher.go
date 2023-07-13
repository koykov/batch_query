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
	askeys := make([]*aerospike.Key, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		var (
			ask *aerospike.Key
			err error
		)
		switch keys[i].(type) {
		case *aerospike.Key:
			ask = keys[i].(*aerospike.Key)
		default:
			if ask, err = aerospike.NewKey(b.Namespace, b.SetName, keys[i]); err != nil {
				return dst, nil
			}
		}
		askeys = append(askeys, ask)
	}
	records, err := b.Client.BatchGet(b.Policy, askeys, b.Bins...)
	if err != nil {
		return dst, err
	}
	for i := 0; i < len(records); i++ {
		dst = append(dst, records[i])
	}
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
