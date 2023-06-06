package batch_query

import "time"

const (
	defaultChunkSize   = 64
	defaultCollectTime = time.Second
)

type BatchQuery struct {
	chunkSize   uint64
	collectTime time.Duration
}

func NewBatchQuery(chunkSize uint64, collectTime time.Duration) (*BatchQuery, error) {
	if collectTime < 0 {
		return nil, ErrNegativeDuration
	}
	q := BatchQuery{
		chunkSize:   chunkSize,
		collectTime: collectTime,
	}
	return &q, nil
}

func (q *BatchQuery) Find(key any) (any, error) {
	_ = key
	return nil, nil
}

func (q *BatchQuery) Close() error {
	return nil
}
