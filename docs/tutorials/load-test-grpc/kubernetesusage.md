---
template: main.html
tags:
- load testing
- benchmarking
- SLOs
- gRPC
- Kubernetes
---

# Benchmark and Validate Kubernetes gRPC Services

Load test, benchmark, and validate a gRPC service that is running within a Kubernetes cluster using the  [`load-test-grpc` experiment](basicusage.md). The gRPC service may be exposed outside the cluster or may be local.

<p align='center'>
  <img alt-text="load-test-http" src="../images/kubernetesusage.png" width="90%" />
</p>

***

--8<-- "docs/tutorials/load-test-grpc/usecases.md"

***

???+ warning "Before you begin"
    ```shell title="Deploy sample gRPC app in Kubernetes cluster"
    kubectl create deploy hello --image=docker.io/grpc/java-example-hostname:latest --port=50051
    kubectl expose deploy hello --port=50051      
    ```

***

## Launch experiment
```shell
iter8 k launch -c load-test-grpc \
--set host="hello:50051" \
--set call="helloworld.Greeter.SayHello" \
--set protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto" \
--set SLOs.grpc/latency/mean=50 \
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
You can set any parameter described in the [basic usage of `load-test-grpc`](basicusage.md) during `iter8 k launch` of the `load-test-grpc` experiment.