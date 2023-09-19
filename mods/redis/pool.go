package redis

import "sync"

var kpool = sync.Pool{New: func() any { return make([]string, 0, 16) }}
