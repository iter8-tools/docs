---
template: main.html
---

# Benchmark and Validate HTTP services

The `slo-validation-istio` collects latency and error-related metrics and validates service-level objectives (SLOs) from an application using Istio.

<!-- <p align='center'>
  <img alt-text="load-test-http" src="../images/http-overview.png" width="90%" />
</p> -->

<!-- ***

--8<-- "docs/tutorials/load-test-http/usecases.md"

*** -->

???+ warning "Before you begin"
    <!-- Run the [httpbin](https://httpbin.org) sample service from a separate terminal.
    ```shell
    docker run -p 80:80 kennethreitz/httpbin
    ```
    You can also use [Podman](https://podman.io) or other alternatives to Docker in the above command. -->

    Create a new Minikube cluster
    ```shell
    minikube start --memory=16384 --cpus=4 -p slo-validation-istio
    ```

    Install Istio
    ```shell
    istioctl install
    ```

    Install Prometheus plugin
    ```shell
    kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.13/samples/addons/prometheus.yaml
    ```

    Run the httpbin sample service
    ```shell
    istioctl kube-inject -f https://raw.githubusercontent.com/Alan-Cha/docs/slo-validation-istio/docs/tutorials/slo-validation-istio/samples/httpbin-v1-deployment.yaml | kubectl apply -f -
    istioctl kube-inject -f https://raw.githubusercontent.com/Alan-Cha/docs/slo-validation-istio/docs/tutorials/slo-validation-istio/samples/httpbin-service.yaml | kubectl apply -f -
    ```

    In a separate terminal window, run a Minikube tunnel 
    ```shell
    minikube tunnel -p slo-validation-istio
    ```

    In another separate terminal window, open the Prometheus API and allow local access to the Prometheus database
    ```shell
    istioctl dashboard prometheus
    ```

    ***

    To send traffic, first run the sleep service
    ```shell
    istioctl kube-inject -f https://raw.githubusercontent.com/Alan-Cha/docs/slo-validation-istio/docs/tutorials/slo-validation-istio/samples/sleep-deployment.yaml | kubectl apply -f -
    ```

    Then, make some requests
    ```shell
    export SLEEP_POD=$(kubectl get pod -l app=sleep -o jsonpath={.items..metadata.name})
    kubectl exec "${SLEEP_POD}" -c sleep -- curl -sS http://httpbin:8000/headers
    ```
***

## Collect metrics
To collect metrics from a service, following must be set:

- `endpoint`, the Prometheus database endpoint
- `destination_workload`, the name of the service
- `destination_workload_namespace`, the namespace of the service 
- `startingTime`, the starting time from which to collect metrics (note the date format)

```shell
iter8 launch -c slo-validation-istio \
--set endpoint=http://localhost:9090 \
--set destination_workload=httpbin-v1 \
--set destination_workload_namespace=default \
--set startingTime="Feb 4\, 2014 at 6:05pm (PST)"
```

## Validate SLOs

The following metrics are collected by default by this experiment:

- `http/request-count`: total number of requests sent
- `http/error-count`: number of error responses
- `http/error-rate`: fraction of error responses
- `http/latency-mean`: mean of observed latency values

All latency metrics have `msec` units.

***

Launch the following experiment. The `--noDownload` flag reuses the Iter8 experiment `charts` folder downloaded during the previous `iter8 launch` invocation.

To validate SLOs in addition to collecting metrics, SLOs must also be provided.

```shell
iter8 launch -c slo-validation-istio \
--noDownload \
--set endpoint=http://localhost:9090 \
--set destination_workload=httpbin-v1 \
--set destination_workload_namespace=default \
--set startingTime="Feb 4\, 2014 at 6:05pm (PST)"
--set SLOs.http/error-rate=0 \
--set SLOs.http/latency-mean=50
```

In the experiment above, the following SLOs are validated.

- error rate is 0
- mean latency is under 50 msec

***

## View experiment report

--8<-- "docs/tutorials/load-test-http/expreport.md"

***

## Assert experiment outcomes

--8<-- "docs/tutorials/load-test-http/assert.md"
