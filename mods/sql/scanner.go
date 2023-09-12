package sql

import "database/sql"

type RecordScanner interface {
	Scan(rows *sql.Rows) (any, error)
}
