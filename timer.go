package batch_query

import (
	"math"
	"sync/atomic"
	"time"
)

type timerSignal uint8

const (
	timerReach timerSignal = iota
	timerReset
	timerStop
)

// Internal timer implementation.
type timer struct {
	t *time.Timer
	c chan timerSignal
	s uint32
}

func newTimer() *timer {
	t := timer{
		c: make(chan timerSignal, 1),
		t: time.NewTimer(math.MaxInt64),
	}
	return &t
}

// Background waiter method.
func (t *timer) wait(query *BatchQuery) {
	t.t.Reset(query.config.CollectInterval)
	for {
		select {
		case signal, ok := <-t.c:
			if !ok {
				return
			}
			switch signal {
			case timerReach:
				query.flush(flushReasonInterval)
			case timerReset:
				t.t.Stop()
				atomic.StoreUint32(&query.flags[flagTimer], 0)
				break
			case timerStop:
				t.t.Stop()
				atomic.StoreUint32(&t.s, 1)
				close(t.c)
				return
			}
		case <-t.t.C:
			t.reach()
			atomic.StoreUint32(&query.flags[flagTimer], 0)
		}
	}
}

// Send time reach signal.
func (t *timer) reach() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerReach
}

// Send reset signal.
func (t *timer) reset() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerReset
}

// Send stop signal.
func (t *timer) stop() {
	if atomic.LoadUint32(&t.s) != 0 {
		return
	}
	t.c <- timerStop
}
