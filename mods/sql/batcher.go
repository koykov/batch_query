package sql

import (
	"context"
	"database/sql"
)

type Batcher struct {
	DB             *sql.DB
	Query          string
	QueryFormatter QueryFormatter
	RecordScanner  RecordScanner
	RecordMatcher  RecordMatcher
}

func (b Batcher) Batch(dst []any, keys []any, ctx context.Context) ([]any, error) {
	if b.DB == nil {
		return dst, ErrNoDB
	}
	if len(b.Query) == 0 {
		return dst, ErrNoQuery
	}
	if b.QueryFormatter == nil {
		return dst, ErrNoQueryFmt
	}
	if b.RecordScanner == nil {
		return dst, ErrNoRecScnr
	}
	if b.RecordMatcher == nil {
		return dst, ErrNoRecMtch
	}

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
		rec, err := b.RecordScanner.Scan(rows)
		if err != nil {
			return dst, err
		}
		dst = append(dst, rec)
	}

	return dst, nil
}

func (b Batcher) MatchKey(key, val any) bool {
	if b.RecordMatcher == nil {
		return false
	}
	return b.RecordMatcher.Match(key, val)
}
