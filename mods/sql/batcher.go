package sql

import (
	"context"
	"database/sql"
)

type Batcher struct {
	DB    *sql.DB
	Query string
}

func (b Batcher) Batch(dst []any, keys []any, ctx context.Context) ([]any, error) {
	query := b.Query // todo: declare and implement QueryFormatter
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
