---
template: main.html
---

# assess

Assess if [service-level objectives (SLOs)](../../getting-started/concepts.md#service-level-objectives) are satisfied by app versions.

## Usage example

In this experiment, the `assess` task validates if the `http/latency-mean` metric has a value that does not exceed 50, and the `http/error-count` metric has a value that does not exceed 0.
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
| SLOs  | struct | [Service-level objectives](../../getting-started/concepts.md#service-level-objectives) that will be validated by this task. This struct contains two fields [upper](#upper) and [lower](#lower). |

### upper

| Name | Type | Description |
| ---- | ---- | ----------- |
| upper  | map[string]float | Map keys are [fully-qualified metric names](../topics/metrics.md#fully-qualified-names) and map values are upper limits of those metrics.  |

### lower

| Name | Type | Description |
| ---- | ---- | ----------- |
| lower  | map[string]float | Map keys are [fully-qualified metric names](../topics/metrics.md#fully-qualified-names) and map values are lower limits of those metrics.  |
