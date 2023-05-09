---
template: main.html
---

# Canary Rollout of a ML Model

This tutorial shows how Iter8 can be used to implement a canary rollout of ML models hosted in a KServe modelmesh serving environment. In a canary rollout, inference requests that match a particular pattern, for example those that have a particular header, are directed to the candidate version of the model. The remaining requests go to the primary, or initial, version of the model. Iter8 enables a canary rollout by automatically configuring the network to distribute inference requests.

After a one time initialization step, the end user merely deploys candidate models, evaluates them and either promotes deletes them. Iter8 automatically handles the underlying network configuration.

![Canary rollout](images/canary.png)

In this tutorial, we use the Istio service mesh to distribute inference requests between different versions of a model.

???+ "Before you begin"
    1. Ensure that you have the [kubectl CLI](https://kubernetes.io/docs/reference/kubectl/).
    2. Have access to a cluster running [ModelMesh Serving](https://github.com/kserve/modelmesh-serving) and [Istio](https://istio.io). For example, you can create a modelmesh-serving [Quickstart](https://github.com/kserve/modelmesh-serving/blob/main/docs/quickstart.md) environment and install a [demo version](https://istio.io/latest/docs/setup/getting-started/) of Istio. Ensure `istioctl` is in your path.

## Install the Iter8 controller

--8<-- "docs/tutorials/integrations/kserve-mm/installiter8controller.md"

## Deploy a primary model

Deploy the primary version of a model using an `InferenceService`:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-0
  namespace: modelmesh-serving
  labels:
    app.kubernetes.io/name: wisdom
    app.kubernetes.io/version: v1
    iter8.tools/watch: "true"
  annotations:
    serving.kserve.io/deploymentMode: ModelMesh
    serving.kserve.io/secretKey: localMinIO
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
EOF
```

??? note "About the primary `InferenceService`"
    Naming the model with the suffix `-0` (and the candidate with the suffix `-1`) simplifies the rollout initialization. However, any names can be specified.
    
    The label `iter8.tools/watch: "true"` lets Iter8 know that it should pay attention to changes to this InferenceService.

Inspect the deployed `InferenceService`:

```shell
kubectl get inferenceservice wisdom-0
```

When the `READY` field becomes `True`, the model is fully deployed.

## Initialize the Canary traffic pattern

Initialize model rollout with a canary traffic pattern as follows:

```shell
cat <<EOF | helm template traffic --repo https://iter8-tools.github.io/hub traffic-templates -f - | kubectl apply -f -
templateName: initialize-rollout
targetEnv: kserve-modelmesh
trafficStrategy: canary
modelName: wisdom
EOF
```

The `initialize-rollout` template (with `trafficStrategy: canary`) configures the Istio service mesh to route all requests to the primary version of the model (`wisdom-0`). Further, it defines a routing policy that will be used by Iter8 when it observes changes in the models. By default, this routing policy sends inference requests with the header `traffic` set to the value `test` to the candidate version of the model and all remaining inference requests to the primary version of the model. For detailed configuration options, see the helm chart.

## Verify network configuration

To verify the network configuration, you can inspect the network configuration:

```shell
kubectl get virtualservice -o yaml wisdom
```

To send inference requests to the model:

1. In a separate terminal, port-forward the ingress gateway:
  ```shell
  kubectl -n istio-system port-forward svc/istio-ingressgateway 8080:80
  ```

2. Download the proto file and a sample input:
  ```shell
  curl -sO https://raw.githubusercontent.com/iter8-tools/docs/v0.13.18/samples/modelmesh-serving/kserve.proto
  curl -sO https://raw.githubusercontent.com/iter8-tools/docs/v0.13.18/samples/modelmesh-serving/grpc_input.json
  ```

3. Send inference requests:
  ```shell
  cat grpc_input.json | \
  grpcurl -plaintext -proto kserve.proto -d @ \
  -authority wisdom.modelmesh-serving \
  localhost:8080 inference.GRPCInferenceService.ModelInfer
  ```

  or, send request with header `traffic: test`:

  ```shell
  cat grpc_input.json | \
  grpcurl -plaintext -proto kserve.proto -d @ \
  -H 'traffic: test' \
  -authority wisdom.modelmesh-serving \
  localhost:8080 inference.GRPCInferenceService.ModelInfer
  ```

Note that the model version responding to each inference request can be determined from the `modelName` field of the response.

## Deploy a candidate model

Deploy a candidate model using a second `InferenceService`:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-1
  namespace: modelmesh-serving
  labels:
    app.kubernetes.io/name: wisdom
    app.kubernetes.io/version: v2
    iter8.tools/watch: "true"
  annotations:
    serving.kserve.io/deploymentMode: ModelMesh
    serving.kserve.io/secretKey: localMinIO
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
EOF
```

??? note "About the candidate `InferenceService`"
    The model name (`wisdom`) and version (`v2`) are recorded using the labels `app.kubernets.io/name` and `app.kubernets.io.version`.

    In this tutorial, the model source (field `spec.predictor.model.storageUri`) is the same as for the primary version of the model. In a real example, this would be different.

## Verify network configuration changes

The deployment of the candidate model triggers an automatic reconfiguration by Iter8. Inspect the `VirtualService` to see that inference requests are now distributed between the primary model and the secondary model according to a header matching rule:

```shell
kubectl get virtualservice wisdom -o yaml
```

Send additional inference requests as described above.

## Promote the candidate model

Promoting the candidate involves redefining the primary `InferenceService` using the new model and deleting the candidate `InferenceService`.

### Redefine the primary `InferenceService`

```shell
cat <<EOF | kubectl replace -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-0
  namespace: modelmesh-serving
  labels:
    app.kubernetes.io/name: wisdom
    app.kubernetes.io/version: v2
    iter8.tools/watch: "true"
  annotations:
    serving.kserve.io/deploymentMode: ModelMesh
    serving.kserve.io/secretKey: localMinIO
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
EOF
```

??? note "What is different?"
    The version label (`app.kubernets.io/version`) was updated. In a real world example, `spec.predictor.model.storageUri` would also be updated.

### Delete the candidate `InferenceService`

```shell
kubectl delete inferenceservice wisdom-1
```

### Verify network configuration changes

Inspect the `VirtualService` to see that the it has been automatically reconfigured to send requests only to the primary model.

## Clean up

Delete the candidate model:

```shell
kubectl delete --force isvc/wisdom-1
```

Delete routing artifacts:

```shell
cat <<EOF | helm template traffic --repo https://iter8-tools.github.io/hub traffic-templates -f - | kubectl delete --force -f -
templateName: initialize-rollout
targetEnv: kserve-modelmesh
trafficStrategy: canary
modelName: wisdom
EOF
```

Delete the primary model:

```shell
kubectl delete --force isvc/wisdom-0
```

Delete the sleep pod:

```shell
kubectl delete --force deploy/sleep configmap/wisdom-input
```

Uninstall the Iter8 controller:

--8<-- "docs/tutorials/integrations/kserve-mm/deleteiter8controller.md"
