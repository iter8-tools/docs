---
template: main.html
---

# http

Generate requests for an HTTP service and and collect [latency and error-related metrics](#metrics).

## Usage example

In this experiment, the `http` task generates requests for `https://httpbin.org/get`, and collects latency and error-related metrics. The metrics are used by the `assess` task to validate SLOs.

Single endpoint:
```bash
iter8 k launch \
--set "tasks={http}" \
--set http.url=https://httpbin.org/get
```

Multiple endpoints:
```bash
iter8 k launch \
--set "tasks={http}" \
--set http.endpoints.get.url=http://httpbin.default/get \
--set http.endpoints.getAnything.url=http://httpbin.default/anything \
--set http.endpoints.post.url=http://httpbin.default/post \
--set http.endpoints.post.payloadStr=hello
```

## Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| url | string (URL) | URL where requests are sent. |
| headers | map[string]string | HTTP headers to use in the requests. |
| numRequests | int | Number of requests to be sent to the app. Default value is 100. |
| duration | string | Duration of this task. Specified in the [Go duration string format](https://pkg.go.dev/maze.io/x/duration#ParseDuration) (example, `5s`). If both duration and numRequests are specified, then duration is ignored. |
| qps | float | qps stands for queries-per-second. Number of requests per second sent to the app. Default value is 8.0. |
| connections | int | Number of parallel connections used to send requests. Default value is 4. |
| payloadURL | string (URL) | URL from which to download the content that will be used as the request payload. If this field is specified, Iter8 will send HTTP POST requests to the app using this content as the payload. |
| payloadStr | string | String data to be used as the request payload. If this field is specified, Iter8 will send HTTP POST requests to the app using this string as the payload. |
| contentType | string | Content type of the payload. This is intended to be used in conjunction with one of the `payload*` fields. If this field is specified, Iter8 will send HTTP POST requests to the app using this as the Content-Type header value. |
| warmupNumRequests | int | Number of requests to be sent in a warmup task (results are ignored). |
| warmupDuration | string | Duration of warmup task (results are ignored). Specified in the [Go duration string format](https://pkg.go.dev/maze.io/x/duration#ParseDuration) (example, 5s). If both warmupDuration and warmupNumRequests are specified, then warmupDuration is ignored. |
| endpoints | map[string]EndPoint | Used to specify multiple endpoints and their configuration. The `string` is the name of the endpoint and the `EndPoint` struct includes all the parameters described above. Load testing and metric collection will be conducted separately for each endpoint. |

## Precedence

Some parameters have a default value, which can be overwritten. In addition, with the `endpoints` parameter, you can test multiple endpoints and configure parameters for each of those endpoint. In these cases, the priority order is the default value, the value set at the base level, and the value set at the endpoint value.

In the following example, all three endpoints will use the default `qps` (queries-per-second) of 8.

```bash
iter8 k launch \
--set "tasks={http}" \
--set http.endpoints.get.url=http://httpbin.default/get \
--set http.endpoints.getAnything.url=http://httpbin.default/anything \
--set http.endpoints.post.url=http://httpbin.default/post \
--set http.endpoints.post.payloadStr=hello
```

In the following example, the `get` and `getAnything` endpoints will use the default `qps` of 8 and the `post` endpoint will use a `qps` of 15.

```bash
iter8 k launch \
--set "tasks={http}" \
--set http.endpoints.get.url=http://httpbin.default/get \
--set http.endpoints.getAnything.url=http://httpbin.default/anything \
--set http.endpoints.post.url=http://httpbin.default/post \
--set http.endpoints.post.payloadStr=hello \
--set http.endpoints.post.qps=15
```

In the following example, all three endpoints will use a `qps` (queries-per-second) of 10.

```bash
iter8 k launch \
--set "tasks={http}" \
--set http.qps=10 \
--set http.endpoints.get.url=http://httpbin.default/get \
--set http.endpoints.getAnything.url=http://httpbin.default/anything \
--set http.endpoints.post.url=http://httpbin.default/post \
--set http.endpoints.post.payloadStr=hello
```

In the following example, the `get` and `getAnything` endpoints will use a `qps` of 10 and the `post` endpoint will use a `qps` of 15.

```bash
iter8 k launch \
--set "tasks={http}" \
--set http.qps=10 \
--set http.endpoints.get.url=http://httpbin.default/get \
--set http.endpoints.getAnything.url=http://httpbin.default/anything \
--set http.endpoints.post.url=http://httpbin.default/post \
--set http.endpoints.post.payloadStr=hello \
--set http.endpoints.post.qps=15
```

***

Further more, set parameters will trickle down to the endpoints.

```bash
iter8 k launch \
--set "tasks={http}" \
--set http.numRequests=50 \
--set http.endpoints.get.url=http://httpbin.default/get \
--set http.endpoints.getAnything.url=http://httpbin.default/anything \
--set http.endpoints.post.url=http://httpbin.default/post \
--set http.endpoints.post.payloadStr=hello
```

In this example, all three endpoints will have a `numRequests` of 50.

## Metrics

This task creates a built-in [provider](../topics/metrics.md#fully-qualified-names) named `http`. The following metrics are collected by this task:

- `http/request-count`: total number of requests sent
- `http/error-count`: number of error responses
- `http/error-rate`: fraction of error responses
- `http/latency-mean`: mean of observed latency values
- `http/latency-stddev`: standard deviation of observed latency values
- `http/latency-min`: min of observed latency values
- `http/latency-max`: max of observed latency values
- `http/latency-pX`: X^th^ percentile latency, for X in `[50.0, 75.0, 90.0, 95.0, 99.0, 99.9]`

All latency metrics have `msec` units.

***

In the case of multiple endpoints, the name of the endpoint will be appended to the name of the provider. For example, if the endpoint name is `httpbin`, then the following metrics would be collected by this task:

- `http-httpbin/request-count`: total number of requests sent
- `http-httpbin/error-count`: number of error responses
- `http-httpbin/error-rate`: fraction of error responses
- `http-httpbin/latency-mean`: mean of observed latency values
- `http-httpbin/latency-stddev`: standard deviation of observed latency values
- `http-httpbin/latency-min`: min of observed latency values
- `http-httpbin/latency-max`: max of observed latency values
- `http-httpbin/latency-pX`: X^th^ percentile latency, for X in `[50.0, 75.0, 90.0, 95.0, 99.0, 99.9]`

To learn more about the names of metrics, please see [here](../topics/metrics.md#fully-qualified-names).