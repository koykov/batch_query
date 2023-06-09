package batch_query

import (
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

	mux   sync.Mutex
	buf   []pair
	timer *time.Timer

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
	q.status = StatusActive
}

func (q *BatchQuery) Find(key any) (any, error) {
	q.once.Do(q.init)
	if status := q.getStatus(); status == StatusClose || status == StatusFail {
		return nil, ErrQueryClosed
	}

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
