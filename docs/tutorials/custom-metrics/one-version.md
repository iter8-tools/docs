---
template: main.html
---

# SLO validation using custom metrics

Validate [SLOs](../../getting-started/concepts.md#service-level-objectives) for an app using [custom metrics](custommetrics.md) provided by a database.

<p align='center'>
  <img alt-text="custom-metrics-one-version" src="../images/one-version.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md).
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
--set custommetrics.templates.istio-prom="https://raw.githubusercontent.com/iter8-tools/iter8/master/custommetrics/istio-prom.tpl" \
--set custommetrics.values.destinationWorkload=httpbin \
--set custommetrics.values.destinationWorkloadNamespace=default \
--set assess.SLOs.upper.istio-prom/error-rate=0 \
--set assess.SLOs.upper.istio-prom/latency-mean=100 \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

???+ note "About this experiment"
    This experiment consists of two [tasks](tasks.md), namely, [custommetrics](custommetrics.md), and [assess](assess.md). 
    
    The [custommetrics](custommetrics.md) task downloads a metrics template named `istio-prom` from a URL, [substitutes the template variables](substitution.md) with `custommetrics.values`, and uses the resulting metrics spec to [fetch metrics](fetchmetrics.md) from Prometheus. Metrics defined by this template include `error-rate` and `latency-mean`; variables used by this template include `destinationWorkload` and `destinationWorkloadNamespace`; they are [documented as part of the template](https://raw.githubusercontent.com/iter8-tools/iter8/master/custommetrics/istio-prom.tpl). 
    
    The [assess](assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, and ii) the mean latency of the app does not exceed 100 msec. 
    
    This is an experiment with [repeated loops](loops.md): the [runner](runner.md) specifies that the experiment should be [run using a Kubernetes cronjob](runner.md), and the `cronjobSchedule` specifies that experiment loops are repeated every minute. Both the experiment tasks are executed in each loop. This enables Iter8 to refresh the metric values and perform SLO validation using the latest metric values in each loop.

***

## Assert experiment outcomes

--8<-- "docs/tutorials/custom-metrics/assert.md"

***

View experiment report and logs, and cleanup as described in [your first experiment](../../getting-started/your-first-experiment.md).
