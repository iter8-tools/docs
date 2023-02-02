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

