package aerospike

import (
	"context"

	as "github.com/aerospike/aerospike-client-go/v7"
)

func fetch(cln *as.Client, pol *as.BatchPolicy, ns, set string, bins []string, dst []any, keys []any, _ context.Context) ([]any, error) {
	askeys := make([]*as.Key, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		var (
			ask *as.Key
			err error
		)
		switch x := keys[i].(type) {
		case *as.Key:
			ask = x
		default:
			if ask, err = as.NewKey(ns, set, keys[i]); err != nil {
				return dst, nil
			}
		}
		askeys = append(askeys, ask)
	}
	records, err := cln.BatchGet(pol, askeys, bins...)
	if err != nil {
		return dst, err
	}
	for i := 0; i < len(records); i++ {
		dst = append(dst, records[i])
	}
	return dst, nil
}

func matchKey(key, val any, ns, set string) bool {
	var ask, asv *as.Key
	switch x := key.(type) {
	case *as.Key:
		ask = x
	default:
		ask, _ = as.NewKey(ns, set, key)
	}
	switch x := val.(type) {
	case *as.Key:
		asv = x
	case *as.Record:
		if x != nil {
			asv = x.Key
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
