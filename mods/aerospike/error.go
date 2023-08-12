package aerospike

import "errors"

var (
	ErrNoNS     = errors.New("no namespace provided")
	ErrNoSet    = errors.New("no set name provided")
	ErrNoPolicy = errors.New("no batch policy provided")
	ErrNoClient = errors.New("no client provided")
)
