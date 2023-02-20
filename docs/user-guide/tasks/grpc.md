---
template: main.html
---

# grpc

Generate requests for a gRPC service and and collect [latency and error-related metrics](#metrics).

## Usage example

In this experiment, the `grpc` task generates call requests for a gRPC service hosted at `hello.default:50051`, defined in the [protobuf](https://developers.google.com/protocol-buffers) file located at `grpc.protoURL`, with a gRPC method named `helloworld.Greeter.SayHello`. Metrics collected by this task are used by the `assess` task to validate SLOs.

Single method:
```bash
iter8 k launch \
--set "tasks={ready,grpc,assess}" \
--set ready.deploy=routeguide \
--set ready.service=routeguide \
--set ready.timeout=60s \
--set grpc.host="routeguide.default:50051" \
--set grpc.call="routeguide.RouteGuide.GetFeature" \
--set grpc.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/unary.json" \
--set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/route_guide/routeguide/route_guide.proto" \
--set assess.SLOs.upper.grpc/error-rate=0 \
--set assess.SLOs.upper.grpc/latency/mean=200 \
--set runner=job
```

Multiple methods:
```bash
iter8 k launch \
--set "tasks={ready,grpc,assess}" \
--set ready.deploy=routeguide \
--set ready.service=routeguide \
--set ready.timeout=60s \
--set grpc.host="routeguide.default:50051" \
--set grpc.endpoints.getFeature.call="routeguide.RouteGuide.GetFeature" \
--set grpc.endpoints.getFeature.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/unary.json" \
--set grpc.endpoints.listFeature.call="routeguide.RouteGuide.ListFeatures" \
--set grpc.endpoints.listFeature.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/server.json" \
--set grpc.endpoints.recourdRoute.call="routeguide.RouteGuide.RecordRoute" \
--set grpc.endpoints.recourdRoute.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/client.json" \
--set grpc.endpoints.routeChat.call="routeguide.RouteGuide.RouteChat" \
--set grpc.endpoints.routeChat.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/bidirectional.json" \
--set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/route_guide/routeguide/route_guide.proto" \
--set assess.SLOs.upper.grpc/getFeature/error-rate=0 \
--set assess.SLOs.upper.grpc/listFeature/error-rate=0 \
--set assess.SLOs.upper.grpc/recourdRoute/error-rate=0 \
--set assess.SLOs.upper.grpc/routeChat/error-rate=0 \
--set runner=job
```

## Parameters

Any field in the [`Config` struct](https://github.com/bojand/ghz/blob/master/runner/config.go) of the [`ghz` runner package](https://github.com/bojand/ghz/tree/master/runner) can be used as a parameter in this task. The JSON tags of the struct fields directly correspond to the names of the parameters of this task. In the [usage example](#usage-example), the parameters `host` and `call` correspond to the `Host` and `Call` fields respectively in the `Config` struct.

In addition, the following fields are defined by this task. 

| Name | Type | Description |
| ---- | ---- | ----------- |
| protoURL | string (URL) | URL where the [protobuf file](https://developers.google.com/protocol-buffers) that defines the gRPC service is located. |
| dataURL | string (URL) | URL where JSON data to be used in call requests is located. |
| binaryDataURL | string (URL) | URL where binary data to be used in call requests is located. |
| metadataURL | string (URL) | URL where the JSON metadata data to be used in call requests is located. |
| warmupNumRequests | int | Number of requests to be sent in a warmup task (results are ignored).  |
| warmupDuration | string | Duration of warmup task (results are ignored). Specified in the [Go duration string format](https://pkg.go.dev/maze.io/x/duration#ParseDuration) (example, 5s). If both warmupDuration and warmupNumRequests are specified, then warmupDuration is ignored. |
| endpoints | map[string]EndPoint | Used to specify multiple endpoints and their configuration. The `string` is the name of the endpoint and the `EndPoint` struct includes all the parameters described above as well as those from the `Config` struct. Load testing and metric collection will be conducted separately for each endpoint. |

## Precedence

With the `endpoints` parameter, it is possible to configure parameters on a task level as well as on a per endpoint level.

```bash
iter8 k launch \
--set "tasks={grpc,assess}" \
--set grpc.host="routeguide.default:50051" \
--set grpc.endpoints.getFeature.call="routeguide.RouteGuide.GetFeature" \
--set grpc.endpoints.getFeature.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/unary.json" \
--set grpc.endpoints.listFeature.call="routeguide.RouteGuide.ListFeatures" \
--set grpc.endpoints.listFeature.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/server.json" \
--set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/route_guide/routeguide/route_guide.proto" \
--set assess.SLOs.upper.grpc/getFeature/error-rate=0 \
--set assess.SLOs.upper.grpc/listFeature/error-rate=0 \
--set runner=job
```

In this example, each endpoint has a different `call` and `dataURL` but they have the same `host`.

If the same task and endpoint parameter is set, then the endpoint parameter will have precedence.

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