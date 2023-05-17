---
template: main.html
---

# Mirrored Rollout of a ML Model

This tutorial shows how Iter8 can be used to implement a mirrored rollout of ML models in a KServe modelmesh serving environment. In a mirrored rollout, all inference requests are sent to the primary version of the model. In addition, a portion of the requests are also sent to the candidate version of the model. The responses from the candidate version are ignored. Iter8 enables a mirrored rollout by automatically configuring the network to distribute inference requests.

After a one time initialization step, the end user merely deploys candidate models, evaluates them, and either promotes or deletes them. Optionally, the end user can modify the percentage of inference requests being sent to the candidate version of the model. Iter8 automatically handles the underlying network configuration.

![Mirror rollout](images/mirror.png)

In this tutorial, we use the Istio service mesh to distribute inference requests between different versions of a model.

???+ "Before you begin"
    1. Ensure that you have the [kubectl CLI](https://kubernetes.io/docs/reference/kubectl/).
    2. Have access to a cluster running [KServe ModelMesh Serving](https://github.com/kserve/modelmesh-serving). For example, you can create a modelmesh-serving [Quickstart](https://github.com/kserve/modelmesh-serving/blob/main/docs/quickstart.md) environment.
    3. Install [Istio](https://istio.io). You can install the [demo profile](https://istio.io/latest/docs/setup/getting-started/).

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
    Naming the model with the suffix `-0` (and the candidate with the suffix `-1`) simplifies the rollout initialization. However, any name can be specified.
    
    The label `iter8.tools/watch: "true"` lets Iter8 know that it should pay attention to changes to this `InferenceService`.

Inspect the deployed `InferenceService`:

```shell
kubectl get inferenceservice wisdom-0
```

When the `READY` field becomes `True`, the model is fully deployed.
    
## Initialize the Mirroring routing policy

Initialize the model rollout with a mirror traffic pattern as follows:

```shell
cat <<EOF | helm template traffic --repo https://iter8-tools.github.io/iter8 traffic-templates -f - | kubectl apply -f -
templateName: initialize-rollout
targetEnv: kserve-modelmesh
trafficStrategy: mirror
modelName: wisdom
EOF
```

The `initialize-rollout` template ( with `trafficStrategy: mirror`) configures the Istio service mesh to route all requests to the primary model (`wisdom-0`). Further, it defines a routing policy that will be used by Iter8 when it observes changes in the models. By default, this routing policy sends all inference requests to the primary version of the model. It also sends them all of them (100%) to the candidate version of the model. Responses from the candidate are ignored.

## Verify network configuration

You can inspect the network configuration:

```shell
kubectl get virtualservice -o yaml wisdom
```

To send inference requests to the model:

=== "From within the cluster"
    1. Create a "sleep" pod in the cluster from which requests can be made:
    ```shell
    curl -s https://raw.githubusercontent.com/iter8-tools/docs/v0.14.1/samples/modelmesh-serving/sleep.sh | sh -
    ```

    2. exec into the sleep pod:
    ```shell
    kubectl exec --stdin --tty "$(kubectl get pod --sort-by={metadata.creationTimestamp} -l app=sleep -o jsonpath={.items..metadata.name} | rev | cut -d' ' -f 1 | rev)" -c sleep -- /bin/sh
    ```

    3. Make inference requests:
    ```shell
    cd demo
    cat wisdom.sh
    . wisdom.sh
    ```

=== "From outside the cluster"
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

    In this tutorial, the model source (field `spec.predictor.model.storageUri`) is the same as for the primary version of the model. In a real world example, this would be different.

## Verify network configuration changes

The deployment of the candidate model triggers an automatic reconfiguration by Iter8. Inspect the `VirtualService` to see that inference requests are now distributed between the primary model and the secondary model:

```shell
kubectl get virtualservice wisdom -o yaml
```

## Modify the percentage of mirrored traffic (optional)

You can modify the percentage of inference requests that are mirrored (send to the candidate version) using the Iter8 `traffic-template` chart. For example, to change the mirrored percentage to 20%, use:

```shell
cat <<EOF | helm template traffic --repo https://iter8-tools.github.io/iter8 traffic-templates -f - | kubectl apply -f -
templateName: modify-weights
targetEnv: kserve-modelmesh
trafficStrategy: mirror
modelName: wisdom
mirrorPercentage: 20
EOF
```

Note that using the `modify-weights` overrides the default behavior for all future candidate deployments.

As above, you can verify the network configuration changes.

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
cat <<EOF | helm template traffic --repo https://iter8-tools.github.io/iter8 traffic-templates -f - | kubectl delete --force -f -
templateName: initialize-rollout
targetEnv: kserve-modelmesh
trafficStrategy: mirror
modelName: wisdom
EOF
```

Delete the primary model:

```shell
kubectl delete --force isvc/wisdom-0
```

Uninstall the Iter8 controller:

--8<-- "docs/tutorials/integrations/kserve-mm/deleteiter8controller.md"
