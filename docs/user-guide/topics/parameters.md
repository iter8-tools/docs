---
template: main.html
---

# Experiment Parameters

Iter8 is built on [Helm](https://helm.sh). Iter8 experiments can be configured with parameters using the same mechanisms provided by Helm for [setting chart values](https://helm.sh/docs/chart_template_guide/values_files/#helm). 

The set of configurable parameters of an experiment includes the parameters of the tasks involved in the experiment. Iter8 uses the convention that the parameters of a task are nested under the name of that task. In the following example, the `url` parameter of the `http` task is nested under the `http` object.

```shell
iter8 k launch \
--set "tasks={http}" \
--set http.url=https://httpbin.org/get
```

All the parameters of a task or an experiment are optional unless indicated otherwise in the documentation of the task or experiment.

## Parameters

Global experiment parameters are described here. Task specific parameters are documented in each task description.

| Name | Type | Description |
| ---- | ---- | ----------- |
| serviceAccountName  | string | Optional name of a service account to use. If specified, it is assumed the service account has the necessary permissions to run an experiment. If not specified, Iter8 will create a service account. |
| logLevel | string | Log level. Must be one of `trace`, `debug`, `info` (default), `warning`, or `error`. |
