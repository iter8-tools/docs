---
template: main.html
---

# Benchmark and Validate HTTP services

The `slo-validation-istio` collects latency and error-related metrics and validates service-level objectives (SLOs) from an application using Istio.

<p align='center'>
  <img alt-text="slo-validation-istio" src="../images/istio-overview.png" width="90%" />
</p>

<!-- ***

--8<-- "docs/tutorials/load-test-http/usecases.md"

*** -->

???+ warning "Before you begin"
    1. [Install Istio](https://istio.io/latest/docs/setup/install/).

    2. [Install Prometheus plugin](https://istio.io/latest/docs/ops/integrations/prometheus/).

    3. [Enable automatic Istio sidecar injection](https://istio.io/latest/docs/setup/additional-setup/sidecar-injection/) for the `default` namespace. This ensures that the pods created in steps 4 and 5 will have the Istio sidecar.
    ```shell
    kubectl label namespace default istio-injection=enabled --overwrite
    ```

    4. Deploy the sample HTTP service in the Kubernetes cluster.
    ```shell
    kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
    kubectl expose deploy httpbin --port=80
    ```

    5. Generate load.
    ```shell
    kubectl run fortio --image=fortio/fortio --command -- fortio load -t 6000s http://httpbin.default/get
    ```
***
## Validate SLOs for a service

The `slo-validation-istio` experiment will regularly check for metrics and validate SLOs.

***

The following metrics are collected by default by this experiment:

- `istio/request-count`: total number of requests sent
- `istio/error-count`: number of error responses
- `istio/error-rate`: fraction of error responses
- `istio/latency-mean`: mean of observed latency values
- `istio/le500ms-latency-percentile`: percentile of requests with less than or equal latency of 500ms

All latency metrics have `msec` units.

***

In order to run the `slo-validation-istio` experiment, the following need to be set.

- `cronjobSchedule`, the [cron schedule](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#cron-schedule-syntax) that determines when the metrics collection and SLO validation is run. It is currently configured to run every minute.
- `providerURL`, the URL of the metrics provider, in this case Prometheus database endpoint
- `destination_workload`, the name of the service
- `destination_workload_namespace`, the namespace of the service 

In addition, the SLOs to be validated also need to be set. In this case, `istio/error-rate` and `istio/latency-mean` will be validated.

```shell
iter8 k launch -c slo-validation-istio \
--set cronjobSchedule="*/1 * * * *" \
--set providerURL=http://prometheus.istio-system:9090 \
--set versionInfo.destination_workload=httpbin \
--set versionInfo.destination_workload_namespace=default \
--set SLOs.istio/error-rate=0 \
--set SLOs.istio/latency-mean=100
```

In the experiment above, the following SLOs are validated.

- error rate is 0
- mean latency is under 100 msec

***

## View experiment report

--8<-- "docs/tutorials/slo-validation-istio/expreport.md"

Because the `cronjobSchedule` has been set to run every minute, the report will change periodically.

***

## Assert experiment outcomes

--8<-- "docs/tutorials/slo-validation-istio/assert.md"
