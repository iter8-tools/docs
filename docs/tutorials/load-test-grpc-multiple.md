---
template: main.html
---

# Load Test Multiple gRPC methods

[Load Test gRPC with SLOs](./load-test-grpc.md) describes how to load test a single method from a gRPC service inside Kubernetes. This tutorial expands on the previous tutorial and describes how to load test multiple endpoints from an HTTP service.

<p align='center'>
  <img alt-text="load-test-grpc" src="../images/grpc.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../getting-started/your-first-experiment.md). Understand the main [concepts](../getting-started/concepts.md) behind Iter8 experiments.
    2. Deploy the sample gRPC service in the Kubernetes cluster.
    ```shell
    kubectl create deployment grpc-example  --image=golang --port=50051 -- bash -c "git clone -b v1.52.0 --depth 1 https://github.com/grpc/grpc-go; cd grpc-go/examples/route_guide; go run server/server.go"
    kubectl expose deployment grpc-example --port=50051
    ```

***

## Launch experiment

```shell
helm install autox autox --repo https://iter8-tools.github.io/hub/ --version 0.1.6 \
--set 'groups.grpc-example.trigger.name=grpc-example' \
--set 'groups.grpc-example.trigger.namespace=default' \
--set 'groups.grpc-example.trigger.group=apps' \
--set 'groups.grpc-example.trigger.version=v1' \
--set 'groups.grpc-example.trigger.resource=deployments' \
--set 'groups.grpc-example.specs.iter8.name=iter8' \
--set 'groups.grpc-example.specs.iter8.values.tasks={ready,http,assess}' \
--set 'groups.grpc-example.specs.iter8.values.ready.deploy=grpc-example' \
--set 'groups.grpc-example.specs.iter8.values.ready.service=grpc-example' \
--set 'groups.grpc-example.specs.iter8.values.ready.timeout=60s' \
--set 'groups.grpc-example.specs.iter8.values.grpc.host=grpc-example.default:81' \
--set 'groups.grpc-example.specs.iter8.values.grpc.endpoints.getFeature.call=routeguide.RouteGuide.GetFeature' \
--set 'groups.grpc-example.specs.iter8.values.grpc.endpoints.getFeature.dataURL=...' \
--set 'groups.grpc-example.specs.iter8.values.grpc.endpoints.listFeatures.call=routeguide.RouteGuide.ListFeatures' \
--set 'groups.grpc-example.specs.iter8.values.grpc.endpoints.listFeatures.dataURL=...' \
--set 'groups.grpc-example.specs.iter8.values.grpc.protoURL=https://raw.githubusercontent.com/grpc/grpc-go/master/examples/route_guide/routeguide/route_guide.proto' \
--set 'groups.grpc-example.specs.iter8.values.assess.SLOs.upper.grpc/getFeature/error-count=0' \
--set 'groups.grpc-example.specs.iter8.values.assess.SLOs.upper.grpc/getFeature/latency-mean=50' \
--set 'groups.grpc-example.specs.iter8.values.assess.SLOs.upper.grpc/listFeatures/error-count=0' \
--set 'groups.grpc-example.specs.iter8.values.assess.SLOs.upper.grpc/listFeatures/latency-mean=100' \
--set 'groups.grpc-example.specs.iter8.version=0.13.0' \
--set 'groups.grpc-example.specs.iter8.values.runner=job'
```

??? note "About this experiment"
    This experiment consists of three [tasks](../getting-started/concepts.md#iter8-experiment), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md). 
    
    The [ready](../user-guide/tasks/ready.md) task checks if the `grpc-example` deployment exists and is available, and the `grpc-example` service exists. 

    The [grpc](../user-guide/tasks/grpc.md) task sends call requests to two methods of the cluster-local gRPC service, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). The two methods are `routeguide.RouteGuide.GetFeature` and `routeguide.RouteGuide.ListFeatures`. Note that each method also has their own `dataURL`.

    The [assess](../user-guide/tasks/assess.md) task verifies if the two methods satisfy the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. 

    The [assess](../user-guide/tasks/assess.md) task verifies if each method satisfy their respective error count and mean latency SLOs. Both methods must have an error count of 0 but the `getFeature` and `listFeatures` methods are allowed a maximum mean latency of 50 and 100 msecs, respectively.
    
    This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

??? note "Some variations and extensions of this experiment"
    1. The [grpc task](../user-guide/tasks/grpc.md) can be configured with load related parameters such as the total number of requests, requests per second, or number of concurrent connections.
    2. The [grpc task](../user-guide/tasks/grpc.md) can be configured to JSON or binary data as payload. You can use this task to test unary or streaming gRPC methods.
    3. The [assess task](../user-guide/tasks/assess.md) can be configured with SLOs for any of [Iter8's built-in grpc load test metrics](../user-guide/tasks/grpc.md#metrics).   

***

Assert experiment outcomes, view experiment report, view experiment logs, and cleanup as described in [your first experiment](../getting-started/your-first-experiment.md).

***

## Cleanup
Remove the Iter8 experiment and the sample app from the Kubernetes cluster and the local Iter8 `charts` folder.

```shell
iter8 k delete
kubectl delete svc/grpc-example
kubectl delete deploy/grpc-example
```