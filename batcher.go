package batch_query

import "context"

type Batcher interface {
	Batch(dst []any, keys []any) ([]any, error)
	BatchContext(dst []any, keys []any, ctx context.Context) ([]any, error)
}
