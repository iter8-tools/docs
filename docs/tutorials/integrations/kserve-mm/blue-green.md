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

## Deploy Initial InferenceService

Deploy the initial version the infernce service:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-primary
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
    Naming the model with the suffix `-primary` (and any candidate with the suffix `-candidate`) simplifies configuration of any experiments. However, it is not strictly required.
    
    The label `iter8.tools/watch: "true"` lets Iter8 know that it should pay attention to changes to the InferenceService.


## Initialize Routing

```shell
cat << EOF | helm template --repo https://iter8-tools.github.io/hub traffic-templates -f - | kubectl apply -f -
targetEnv: kserve-modelmesh
externalAccess: true
trafficStrategy: blue-green
modelName: wisdom
- versions: wisdom-primary
- versions: wisdom-candidate
EOF
```

???+ warning "Question"
    Did we intend to have 2 steps for externalizing and starting

??? note "Initialization details"
    The initialization step does two things:

    1. Configure Istio 
    2. Associates the models associated with the name (prefix) `wisdom` with the blue-green rollout strategy.

    By default, traffic will be split 50-50 between the primary and candidate versions. 

## Deploy a Candidate Model

cat <<EOF | kubectl apply -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-candidate
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
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: wisdom-candidate-weight
  labels:
    iter8.tools/watch: "true"
EOF

???+ warning "Question"
    How does the user know to create the configmap?

## Optionally modify inference requets distribution

Modify the weight distribution to the candidate model; the weight distributed to the primary model can be determined by this. For example, to send 80% of the inference requests to the candidate model:

```shell
kubectl annotate --overwrite \
configmap wisdom-candidate-weight-config \
iter8.tools/weight='80'
```

???+ warning "Question"
    How does the user know the name of the configmap?

## Promote the Candidate Model

Promotion is two steps: (a) redefine the primary model using the new version and (b) delete the candidate version:

```shell
cat <<EOF | kubectl replace -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-primary
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
kubectl delete isvc/wisdom-candidate cm/wisdom-candidate-weight
```

???+ warning "Question"
    How does the user know to delete the configmap?

## 

??? note "Some variations and extensions of this experiment" 
    1. Modify the inference weight distribution
    2. 

## Clean up

Delete all artifacts:

```shell
kubectl delete isvc/wisdom-candidate cm/wisdom-candidate-weight-config

cat << EOF | helm template --repo https://iter8-tools.github.io/hub traffic-templates -f - | kubectl delete -f -
targetEnv: blue-green
modelName: wisdom
- versions: wisdom-primary
- versions: wisdom-candidate
EOF

kubectl delete isvc/wisdom-primary
```

Delete the Iter8 controller:

```shell
helm delete iter8-traffic
```