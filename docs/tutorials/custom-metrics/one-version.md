---
template: main.html
---

# SLO validation using custom metrics (single version)

Validate [SLOs](../../getting-started/concepts.md#service-level-objectives) for an app by fetching the app's metrics from a database (like Prometheus). This is a [multi-loop](../../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../../getting-started/concepts.md#kubernetes-experiments).

<p align='center'>
  <img alt-text="custom-metrics-one-version" src="../images/one-version.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md). Understand the main [concepts](../../getting-started/concepts.md) behind Iter8 experiments.
    2. [Install Istio](https://istio.io/latest/docs/setup/install/).
    3. [Install Istio's Prometheus add-on](https://istio.io/latest/docs/ops/integrations/prometheus/).
    4. [Enable automatic Istio sidecar injection](https://istio.io/latest/docs/setup/additional-setup/sidecar-injection/) for the `default` namespace. This ensures that the pods created in steps 5 and 6 will have the Istio sidecar.
    ```shell
    kubectl label namespace default istio-injection=enabled --overwrite
    ```
    5. Deploy the sample HTTP service in the Kubernetes cluster.
    ```shell
    kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
    kubectl expose deploy httpbin --port=80
    ```
    6. Generate load.
    ```shell
    kubectl run fortio --image=fortio/fortio --command -- fortio load -t 6000s http://httpbin.default/get
    ```

***

## Launch experiment

```shell
iter8 k launch \
--set "tasks={custommetrics,assess}" \
--set custommetrics.templates.istio-prom="https://raw.githubusercontent.com/iter8-tools/hub/main/templates/custommetrics/istio-prom.tpl" \
--set custommetrics.values.labels.namespace=default \
--set custommetrics.values.labels.destination_app=httpbin \
--set assess.SLOs.upper.istio-prom/error-rate=0 \
--set assess.SLOs.upper.istio-prom/latency-mean=100 \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

??? note "About this experiment"
    This experiment consists of two [tasks](../../getting-started/concepts.md#iter8-experiment), namely, [custommetrics](../../user-guide/tasks/custommetrics.md), and [assess](../../user-guide/tasks/assess.md). 
    
    The `custommetrics` task in this experiment [works](../../user-guide/tasks/custommetrics.md#how-it-works) by downloading a [provider template](../../user-guide/tasks/custommetrics.md#provider-template) named `istio-prom` from a URL, [substituting the template variables with values](../../user-guide/tasks/custommetrics.md#computing-variable-values), using the resulting [provider spec](../../user-guide/tasks/custommetrics.md#provider-spec) to query Prometheus for metrics, and [processing the response from Prometheus](../../user-guide/tasks/custommetrics.md#processing-response) to extract the metric values. Metrics defined by this template include `error-rate` and `latency-mean`; the [Prometheus labels](https://istio.io/latest/docs/reference/config/metrics/#labels) used by this template are stored in `labels`; all the metrics and variables associated with this template are [documented as part of the template](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/custommetrics/istio-prom.tpl). 
    
    The [assess](../../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, and ii) the mean latency of the app does not exceed 100 msec. 

    This is a [multi-loop](../../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../../getting-started/concepts.md#kubernetes-experiments). Hence, its [runner](../../getting-started/concepts.md#runners) value is set to `cronjob`. The `cronjobSchedule` expression specifies that each experiment loop (i.e., the sequence of tasks in the experiment) is scheduled for execution periodically once every minute. This enables Iter8 to refresh the metric values and perform SLO validation using the latest metric values during each loop.

***

## Assert experiment outcomes

--8<-- "docs/tutorials/custom-metrics/assert.md"

***

View experiment report and logs, and cleanup as described in [your first experiment](../../getting-started/your-first-experiment.md).

***

??? note "Some variations and extensions of this experiment"
    1. Perform [SLO validation for multiple versions of an app using custom metrics](two-or-more-versions.md).
    2. Define and use your own provider templates. This enables you to use any app-specific metrics from any database as part of Iter8 experiments. Read the [documentation for the `custommetrics` task](../../user-guide/tasks/custommetrics.md) to learn more.
    3. Alter the `cronjobSchedule` expression so that experiment loops are repeated at a frequency of your choice. Use use [https://crontab.guru](https://crontab.guru) to learn more about `cronjobSchedule` expressions.

