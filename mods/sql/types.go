package sql

import "database/sql"

type PlaceholderType uint8

const (
	PlaceholderMySQL PlaceholderType = iota
	PlaceholderPgSQL
)

type QueryFormatter interface {
	Format(query string, args []any) (string, error)
}

type RecordScanner interface {
	Scan(rows *sql.Rows) (any, error)
}

type RecordMatcher interface {
	Match(key, value any) bool
}
