---
template: main.html
---

# Load Test a KServe Model (via gRPC)

This tutorial shows how easy it is to run a load test for KServe when using gRPC to make requests. We use a sklearn model to demonstrate. The same approach works for any model type. 

???+ "Before you begin"
    1. Try [your first experiment](../../../getting-started/your-first-experiment.md). Understand the main [concepts](../../../getting-started/concepts.md) behind Iter8 experiments.
    2. Ensure that you have the [kubectl](https://kubernetes.io/docs/reference/kubectl/) CLI.
    3. Have access to a cluster running [KServe](https://kserve.github.io/website). You can create a [KServe Quickstart](https://kserve.github.io/website/0.10/get_started/#before-you-begin) environment as follows:
    ```shell
    curl -s "https://raw.githubusercontent.com/kserve/kserve/release-0.10/hack/quick_install.sh" | bash
    ```
    4. Have Grafana available. For example, Grafana can be installed on your cluster as follows:
    ```shell
    kubectl create deploy grafana --image=grafana/grafana
    kubectl expose deploy grafana --port=3000
    ```

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

***

## Launch Experiment

Launch the Iter8 experiment inside the Kubernetes cluster:

```shell
GRPC_HOST=$(kubectl get isvc sklearn-irisv2 -o jsonpath='{.status.components.predictor.address.url}' | sed 's#.*//##')
GRPC_PORT=80
```

```shell
iter8 k launch \
--set "tasks={ready,grpc}" \
--set ready.isvc=sklearn-irisv2 \
--set ready.timeout=180s \
--set grpc.protoURL=https://raw.githubusercontent.com/kserve/kserve/master/docs/predict-api/v2/grpc_predict_v2.proto \
--set grpc.host=${GRPC_HOST}:${GRPC_PORT} \
--set grpc.call=inference.GRPCInferenceService.ModelInfer \
--set grpc.dataURL=https://gist.githubusercontent.com/kalantar/6e9eaa03cad8f4e86b20eeb712efef45/raw/56496ed5fa9078b8c9cdad590d275ab93beaaee4/sklearn-irisv2-input-grpc.json
```

??? note "About this experiment"
    This experiment consists of two [tasks](../../../getting-started/concepts.md#design), namely, [ready](../../../user-guide/tasks/ready.md) and [grpc](../../../user-guide/tasks/grpc.md). 
    
    The [ready](../../../user-guide/tasks/ready.md) task checks if the `sklearn-irisv2` InferenceService exists and is `Ready`. 

    The [grpc](../../../user-guide/tasks/grpc.md) task sends call requests to the `inference.GRPCInferenceService.ModelInfer` method of the cluster-local gRPC service with host address `${GRPC_HOST}:${GRPC_PORT}`, and collects Iter8's built-in gRPC load test metrics.

***

View the experiment results by using the Iter8 Grafana dashboard, as described in [your first experiment](../../../getting-started/your-first-experiment.md).

??? note "Some variations and extensions of this experiment" 
    1. The [grpc task](../../../user-guide/tasks/grpc.md) can be configured with load related parameters such as the number of requests, requests per second, or number of concurrent connections.

## Clean up

```shell
iter8 k delete
kubectl delete inferenceservice sklearn-irisv2
```