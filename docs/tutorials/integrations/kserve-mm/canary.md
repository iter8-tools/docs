---
template: main.html
---

# Model Rollouts: Blue-Green, Canary, and Mirroring

???+ warning "TODO"
    1. fix template chart location (all calls) to use repository
    2. do we have an example where we really have two versions with different values for `storageUri`
    3. make this work with mesh gateway
    4. make this work for multiple models
    5. picture(s)
    6. template for cleanup?
    7. should we include sleep.sh, execintosleep.sh? These are tied to the example model. Is this a standard example that will be available long term?

This tutorial shows how Iter8 can be used to implement a blue-green rollout of ML models. In a blue-green rollout, a percentage of inference requests are directed to a candidate version. The remaining requests go to the primary, or initial, version. Iter8 enables a blue-green rollout by automatically configuring the network to distribute inference requests.

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

??? note "Some observations on the InferenceService"
    **Something about target namespace?**

    Naming the model with the suffix `-0` (and the candidate with the suffix `-1`) simplifies configuration of any experiments. Iter8 assumes this convention by default. However, any names can be specified.
    
    The label `iter8.tools/watch: "true"` lets Iter8 know that it should pay attention to changes to this InferenceService.

??? note "Verifying Deployment of InferenceService"
    Inspect the output of `kubectl get` to see that the `READY` field becomes `True`.
    
    ```shell
    kubectl -n modelmesh-serving get inferenceservice wisdom-0
    ```

## Initialize Rollout Traffic Pattern

=== "Blue-Green"
    Initialize the model rollout with a blue-green traffic pattern as follows.

    ```shell
    cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
    templateName: initialize
    targetEnv: kserve-modelmesh
    trafficStrategy: blue-green
    modelName: wisdom
    modelVersions:
    - weight: 50
    - weight: 50
    EOF
    ```
    ??? note "What Happended?"
        The `initialize` template for `trafficStrategy: blue-green` does two things. First, it configures the Istio service mesh to route all requests to the primary model (`wisdom-0`). Second, it defines a routing policy that will be used by Iter8 when it observes new/candidate versions (or their promotion). This routing policy splits inferences requests 50-50 between the primary and candidate versions when both are present.

=== "Canary"
    Initialize the model rollout with a canary traffic pattern as follows.

    ```shell
    cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - | kubectl apply -f -
    templateName: initialize
    targetEnv: kserve-modelmesh
    trafficStrategy: canary
    modelName: wisdom
    modelVersions:
    - name: wisdom-0
    - name: wisdom-1
      match:
      - headers:
          traffic:
            exact: test
    EOF
    ```

    ???+ warning "Question"
        Canary weight in list of variants should be 100?

    ??? note "What Happended?"
        The `initialize` template for `trafficStrategy: canary`. First, it configures the Istio service mesh to route all requests to the primary model (`wisdom-0`). Second, it defines a routing policy that will be used by Iter8 when it observes new/candidate versions (or their promotion). This routing policy sends traffic matching a particular pattern to the candidate version while remaining traffic is sent to the primary version.

=== "Mirroring"
    Initialize the model rollout with a mirroring traffic pattern as follows.

    ```shell
    cat <<EOF | helm template traffic ../../../../hub/charts/traffic-templates -f - > init.yaml
    templateName: initialize
    targetEnv: kserve-modelmesh
    trafficStrategy: mirror
    modelName: wisdom
    modelVersions:
      # can only mirror to a single candidate version
      - name: wisdom-0
      - name: wisdom-1
    EOF
    ```

    ??? note "What Happended?"
        The `initialize` template for `trafficStrategy: mirror`. First, it configures the Istio service mesh to route all requests to the primary model (`wisdom-0`). Second, it defines a routing policy that will be used by Iter8 when it observes new/candidate versions (or their promotion). This routing policy continutes to send all traffic to the primary version, but also mirrors it to the candidate version. Responses from the candidate are ignored.

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

??? note "Comments on the candiate `InferenceService`"
    The model name (`wisdom`) and version (`v2`) are recorded using the labels `app.kubernets.io/name` and `app.kubernets.io.version`.

    In this tutorial, the model source (see the field `spec.predictor.model.storageUri`) is the same as for the primary. In a real example, this would be different.

## Modify Inference Request Distribution

??? note "Optional; Applies only to the blue-green traffic pattern"

You can modify the weight distribution of inference requests by annotating the model objects with the desired weights. The Iter8 `traffic-template` chart can be used to simplify doing so:

```shell
kubectl annotate --override isvc wisdom-0 iter8.tools/weight='20'
kubectl annotate --override isvc wisdom-1 iter8.tools/weight='80'
```

??? note "Verifying the change"
    You can verify the change in inference request distribution by inspecting the `VirtualService`:

    ```shell
    kubectl get virtualservice wisdom -o yaml
    ```

    As you continue to send requests, you should see the distribution change (see the `modelName` field in inference responses).

## Promote the Candidate Model

Promoting the candidate is a two step process. First, redefine the primary `InferenceService` using the new model. Second, delete the candidate `InferenceService`.

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
    The version label `app.kubernets.io/version` has been updated and the annotation `iter8/weight` was removed. In practise, the field `spec.predictor.model.storageUri` would also be updated.

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