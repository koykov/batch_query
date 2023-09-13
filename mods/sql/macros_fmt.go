package sql

import (
	"strconv"
	"strings"
)

const defaultMacros = "::args::"

type MacrosQueryFormatter struct {
	QueryMacros     string
	PlaceholderType PlaceholderType
}

func (qf MacrosQueryFormatter) Format(query string, args []any) (string, error) {
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
	return strings.Replace(query, defaultMacros, string(buf), 1), nil
}
