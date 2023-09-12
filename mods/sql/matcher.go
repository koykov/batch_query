package sql

type RecordMatcher interface {
	Match(key, value any) bool
}
