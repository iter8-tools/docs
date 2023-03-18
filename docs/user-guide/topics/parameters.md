---
template: main.html
---

# Experiment Parameters

Iter8 is built on [Helm](https://helm.sh). Iter8 experiments can be configured with parameters using the same mechanisms provided by Helm for [setting chart values](https://helm.sh/docs/chart_template_guide/values_files/#helm). 

The set of configurable parameters of an experiment includes the parameters of the tasks involved in the experiment. Iter8 uses the convention that the parameters of a task are nested under the name of that task. In the following example, the `url` parameter of the `http` task is nested under the `http` object, and the `SLOs` parameter of the `assess` task is nested under the `assess` object.

```shell
iter8 k launch \
--set "tasks={http,assess}" \
--set http.url=https://httpbin.org/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0
--set runner=job
```

All the parameters of a task or an experiment are optional unless indicated otherwise in the documentation of the task or experiment.

## Parameters

Global experiment parameters are described here. Task specific parameters are documented in each task description.

| Name | Type | Description |
| ---- | ---- | ----------- |
| serviceAccountName  | string | Optional name of a service account to use. If specified, it is assumed the service account has the necessary permissions to run an experiment. If not specified, Iter8 will create a service account. |
| runner | string | One of `job` or `cronjob` indicating whether the experiment has a single loop or multiple loops, respectively. |
| cronjobSchedule | string | Schedule for a multi-loop experiment. Required if `runner` is `cronjob`. Ignored otherwise. Expressed using Unix cronjob notation. |
| logLevel | string | Log level. Must be one of `trace`, `debug`, `info` (default), `warning`, or `error`. |
