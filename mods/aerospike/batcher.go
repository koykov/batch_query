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

func (b Batcher) Batch(dst []any, keys []any, _ context.Context) ([]any, error) {
	if len(b.Namespace) == 0 {
		return dst, ErrNoNS
	}
	if len(b.SetName) == 0 {
		return dst, ErrNoSet
	}
	if b.Policy == nil {
		return dst, ErrNoPolicy
	}
	if b.Client == nil {
		return dst, ErrNoClient
	}
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
	var ask, asv *aerospike.Key
	switch key.(type) {
	case *aerospike.Key:
		ask = key.(*aerospike.Key)
	default:
		ask, _ = aerospike.NewKey(b.Namespace, b.SetName, key)
	}
	switch val.(type) {
	case *aerospike.Key:
		asv = val.(*aerospike.Key)
	case *aerospike.Record:
		if raw := val.(*aerospike.Record); raw != nil {
			asv = val.(*aerospike.Record).Key
		}
	default:
		return false
	}
	if ask != nil && asv == nil {
		return false
	}
	if ask == nil && asv != nil {
		return false
	}
	return ask.Equals(asv)
}
