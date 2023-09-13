package sql

import (
	"strconv"
	"testing"
)

var stages = []struct {
	query  string
	keys   []any
	pt     PlaceholderType
	expect string
}{
	{
		"select * from users where uid in (::args::)",
		[]any{1, 2, 3, 4, 5, 6, 7},
		PlaceholderMySQL,
		"select * from users where uid in (?,?,?,?,?,?,?)",
	},
	{
		"select * from users where uid in (::args::)",
		[]any{1, 2, 3, 4, 5, 6, 7},
		PlaceholderPgSQL,
		"select * from users where uid in ($1,$2,$3,$4,$5,$6,$7)",
	},
}

func TestMacrosQueryFormatter(t *testing.T) {
	for i, stage := range stages {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			qf := MacrosQueryFormatter{PlaceholderType: stage.pt}
			query, _ := qf.Format(stage.query, stage.keys)
			if query != stage.expect {
				t.FailNow()
			}
		})
	}
}

func BenchmarkMacrosQueryFormatter(b *testing.B) {
	for i, stage := range stages {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.ReportAllocs()
			for j := 0; j < b.N; j++ {
				qf := MacrosQueryFormatter{PlaceholderType: stage.pt}
				query, _ := qf.Format(stage.query, stage.keys)
				if query != stage.expect {
					b.FailNow()
				}
			}
		})
	}
}
