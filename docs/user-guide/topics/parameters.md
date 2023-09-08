---
template: main.html
---

# Performance test parameters

Iter8 is built on [Helm](https://helm.sh). Performance tests can be configured with parameters using the same mechanisms provided by Helm for [setting chart values](https://helm.sh/docs/chart_template_guide/values_files/#helm). 

The set of configurable parameters for a performance test includes the parameters of the tasks involved in the test. Iter8 uses the convention that the parameters of a task are nested under the name of that task. In the following example, the `url` parameter of the `http` task is nested under the `http` object.

```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 0.17 httpbin-test iter8 \
--set "tasks={http}" \
--set http.url=https://httpbin.org/get
```

All the parameters of the performance test (including of all included tasks) are optional unless otherwise documented.

## Parameters

Global performance test parameters are described here. Task specific parameters are documented in each task description.

| Name | Type | Description |
| ---- | ---- | ----------- |
| serviceAccountName  | string | Optional name of a service account to use. If specified, it is assumed the service account has the necessary permissions to run a performance test. If not specified, Iter8 will create a service account. |
| logLevel | string | Log level. Must be one of `trace`, `debug`, `info` (default), `warning`, or `error`. |
