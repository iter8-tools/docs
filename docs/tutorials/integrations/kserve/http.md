---
template: main.html
---

# Load Test a KServe Model (via HTTP)

This tutorial shows how easy it is to run a load test for KServe when using HTTP to make requests. We use a sklearn model to demonstrate. The same approach works for any model type. 

???+ warning "Before you begin"
    1. Try [your first experiment](../../../getting-started/your-first-experiment.md). Understand the main [concepts](../../../getting-started/concepts.md) behind Iter8 experiments.
    2. Ensure that you have the [kubectl](https://kubernetes.io/docs/reference/kubectl/) CLI.
    3. Have access to a cluster running [KServe](https://kserve.github.io/website). You can create a [KServe Quickstart](https://kserve.github.io/website/0.10/get_started/#before-you-begin) environment as follows:
    ```shell
    curl -s "https://raw.githubusercontent.com/kserve/kserve/release-0.10/hack/quick_install.sh" | bash
    ```

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

***

## Launch Experiment

Launch an Iter8 experiment inside the Kubernetes cluster:

```shell
iter8 k launch \
--set "tasks={ready,http}" \
--set ready.isvc=sklearn-irisv2 \
--set ready.timeout=180s \
--set http.url=http://sklearn-irisv2.default.svc.cluster.local/v2/models/sklearn-irisv2/infer \
--set http.payloadURL=https://gist.githubusercontent.com/kalantar/d2dd03e8ebff2c57c3cfa992b44a54ad/raw/97a0480d0dfb1deef56af73a0dd31c80dc9b71f4/sklearn-irisv2-input.json \
--set http.contentType="application/json"
```

??? note "About this experiment"
    This experiment consists of two [tasks](../../../getting-started/concepts.md#design), namely, [ready](../../../user-guide/tasks/ready.md) and [http](../../../user-guide/tasks/http.md). 
    
    The [ready](../../../user-guide/tasks/ready.md) task checks if the `sklearn-irisv2` InferenceService exists and is `Ready`. 

    The [http](../../../user-guide/tasks/http.md) task sends requests to the cluster-local HTTP service whose URL exposed by the InferenceService, `http://sklearn-irisv2.default.svc.cluster.local/v2/models/sklearn-irisv2/infer`, and collects [Iter8's built-in HTTP load test metrics](../../../user-guide/tasks/http.md#metrics).

***

View the experiment results by using the Iter8 Grafana dashboard, as described in [your first experiment](../getting-started/your-first-experiment.md).

??? note "Some variations and extensions of this experiment" 
    1. The [http task](../../../user-guide/tasks/http.md) can be configured with load related parameters such as the number of requests, queries per second, or number of parallel connections.

## Clean up

```shell
iter8 k delete
kubectl delete inferenceservice sklearn-irisv2
```

