---
template: main.html
---

# Load test gRPC with SLOs

Load test a gRPC service inside Kubernetes and validate its [SLOs](slos.md).

<p align='center'>
  <img alt-text="load-test-grpc" src="../images/grpc.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../getting-started/your-first-experiment.md).
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

???+ note "About this experiment"
    This experiment consists of three [tasks](tasks.md), namely, [ready](ready.md), [grpc](grpc.md), and [assess](assess.md). The [ready](ready.md) task checks if the `hello` deployment exists and is available, and the `hello` service exists. The [grpc](grpc.md) task calls the `helloworld.Greeter.SayHello` method of the cluster-local gRPC service hosted at `hello.default:50051`, and collects [Iter8's built-in gRPC load test metrics](built-in.md). The [assess](assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. The [runner](runner.md) value specifies that the experiment should be [run using a Kubernetes job](runner.md).

***

Assert experiment outcomes, view experiment report, view experiment logs, and cleanup as described in [your first experiment](../getting-started/your-first-experiment.md). As in the case of the HTTP load test, you can configure the load profile and SLOs in a flexible manner. Please refer to the [grpc task](grpc.md) for more details.