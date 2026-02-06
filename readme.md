# Batch Query

Batch query collects single requests to batches and process that allows to reduce pressure to some resources (network, 
database, ...).

## Usage

Example of usage allows in [demo application](https://github.com/koykov/demo/tree/master/batch_query).

## Metrics

Currently, package contains two built-in implementations of [MetricsWriter](metrics.go):
* [Prometheus implementation](https://github.com/koykov/batch_query/tree/master/metrics/prometheus)
* [VictoriaMetrics](https://github.com/koykov/batch_query/tree/master/metrics/victoria)

Feel free to implement your own implementation.

## Config params

See [source code](config.go)

## Modules

Currently, supports modules:
* [Aerospike](mods/aerospike)
* [Redis](mods/redis)
* [SQL](mods/sql)
