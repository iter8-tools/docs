---
template: main.html
---

# Load Test Multiple HTTP endpoints

[Your first experiment](../getting-started/your-first-experiment.md) describes how to load test a single endpoint from an HTTP service inside Kubernetes. This tutorial expands on the previous tutorial and describes how to load test multiple endpoints from an HTTP service.

![load-test-http](../getting-started/images/kubernetes.png)

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
--set "tasks={ready,http}" \
--set ready.deploy=httpbin \
--set ready.service=httpbin \
--set ready.timeout=60s \
--set http.endpoints.get.url=http://httpbin.default/get \
--set http.endpoints.getAnything.url=http://httpbin.default/anything \
--set http.endpoints.post.url=http://httpbin.default/post \
--set http.endpoints.post.payloadStr=hello
```

??? note "About this experiment"
    This experiment consists of two [tasks](../getting-started/concepts.md#design), namely, [ready](../user-guide/tasks/ready.md) and [http](../user-guide/tasks/http.md).
    
    The [ready](../user-guide/tasks/ready.md) task checks if the `httpbin` deployment exists and is available, and the `httpbin` service exists. 
    
    The [http](../user-guide/tasks/http.md) task sends requests to three endpoints from the cluster-local HTTP service, and collects [Iter8's built-in HTTP load test metrics](../user-guide/tasks/http.md#metrics). The three endpoints are `http://httpbin.default/get`, `http://httpbin.default/anything`, and `http://httpbin.default/post`. The last endpoint also has a payload string `hello`.

***

View the experiment results by using the Iter8 Grafana dashboard, as described in [your first experiment](../getting-started/your-first-experiment.md).

***

## Cleanup
Remove the Iter8 experiment and the sample app from the Kubernetes cluster and the local Iter8 `charts` folder.

```shell
iter8 k delete
kubectl delete svc/httpbin
kubectl delete deploy/httpbin
```