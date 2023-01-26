---
template: main.html
---

# SLO validation using custom metrics (multiple versions)

Validate [SLOs](../../getting-started/concepts.md#service-level-objectives) for multiple versions of an app by fetching metrics for each app version from a database (like Prometheus). This is a [multi-loop](../../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../../getting-started/concepts.md#kubernetes-experiments).

<p align='center'>
  <img alt-text="custom-metrics-two-or-more-versions" src="../images/two-or-more-versions.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md). Understand the main [concepts](../../getting-started/concepts.md) behind Iter8 experiments. Try [an SLO validation experiment using custom metrics for a single version of an app](one-version.md).
    2. [Complete the Istio traffic mirroring tutorial](https://istio.io/latest/docs/tasks/traffic-management/mirroring/), specifically, the [Before you begin](https://istio.io/latest/docs/tasks/traffic-management/mirroring/#before-you-begin) section, the [Creating a default routing policy](https://istio.io/latest/docs/tasks/traffic-management/mirroring/#creating-a-default-routing-policy) section, and the [Mirroring traffic to v2](https://istio.io/latest/docs/tasks/traffic-management/mirroring/#mirroring-traffic-to-v2) section. Omit the the `Cleaning up` step (you can clean up once you are done with this tutorial).
    3. [Install Istio's Prometheus add-on](https://istio.io/latest/docs/ops/integrations/prometheus/).
    4. Generate load.
    ```shell
    kubectl run fortio --image=fortio/fortio --command -- fortio load -t 6000s http://httpbin.default:8000/get
    ```
***

## Launch experiment

```shell
iter8 k launch \
--set "tasks={custommetrics,assess}" \
--set custommetrics.templates.istio-prom="https://raw.githubusercontent.com/iter8-tools/hub/iter8-0.13.0/templates/custommetrics/istio-prom.tpl" \
--set custommetrics.values.labels.namespace=default \
--set custommetrics.values.labels.destination_app=httpbin \
--set custommetrics.values.labels.reporter=destination \
--set 'custommetrics.versionValues[0].labels.destination_version=v1' \
--set 'custommetrics.versionValues[1].labels.destination_version=v2' \
--set assess.SLOs.upper.istio-prom/error-rate=0 \
--set assess.SLOs.upper.istio-prom/latency-mean=100 \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

??? note "About this experiment"
    This experiment extends the [SLO validation experiment using custom metrics for a single app version](one-version.md). There are two versions of the app in this experiment. Variable values that are specific to the first version are specified under `custommetrics.versionValues[0]`, while those that are specific to the second version are specified under `custommetrics.versionValues[1]`. For the first version, Iter8 merges `custommetrics.values` with `custommetrics.versionValues[0]` (the latter takes precedence), and uses the result for template variable substitution. Similarly, for the second version, Iter8 merges `custommetrics.values` with `custommetrics.versionValues[1]` (the latter takes precedence), and uses the result for template variable substitution.

***

Assert experiment outcomes, view experiment report, view experiment logs, and cleanup as described in [this experiment](../../tutorials/custom-metrics/one-version.md).

***

??? note "Some variations and extensions of this experiment"
    1. Define and use your own provider templates. This enables you to use any app-specific metrics from any database as part of Iter8 experiments. Read the [documentation for the `custommetrics` task](../../user-guide/tasks/custommetrics.md) to learn more.
    2. Alter the `cronjobSchedule` expression so that experiment loops are repeated at a frequency of your choice. Use use [https://crontab.guru](https://crontab.guru) to learn more about `cronjobSchedule` expressions.

    