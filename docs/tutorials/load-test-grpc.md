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
    kubectl create deployment routeguide  --image=golang --port=50051 -- bash -c "git clone -b v1.52.0 --depth 1 https://github.com/grpc/grpc-go; cd grpc-go/examples/route_guide; go run server/server.go"
    kubectl expose deployment routeguide --port=50051
    ```

***

## Launch experiment

=== "Unary example"
    ```shell
    iter8 k launch \
    --set "tasks={ready,grpc,assess}" \
    --set ready.deploy=routeguide \
    --set ready.service=routeguide \
    --set ready.timeout=60s \
    --set grpc.host="routeguide.default:50051" \
    --set grpc.call="routeguide.RouteGuide.GetFeature" \
    --set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/route_guide/routeguide/route_guide.proto" \
    --set grpc.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/unary.json" \
    --set assess.SLOs.upper.grpc/error-rate=0 \
    --set assess.SLOs.upper.grpc/latency/mean=200 \
    --set assess.SLOs.upper.grpc/latency/p'97\.5'=800 \
    --set runner=job
    ```

    ??? note "About this experiment"
        This experiment consists of three [tasks](../getting-started/concepts.md#iter8-experiment), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md). 
        
        The [ready](../user-guide/tasks/ready.md) task checks if the `routeguide` deployment exists and is available, and the `routeguide` service exists. 
        
        The [grpc](../user-guide/tasks/grpc.md) task sends call requests to the `routeguide.RouteGuide.GetFeature` method of the cluster-local gRPC service with host address `routeguide.default:50051`, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). 
        
        The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. 
        
        This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

=== "Server streaming example"
    ```shell
    iter8 k launch \
    --set "tasks={ready,grpc,assess}" \
    --set ready.deploy=routeguide \
    --set ready.service=routeguide \
    --set ready.timeout=60s \
    --set grpc.host="routeguide.default:50051" \
    --set grpc.call="routeguide.RouteGuide.ListFeatures" \
    --set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/route_guide/routeguide/route_guide.proto" \
    --set grpc.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/server.json" \
    --set assess.SLOs.upper.grpc/error-rate=0 \
    --set assess.SLOs.upper.grpc/latency/mean=200 \
    --set assess.SLOs.upper.grpc/latency/p'97\.5'=800 \
    --set runner=job
    ```

    ??? note "About this experiment"
        This experiment consists of three [tasks](../getting-started/concepts.md#iter8-experiment), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md). 
        
        The [ready](../user-guide/tasks/ready.md) task checks if the `routeguide` deployment exists and is available, and the `routeguide` service exists. 
        
        The [grpc](../user-guide/tasks/grpc.md) task sends call requests to the `routeguide.RouteGuide.ListFeatures` method of the cluster-local gRPC service with host address `routeguide.default:50051`, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). 
        
        The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. 
        
        This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

=== "Client streaming example"
    ```shell
    iter8 k launch \
    --set "tasks={ready,grpc,assess}" \
    --set ready.deploy=routeguide \
    --set ready.service=routeguide \
    --set ready.timeout=60s \
    --set grpc.host="routeguide.default:50051" \
    --set grpc.call="routeguide.RouteGuide.RecordRoute" \
    --set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/route_guide/routeguide/route_guide.proto" \
    --set grpc.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/client.json" \
    --set assess.SLOs.upper.grpc/error-rate=0 \
    --set assess.SLOs.upper.grpc/latency/mean=200 \
    --set assess.SLOs.upper.grpc/latency/p'97\.5'=800 \
    --set runner=job
    ```

    ??? note "About this experiment"
        This experiment consists of three [tasks](../getting-started/concepts.md#iter8-experiment), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md). 
        
        The [ready](../user-guide/tasks/ready.md) task checks if the `routeguide` deployment exists and is available, and the `routeguide` service exists. 
        
        The [grpc](../user-guide/tasks/grpc.md) task sends call requests to the `routeguide.RouteGuide.RecordRoute` method of the cluster-local gRPC service with host address `routeguide.default:50051`, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). 
        
        The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. 
        
        This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

=== "Bi-directional example"
    ```shell
    iter8 k launch \
    --set "tasks={ready,grpc,assess}" \
    --set ready.deploy=routeguide \
    --set ready.service=routeguide \
    --set ready.timeout=60s \
    --set grpc.host="routeguide.default:50051" \
    --set grpc.call="routeguide.RouteGuide.RouteChat" \
    --set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/route_guide/routeguide/route_guide.proto" \
    --set grpc.dataURL="https://raw.githubusercontent.com/iter8-tools/docs/main/samples/grpc-payload/bidirectional.json" \
    --set assess.SLOs.upper.grpc/error-rate=0 \
    --set assess.SLOs.upper.grpc/latency/mean=200 \
    --set assess.SLOs.upper.grpc/latency/p'97\.5'=800 \
    --set runner=job
    ```

    ??? note "About this experiment"
        This experiment consists of three [tasks](../getting-started/concepts.md#iter8-experiment), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md). 
        
        The [ready](../user-guide/tasks/ready.md) task checks if the `routeguide` deployment exists and is available, and the `routeguide` service exists. 
        
        The [grpc](../user-guide/tasks/grpc.md) task sends call requests to the `routeguide.RouteGuide.RouteChat` method of the cluster-local gRPC service with host address `routeguide.default:50051`, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). 
        
        The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. 
        
        This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

??? note "Some variations and extensions of this experiment"
    1. The [grpc task](../user-guide/tasks/grpc.md) can be configured with load related parameters such as the total number of requests, requests per second, or number of concurrent connections.
    2. The [assess task](../user-guide/tasks/assess.md) can be configured with SLOs for any of [Iter8's built-in grpc load test metrics](../user-guide/tasks/grpc.md#metrics).   

***

Assert experiment outcomes, view experiment report, view experiment logs, and cleanup as described in [your first experiment](../getting-started/your-first-experiment.md).

***

## Cleanup
Remove the Iter8 experiment and the sample app from the Kubernetes cluster and the local Iter8 `charts` folder.
```shell
iter8 k delete
kubectl delete svc/routeguide
kubectl delete deploy/routeguide
```