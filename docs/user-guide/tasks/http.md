---
template: main.html
---

# http

Generate requests for an HTTP service and and collect [latency and error-related metrics](#metrics).

## Usage example

In this experiment, the `http` task generates requests for `https://httpbin.org/get`, and collects latency and error-related metrics. Metrics collected by this task are viewable with an Iter8 dashboard.

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

## Grafana Dashboard

The results of the `http` task is visualized using the `http` Iter8 Grafana dashboard. The dashboard can be found [here](https://raw.githubusercontent.com/iter8-tools/iter8/v0.16.0/grafana/http.json).

To use the dashboard:

1. Open Grafana in a browser. 
2. Add a new data JSON API data source with the following parameters
    * URL: `<link to Grafana service>/httpDashboard`
    * Query string: `namespace=<namespace of experiment>&experiment=<name of experiment>`
3. Import the `http` Iter8 Grafana dashboard
    * Copy and paste the contents of this [link](https://raw.githubusercontent.com/iter8-tools/iter8/v0.16.0/grafana/http.json) into the text box

You will see a visualization of the experiment like the following:

![`http` Iter8 dashboard](images/httpdashboard.png)

For multiple endpoints, the visualization will look like the following:

![`http` Iter8 dashboard with multiple endpoints](images/httpmultipledashboard.png)
