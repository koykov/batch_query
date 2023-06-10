package batch_query

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Status uint32

const (
	StatusNil Status = iota
	StatusFail
	StatusActive
	StatusThrottle
	StatusClose
)

type BatchQuery struct {
	once   sync.Once
	config *Config
	status Status

	chunkSize   uint64
	collectTime time.Duration

	mux    sync.Mutex
	buf    []pair
	c      chan []pair
	timer  *time.Timer
	cancel context.CancelFunc

	err error
}

func New(conf *Config) (*BatchQuery, error) {
	if conf == nil {
		return nil, ErrNoConfig
	}
	q := BatchQuery{config: conf.Copy()}
	q.once.Do(q.init)
	return &q, q.err
}

func (q *BatchQuery) init() {
	if q.config == nil {
		q.err = ErrNoConfig
		q.status = StatusFail
		return
	}

	q.config = q.config.Copy()
	c := q.config

	if c.ChunkSize == 0 {
		c.ChunkSize = defaultChunkSize
	}
	if c.CollectInterval <= 0 {
		c.CollectInterval = defaultCollectInterval
	}
	if c.Workers == 0 {
		q.err = ErrNoWorkers
		q.status = StatusFail
		return
	}
	if c.Batcher == nil {
		q.err = ErrNoBatcher
		q.status = StatusFail
		return
	}

	var ctx context.Context
	ctx, q.cancel = context.WithCancel(context.Background())
	for i := uint(0); i < c.Workers; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case p := <-q.c:
					// Prepare keys.
					keys := make([]any, 0, len(p))
					for i := 0; i < len(p); i++ {
						keys = append(keys, p[i].key)
					}
					// Exec batch operation.
					dst := make([]any, 0, len(p))
					var err error
					dst, err = q.config.Batcher.Batch(dst, keys, ctx)
					if err != nil {
						// Report about error encountered.
						for i := 0; i < len(p); i++ {
							p[i].c <- tuple{err: err}
						}
						continue
					}
					// todo send response values to corresponding channels
				case <-ctx.Done():
					return
				}
			}
		}(ctx)
	}

	q.status = StatusActive
}

func (q *BatchQuery) Find(key any) (any, error) {
	q.once.Do(q.init)
	if status := q.getStatus(); status == StatusClose || status == StatusFail {
		return nil, ErrQueryClosed
	}

	c := make(chan tuple, 1)
	q.find(key, c)
	rec := <-c
	close(c)
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
	if q.getStatus() == StatusClose {
		return ErrQueryClosed
	}
	return nil
}

func (q *BatchQuery) setStatus(status Status) {
	atomic.StoreUint32((*uint32)(&q.status), uint32(status))
}

func (q *BatchQuery) getStatus() Status {
	return Status(atomic.LoadUint32((*uint32)(&q.status)))
}

type pair struct {
	key any
	c   chan tuple
}

var _ = New
var _, _ = StatusNil, StatusThrottle
