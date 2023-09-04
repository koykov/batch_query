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
	flagTimer = iota
	flagNoMetrics
)
const (
	ctxInt uint8 = iota
	ctxTO
)

// BatchQuery is an implementation of query that collects single request to resource (database, network, ...) to batches
// and thus reduces pressure to resource.
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

// New makes new query instance and initialize it according config params.
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

	// Make a copy of config instance to protect queue from changing params after start.
	q.config = q.config.Copy()
	c := q.config

	// Check config params.
	if c.BatchSize == 0 {
		c.BatchSize = defaultBatchSize
	}
	if c.CollectInterval <= 0 {
		c.CollectInterval = defaultCollectInterval
	}
	if c.TimeoutInterval <= 0 {
		c.TimeoutInterval = defaultTimeoutInterval
	}
	if c.TimeoutInterval < c.CollectInterval {
		q.err = ErrBadIntervals
		q.status = StatusFail
		return
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

	// Run internal workers.
	var ctx context.Context
	ctx, q.cancel = context.WithCancel(context.Background())
	for i := uint(0); i < c.Workers; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case p := <-q.c:
					q.mw().BufferOut()
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

// Fetch add single request to current batch using default timeout interval.
func (q *BatchQuery) Fetch(key any) (any, error) {
	return q.FetchTimeout(key, q.config.TimeoutInterval)
}

// FetchContext add single request to current batch with context.
func (q *BatchQuery) FetchContext(key any, ctx context.Context) (any, error) {
	q.once.Do(q.init)
	if status := q.getStatus(); status == StatusClose || status == StatusFail {
		return nil, ErrQueryClosed
	}

	return q.fetch(key, ctx, ctxInt)
}

// FetchTimeout add single request to current batch using given timeout interval.
func (q *BatchQuery) FetchTimeout(key any, timeout time.Duration) (any, error) {
	if timeout <= 0 {
		return nil, ErrTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	_ = cancel
	return q.fetch(key, ctx, ctxTO)
}

// FetchDeadline add single request to current batch using given deadline.
func (q *BatchQuery) FetchDeadline(key any, deadline time.Time) (any, error) {
	timeout := -time.Since(deadline)
	return q.FetchTimeout(key, timeout)
}

func (q *BatchQuery) fetch(key any, ctx context.Context, ctxt uint8) (any, error) {
	q.mw().Fetch()
	c := make(chan tuple, 1)
	now := q.now()
	q.fetch_(key, c)
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
	case <-ctx.Done():
		switch ctxt {
		case ctxTO:
			q.mw().Timeout()
			return nil, ErrTimeout
		case ctxInt:
			fallthrough
		default:
			q.mw().Interrupt()
			return nil, ErrInterrupt
		}
	}
}

func (q *BatchQuery) fetch_(key any, c chan tuple) {
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

// Close gracefully stops the query.
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

// ForceClose closes the query and immediately flush all batches.
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
	if atomic.LoadUint32(&q.flags[flagNoMetrics]) == 1 {
		return
	}
	t = time.Now()
	return
}

var _ = New
var _, _ = StatusNil, StatusThrottle
