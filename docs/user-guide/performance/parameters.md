---
template: main.html
---

# Performance test parameters

Iter8 performance tests are run using [Helm](https://helm.sh). Tests can be configured with parameters in the same manner as any other Helm command -- using [chart values](https://helm.sh/docs/chart_template_guide/values_files/#helm) and the `--set` option.

The `tasks` parameter is used to identify the sequence of tasks that should be executed. Each task has its own parameters. Iter8 uses the convention that the parameters of a task are nested under the name of that task. In the following example, the `url` parameter of the `http` task is nested under the `http` object.

```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 0.18 httpbin-test iter8 \
--set "tasks={http}" \
--set http.url=https://httpbin.org/get
```

All the parameters of the performance test (including of all included tasks) are optional unless otherwise documented. The task-specific parameters are documented in the task documentation.

Currently available tasks are:

- [`http`](tasks/http.md) - generate synthetic load and capture performance metrics for HTTP endpoints
- [`grpc`](tasks/grpc.md) - generate synthetic load and capture performance metrics for gRPC methods
- [`ready`](tasks/ready.md) - check readiness of an object
- [`github`](tasks/github.md) - send a GitHub notification
- [`slack`](tasks/slack.md) - send a Slack notification

## Global parameters

In addition to task-specific parameters, the following global parameters are available:

| Name | Type | Description |
| ---- | ---- | ----------- |
| serviceAccountName  | string | Optional name of a service account to use. If specified, it is assumed the service account has the necessary permissions to run a performance test. If not specified, Iter8 will create a service account. |
| logLevel | string | Log level. Must be one of `trace`, `debug`, `info` (default), `warning`, or `error`. |
