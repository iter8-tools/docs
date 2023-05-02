---
template: main.html
---

# Blue-Green Rollout of a ML Model

This tutorial shows how Iter8 can be used to implement a blue-green rollout of ML models hosted in a KServe modelmesh serving environment. In a blue-green rollout, a percentage of inference requests are directed to a candidate version of the model. The remaining requests go to the primary, or initial, version of the model. Iter8 enables a blue-green rollout by automatically configuring the network to distribute inference requests.

After a one time initialization step, the end user merely deploys candidate models, evaluates them and either promotes or deletes them. Optionally, the end user can modify the percentage of inference requests being sent to the candidate model. Iter8 automatically handles all underlying network configuration.

![Blue-Green rollout](images/blue-green.png)

In this tutorial, we use the Istio service mesh to distribute inference requests between different versions of a model.

???+ "Before you begin"
    1. Ensure that you have the [kubectl CLI](https://kubernetes.io/docs/reference/kubectl/).
    2. Have access to a cluster running [KServe ModelMesh Serving](https://github.com/kserve/modelmesh-serving) and [Istio](https://istio.io). For example, you can create a modelmesh-serving [Quickstart](https://github.com/kserve/modelmesh-serving/blob/main/docs/quickstart.md) environment and install a [demo version](https://istio.io/latest/docs/setup/getting-started/) of Istio. Ensure `istioctl` is in your path.

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
    
## Initialize the Blue-Green traffic pattern

Initialize model rollout with a blue-green traffic pattern as follows:

```shell
cat <<EOF | helm template traffic ../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: initialize
targetEnv: kserve-modelmesh
trafficStrategy: blue-green
modelName: wisdom
EOF
```

The `initialize` template (with `trafficStrategy: blue-green`) configures the Istio service mesh to route all requests to the primary version of the model (`wisdom-0`). Further, it defines the routing policy that will be used by Iter8 when it observes changes in the models. By default, this routing policy splits inference requests 50-50 between the primary and candidate versions. For detailed configuration options, see the helm chart.

## Verify network configuration

You can inspect the network configuration:

```shell
kubectl get virtualservice -o yaml wisdom
```

You can also run tests by sending inference requests from a pod in the cluster. For the models in this tutorial you can deploy a pod with the necessary artifacts as follows:

```shell
curl -s https://raw.githubusercontent.com/kalantar/docs/mm/samples/modelmesh-serving/sleep.sh | \
sh - 
```

In a separate terminal, exec into the pod:

```shell
curl -sO https://raw.githubusercontent.com/kalantar/docs/mm/samples/modelmesh-serving/execintosleep.sh
source execintosleep.sh
```

The necessary artifacts are in the directory wisdom:

```shell
cd wisdom
ls -l
```

Run inference requests using `grpcurl` via the scripts `wisdom.sh` and `wisdom-test.sh`:

```shell
. wisdom.sh
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

You can also send inference requests from the sleep pod in the cluster to verify the distribution:

```shell
. wisdom.sh
```

## Modify inference request distribution (optional)

You can modify the weight distribution of inference requests using the Iter8 `traffic-template` chart:

```shell
cat <<EOF | helm template traffic ../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: modify-weights
targetEnv: kserve-modelmesh
trafficStrategy: blue-green
modelName: wisdom
modelVersions:
  - weight: 20
  - weight: 80
EOF
```

Note that using the `modify-weights` modifies the default traffic split for all future candidate deployments.

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

Inspect the `VirtualService` to see that the it has been automaticaly reconfigured to send requests only to the primary model.

## Clean up

Delete the candidate model:

```shell
kubectl delete --force isvc/wisdom-1
```

Delete routing artifacts:

```shell
cat <<EOF | helm template traffic ../../../hub/charts/traffic-templates -f - | kubectl delete --force -f -
templateName: initialize
targetEnv: kserve-modelmesh
trafficStrategy: blue-green
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
