---
template: main.html
---

# A/B/n Experiments with Rewards

This tutorial describes how to to use Iter8 to evaluate two or more versions on an application or ML model to identify the "best" version according to some reward metrics(s).

## Assumptions

We assume that you have deployed multiple versions of an application (or ML model) with the following characteristics:

- There is a way to route user traffic to the deployed versions. This might be done using the Iter8 SDK, the Iter8 traffic control features, or some other mechanism.
- Metrics, including reward metrics, are being exported to a metrics store such as Prometheus.
- Metrics can be retrieved from the metrics store by application (model) version.

In this tutorial, we mock a Prometheus service and demonstrate how to write an Iter8 experiment that evaluates reward metrics.

## Mock Prometheus

For simplicity, we use [mockoon](https://mockoon.com/) to create a mocked Prometheus service instead of deploying Prometheus itself:

```shell
kubectl create deploy prometheus-mock \
--image mockoon/cli:latest \
--port 9090 \
-- mockoon-cli start --daemon-off \
--port 9090 \
--data https://gist.githubusercontent.com/kalantar/2efbbdc4ac753e47399e9e1976655397/raw/e4e4a690a6e98eb9f230b6b9978b15ae2f3c7c8e/model-prometheus-abn-tutorial.json
kubectl expose deploy prometheus-mock --port 9090
```

## Define template

It is necessary to create a [_provider specification_](../../user-guide/tasks/custommetrics.md#provider-spec) that describes, for each metric, how Iter8 should fetch the metric value from the metrics store. Common to all metrics, is information about the provider url, the HTTP method to be used and any common headers. Further, for each metric, there is:
- metadata, such as name, type and description, 
- HTTP query parameters, and 
- a jq expression describing how to extract the metric value from the response.

For example, a specificaytion for the mean latency metric is: 

```
metric:
- name: latency-mean
  type: gauge
  description: |
    Mean latency
  params:
  - name: query
    value: |
      (sum(last_over_time(revision_app_request_latencies_sum{
        {{- template "labels" . }}
      }[{{ .elapsedTimeSeconds }}s])) or on() vector(0))/(sum(last_over_time(revision_app_request_latencies_count{
        {{- template "labels" . }}
      }[{{ .elapsedTimeSeconds }}s])) or on() vector(0))
  jqExpression: .data.result[0].value[1] | tonumber
```

Note that the template is parameterized. Values are provided by the Iter8 experiment at runtime.

A sample provider specification for our reward metrics is provider [here](https://gist.githubusercontent.com/kalantar/80c9efc0fd4cc34572d893cc82bdc4d2/raw/f3629aa62cdc9fd7e39ee2b6b113a8bf7b6b4463/model-prometheus-abn-tutorial.tpl).

It describes the following metrics:

- request-count
- latency-mean
- profit-mean

## Launch experiment

```shell
iter8 k launch \
--set "tasks={custommetrics,assess}" \
--set custommetrics.templates.model-prometheus="https://gist.githubusercontent.com/kalantar/80c9efc0fd4cc34572d893cc82bdc4d2/raw/f3629aa62cdc9fd7e39ee2b6b113a8bf7b6b4463/model-prometheus-abn-tutorial.tpl" \
--set custommetrics.values.labels.model_name=wisdom \
--set 'custommetrics.versionValues[0].labels.mm_vmodel_id=wisdom-1' \
--set 'custommetrics.versionValues[1].labels.mm_vmodel_id=wisdom-2' \
--set assess.SLOs.upper.model-prometheus/latency-mean=50 \
--set "assess.rewards.max={model-prometheus/profit-mean}" \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

This experiment executes in a loop (`runner=cronjob`), once every minute (see `cronjobSchedule`). It uses the `custommetrics` task to read metrics from the (mocked) Prometheus provider. Finally, the `assess` task verifies that the `latency-mean` is below 50 msec and identifies which version provides the greatest reward; that is, the greater mean profit.

## Inspect experiment report
--8<-- "docs/getting-started/expreport.md"

Because the experiment loops, the reported results will change over time.

***

## Cleanup

Delete the experiment:

```shell
iter8 k delete
```

Terminate the mocked Prometheus service:

```shell
kubectl delete deploy/prometheus-mock svc/prometheus-mock
```