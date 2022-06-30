---
template: main.html
---

# Experiment Parameters

Iter8 experiments can be configured with parameters using the mechanism provided by [Helm](https://helm.sh) for [setting chart values](https://helm.sh/docs/chart_template_guide/values_files/#helm). An example is as follows.

```shell
iter8 launch \
--set "tasks={http,assess}" \
--set http.url=https://httpbin.org/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0
```

The set of configurable parameters of an experiment includes the parameters of the tasks involved in the experiment. When configuring an experiment, Iter8 uses the convention that the parameters of a task are nested under the name of that task. In the above example, the `url` parameter of the `http` task is nested under the `http` object, and the `SLOs` parameter of the `assess` task is nested under the `assess` object.