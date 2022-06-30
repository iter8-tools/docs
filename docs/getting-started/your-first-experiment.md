---
template: main.html
---

# Your First Experiment

Perform your first [Iter8 experiment](concepts.md#iter8-experiment) by load testing an HTTP service inside Kubernetes and validating its [service-level objectives (SLOs)](concepts.md#service-level-objectives).

<p align='center'>
  <img alt-text="load-test-http" src="../images/http.png" />
</p>

???+ warning "Before you begin"
    1. Ensure that you have a Kubernetes cluster and the [`kubectl` CLI](https://kubernetes.io/docs/reference/kubectl/). You may run a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
    2. Deploy the sample HTTP service in the Kubernetes cluster.
    ```shell
    kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
    kubectl expose deploy httpbin --port=80
    ```

***

## Install Iter8 CLI
--8<-- "docs/getting-started/installbrewbins.md"

***

## Launch experiment
Launch the Iter8 experiment inside the Kubernetes cluster.

```shell
iter8 k launch \
--set "tasks={ready,http,assess}" \
--set ready.deploy=httpbin \
--set ready.service=httpbin \
--set ready.timeout=60s \
--set http.url=http://httpbin.default/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0 \
--set runner=job
```

???+ note "About this experiment"
    This experiment consists of three [tasks](../getting-started/concepts.md#tasks), namely, [ready](ready.md), [http](http.md), and [assess](assess.md). The [ready](ready.md) task checks if the `httpbin` deployment exists and is available, and the `httpbin` service exists. The [http](http.md) task sends requests to the cluster-local HTTP service whose URL is `http://httpbin.default/get`, and collects [Iter8's built-in HTTP load test metrics](built-in.md). The [assess](assess.md) task verifies if the app satisfies the specified SLOs: i) the mean latency of the service does not exceed 50 msec, and ii) there are no errors (4xx or 5xx response codes) in the responses. The [runner](../getting-started/concepts.md#runners) value specifies that the experiment should be [run using a Kubernetes job](../getting-started/concepts.md#runners).

***

## Assert experiment outcomes
Assert that the experiment completed without failures, and all SLOs are satisfied. The timeout flag below specifies a period of 120 sec for assert conditions to be satisfied.

```shell
iter8 k assert -c completed -c nofailure -c slos --timeout 120s
```

If the assert conditions are satisfied, the above command exits with code 0; else, it exits with code 1. Assertions are especially useful inside CI/CD/GitOps pipelines. Depending on the exit code of the assert command, your pipeline can branch into different actions.

***

## View experiment report
--8<-- "docs/getting-started/expreport.md"

***

## View experiment logs
Logs are useful when debugging an experiment.

```shell
iter8 k log
```

--8<-- "docs/getting-started/logs.md"

***

## Cleanup
Remove the Kubernetes resource objects created during the Iter8 experiment launch.
```shell
iter8 k delete
```

***

Congratulations! :tada: You completed your first Iter8 experiment.

???+ tip "Variations"
    1. The [http task](http.md) can be configured with [load related parameters](loadprofile.md) such as the number of requests, queries per second, or number of parallel connections.
    2. The [http task](http.md) can be configured to [send various types of content as payload](payload.md).
    3. The [assess task](assess.md) can be configured with SLOs for any of [Iter8's built-in HTTP load test metrics](built-in.md).
    4. This experiment can also be run in your [local environment](local.md) or run within a [GitHub Actions pipeline](githubactions.md).
