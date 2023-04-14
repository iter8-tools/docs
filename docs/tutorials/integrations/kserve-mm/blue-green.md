---
template: main.html
---

# Blue-Green Model Rollouts

This tutorial shows how Iter8 can be used to implement a blue-green rollout of ML models. In a blue-green rollout, a percentage of inference requests are directed to a candidate version. The remaining requests go to the primary, or initial, version. Iter8 enables a blue-green rollout by automatically configuring the network to distribute inference requests.

After a one time initialization step, the end user merely deploys candidate models, evaluates them and either promotes the candiate or deletes it. Optionally, the end user can modify the percentage of inference requests being sent to the candidate. Iter8 automatically handles the underlying network configuration.

In this tutorial, we use the Istio service mesh to distribute inference requests between different versions of a model.

???+ "Before you begin"
    1. Install the Iter8 controller.
    2. Ensure that you have the [kubectl](https://kubernetes.io/docs/reference/kubectl/) CLI.
    3. Have access to a cluster running [ModleMesh Serving](https://github.com/kserve/modelmesh-serving) and [Istio](). For example, you can create a [Quickstart](https://github.com/kserve/modelmesh-serving/blob/main/docs/quickstart.md) environment.

## Install the Iter8 controller

```shell
curl -s https://raw.githubusercontent.com/iter8-tools/iter8/v0.13.13/testdata/controllers/config.yaml | \
helm install --repo https://iter8-tools.github.io/hub iter8-traffic traffic --values -
```

??? note "About controller configuration"
    The configuration file specifies a list resources types to watch. The default list is suitable for KServe modelmesh supported ML models. 

## Configure External Routing

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: external
targetEnv: kserve-modelmesh
EOF
```

???+ warning "TODO"
    fix template chart location (all calls)

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

??? note "Some observations on the InferenceService"
    Naming the model with the suffix `-0` (and the candidate with the suffix `-1`) simplifies configuration of any experiments. Iter8 assumes this convention by default. However, any names can be specified.
    
    The label `iter8.tools/watch: "true"` lets Iter8 know that it should pay attention to changes to this InferenceService.

## Initialize Routing

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: bg-initialize
targetEnv: kserve-modelmesh
trafficStrategy: blue-green
modelName: wisdom
modelVersions:
- weight: 50
- weight: 50
EOF
```

??? note "Initialization details"
    The initialization step does two things:

    1. Configure Istio to route requests to the primary model.

    2. Associates the `modelName` (`wisdom`) models (`wisdom-0` and `wisdom-1`) with the blue-green rollout strategy. The blue-green traffic strategy is configured to send 50% of requests to each of a primary and candidate version when both are present.

## Deploy a Candidate Model

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

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: deploy-candidate
targetEnv: kserve-modelmesh
trafficStrategy: blue-green
modelName: wisdom
EOF
```

## Optionally modify inference request distribution

Modify the distribution weight for inference requests.

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: modify-weights
targetEnv: kserve-modelmesh
trafficStrategy: blue-green
modelName: wisdom
modelVersions:
  - weight: 20
  - weight: 80
EOF
```

## Promote the Candidate Model

Promotion is two steps: (a) redefine the primary model using the new version and (b) delete the candidate version.

### Redefine the deployed primary model

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

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
templateName: promote-candidate
targetEnv: kserve-modelmesh
trafficStrategy: blue-green
modelName: wisdom
EOF
```

### Delete candidate model

```shell
kubectl delete isvc wisdom-1
```

```shell
cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl delete -f -
templateName: deploy-candidate
targetEnv: kserve-modelmesh
trafficStrategy: blue-green
modelName: wisdom
EOF
```

## 

??? note "Some variations and extensions of this experiment" 
    1. Modify the inference weight distribution
    2. 

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

???+ warning "TODO"
    Should we use a template for cleanup too?

Uninstall the Iter8 controller:

```shell
helm delete iter8-traffic
```