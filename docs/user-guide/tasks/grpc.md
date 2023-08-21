---
template: main.html
---

# grpc

Generate requests for a gRPC service and and collect [latency and error-related metrics](#metrics).

## Usage example

In this experiment, the `grpc` task generates call requests for a gRPC service hosted at `hello.default:50051`, defined in the [protobuf](https://developers.google.com/protocol-buffers) file located at `grpc.protoURL`, with a gRPC method named `helloworld.Greeter.SayHello`. Metrics collected by this task are viewable with an Iter8 dashboard.

Single method:
```bash
iter8 k launch \
--set "tasks={grpc}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.protoURL=https://raw.githubusercontent.com/grpc/grpc-go/v1.52.0/examples/route_guide/routeguide/route_guide.proto \
--set grpc.call=routeguide.RouteGuide.GetFeature \
--set grpc.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json
```

Multiple methods:
```bash
iter8 k launch \
--set "tasks={grpc}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.protoURL=https://raw.githubusercontent.com/grpc/grpc-go/v1.52.0/examples/route_guide/routeguide/route_guide.proto \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeatures.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeatures.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json \
--set grpc.endpoints.routeChat.call=routeguide.RouteGuide.RouteChat \
--set grpc.endpoints.routeChat.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/bidirectional.json
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
--set "tasks={grpc}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.protoURL=https://raw.githubusercontent.com/grpc/grpc-go/v1.52.0/examples/route_guide/routeguide/route_guide.proto \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeatures.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeatures.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json
```

In the following example, the `getFeature` and `listFeatures` endpoints will use the default `timeout` of `20s` and the `recordRoute` endpoint will use a `timeout` of `30s`.

```bash
iter8 k launch \
--set "tasks={grpc}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.protoURL=https://raw.githubusercontent.com/grpc/grpc-go/v1.52.0/examples/route_guide/routeguide/route_guide.proto \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeatures.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeatures.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json
```

In the following example, all three endpoints will use a `qps` of `40s`.

```bash
iter8 k launch \
--set "tasks={grpc}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.protoURL=https://raw.githubusercontent.com/grpc/grpc-go/v1.52.0/examples/route_guide/routeguide/route_guide.proto \
--set grpc.timeout=40s \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeatures.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeatures.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json
```

In the following example, the `getFeature` and `listFeatures` endpoints will use a `timeout` of `40s` and the `listFeatures` endpoint will use a `timeout` of `30s`.

```bash
iter8 k launch \
--set "tasks={grpc}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.protoURL=https://raw.githubusercontent.com/grpc/grpc-go/v1.52.0/examples/route_guide/routeguide/route_guide.proto \
--set grpc.timeout=40s \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeatures.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeatures.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.listFeatures.timeout=30s \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json \
--set grpc.endpoints.recordRoute.timeout=30s
```

***

Further more, set parameters will trickle down to the endpoints.

```bash
iter8 k launch \
--set "tasks={grpc}" \
--set grpc.host=routeguide.default:50051 \
--set grpc.protoURL=https://raw.githubusercontent.com/grpc/grpc-go/v1.52.0/examples/route_guide/routeguide/route_guide.proto \
--set grpc.skipFirst=5 \
--set grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature \
--set grpc.endpoints.getFeature.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/unary.json \
--set grpc.endpoints.listFeatures.call=routeguide.RouteGuide.ListFeatures \
--set grpc.endpoints.listFeatures.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/server.json \
--set grpc.endpoints.listFeatures.timeout=30s \
--set grpc.endpoints.recordRoute.call=routeguide.RouteGuide.RecordRoute \
--set grpc.endpoints.recordRoute.dataURL=https://raw.githubusercontent.com/iter8-tools/docs/v0.13.13/samples/grpc-payload/client.json
```

In this example, all three endpoints will have a `skipFirst` of 5.

## Grafana Dashboard

The results of the `grpc` task is visualized using the `grpc` Iter8 Grafana dashboard. The dashboard can be found [here](https://raw.githubusercontent.com/iter8-tools/iter8/v0.16.0).

To use the dashboard:

1. Open Grafana in a browser. 
2. Add a new data JSON API data source with the following parameters
    * URL: `<link to Grafana service>/grpcDashboard`
    * Query string: `namespace=<namespace of experiment>&experiment=<name of experiment>`
3. Import the `grpc` Iter8 Grafana dashboard
    * Copy and paste the contents of this [link](https://raw.githubusercontent.com/iter8-tools/iter8/v0.16.0) into the text box

You will see a visualization of the experiment like the following:

![`grpc` Iter8 dashboard](images/grpcdashboard.png)

For multiple endpoints, the visualization will look like the following:

![`grpc` Iter8 dashboard with multiple endpoints](images/grpcmultipledashboard.png)