---
template: main.html
---

---
template: main.html
---

# custommetrics

Fetch metrics from databases (like Prometheus).

## Usage Example
In this example, the `custommetrics` task fetches metrics from the Prometheus database that is created by [Istio's Prometheus add-on](https://istio.io/latest/docs/ops/integrations/prometheus/). 

```shell
iter8 k launch \
--set "tasks={custommetrics,assess}" \
--set custommetrics.templates.istio-prom="https://raw.githubusercontent.com/iter8-tools/iter8/master/custommetrics/istio-prom.tpl" \
--set custommetrics.values.destinationWorkload=httpbin \
--set custommetrics.values.destinationWorkloadNamespace=default \
--set assess.SLOs.upper.istio-prom/error-rate=0 \
--set assess.SLOs.upper.istio-prom/latency-mean=100 \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

[This section](#concepts) describes the concepts of [provider spec](#provider-spec) and [provider template](#provider-template) that are central to this task. [This section](#how-it-works) describes how this task works under the hood.

## Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| templates  | map[string]string | A map where each key is the name of a provider, and the corresponding value is a URL containing the [provider template](#provider-template). |
| values  | map[string]interface{} | A map that contains the values for variables in [provider templates](#provider-template). When there are two or more app versions, this map contains values that are common to all versions. |
| versionValues  | []map[string]interface{} | An array that contains version-specific values for variables in [provider templates](#provider-template). While fetching metrics for version `i`, the task merges `values` with `versionValues[i]` (latter takes precedence), and the merged map contains the values for variables in provider templates. |

## Concepts

### Provider spec

Iter8 needs the information following in order to fetch metrics from a database.

1. The HTTP URL where the database can be queried.
2. The HTTP headers and method (GET/POST) to be used while querying the database.
3. For each metric to be fetched from the database:
    * The specific HTTP query to be used, in particular, the HTTP query parameters and body (if any).
    * The logic for parsing the query response and retrieving the metric value.

The above information is encapsulated by `ProviderSpec`, a data structure which Iter8 associates with each provider in this task.

???+ tip "Golang type definitions for ProviderSpec and Metric"
    ```go linenums="1"
    type ProviderSpec struct {
      // URL is the database endpoint
      URL string `json:"url" yaml:"url"`
      // Method is the HTTP method that needs to be used
      Method string `json:"method" yaml:"method"`
      // Headers is the set of HTTP headers that need to be sent
      Headers map[string]string `json:"headers" yaml:"headers"`
      // Metrics is the set of metrics that can be obtained
      Metrics []Metric `json:"metrics" yaml:"metrics"`
    }

    type Metric struct {
      // Name is the name of the metric
      Name string `json:"name" yaml:"name"`
      // Description is the description of the metric
      Description *string `json:"description,omitempty" yaml:"description,omitempty"`
      // Type is the type of the metric, either gauge or counter
      Type string `json:"type" yaml:"type"`
      // Units is the unit of the metric, which can be omitted for unitless metrics
      Units *string `json:"units,omitempty" yaml:"units,omitempty"`
      // Params is the set of HTTP parameters that need to be sent
      Params *[]HTTPParam `json:"params,omitempty" yaml:"params,omitempty"`
      // Body is the HTTP request body that needs to be sent
      Body *string `json:"body,omitempty" yaml:"body,omitempty"`
      // JqExpression is the jq expression that can extract the value from the HTTP
      // response
      JqExpression string `json:"jqExpression" yaml:"jqExpression"`
    }

    type HTTPParam struct {
      // Name is the name of the HTTP parameter
      Name string `json:"name" yaml:"name"`
      // Value is the value of the HTTP parameter
      Value string `json:"value" yaml:"value"`
    }
    ```

### Provider template

An initial idea would be for users to supply one or more [provider specs](#provider-spec), so that Iter8 can construct the metric queries. Iter8 builds on this idea by letting users supply one or more [Golang templates](https://pkg.go.dev/text/template) for provider specs. When a provider template is combined with [values](#computing-variable-values), it generates a [provider spec](#provider-spec) in YAML format that can be used by Iter8.

??? tip "`istio-prom` provider template in the usage example"
    ```yaml linenums="1"
    # This file provides templated metric specifications that enable
    # Iter8 to retrieve metrics from Istio's Prometheus add-on.
    # 
    # For a list of metrics supported out-of-the-box by the Istio Prom add-on, 
    # please see https://istio.io/latest/docs/reference/config/metrics/
    #
    # Iter8 substitutes the placeholders in this file with values, 
    # and uses the resulting metric specs to query Prometheus.
    # The placeholders are as follows.
    # 
    # reporter                        string  optional
    # destinationWorkload             string  required
    # destinationWorkloadNamespace    string  required
    # elapsedTimeSeconds              int     implicit
    # startingTime                    string  optional
    # latencyPercentiles              []int   optional
    #
    # For descriptions of reporter, destinationWorkload, and destinationWorkloadNamespace, 
    # please see https://istio.io/latest/docs/reference/config/metrics/
    #
    # elapsedTimeSeconds: this should not be specified directly by the user. 
    # It is implicitly computed by Iter8 according to the following formula
    # elapsedTimeSeconds := (time.Now() - startingTime).Seconds()
    # 
    # startingTime: By default, this is the time at which the Iter8 experiment started.
    # The user can explicitly specify the startingTime for each app version
    # (for example, the user can set the startingTime to the creation time of the app version)
    #
    # latencyPercentiles: Each item in this slice will create a new metric spec.
    # For example, if this is set to [50,75,90,95],
    # then, latency-p50, latency-p75, latency-p90, latency-p95 metric specs are created.

    {{- define "istio-prom-reporter"}}
    {{- if .reporter }}
    reporter="{{ .reporter }}",
    {{- end }}
    {{- end }}

    {{- define "istio-prom-dest"}}
    {{ template "istio-prom-reporter" . }}
    destination_workload="{{ .destinationWorkload }}",
    destination_workload_namespace="{{ .destinationWorkloadNamespace }}"
    {{- end }}

    # url is the HTTP endpoint where the Prometheus service installed by Istio's Prom add-on
    # can be queried for metrics

    url: {{ default .istioPromURL "http://prometheus.istio-system:9090/api/v1/query" }}
    provider: istio-prom
    method: GET
    metrics:
    - name: request-count
      type: counter
      description: |
        Number of requests
      params:
      - name: query
        value: |
          sum(last_over_time(istio_requests_total{
            {{ template "istio-prom-dest" . }}
          }[{{ .elapsedTimeSeconds }}s])) or on() vector(0)
      jqExpression: .data.result[0].value[1] | tonumber
    - name: error-count
      type: counter
      description: |
        Number of unsuccessful requests
      params:
      - name: query
        value: |
          sum(last_over_time(istio_requests_total{
            response_code=~'5..',
            {{ template "istio-prom-dest" . }}
          }[{{ .elapsedTimeSeconds }}s])) or on() vector(0)
      jqExpression: .data.result[0].value[1] | tonumber
    - name: error-rate
      type: gauge
      description: |
        Fraction of unsuccessful requests
      params:
      - name: query
        value: |
          (sum(last_over_time(istio_requests_total{
            response_code=~'5..',
            {{ template "istio-prom-dest" . }}
          }[{{ .elapsedTimeSeconds }}s])) or on() vector(0))/(sum(last_over_time(istio_requests_total{
            {{ template "istio-prom-dest" . }}
          }[{{ .elapsedTimeSeconds }}s])) or on() vector(0))
      jqExpression: .data.result.[0].value.[1]
    - name: latency-mean
      type: gauge
      description: |
        Mean latency
      params:
      - name: query
        value: |
          (sum(last_over_time(istio_request_duration_milliseconds_sum{
            {{ template "istio-prom-dest" . }}
          }[{{ .elapsedTimeSeconds }}s])) or on() vector(0))/(sum(last_over_time(istio_requests_total{
            {{ template "istio-prom-dest" . }}
          }[{{ .elapsedTimeSeconds }}s])) or on() vector(0))
      jqExpression: .data.result[0].value[1] | tonumber
    {{- range $i, $p := .latencyPercentiles }}
    - name: latency-p{{ $p }}
      type: gauge
      description: |
        {{ $p }} percentile latency
      params:
      - name: query
        value: |
          histogram_quantile(0.{{ $p }}, sum(rate(istio_request_duration_milliseconds_bucket{
            {{ template "istio-prom-dest" $ }}
          }[{{ .elapsedTimeSeconds }}s])) by (le))
      jqExpression: .data.result[0].value[1] | tonumber
    {{- end }}    
    ```    

## How it works

The `custommetrics` tasks works as illustrated in the flowchart below.

```mermaid
graph TD
  A([Start]) --> B([Get provider template]);
  B --> C([Compute variable values]);
  C --> D([Create provider spec by combining provider template with values]);
  D --> E([Query database]);
  E --> F([Process response]);
  F --> G([Update metric value in experiment]);
  G --> H{Done with all metrics?};
  H ---->|No| E;
  H ---->|Yes| I{Done with all versions?};
  I ---->|No| C;
  I ---->|Yes| J([End]);
```

In order to create provider templates, users need an understanding of the `Compute variable values` and `Process response` steps in the above flowchart. We describe these steps below.


### Computing variable values

=== "One version"

    See [usage example](#usage-example).

=== "Two or more versions"

    Template variable substitution for two or more versions.

    ```shell
    iter8 k launch \
    --set "tasks={custommetrics,assess}" \
    --set custommetrics.templates.istio-prom="https://raw.githubusercontent.com/iter8-tools/iter8/master/custommetrics/istio-prom.tpl" \
    --set custommetrics.values.destinationWorkloadNamespace=default \
    --set custommetrics.values.reporter=destination \
    --set custommetrics.versionValues[0].destinationWorkload=httpbin-v1 \
    --set custommetrics.versionValues[1].destinationWorkload=httpbin-v2 \
    --set assess.SLOs.upper.istio-prom/error-rate=0 \
    --set assess.SLOs.upper.istio-prom/latency-mean=100 \
    --set runner=cronjob \
    --set cronjobSchedule="*/1 * * * *"
    ```

=== "Elapsed time"
    ```mermaid
    graph TD
      A([Start]) --> B{startingTime parameter supplied?};
      B ---->|Yes| C([elapsedTimeSeconds := currentTime - startingTime]);
      B ---->|No| D([startingTime := time when experiment was launched]);
      D --> C;
      C --> E([End]);
    ```

    The placeholder `elapsedTimeSeconds` is substituted based on the start of the experiment or `startingTime`, if provided in the CLI. If `startingTime` is provided, then the time should follow [RFC 3339](https://www.rfc-editor.org/rfc/rfc3339) (for example: `2020-02-01T09:44:40Z` or `2020-02-01T09:44:40.954641934Z`).

### Processing response

The metrics provider is expected to respond to Iter8's HTTP request for a metric with a JSON object. The format of this JSON object is provider-specific. Iter8 uses [jq](https://stedolan.github.io/jq/) to extract the metric value from the JSON response of the provider. The `jqExpression` used by Iter8 is supplied as part of the metric definition. When the `jqExpression` is applied to the JSON response, it is expected to yield a number.


=== "Prometheus response example"
    The format of the Prometheus JSON response is [defined here](https://prometheus.io/docs/prometheus/latest/querying/api/#format-overview). A sample Prometheus response is as follows.
    ```json linenums="1"
    {
      "status": "success",
      "data": {
        "resultType": "vector",
        "result": [
          {
            "value": [1556823494.744, "21.7639"]
          }
        ]
      }
    }    
    ```

=== "Prometheus jqExpression example"
    Consider the `jqExpression` defined in the [sample Prometheus metric](#metrics-withwithout-auth). Let us apply it to the [sample JSON response from Prometheus](#json-response).
    ```shell
    echo '{
      "status": "success",
      "data": {
        "resultType": "vector",
        "result": [
          {
            "value": [1556823494.744, "21.7639"]
          }
        ]
      }
    }' | jq ".data.result[0].value[1] | tonumber"
    ```
    Executing the above command results yields `21.7639`, a number, as required by Iter8. 

    > **Note:** The shell command above is for illustration only. Iter8 uses Python bindings for `jq` to evaluate the `jqExpression`.


# Custom Metrics

Custom Iter8 metrics enable you to use data from any database for evaluating app/ML model versions within Iter8 experiments. This document describes how you can define custom Iter8 metrics and (optionally) supply authentication information that may be required by the metrics provider.

Metric providers differ in the following aspects.

* HTTP request authentication method: no authentication, basic auth, API keys, or bearer token
* HTTP request method: GET or POST
* Format of HTTP parameters and/or JSON body used while querying them
* Format of the JSON response returned by the provider
* The logic used by Iter8 to extract the metric value from the JSON response

The examples in this document focus on Prometheus, NewRelic, Sysdig, and Elastic. However, the principles illustrated here will enable you to use metrics from any provider in experiments.

## Metrics with/without auth

> **Note:** Metrics are defined by you, the **Iter8 end-user**.

=== "Prometheus"

    Prometheus does not support any authentication mechanism *out-of-the-box*. However,
    Prometheus can be setup in conjunction with a reverse proxy, which in turn can support HTTP request authentication, as described [here](https://prometheus.io/docs/guides/basic-auth/).

    === "No Authentication"
        The following is an example of an Iter8 metric with Prometheus as the provider. This example assumes that Prometheus can be queried by Iter8 without any authentication.

        ```yaml linenums="1"
        url: http://127.0.0.1:9090/api/v1/query
        provider: istio
        method: GET
        metrics:
        - name: request-count
          type: counter
          description: |
            Number of requests
          params:
          - name: query
            value: |
              sum(last_over_time(istio_requests_total{
                  reporter="source",    
                {{- if .destination_workload }}
                  destination_workload="{{.destination_workload}}",
                {{- end }}
                {{- if .destination_workload_namespace }}
                  destination_workload_namespace="{{.destination_workload_namespace}}",
                {{- end }}
              }[{{.elapsedTimeSeconds}}s])) or on() vector(0)
          jqExpression: .data.result[0].value[1]
        ```

    ??? hint "Brief explanation of the `request-count` metric"

        1. The HTTP query used by Iter8 contains a single query parameter named `query` as [required by Prometheus](https://prometheus.io/docs/prometheus/latest/querying/api/). The value of this parameter is derived by [substituting the placeholders](#placeholder-substitution) in the value string.
        2. The `url` field provides the URL of the Prometheus service.
        3. The `method` field provides the HTTP method, in this case `GET`.
        4. The `jqExpression` enables Iter8 to extract the metric value from the JSON response returned by Prometheus.

