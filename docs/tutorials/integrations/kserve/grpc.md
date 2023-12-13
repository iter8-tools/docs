---
template: main.html
---

# Load test a KServe model (via gRPC)

This tutorial shows how easy it is to run a load test for KServe when using gRPC to make requests. We use a sklearn model to demonstrate. The same approach works for any model type. 

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

Create an InferenceService which exposes a gRPC port. The following serves the sklearn [irisv2 model](https://kserve.github.io/website/0.10/modelserving/v1beta1/sklearn/v2/#deploy-with-inferenceservice):

```shell
cat <<EOF | kubectl create -f -
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
      protocolVersion: v2
      storageUri: "gs://seldon-models/sklearn/mms/lr_model"
      ports:
      - containerPort: 9000
        name: h2c
        protocol: TCP
EOF
```

## Launch performance test

```shell
GRPC_HOST=$(kubectl get isvc sklearn-irisv2 -o jsonpath='{.status.components.predictor.address.url}' | sed 's#.*//##')
GRPC_PORT=80
```

```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 0.18 model-test iter8 \
--set "tasks={ready,grpc}" \
--set ready.isvc=sklearn-irisv2 \
--set ready.timeout=180s \
--set grpc.protoURL=https://raw.githubusercontent.com/kserve/kserve/master/docs/predict-api/v2/grpc_predict_v2.proto \
--set grpc.host=${GRPC_HOST}:${GRPC_PORT} \
--set grpc.call=inference.GRPCInferenceService.ModelInfer \
--set grpc.dataURL=https://gist.githubusercontent.com/kalantar/6e9eaa03cad8f4e86b20eeb712efef45/raw/56496ed5fa9078b8c9cdad590d275ab93beaaee4/sklearn-irisv2-input-grpc.json
```

??? note "About this performance test"
    This performance test consists of two [tasks](../../../getting-started/concepts.md#design), namely, [ready](../../../user-guide/performance/tasks/ready.md) and [grpc](../../../user-guide/performance/tasks/grpc.md). 
    
    The [ready](../../../user-guide/performance/tasks/ready.md) task checks if the `sklearn-irisv2` InferenceService exists and is `Ready`. 

    The [grpc](../../../user-guide/performance/tasks/grpc.md) task sends call requests to the `inference.GRPCInferenceService.ModelInfer` method of the cluster-local gRPC service with host address `${GRPC_HOST}:${GRPC_PORT}`, and collects Iter8's built-in gRPC load test metrics.

## View results using Grafana
Inspect the metrics using Grafana. If Grafana is deployed to your cluster, port-forward requests as follows:

```shell
kubectl port-forward service/grafana 3000:3000
```

Open Grafana in a browser by going to [http://localhost:3000](http://localhost:3000) and login. The default username/password are `admin`/`admin`. 

[Add a JSON API data source](http://localhost:3000/connections/datasources/marcusolsson-json-datasource) `model-test` with the following parameters:

* URL: `http://iter8.default:8080/grpcDashboard` 
* Query string: `namespace=default&test=model-test`

[Create a new dashboard](http://localhost:3000/dashboards) by *import*. Paste the contents of the [`grpc` Grafana dashboard](https://raw.githubusercontent.com/iter8-tools/iter8/v0.18.3/grafana/grpc.json) into the text box and *load* it. Associate it with the JSON API data source defined above.

The Iter8 dashboard will look like the following:

![`grpc` Iter8 dashboard](../../../user-guide/performance/tasks/images/grpcdashboard.png)

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
    1. The [grpc task](../../../user-guide/performance/tasks/grpc.md) can be configured with load related parameters such as the number of requests, requests per second, or number of concurrent connections.