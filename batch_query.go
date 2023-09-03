package batch_query

import (
	"context"
	"math"
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
const (
	flagTimer     = 0
	flagNoMetrics = 1
)

type BatchQuery struct {
	once   sync.Once
	config *Config
	status Status
	flags  [2]uint32

	mux    sync.Mutex
	buf    []pair
	c      chan []pair
	idx    uint64
	timer  *timer
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

	if c.BatchSize == 0 {
		c.BatchSize = defaultBatchSize
	}
	if c.CollectInterval <= 0 {
		c.CollectInterval = defaultCollectInterval
	}
	if c.TimeoutInterval <= 0 {
		c.TimeoutInterval = defaultTimeoutInterval
	}
	q.timer = newTimer()

	if c.Workers == 0 {
		q.err = ErrNoWorkers
		q.status = StatusFail
		return
	}
	if c.Buffer == 0 {
		c.Buffer = defaultBuffer
	}
	if c.Batcher == nil {
		q.err = ErrNoBatcher
		q.status = StatusFail
		return
	}

	if c.MetricsWriter == nil {
		c.MetricsWriter = DummyMetrics{}
		atomic.StoreUint32(&q.flags[flagNoMetrics], 1)
	}

	q.c = make(chan []pair, q.config.Buffer)
	q.idx = math.MaxUint64

	var ctx context.Context
	ctx, q.cancel = context.WithCancel(context.Background())
	for i := uint(0); i < c.Workers; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case p := <-q.c:
					idx := atomic.AddUint64(&q.idx, 1)
					// Prepare keys.
					keys := make([]any, 0, len(p))
					for i := 0; i < len(p); i++ {
						keys = append(keys, p[i].key)
					}
					if l := q.l(); l != nil {
						l.Printf("batch #%d of %d keys\n", idx, len(keys))
					}
					// Exec batch operation.
					dst := make([]any, 0, len(p))
					var err error
					now := q.now()
					dst, err = q.config.Batcher.Batch(dst, keys, ctx)
					if err != nil {
						if l := q.l(); l != nil {
							l.Printf("batch #%d failed due to error: %s\n", idx, err.Error())
						}
						q.mw().BatchFail()
						// Report about error encountered.
						for i := 0; i < len(p); i++ {
							p[i].c <- tuple{err: err}
						}
						continue
					}
					q.mw().BatchOK(q.now().Sub(now))
					var s, r int
					// Send values to corresponding channels.
					for i := 0; i < len(dst); i++ {
						for j := 0; j < len(p); j++ {
							if p[j].done {
								continue
							}
							if p[j].done = q.config.Batcher.MatchKey(p[j].key, dst[i]); p[j].done {
								p[j].c <- tuple{val: dst[i]}
								close(p[j].c)
								s++
								continue
							}
						}
					}
					// Check rest of keys.
					for i := 0; i < len(p); i++ {
						if !p[i].done {
							p[i].c <- tuple{err: ErrNotFound}
							close(p[i].c)
							r++
						}
					}
					if l := q.l(); l != nil {
						l.Printf("batch #%d finish with %d success jobs, %d jobs unresponded\n", idx, s, r)
					}
				case <-ctx.Done():
					return
				}
			}
		}(ctx)
	}

	q.setStatus(StatusActive)
}

func (q *BatchQuery) Find(key any) (any, error) {
	return q.FindTimeout(key, q.config.TimeoutInterval)
}

func (q *BatchQuery) FindTimeout(key any, timeout time.Duration) (any, error) {
	if timeout <= 0 {
		return nil, ErrTimeout
	}

	q.once.Do(q.init)
	if status := q.getStatus(); status == StatusClose || status == StatusFail {
		return nil, ErrQueryClosed
	}

	q.mw().Fetch()
	c := make(chan tuple, 1)
	now := q.now()
	q.find(key, c)
	select {
	case rec := <-c:
		switch {
		case rec.err != nil && rec.err == ErrNotFound:
			q.mw().NotFound()
		case rec.err != nil && rec.err != ErrNotFound:
			q.mw().Fail()
		default:
			q.mw().OK(q.now().Sub(now))
		}
		return rec.val, rec.err
	case <-time.After(timeout):
		q.mw().Timeout()
		return nil, ErrTimeout
	}
}

func (q *BatchQuery) FindDeadline(key any, deadline time.Time) (any, error) {
	timeout := -time.Since(deadline)
	return q.FindTimeout(key, timeout)
}

func (q *BatchQuery) find(key any, c chan tuple) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.buf = append(q.buf, pair{key: key, c: c})
	if uint64(len(q.buf)) == q.config.BatchSize {
		q.flushLF(flushReasonSize)
		return
	}
	if atomic.LoadUint32(&q.flags[flagTimer]) == 0 {
		atomic.StoreUint32(&q.flags[flagTimer], 1)
		go q.timer.wait(q)
	}
}

func (q *BatchQuery) Close() error {
	if q.getStatus() == StatusClose {
		return ErrQueryClosed
	}
	q.setStatus(StatusClose)
	q.mux.Lock()
	defer q.mux.Unlock()
	q.flushLF(flushReasonForce)
	close(q.c)
	q.cancel()
	if l := q.l(); l != nil {
		l.Printf("caught close signal\n")
	}
	return nil
}

func (q *BatchQuery) ForceClose() error {
	if q.getStatus() == StatusClose {
		return ErrQueryClosed
	}
	q.setStatus(StatusClose)
	q.mux.Lock()
	defer q.mux.Unlock()
	close(q.c)
	var c int
	for x := range q.c {
		for _, p := range x {
			p.c <- tuple{err: ErrInterrupt}
			c++
		}
	}
	q.cancel()
	if l := q.l(); l != nil {
		l.Printf("caught force close signal, %d jobs interrupted\n", c)
	}
	return nil
}

func (q *BatchQuery) Error() error {
	return q.err
}

func (q *BatchQuery) setStatus(status Status) {
	atomic.StoreUint32((*uint32)(&q.status), uint32(status))
}

func (q *BatchQuery) getStatus() Status {
	return Status(atomic.LoadUint32((*uint32)(&q.status)))
}

func (q *BatchQuery) mw() MetricsWriter {
	return q.config.MetricsWriter
}

func (q *BatchQuery) l() Logger {
	return q.config.Logger
}

func (q *BatchQuery) now() (t time.Time) {
	if atomic.LoadUint32(&q.flags[flagNoMetrics]) == 0 {
		return
	}
	t = time.Now()
	return
}

type pair struct {
	key  any
	c    chan tuple
	done bool
}

var _ = New
var _, _ = StatusNil, StatusThrottle
