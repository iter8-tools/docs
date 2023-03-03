---
template: main.html
---

# grpc

Generate requests for a gRPC service and and collect [latency and error-related metrics](#metrics).

## Usage example

In this experiment, the `grpc` task generates call requests for a gRPC service hosted at `hello.default:50051`, defined in the [protobuf](https://developers.google.com/protocol-buffers) file located at `grpc.protoURL`, with a gRPC method named `helloworld.Greeter.SayHello`. Metrics collected by this task are used by the `assess` task to validate SLOs.

```
iter8 k launch \
--set "tasks={grpc,assess}" \
--set grpc.host="hello.default:50051" \
--set grpc.call="helloworld.Greeter.SayHello" \
--set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto" \
--set assess.SLOs.upper.grpc/error-rate=0 \
--set assess.SLOs.upper.grpc/latency/p'97\.5'=800 \
--set runner=job
```

## Parameters

Any field in the [`Config` struct](https://github.com/bojand/ghz/blob/master/runner/config.go) of the [`ghz` runner package](https://github.com/bojand/ghz/tree/master/runner) can be used as a parameter in this task. The JSON tags of the struct fields directly correspond to the names of the parameters of this task. In the [usage example](#usage-example), the parameters `host` and `call` correspond to the `Call` and `Host` fields respectively in the `Config` struct.

In addition, the following fields are defined by this task. 

| Name | Type | Description |
| ---- | ---- | ----------- |
| protoURL  | string (URL) | URL where the [protobuf file](https://developers.google.com/protocol-buffers) that defines the gRPC service is located. |
| dataURL  | string (URL) | URL where JSON data to be used in call requests is located. |
| binaryDataURL  | string (URL) | URL where binary data to be used in call requests is located. |
| metadataURL  | string (URL) | URL where the JSON metadata data to be used in call requests is located. |
| warmupNumRequests | int | Number of requests to be sent in a warmup task (results are ignored).  |
| warmupDuration | string | Duration of warmup task (results are ignored). Specified in the [Go duration string format](https://pkg.go.dev/maze.io/x/duration#ParseDuration) (example, 5s). If both warmupDuration and warmupNumRequests are specified, then warmupDuration is ignored. |

## Metrics

This task creates a built-in [provider](../topics/metrics.md#fully-qualified-names) named `grpc`. The following metrics are collected by this task:

- `grpc/request-count`: total number of requests sent
- `grpc/error-count`: number of error responses
- `grpc/error-rate`: fraction of error responses

The following latency metrics are also supported by this task.

- `grpc/latency/mean`: mean latency
- `grpc/latency/stddev`: standard deviation of latency
- `grpc/latency/min`: min latency
- `grpc/latency/max`: max latency
- `grpc/latency/pX`: X^th^ percentile latency, for any X in the range 0.0 to 100.0

All latency metrics have `msec` units.