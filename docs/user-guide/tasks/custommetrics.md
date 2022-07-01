---
template: main.html
---

---
template: main.html
---

# custommetrics

Fetch metrics from databases (like Prometheus).

## Usage Example
In this example, the `custommetrics` task fetches metrics from the Prometheus database created using [Istio's Prometheus add-on](https://istio.io/latest/docs/ops/integrations/prometheus/). 

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

[This section](#concepts) describes the concepts of [provider spec](#provider-spec) and [provider template](#provider-template), [this section](#how-it-works) describes how the task works behind the scenes, and [this section](#customizing) describes how you can customize this task with your own metric providers.


## Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| templates  | map[string]string | A map where each key is the name of a provider, and the corresponding value is a URL containing the [provider template](#provider-template). |
| values  | map[string]interface{} | A map that contains the values for variables in [provider templates](#provider-template). When there are two or more app versions, this map contains values that are common to all versions. |
| versionValues  | []map[string]interface{} | An array that contains version-specific values for variables in [provider templates](#provider-template). While fetching metrics for version `i`, the task merges `values` with `versionValues[i]` (latter takes precedence), and the merged map contains the values for variables in provider templates. |

## Concepts

### Provider spec

Golang struct here... and a yaml under twistie here.

### Provider template

Provider template link here. Describe the variables.

## How it works

Flow chart here ... 

### Template variable values

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


### Elapsed time

### Response handling

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

### Error handling
Errors may occur during Iter8's metric queries due to a number of reasons (for example, due to an invalid metrics template template, or due to the metrics database being unreachable). If errors are encountered during the attempt to retrieve metric values, Iter8 will mark the corresponding metric as unavailable.

## Customizing

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

## Placeholder substitution

> **Note:** This step is automated by **Iter8**.

=== "Prometheus"
  The placeholder `elapsedTimeSeconds` is substituted based on the start of the experiment or `startingTime`, if provided in the CLI. If `startingTime` is provided, then the time should follow [RFC 3339](https://www.rfc-editor.org/rfc/rfc3339) (for example: `2020-02-01T09:44:40Z` or `2020-02-01T09:44:40.954641934Z`).