# Batch Query

Batch query collects single requests to batches and process that allows to reduce pressure to some resources (network, 
database, ...).

## Usage

Example of usage allows in [demo application](https://github.com/koykov/demo/tree/master/batch_query).

## Metrics

Check [Prometheus implementation](https://github.com/koykov/batch_query/tree/master/metrics/prometheus) of metrics writer.

## Config params

See [source code](config.go)

## Modules

Currently, supports modules:
* [Aerospike](mods/aerospike)
* [Redis](mods/redis)
* [SQL](mods/sql)
