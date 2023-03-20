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
--set "tasks={grpc,assess}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.call=routeguide.RouteGuide.GetFeature \
--set grpc.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.protoURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/route_guide/routeguide/route_guide.proto \
--set assess.SLOs.upper.grpc/error-rate=0 \
--set assess.SLOs.upper.grpc/latency/mean=200 \
--set runner=job
```

Multiple methods:
```bash
iter8 k launch \
--set "tasks={grpc,assess}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeature.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json \
--set grpc.endpoints.routeChat.call=routeguide.RouteGuide.RouteChat \
--set grpc.endpoints.routeChat.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/bidirectional.json \
--set grpc.protoURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/route_guide/routeguide/route_guide.proto \
--set assess.SLOs.upper.grpc-getFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-listFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-recordRoute/error-rate=0 \
--set assess.SLOs.upper.grpc-routeChat/error-rate=0 \
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

Some parameters have a default value, which can be overwritten. In addition, with the `endpoints` parameter, you can test multiple endpoints and configure parameters for each of those endpoint. In these cases, the priority order is the default value, the value set at the base level, and the value set at the endpoint value.

In the following example, all three endpoints will use the default `timeout` of `20s` (from `Config` struct).

```bash
iter8 k launch \
--set "tasks={grpc,assess}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeature.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json \
--set grpc.protoURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/route_guide/routeguide/route_guide.proto \
--set assess.SLOs.upper.grpc-getFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-listFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-recordRoute/error-rate=0 \
--set runner=job
```

In the following example, the `getFeature` and `listFeature` endpoints will use the default `timeout` of `20s` and the `recordRoute` endpoint will use a `timeout` of `30s`.

```bash
iter8 k launch \
--set "tasks={grpc,assess}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeature.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json \
--set grpc.endpoints.recordRoute.timeout=30s \
--set grpc.protoURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/route_guide/routeguide/route_guide.proto \
--set assess.SLOs.upper.grpc-getFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-listFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-recordRoute/error-rate=0 \
--set runner=job
```

In the following example, all three endpoints will use a `qps` of `40s`.

```bash
iter8 k launch \
--set "tasks={grpc,assess}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.timeout=40s \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeature.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json \
--set grpc.protoURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/route_guide/routeguide/route_guide.proto \
--set assess.SLOs.upper.grpc-getFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-listFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-recordRoute/error-rate=0 \
--set runner=job
```

In the following example, the `getFeature` and `listFeature` endpoints will use a `timeout` of `40s` and the `listFeature` endpoint will use a `timeout` of `30s`.

```bash
iter8 k launch \
--set "tasks={grpc,assess}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.timeout=40s \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeature.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.listFeature.timeout=30s \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json \
--set grpc.endpoints.recordRoute.timeout=30s \
--set grpc.protoURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/route_guide/routeguide/route_guide.proto \
--set assess.SLOs.upper.grpc-getFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-listFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-recordRoute/error-rate=0 \
--set runner=job
```

***

Further more, set parameters will trickle down to the endpoints.

```bash
iter8 k launch \
--set "tasks={grpc,assess}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.skipFirst=5 \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeature.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.listFeature.timeout=30s \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json \
--set grpc.protoURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/route_guide/routeguide/route_guide.proto \
--set assess.SLOs.upper.grpc-getFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-listFeature/error-rate=0 \
--set assess.SLOs.upper.grpc-recordRoute/error-rate=0 \
--set runner=job
```

In this example, all three endpoints will have a `skipFirst` of 5.

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

***

In the case of multiple endpoints, the name of the endpoint will be appended to the name of the provider. For example, if the endpoint name is `routeguide`, then the following metrics would be collected by this task:

- `grpc-routeguide/request-count`: total number of requests sent
- `grpc-routeguide/error-count`: number of error responses
- `grpc-routeguide/error-rate`: fraction of error responses
- `grpc-routeguide/latency/mean`: mean latency
- `grpc-routeguide/latency/stddev`: standard deviation of latency
- `grpc-routeguide/latency/min`: min latency
- `grpc-routeguide/latency/max`: max latency
- `grpc-routeguide/latency/pX`: X^th^ percentile latency, for any X in the range 0.0 to 100.0

To learn more about the names of metrics, please see [here](../topics/metrics.md#fully-qualified-names).