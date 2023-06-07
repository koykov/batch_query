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
	once   sync.Once
	config *Config

	chunkSize   uint64
	collectTime time.Duration

	mux   sync.Mutex
	buf   []pair
	timer *time.Timer
}

func New(conf *Config) (*BatchQuery, error) {
	q := BatchQuery{config: conf}
	return &q, nil
}

// DEPRECATED
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
	q.once.Do(func() {
		if q.chunkSize == 0 {
			q.chunkSize = defaultChunkSize
		}
		if q.collectTime == 0 {
			q.collectTime = defaultCollectTime
		}
	})
	c := make(chan tuple, 1)
	q.find(key, c)
	rec, ok := <-c
	if !ok {
		return nil, ErrNotFound
	}
	return rec.val, rec.err
}

func (q *BatchQuery) find(key any, c chan tuple) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.buf = append(q.buf, pair{key: key, c: c})
	if uint64(len(q.buf)) == q.chunkSize {
		cpy := append([]pair(nil), q.buf...)
		_ = cpy
		// ...
		q.buf = q.buf[:0]
		return
	}
}

func (q *BatchQuery) Close() error {
	return nil
}

type pair struct {
	key any
	c   chan tuple
}

var _, _ = New, NewBatchQuery
