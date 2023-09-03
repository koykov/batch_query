# Batch Query

Batch query collects single requests to batches and process that allows to reduce pressure to some resources (network, 
database, ...).

## Usage

Example of usage allows in [demo application](https://github.com/koykov/demo/tree/master/batch_query).

## Metrics

Check [Prometheus implementation](https://github.com/koykov/metrics_writers/tree/master/batch_query) of metrics writer.

## Config params

todo ...

## Modules

Currently, supports only [Aerospike](mods/aerospike) module.
