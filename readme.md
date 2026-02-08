# Batch Query

This library helps to combine many small, similar queries to various resources (databases, caches, networks) into batches,
thereby reducing overhead on connections and data transport.

The library initially developed in response to permanent issues between a high-load service and an Aerospike server.
Aerospike's own metrics showed it was underloaded, while the application's metrics indicated huge query latencies to Aerospike.
The library's creation was also inspired by the
[singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight) package - the main idea of reducing the number of small queries,
as a logical continuation, led to the thought of eliminating them totally using batching.

The solution's concept is very simple: collect small queries until either an amount limit (maximum batch size) is reached
or the batch collection time expires (from the moment the first small query enters to the batch).

## Usage

The central component of the library is the `BatchQuery` structure, which does all the work. Before use, it must be configured,
and the [`Config`](config.go) structure serves this purpose. Let's examine its fields:
* `BatchSize` - how many small queries a batch can contain. Optional field, default value is `64`.
* `CollectInterval` - the maximum duration for collecting a batch. Starts counting from the moment the first query enters the batch. Default value is `1` second.
* `TimeoutInterval` - a limit on collection, sending the batch request, and post-processing. Must be greater than `CollectInterval`.
* `Batcher` - an abstraction for a specific storage, see description below. Mandatory parameter.
* `Buffer` - size of storage for collected batches, ready to be sent and processed.
* `Workers` - number of workers for sending/processing batches. They read from the buffer (see `Buffer`).
* `MetricsWriter` - abstraction for a specific TSDB solution.
* `Logger` - abstraction for an internal process logger. Useful for debugging, not recommended for production.

Thus, a usage example looks like this:
```go
package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/koykov/batch_query"
	promw "github.com/koykov/batch_query/metrics/prometheus"
	"github.com/koykov/batch_query/mods/aerospike"
)

func main() {
	policy := as.NewBatchPolicy()
	client, _ := as.NewClientWithPolicy(as.NewClientPolicy(), "localhost", 3000)

	// Prepare config for query.
	conf := batch_query.Config{
		BatchSize:       100,
		CollectInterval: 500 * time.Microsecond,
		TimeoutInterval: 5 * time.Millisecond,
		Workers:         10,
		Buffer:          4,
		// Declare Aerospike batcher with specific params.
		Batcher: aerospike.Batcher{
			Namespace: "my_ns",
			SetName:   "my_set",
			Bins:      []string{"bin01", "bin02", "binNN"},
			Policy:    policy,
			Client:    client,
		},
		// Declare writer to export metrics.
		MetricsWriter: promw.NewWriter("my_query", promw.WithPrecision(time.Millisecond)),
		// Declare logger for debugging purposes.
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	// Initialize the query.
	bq, _ := batch_query.New(&conf)

	// Start 10k goroutines to fetch small keys.
	for i := 0; i < 10000; i++ {
		go func() {
			for {
				resp, _ := bq.Fetch(i + rand.Intn(i))
				_ = resp.(*as.Record)
			}
		}()
	}

	c := make(chan struct{})
	<-c
}
```
In this example, 10k goroutines read single keys, which will be combined into batches and processed in bulk by the query.
At the same time, each goroutine will receive a response specifically for the key it requested, or an error.

## Modules

Currently, the library supports three data storages via the [Batcher](batcher.go) abstraction:
* [Aerospike](mods/aerospike)
* [Redis](mods/redis)
* [SQL](mods/sql)

The interface itself is quite simple, and if necessary, it's fairly straightforward to write your own version for the required storage.

## Metrics

To evaluate the query's efficiency and/or tune configuration parameters, you can set a component for writing and exporting
metrics via the [MetricsWriter](metrics.go) abstraction. Currently, the following TSDBs are supported:
* [Prometheus](https://github.com/koykov/batch_query/tree/master/metrics/prometheus)
* [VictoriaMetrics](https://github.com/koykov/batch_query/tree/master/metrics/victoria)

Using them is very simple - you need to set a unique queue name and, optionally, the timestamp precision
(by default, one nanosecond, but it's more reasonable to set one millisecond, see the usage example).
