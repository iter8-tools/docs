---
template: main.html
---

# Mirrored Rollout of a ML Model

This tutorial shows how Iter8 can be used to implement a mirrored rollout of ML models in a KServe modelmesh serving environment. In a mirrored rollout, all inference requests are sent to the primary version of the model. In addition, a portion of the requests are also sent to the candidate version of the model. The responses from the candidate version are ignored. Iter8 enables a mirrored rollout by automatically configuring the network to distribute inference requests.

After a one time initialization step, the end user merely deploys candidate models, evaluates them and either promotes or deletes them. Optionally, the end user can modify the percentage of inference requests being sent to the candidate version of the model. Iter8 automatically handles the underlying network configuration.

In this tutorial, we use the Istio service mesh to distribute inference requests between different versions of a model.

???+ "Before you begin"
    1. Ensure that you have the [kubectl CLI](https://kubernetes.io/docs/reference/kubectl/).
    2. Have access to a cluster running [ModelMesh Serving](https://github.com/kserve/modelmesh-serving) and [Istio](https://istio.io). For example, you can create a modelmesh-serving [Quickstart](https://github.com/kserve/modelmesh-serving/blob/main/docs/quickstart.md) environment and install a [demo version](https://istio.io/latest/docs/setup/getting-started/) of Istio.

## Install the Iter8 controller

The Iter8 controller can be installed using a helm chart as follows:

```shell
helm install --repo https://iter8-tools.github.io/hub iter8-traffic traffic
```

## Configure external routing (optional)

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: external
targetEnv: kserve-modelmesh
EOF
```

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
kubectl get inferenceservice wisdom-0 -o yaml
```

When the `READY` field becomes `True`, the model is fully deployed.
    
## Initialize the Mirror traffic pattern

Initialize the model rollout with a mirror traffic pattern as follows:

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: initialize
targetEnv: kserve-modelmesh
trafficStrategy: mirror
modelName: wisdom
EOF
```

The `initialize` template ( with `trafficStrategy: mirror`) configures the Istio service mesh to route all requests to the primary model (`wisdom-0`). Further, it defines a routing policy that will be used by Iter8 when it observes changes in the models. By default, this routing policy sends all inference requests to the primary version of the model. It also sends them all of them (100 %) to the candidate version of the model. Responses from the candidate are ignored.

## Verify network configuration

You can inspect the network configuration:

```shell
get virtualservice -o yaml wisdom
```

You can also run tests by sending inference requests from a pod in the cluster. For the models in this tutorial you can deploy a pod with the necessary artifacts as follows:

```shell
curl -s https://raw.githubusercontent.com/iter8-tools/doc/master/samples/controllers/canary-mm/sleep.sh | \
sh - 
```

In a separate terminal, exec into the pod:

```shell
curl -sO https://raw.githubusercontent.com/kalantar/iter8/mm-demos/testdata/controllers/canary-mm/execintosleep.sh | \
sh -
```

The necessary artifacts are in the directory wisdom:

```shell
cd widsom
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

## Modify the percentage of mirrored traffic (optional)

You can modify the percentage of inference requests that are mirrored (send to the candidate version) using the Iter8 `traffic-template` chart. For example, to change the mirrored percentage to 20%, use:

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: modify-weights
targetEnv: kserve-modelmesh
trafficStrategy: mirror
modelName: wisdom
mirrorPercentage: 20
EOF
```

Note that using the `modify-weights` modifies the default behavior for all future candidate deployments.

??? note "Verifying Network Configuration Change

    You can verify the change in inference request distribution by inspecting the `VirtualService`:

    ```shell
    kubectl get virtualservice wisdom -o yaml
    ```

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

## Clean up

Delete the candidate model:

```shell
kubectl delete --force isvc/wisdom-1
```

Delete routing artifacts:

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl delete --force -f -
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

Delete artifacts created to configure external routing (if created):

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl delete --force -f -
templateName: external
targetEnv: kserve-modelmesh
EOF
```

Delete the sleep pod:

```shell
kubectl delete --force deploy/sleep configmap/wisdom-input
```

Uninstall the Iter8 controller:

```shell
helm delete iter8-traffic
```
