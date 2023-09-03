package batch_query

import "context"

// Batcher describes object to process batches.
type Batcher interface {
	// Batch processes collected batch.
	Batch(dst []any, keys []any, ctx context.Context) ([]any, error)
	// MatchKey checks if key corresponds to val.
	MatchKey(key, val any) bool
}
