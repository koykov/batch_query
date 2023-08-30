package batch_query

import (
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
	t := timer{c: make(chan timerSignal, 1)}
	return &t
}

// Background waiter method.
func (t *timer) wait(query *BatchQuery) {
	t.t = time.AfterFunc(query.config.CollectInterval, func() {
		t.reach()
		query.SetBit(flagTimer, false)
	})
	for {
		signal, ok := <-t.c
		if !ok {
			return
		}
		switch signal {
		case timerReach:
			query.flush(flushReasonInterval)
		case timerReset:
			t.t.Stop()
			query.SetBit(flagTimer, false)
			break
		case timerStop:
			t.t.Stop()
			atomic.StoreUint32(&t.s, 1)
			close(t.c)
			return
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
