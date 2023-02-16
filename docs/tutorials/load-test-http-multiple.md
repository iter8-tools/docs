---
template: main.html
---

# Load Test Multiple HTTP endpoints

[Your first experiment](../getting-started/your-first-experiment.md) describes how to load test a single endpoint from an HTTP service inside Kubernetes. This tutorial expands on the previous tutorial and describes how to load test multiple endpoints from an HTTP service.

<p align='center'>
  <img alt-text="load-test-http" src="../../getting-started/images/http.png" />
</p>

???+ warning "Before you begin"
    1. Ensure that you have a Kubernetes cluster and the [`kubectl` CLI](https://kubernetes.io/docs/reference/kubectl/). You can create a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
    2. Deploy the sample HTTP service in the Kubernetes cluster.
    ```shell
    kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
    kubectl expose deploy httpbin --port=80
    ```

***

## Launch experiment
Launch the Iter8 experiment inside the Kubernetes cluster.

```bash
iter8 k launch \
--set "tasks={ready,http,assess}" \
--set ready.deploy=httpbin \
--set ready.service=httpbin \
--set ready.timeout=60s \
--set http.endpoints.get.url=http://httpbin.default/get \
--set http.endpoints.getAnything.url=http://httpbin.default/anything \
--set http.endpoints.post.url=http://httpbin.default/post \
--set http.endpoints.post.payloadStr=hello \
--set http.endpoints.post.qps=10 \
--set assess.SLOs.upper.http/get/error-count=0 \
--set assess.SLOs.upper.http/get/latency-mean=50 \
--set assess.SLOs.upper.http/getAnything/error-count=0 \
--set assess.SLOs.upper.http/getAnything/latency-mean=100 \
--set assess.SLOs.upper.http/post/error-count=0 \
--set assess.SLOs.upper.http/post/latency-mean=150 \
--set runner=job
```

??? note "About this experiment"
    This experiment consists of three [tasks](../getting-started/concepts.md#iter8-experiment), namely, [ready](../user-guide/tasks/ready.md), [http](../user-guide/tasks/http.md), and [assess](../user-guide/tasks/assess.md). 
    
    The [ready](../user-guide/tasks/ready.md) task checks if the `httpbin` deployment exists and is available, and the `httpbin` service exists. 
    
    The [http](../user-guide/tasks/http.md) task sends requests to three endpoints from the cluster-local HTTP service, and collects [Iter8's built-in HTTP load test metrics](../user-guide/tasks/http.md#metrics). The three endpoints are `http://httpbin.default/get`, `http://httpbin.default/anything`, and `http://httpbin.default/post`. The last endpoint also has a payload string `hello`. Furthermore, by default, the endpoints are queried at a rate of 8 qps (queries-per-second), but the last endpoint will be queried at 10 qps.
    
    The [assess](../user-guide/tasks/assess.md) task verifies if each endpoint satisfies their respective error count and mean latency SLOs. All three must have an error count of 0 but the `get`, `getAnything`, and `post` endpoints are allowed a maximum mean latency of 50, 75, and 100 msecs, respectively.
    
    This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

***

Assert experiment outcomes, view experiment report, view experiment logs, and cleanup as described in [your first experiment](../getting-started/your-first-experiment.md).

***

## Cleanup
Remove the Iter8 experiment and the sample app from the Kubernetes cluster and the local Iter8 `charts` folder.

```shell
iter8 k delete
kubectl delete svc/httpbin
kubectl delete deploy/httpbin
```