---
template: main.html
---

# Load test a KServe model (via HTTP)

This tutorial shows how easy it is to run a load test for KServe when using HTTP to make requests. We use a sklearn model to demonstrate. The same approach works for any model type. 

???+ warning "Before you begin"
    1. Ensure that you have the [`kubectl`](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs installed.
    2. Have access to a cluster running [KServe](https://kserve.github.io/website). You can create a [KServe Quickstart](https://kserve.github.io/website/0.11/get_started/#before-you-begin) environment as follows:
    ```shell
    curl -s "https://raw.githubusercontent.com/kserve/kserve/release-0.11/hack/quick_install.sh" | bash
    ```
    If using a local cluster (for example, [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/)), we recommend providing the cluster with at least 16GB of memory.
    4. Have Grafana available. For example, Grafana can be installed on your cluster as follows:
    ```shell
    kubectl create deploy grafana --image=grafana/grafana
    kubectl expose deploy grafana --port=3000
    ```

## Install the Iter8 controller

--8<-- "docs/getting-started/install.md"

## Deploy an InferenceService

Create an InferenceService which exposes an HTTP port. The following serves the sklearn [irisv2 model](https://kserve.github.io/website/0.10/modelserving/v1beta1/sklearn/v2/#deploy-with-inferenceservice):

```shell
cat <<EOF | kubectl apply -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-irisv2"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      runtime: kserve-mlserver
      storageUri: "gs://seldon-models/sklearn/mms/lr_model"
EOF
```

## Launch performance test

```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 0.18 model-test iter8 \
--set "tasks={ready,http}" \
--set ready.isvc=sklearn-irisv2 \
--set ready.timeout=180s \
--set http.url=http://sklearn-irisv2.default.svc.cluster.local/v2/models/sklearn-irisv2/infer \
--set http.payloadURL=https://gist.githubusercontent.com/kalantar/d2dd03e8ebff2c57c3cfa992b44a54ad/raw/97a0480d0dfb1deef56af73a0dd31c80dc9b71f4/sklearn-irisv2-input.json \
--set http.contentType="application/json"
```

??? note "About this performance test"
    This performance test consists of two [tasks](../../../getting-started/concepts.md#design), namely, [ready](../../../user-guide/tasks/ready.md) and [http](../../../user-guide/tasks/http.md). 
    
    The [ready](../../../user-guide/tasks/ready.md) task checks if the `sklearn-irisv2` InferenceService exists and is `Ready`. 

    The [http](../../../user-guide/tasks/http.md) task sends requests to the cluster-local HTTP service whose URL exposed by the InferenceService, `http://sklearn-irisv2.default.svc.cluster.local/v2/models/sklearn-irisv2/infer`, and collects [Iter8's built-in HTTP load test metrics](../../../user-guide/tasks/http.md#metrics).

## View results using Grafana
Inspect the metrics using Grafana. If Grafana is deployed to your cluster, port-forward requests as follows:

```shell
kubectl port-forward service/grafana 3000:3000
```

Open Grafana in a browser by going to [http://localhost:3000](http://localhost:3000) and login. The default username/password are `admin`/`admin`.

[Add a JSON API data source](http://localhost:3000/connections/datasources/marcusolsson-json-datasource) `model-test` with the following parameters:

* URL: `http://iter8.default:8080/httpDashboard` 
* Query string: `namespace=default&test=model-test`

[Create a new dashboard](http://localhost:3000/dashboards) by *import*. Paste the contents of the [`http` Grafana dashboard](https://raw.githubusercontent.com/iter8-tools/iter8/v0.18.3/grafana/http.json) into the text box and *load* it. Associate it with the JSON API data source defined above.

The Iter8 dashboard will look like the following:

![`http` Iter8 dashboard](../../../user-guide/tasks/images/httpdashboard.png)

## Cleanup

```shell
helm delete model-test
kubectl delete inferenceservice sklearn-irisv2
```

### Uninstall the Iter8 controller

--8<-- "docs/getting-started/uninstall.md"

If you installed Grafana, you can delete it as follows:

```shell
kubectl delete svc/grafana deploy/grafana
```

??? note "Some variations and extensions of this performance test"
    1. The [http task](../../../user-guide/tasks/http.md) can be configured with load related parameters such as the number of requests, queries per second, or number of parallel connections.
    2. The [http task](../../../user-guide/tasks/http.md) can be configured to send various types of content as payload.
