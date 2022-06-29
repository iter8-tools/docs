---
template: main.html
---

# SLO validation using custom metrics

Validate [SLOs](slos.md) for multiple versions of an app using [custom metrics](custommetrics.md) provided by a database.

<p align='center'>
  <img alt-text="custom-metrics-two-or-more-versions" src="../images/two-or-more-versions.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md).
    2. Try the [custom metrics experiment with one app version](one-version.md).
    3. [Follow the Istio traffic mirroring tutorial](https://istio.io/latest/docs/tasks/traffic-management/mirroring/).
    4. [Install Prometheus plugin](https://istio.io/latest/docs/ops/integrations/prometheus/).
    5. Generate load.
    ```shell
    kubectl run fortio --image=fortio/fortio --command -- fortio load -t 6000s http://httpbin.default:8000/get
    ```
***

## Launch experiment

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

???+ note "About this experiment"
    This experiment extends the [custom metrics experiment with one app version](one-version.md). There are two versions of the app in this experiment. Variable values that are specific to the first version are specified under `custommetrics.versionValues[0]`, while those that are specific to the second version are specified under `custommetrics.versionValues[1]`. For the first version, Iter8 merges `custommetrics.values` with `custommetrics.versionValues[0]` (the latter takes precedence), and uses the result for template variable substitution. Similarly, for the second version, Iter8 merges `custommetrics.values` with `custommetrics.versionValues[1]` (the latter takes precedence), and uses the result for template variable substitution.
    