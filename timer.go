package batch_query

import (
	"math"
	"sync/atomic"
	"time"
)

type timerSignal uint8

const (
	timerSignalReach timerSignal = iota
	timerSignalPause
	timerSignalResume
	timerSignalStop
)

const (
	timerStatusActive = iota
	timerStatusPaused
	timerStatusStopped
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
			case timerSignalReach:
				atomic.StoreUint32(&t.s, timerStatusPaused)
				t.t.Stop()
				query.flush(flushReasonInterval)
				break
			case timerSignalPause:
				atomic.StoreUint32(&t.s, timerStatusPaused)
				t.t.Stop()
				break
			case timerSignalResume:
				atomic.StoreUint32(&t.s, timerStatusActive)
				t.t.Reset(query.config.CollectInterval)
				break
			case timerSignalStop:
				atomic.StoreUint32(&t.s, timerStatusStopped)
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
	if atomic.LoadUint32(&t.s) != timerStatusActive {
		return
	}
	t.send(timerSignalReach)
}

// Send pause signal.
func (t *timer) pause() {
	if atomic.LoadUint32(&t.s) != timerStatusActive {
		return
	}
	t.send(timerSignalPause)
}

// Check timer is paused.
func (t *timer) paused() bool {
	return atomic.LoadUint32(&t.s) == timerStatusPaused
}

// Send resume signal.
func (t *timer) resume() {
	if atomic.CompareAndSwapUint32(&t.s, timerStatusPaused, timerStatusActive) {
		t.send(timerSignalResume)
	}
}

// Send stop signal.
func (t *timer) stop() {
	if atomic.LoadUint32(&t.s) == timerStatusStopped {
		return
	}
	t.send(timerSignalStop)
}

func (t *timer) send(sig timerSignal) bool {
	select {
	case t.c <- sig:
		return true
	default:
		return false
	}
}
