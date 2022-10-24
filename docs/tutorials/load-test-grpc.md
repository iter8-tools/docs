---
template: main.html
---

# Load Test gRPC with SLOs

Load test a Kubernetes gRPC service and validate its [service-level objectives (SLOs)](../getting-started/concepts.md#service-level-objectives). This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments).

<p align='center'>
  <img alt-text="load-test-grpc" src="../images/grpc.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../getting-started/your-first-experiment.md). Understand the main [concepts](../getting-started/concepts.md) behind Iter8 experiments.
    2. Deploy the sample gRPC service in the Kubernetes cluster.
    ```shell
    kubectl create deploy hello --image=docker.io/grpc/java-example-hostname:latest --port=50051
    kubectl expose deploy hello --port=50051
    ```

***

## Launch experiment

```shell
iter8 k launch \
--set "tasks={ready,grpc,assess}" \
--set ready.deploy=hello \
--set ready.service=hello \
--set ready.timeout=60s \
--set grpc.host="hello.default:50051" \
--set grpc.call="helloworld.Greeter.SayHello" \
--set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto" \
--set assess.SLOs.upper.grpc/error-rate=0 \
--set assess.SLOs.upper.grpc/latency/mean=200 \
--set assess.SLOs.upper.grpc/latency/p'97\.5'=800 \
--set runner=job
```

??? note "About this experiment"
    This experiment consists of three [tasks](../getting-started/concepts.md#iter8-experiment), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md). 
    
    The [ready](../user-guide/tasks/ready.md) task checks if the `hello` deployment exists and is available, and the `hello` service exists. 
    
    The [grpc](../user-guide/tasks/grpc.md) task sends call requests to the `helloworld.Greeter.SayHello` method of the cluster-local gRPC service with host address `hello.default:50051`, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). 
    
    The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. 
    
    This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

??? note "Some variations and extensions of this experiment"
    1. The [grpc task](../user-guide/tasks/grpc.md) can be configured with load related parameters such as the total number of requests, requests per second, or number of concurrent connections.
    2. The [grpc task](../user-guide/tasks/grpc.md) can be configured to JSON or binary data as payload. You can use this task to test unary or streaming gRPC methods.
    3. The [assess task](../user-guide/tasks/assess.md) can be configured with SLOs for any of [Iter8's built-in grpc load test metrics](../user-guide/tasks/grpc.md#metrics).
    4. This experiment can also be run in your [local environment](../tutorials/integrations/local.md) or run within a [GitHub Actions pipeline](../tutorials/integrations/ghactions.md).    

***

Assert experiment outcomes, view experiment report, view experiment logs, and cleanup as described in [your first experiment](../getting-started/your-first-experiment.md).

***

## Cleanup
Remove the Iter8 experiment and the sample app from the Kubernetes cluster and the local Iter8 `charts` folder.
```shell
iter8 k delete
kubectl delete svc/hello
kubectl delete deploy/hello
rm -rf charts
```