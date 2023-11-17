---
template: main.html
---

# Load test HTTP endpoint

Run your first performance test by load testing a Kubernetes HTTP service and visualizing the performance with an Iter8 Grafana dashboard.

![Load test HTTP](images/kubernetesusage.png)

???+ warning "Before you begin"
    1. Ensure that you have a Kubernetes cluster and the [`kubectl`](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs. If using a local cluster (for example, [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/)), we recommend providing the cluster with at least 16GB of memory.
    2. Deploy the sample HTTP service in the Kubernetes cluster.
    ```shell
    kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
    kubectl expose deploy httpbin --port=80
    ```
    3. Have Grafana available. For example, Grafana can be installed on your cluster as follows:
    ```shell
    kubectl create deploy grafana --image=grafana/grafana
    kubectl expose deploy grafana --port=3000
    ```

***

## Install the Iter8 controller

--8<-- "docs/getting-started/install.md"

## Launch performance test

=== "GET example"
    ```shell
    helm upgrade --install \
    --repo https://iter8-tools.github.io/iter8 --version 0.18 httpbin-test iter8 \
    --set "tasks={ready,http}" \
    --set ready.deploy=httpbin \
    --set ready.service=httpbin \
    --set ready.timeout=60s \
    --set http.url=http://httpbin.default/get
    ```

=== "POST example"
    ```shell
    helm upgrade --install \
    --repo https://iter8-tools.github.io/iter8 --version 0.18 httpbin-test iter8 \
    --set "tasks={ready,http}" \
    --set ready.deploy=httpbin \
    --set ready.service=httpbin \
    --set ready.timeout=60s \
    --set http.url=http://httpbin.default/post \
    --set http.payloadStr=hello
    ```

??? note "About this performance test"
    This performance test consists of two [tasks](concepts.md#design), namely, [ready](../user-guide/tasks/ready.md) and [http](../user-guide/tasks/http.md). 
    
    The [ready](../user-guide/tasks/ready.md) task checks if the `httpbin` deployment exists and is available, and the `httpbin` service exists. 
    
    The [http](../user-guide/tasks/http.md) task sends requests to the cluster-local HTTP service using the specified `url`, and collects [Iter8's built-in HTTP load test metrics](../user-guide/tasks/http.md#metrics). This tasks supports both GET and POST requests, and for POST requests, a payload can be provided by using either `payloadStr` or `payloadURL`.

## View results using Grafana
Inspect the metrics using Grafana. If Grafana is deployed to your cluster, port-forward requests as follows:

```shell
kubectl port-forward service/grafana 3000:3000
```

Open Grafana in a browser by going to [http://localhost:3000](http://localhost:3000) and login. The default username/password are `admin`/`admin`.

[Add a JSON API data source](http://localhost:3000/connections/datasources/marcusolsson-json-datasource) `httpbin-test` with the following parameters:

* URL: `http://iter8.default:8080/httpDashboard` 
* Query string: `namespace=default&test=httpbin-test`

[Create a new dashboard](http://localhost:3000/dashboards) by *import*. Paste the contents of the [`http` Grafana dashboard](https://raw.githubusercontent.com/iter8-tools/iter8/v0.18.3/grafana/http.json) into the text box and *load* it. Associate it with the JSON API data source defined above.

The Iter8 dashboard will look like the following:

![`http` Iter8 dashboard](../user-guide/tasks/images/httpdashboard.png)

## View logs
Logs are useful for debugging. To see the test logs:

```shell
kubectl logs -l iter8.tools/test=httpbin-test
```

--8<-- "docs/getting-started/logs.md"

***

## Cleanup
Remove the performance test and the sample app from the Kubernetes cluster.
```shell
helm delete httpbin-test
kubectl delete svc/httpbin
kubectl delete deploy/httpbin
```

### Uninstall the Iter8 controller

--8<-- "docs/getting-started/uninstall.md"

If you installed Grafana, you can delete it as follows:

```shell
kubectl delete svc/grafana deploy/grafana
```

***

Congratulations! :tada: You completed your first performance test with Iter8.

***

??? note "Some variations and extensions of this performance test"
    1. The [http task](../user-guide/tasks/http.md) can be configured with load related parameters such as the number of requests, queries per second, or number of parallel connections.
    2. The [http task](../user-guide/tasks/http.md) can be configured to send various types of content as payload.
