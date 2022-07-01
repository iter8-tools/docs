---
template: main.html
---

# Load Test gRPC with SLOs

Load test a gRPC service inside Kubernetes and validate its [SLOs](../getting-started/concepts.md#service-level-objectives).

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
    This experiment consists of three [tasks](../getting-started/concepts.md#tasks), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md). The [ready](../user-guide/tasks/ready.md) task checks if the `hello` deployment exists and is available, and the `hello` service exists. The [grpc](../user-guide/tasks/grpc.md) task calls the `helloworld.Greeter.SayHello` method of the cluster-local gRPC service hosted at `hello.default:50051`, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. The [runner](../getting-started/concepts.md#runners) value specifies that the experiment should be [run using a Kubernetes job](../getting-started/concepts.md#runners).

***

Assert experiment outcomes, view experiment report, view experiment logs, and cleanup as described in [your first experiment](../getting-started/your-first-experiment.md). As in the case of the HTTP load test, you can configure the load profile and SLOs in a flexible manner. Please refer to the [grpc task](../user-guide/tasks/grpc.md) for more details.