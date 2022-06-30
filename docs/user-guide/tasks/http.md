---
template: main.html
---

# http

Generate requests for an HTTP service and and collect [latency and error-related metrics](#metrics).

## Usage example
In this experiment, the `http` task generates requests for `https://httpbin.org/get`, and collects latency and error-related metrics. The metrics are used by the `assess` task to validate SLOs.

```
iter8 launch \
--set "tasks={http,assess}" \
--set http.url=https://httpbin.org/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0
```

## Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| url  | string (URL) | HTTP URL where requests are sent. |
| headers  | map[string]string | HTTP headers to use in the requests. |
| numRequests  | int | Number of requests to be sent to the app. Default value is 100. |
| duration  | string | Duration of this task. Specified in the [Go duration string format](https://pkg.go.dev/maze.io/x/duration#ParseDuration) (example, `5s`). If both duration and numRequests are specified, then duration is ignored. |
| qps  | float | qps stands for queries-per-second. Number of requests per second sent to the app. Default value is 8.0. |
| connections  | int | Number of parallel connections used to send requests. Default value is 4. |
| payloadURL  | string (URL) | URL from which to download the content that will be used as the request payload. If this field is specified, Iter8 will send HTTP POST requests to the app using this content as the payload. |
| payloadStr  | string | String data to be used as the request payload. If this field is specified, Iter8 will send HTTP POST requests to the app using this string as the payload. |
| contentType  | string | Content type of the payload. This is intended to be used in conjunction with one of the `payload*` fields. If this field is specified, Iter8 will send HTTP POST requests to the app using this as the Content-Type header value. |


## Metrics

The following metrics are collected by this task:

- `http/request-count`: total number of requests sent
- `http/error-count`: number of error responses
- `http/error-rate`: fraction of error responses
- `http/latency-mean`: mean of observed latency values
- `http/latency-stddev`: standard deviation of observed latency values
- `http/latency-min`: min of observed latency values
- `http/latency-max`: max of observed latency values
- `http/latency-pX`: X^th^ percentile latency, for X in `[50.0, 75.0, 90.0, 95.0, 99.0, 99.9]`

All latency metrics have `msec` units. Note that the name of the provider for all the metrics collected by this task is `http`.
