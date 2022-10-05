---
template: main.html
---

# abnmetrics

Fetch metrics from the Iter8 A/B(/n) service.

## Usage example
In this experiment, the `abnmetrics` task fetches metrics from the A/B(/n) service for the application `default/backend`. The task is run periodically (as defined by `cronjobSchedule`).

```
iter8 launch \
--set "tasks={abnmetrics}" \
--set http.url=https://httpbin.org/get \
--set abnmetrics.application=default/backend
-set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

## Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| endpont  | string | Endpoint of the A/B(/n) service. Defaults to `iter8-abn:50051` |
| application | string | Application name in form `namespace/name` |
