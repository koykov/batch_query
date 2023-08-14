package aerospike

import (
	"context"

	as "github.com/aerospike/aerospike-client-go"
)

type Batcher struct {
	Namespace string
	SetName   string
	Bins      []string
	Policy    *as.BatchPolicy
	Client    *as.Client
}

func (b Batcher) Batch(dst []any, keys []any, _ context.Context) ([]any, error) {
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
	askeys := make([]*as.Key, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		var (
			ask *as.Key
			err error
		)
		switch keys[i].(type) {
		case *as.Key:
			ask = keys[i].(*as.Key)
		default:
			if ask, err = as.NewKey(b.Namespace, b.SetName, keys[i]); err != nil {
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
	var ask, asv *as.Key
	switch key.(type) {
	case *as.Key:
		ask = key.(*as.Key)
	default:
		ask, _ = as.NewKey(b.Namespace, b.SetName, key)
	}
	switch val.(type) {
	case *as.Key:
		asv = val.(*as.Key)
	case *as.Record:
		if raw := val.(*as.Record); raw != nil {
			asv = raw.Key
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
