---
template: main.html
---

# Canary Rollout of a Model

???+ warning "TODO"
    1. fix template chart location (all calls) to use repository
    2. do we have an example where we really have two versions with different values for `storageUri`
    3. make this work with mesh gateway
    4. make this work for multiple models
    5. picture(s)
    6. template for cleanup?
    7. should we include sleep.sh, execintosleep.sh? These are tied to the example model. Is this a standard example that will be available long term?

This tutorial shows how Iter8 can be used to implement a canary rollout of ML models. In a canary rollout, inference requests that match a particular pattern, for example those that have a particular header, are directed to the candidate version. The remaining requests go to the primary, or initial, version. Iter8 enables a canary rollout by automatically configuring the network to distribute inference requests.

After a one time initialization step, the end user merely deploys candidate models, evaluates them and either promotes the candiate or deletes it. Optionally, the end user can modify the percentage of inference requests being sent to the candidate. Iter8 automatically handles the underlying network configuration.

In this tutorial, we use the Istio service mesh to distribute inference requests between different versions of a model.

???+ "Before you begin"
    1. Install the Iter8 controller.
    2. Ensure that you have the [kubectl CLI](https://kubernetes.io/docs/reference/kubectl/).
    3. Have access to a cluster running [ModleMesh Serving](https://github.com/kserve/modelmesh-serving) and [Istio](https://istio.io). For example, you can create a modelmesh-serving [Quickstart](https://github.com/kserve/modelmesh-serving/blob/main/docs/quickstart.md) environment and install a [demo version](https://istio.io/latest/docs/setup/getting-started/) of Istio.

## Install the Iter8 controller

```shell
curl -s https://raw.githubusercontent.com/iter8-tools/iter8/v0.13.13/testdata/controllers/config.yaml | \
helm install --repo https://iter8-tools.github.io/hub iter8-traffic traffic -f -
```

??? note "About controller configuration"
    The configuration file specifies a list resources types to watch. The default list is suitable for KServe modelmesh supported ML models. 

## Configure External Routing (Optional)

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: external
targetEnv: kserve-modelmesh
EOF
```

## Deploy Initial InferenceService

Deploy the initial version the infernce service:

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
    *Something about target namespace?*

    Naming the model with the suffix `-0` (and the candidate with the suffix `-1`) simplifies configuration of any experiments. Iter8 assumes this convention by default. However, any names can be specified.
    
    The label `iter8.tools/watch: "true"` lets Iter8 know that it should pay attention to changes to this InferenceService.

??? note "Verifying Deployment of InferenceService"
    Inspect the output of `kubectl get` to see that the `READY` field becomes `True`.
    
    ```shell
    kubectl -n modelmesh-serving get inferenceservice wisdom-0
    ```

## Initialize Rollout Traffic Pattern

Initialize the model rollout with a canary traffic pattern as follows.

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: initialize
targetEnv: kserve-modelmesh
trafficStrategy: canary
modelName: wisdom
EOF
```

The `initialize` template for `trafficStrategy: canary` configures the Istio service mesh to route all requests to the primary model (`wisdom-0`). Further, it defines a routing policy that will be used by Iter8 when it observes new candidate versions (or their promotion). By default, this routing policy sends inference requests with the header `traffic` set to the value `test` to the candidate version while remaining inference requets are sent to the primary version.

??? note "Verifying Network Configuration"
    In all cases, when just the primary model is deployed, all inference requests should be sent to it. You can test using a tool such as `grpcurl`.

## Deploy a Candidate Model

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

??? note "About the candiate `InferenceService`"
    The model name (`wisdom`) and version (`v2`) are recorded using the labels `app.kubernets.io/name` and `app.kubernets.io.version`.

    In this tutorial, the model source (see the field `spec.predictor.model.storageUri`) is the same as for the primary. In a real example, this would be different.

## Promote the Candidate Model

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

??? note "What changed?"
    The version label (`app.kubernets.io/version`) and `spec.predictor.model.storageUri` are updated.

### Delete the candidate `InferenceService`

```shell
kubectl delete inferenceservice wisdom-1
```

??? note "Some variations and extensions to try" 
    1. 

## Clean up

Delete all artifacts:

```shell
kubectl delete \
  isvc/wisdom-0 \
  isvc/wisdom-1 \
  virtualservice.networking.istio.io/wisdom \
  gateway.networking.istio.io/mm-external-gateway \
  service/mm-external \
  configmap/wisdom-0-weight-config \
  configmap/wisdom-1-weight-config \
  configmap/wisdom-routemap \
  deployment/sleep \
  configmap/wisdom-input
```

Uninstall the Iter8 controller:

```shell
helm delete iter8-traffic
```