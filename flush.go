package batch_query

type flushReason uint8

const (
	flushReasonSize flushReason = iota
	flushReasonInterval
	flushReasonForce
)

func (r flushReason) String() string {
	switch r {
	case flushReasonSize:
		return "size"
	case flushReasonInterval:
		return "interval"
	case flushReasonForce:
		return "force"
	default:
		return "unknown"
	}
}

func (q *BatchQuery) flush(reason flushReason) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.flushLF(reason)
}

func (q *BatchQuery) flushLF(reason flushReason) {
	if reason != flushReasonInterval {
		q.timer.stop()
	}
	cpy := append([]pair(nil), q.buf...)
	q.buf = q.buf[:0]
	q.mw().BatchIn()
	if l := q.l(); l != nil {
		l.Printf("flush by reason '%s'\n", reason.String())
	}
	q.c <- cpy
}
