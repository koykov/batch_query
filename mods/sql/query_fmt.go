package sql

import (
	"strconv"
	"strings"
)

type QueryFormatter interface {
	Format(query string, args []any) (string, error)
}

type PlaceholderType uint8

const (
	PlaceholderMySQL PlaceholderType = iota
	PlaceholderPgSQL
)

const QuerySubstring = "::args::"

type SubstringQueryFormatter struct {
	QuerySubstring  string
	PlaceholderType PlaceholderType
}

func (qf SubstringQueryFormatter) Format(query string, args []any) (string, error) {
	buf := make([]byte, 0, len(args)*5)
	for i := 0; i < len(args); i++ {
		if i > 0 {
			buf = append(buf, ',')
		}
		switch qf.PlaceholderType {
		case PlaceholderMySQL:
			buf = append(buf, '?')
		case PlaceholderPgSQL:
			buf = append(buf, '$')
			buf = strconv.AppendInt(buf, int64(i+1), 10)
		default:
			return "", ErrUnknownPlaceholderType
		}
	}
	return strings.Replace(query, QuerySubstring, string(buf), 1), nil
}
