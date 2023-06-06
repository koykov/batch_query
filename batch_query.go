package batch_query

import (
	"sync"
	"time"
)

const (
	defaultChunkSize   = 64
	defaultCollectTime = time.Second
)

type BatchQuery struct {
	chunkSize   uint64
	collectTime time.Duration

	mux   sync.Mutex
	buf   []pair
	timer *time.Timer
}

func NewBatchQuery(chunkSize uint64, collectTime time.Duration) (*BatchQuery, error) {
	if chunkSize == 0 {
		return nil, ErrUselessChunkSize
	}
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
	c := make(chan tuple, 1)
	q.find(key, c)
	rec, ok := <-c
	if !ok {
		return nil, ErrNotFound
	}
	return rec.val, rec.err
}

func (q *BatchQuery) find(key any, c chan tuple) {

}

func (q *BatchQuery) Close() error {
	return nil
}

type pair struct {
	key any
	c   chan tuple
}
