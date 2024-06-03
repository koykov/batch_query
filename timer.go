package batch_query

import (
	"math"
	"sync/atomic"
	"time"
)

type timerSignal uint8

const (
	timerReach timerSignal = iota
	timerPause
	timerResume
	timerStop
)

const (
	timerActive = iota
	timerPaused
	timerStopped
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

func (t *timer) observe(query *BatchQuery) {
	t.t.Reset(query.config.CollectInterval)
	for {
		select {
		case signal, ok := <-t.c:
			if !ok {
				return
			}
			switch signal {
			case timerReach:
				atomic.StoreUint32(&t.s, timerPaused)
				t.t.Stop()
				query.flush(flushReasonInterval)
			case timerPause:
				atomic.StoreUint32(&t.s, timerPaused)
				t.t.Stop()
				break
			case timerResume:
				atomic.StoreUint32(&t.s, timerActive)
				t.t.Reset(query.config.CollectInterval)
				break
			case timerStop:
				atomic.StoreUint32(&t.s, timerStopped)
				t.t.Stop()
				close(t.c)
				return
			}
		case <-t.t.C:
			t.reach()
		}
	}
}

// Send time reach signal.
func (t *timer) reach() {
	if atomic.LoadUint32(&t.s) != timerActive {
		return
	}
	t.c <- timerReach
}

// Send pause signal.
func (t *timer) pause() {
	if atomic.LoadUint32(&t.s) != timerActive {
		return
	}
	t.c <- timerPause
}

// Check timer is paused.
func (t *timer) paused() bool {
	return atomic.LoadUint32(&t.s) == timerPaused
}

// Send resume signal.
func (t *timer) resume() {
	if atomic.CompareAndSwapUint32(&t.s, timerPaused, timerActive) {
		t.c <- timerResume
	}
}

// Send stop signal.
func (t *timer) stop() {
	if atomic.LoadUint32(&t.s) == timerStopped {
		return
	}
	t.c <- timerStop
}
