---
template: main.html
---

# Benchmark and Validate HTTP services with traffic mirroring

Traffic mirroring is a method of duplicating traffic to another service. This other service can be a new version. Traffic mirroring will allow you to safely test how the new version will behave using live traffic without any risk of affecting end users.

In this tutorial, you will learn how to set up a new version and traffic mirroring with Istio and use Iter8 to test the performance of the new version. 

<p align='center'>
  <img alt-text="slo-validation-istio" src="../images/istio-traffic-overview.png" width="90%" />
</p>

<!-- ***

--8<-- "docs/tutorials/load-test-http/usecases.md"

*** -->

???+ warning "Before you begin"
    <!-- Run the [httpbin](https://httpbin.org) sample service from a separate terminal.
    ```shell
    docker run -p 80:80 kennethreitz/httpbin
    ```
    You can also use [Podman](https://podman.io) or other alternatives to Docker in the above command. -->

    1. [Follow the Istio traffic mirroring tutorial](https://istio.io/latest/docs/tasks/traffic-management/mirroring/).

    2. [Install Prometheus plugin](https://istio.io/latest/docs/ops/integrations/prometheus/).

    3. Generate load.
    ```shell
    kubectl run fortio --image=fortio/fortio --command -- fortio load -t 6000s http://httpbin.default:8000/get
    ```
***

## Collect metrics from the mirrored service
To collect metrics, following must be set:

- `cronjobSchedule`, the [cron schedule](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#cron-schedule-syntax) that determines when the metrics collection and SLO validation is run. It is currently configured to run every minute.
- `providerURL`, the URL of the metrics provider, in this case Prometheus database endpoint
- `destination_workload`, the name of the service
- `destination_workload_namespace`, the namespace of the service 

```shell
iter8 k launch -c slo-validation-istio \
--set reporter=destination \
--set cronjobSchedule="*/1 * * * *" \
--set providerURL=http://prometheus.istio-system:9090/api/v1/query \
--set versionInfo.destination_workload=httpbin-v2 \
--set versionInfo.destination_workload_namespace=default
```

## Validate SLOs for the mirrored

The following metrics are collected by default by this experiment:

- `istio/request-count`: total number of requests sent
- `istio/error-count`: number of error responses
- `istio/error-rate`: fraction of error responses
- `istio/latency-mean`: mean of observed latency values

All latency metrics have `msec` units.

***

To validate SLOs in addition to collecting metrics, SLOs must be provided. The `--noDownload` flag reuses the Iter8 experiment `charts` folder downloaded during the previous `iter8 launch` invocation.

```shell
iter8 k launch -c slo-validation-istio \
--set reporter=destination \
--set cronjobSchedule="*/1 * * * *" \
--set providerURL=http://prometheus.istio-system:9090/api/v1/query \
--set versionInfo.destination_workload=httpbin-v2 \
--set versionInfo.destination_workload_namespace=default \
--set SLOs.istio/error-rate=0 \
--set SLOs.istio/latency-mean=100 \
--noDownload
```

In the experiment above, the following SLOs are validated.

- error rate is 0
- mean latency is under 100 msec

***

## View experiment report

--8<-- "docs/tutorials/slo-validation-istio/expreport.md"

***

## Assert experiment outcomes

--8<-- "docs/tutorials/slo-validation-istio/assert.md"
