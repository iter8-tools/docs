---
template: main.html
tags:
- load testing
- benchmarking
- SLOs
- HTTP
- Kubernetes
---

# Benchmark and Validate Kubernetes HTTP Services

Load test, benchmark, and validate an HTTP service that is running within a Kubernetes cluster using the  [`load-test-http` experiment](basicusage.md). The HTTP service may be exposed outside the cluster or may be local.

<p align='center'>
  <img alt-text="load-test-http" src="../images/kubernetesusage.png" width="90%" />
</p>

***

--8<-- "docs/tutorials/load-test-http/usecases.md"

***

???+ warning "Before you begin"
    ```shell title="Deploy sample HTTP app in Kubernetes cluster"
    kubectl create deploy hello --image=kennethreitz/httpbin --port=80
    kubectl expose deploy hello --port=80 
    ```

***

## Launch experiment
```shell
iter8 k launch -c load-test-http \
--set url=http://hello.default \
--set SLOs.http/latency/mean=50 \
--set readiness.service="hello" \
--set readiness.deploy="hello"
```

Also refer to [readiness check](../../user-guide/topics/readiness.md).

***

## Assert outcomes
```shell
iter8 k assert -c completed -c nofailure -c slos --timeout=30s
```

***

## View report
=== "Text"
    ```shell
    iter8 k report
    ```

=== "HTML"
    ```shell
    iter8 k report -o html > report.html # view in a browser
    ```

***

## Get experiment logs
```shell
iter8 k logs
```

Refer to the log level flag during `iter8 k launch`.

***

## Set values
You can set any parameter described in the [basic usage of `load-test-http`](basicusage.md) during `iter8 k launch` of the `load-test-http` experiment.

***

## Cleanup
Cleanup all Kubernetes resources created by the Iter8 experiment in the cluster.

```shell
iter8 k cleanup
```
