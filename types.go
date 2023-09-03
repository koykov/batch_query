package batch_query

// pair represents internal request in batches.
// See tuple type.
type pair struct {
	key  any
	c    chan tuple
	done bool
}

// tuple represents internal response to single request.
// See pair type.
type tuple struct {
	val any
	err error
}
