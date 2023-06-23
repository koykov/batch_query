package batch_query

type Logger interface {
	Printf(format string, v ...any)
	Print(v ...any)
	Println(v ...any)
}
