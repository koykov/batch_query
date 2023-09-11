package sql

import (
	"context"
	"database/sql"
)

type Batcher struct {
	DB             *sql.DB
	Query          string
	QueryFormatter QueryFormatter
}

type RecordBuilder interface {
	Build(args []any) (any, error)
}

func (b Batcher) Batch(dst []any, keys []any, ctx context.Context) ([]any, error) {
	query, err := b.QueryFormatter.Format(b.Query, keys)
	if err != nil {
		return dst, err
	}
	rows, err := b.DB.QueryContext(ctx, query, keys...)
	if err != nil {
		return dst, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		// todo: declare and implement Scanner (RecordBuilder?)
		// if err = rows.Scan(); err != nil {
		// 	return dst, err
		// }
	}

	return dst, nil
}

func (b Batcher) MatchKey(key, val any) bool {
	// todo: implement me
	_, _ = key, val
	return false
}
