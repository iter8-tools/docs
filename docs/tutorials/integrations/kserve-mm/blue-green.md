---
template: main.html
---

# Blue-green release of a ML model

This tutorial shows how Iter8 can be used to release ML models hosted in a KServe ModelMesh environment using a blue-green rollout strategy. 
In a blue-green rollout, a percentage of requests are directed to a candidate version of the model. 
This percentage can be changed over time. 
The user declaratively describes the desired application state at any given moment. 
An Iter8 `release` chart assists users who describe the application state at any given moment. 
The chart provides the configuration needed for Iter8 to automatically deploy application versions and configure the routing to implement the blue-green rollout strategy.

![Blue-green rollout](../../images/blue-green.png)

???+ warning "Before you begin"
    1. Ensure that you have the [`kubectl`](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs.
    2. Have access to a cluster running [KServe ModelMesh Serving](https://github.com/kserve/modelmesh-serving). For example, you can create a modelmesh-serving [Quickstart](https://github.com/kserve/modelmesh-serving/blob/release-0.11/docs/quickstart.md) environment.  If using the Quickstart environment, change your default namespace to `modelmesh-serving`: 
    ```shell
    kubectl config set-context --current --namespace=modelmesh-serving
    ```
    3. Install [Istio](https://istio.io). It suffices to install the [demo profile](https://istio.io/latest/docs/setup/getting-started/), for example by using: 
    ```shell
    istioctl install --set profile=demo -y
    ```
## Install the Iter8 controller

--8<-- "docs/getting-started/install.md"

## Deploy initial version

Deploy the initial version of the model using the Iter8 `release` chart by identifying the environment into which it should be deployed, a list of the versions to be deployed (only one here), and the rollout strategy to be used:

```shell
cat <<EOF | helm upgrade --install wisdom --repo https://iter8-tools.github.io/iter8 release --version 0.18 -f -
environment: kserve-modelmesh-istio
application: 
  metadata:
    labels:
      app.kubernetes.io/name: wisdom
    annotations:
      serving.kserve.io/secretKey: localMinIO
  modelFormat: sklearn
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
  strategy: blue-green
EOF
```

Wait for the backend model to be ready:

```shell
kubectl wait --for condition=ready isvc/wisdom-0 --timeout=600s
```

??? note "What happens?"
    - Because `environment` is set to `kserve-modelmesh-istio`,  an `InferenceService` object is created.
    - The namespace `default` is inherited from the Helm release namespace since it is not specified in the version or in `application.metadata`.
    - The name `wisdom-0` is derived from the Helm release name since it is not specified in the version or in `application.metadata`. The name is derived by appending the index of the version in the list of versions; `-0` in this case.
    - Alternatively, an `inferenceServiceSpecification` could have been provided.

    To support routing, a `ServiceEntry` named `default/wisdom` is deployed. Further, an Iter8 [routemap](../../../user-guide/topics/routemap.md) is created.

Once the `InferenceService` is ready, the Iter8 controller automatically configures the routing by creating an Istio `VirtualService`. It is configured to route all inference requests to the only deployed version, `wisdom-0`.

### Verify routing

You can send verify the routing configuration by inspecting the `VirtualService`:

```shell
kubectl get virtualservice wisdom -o yaml
```

You can also send inference requests from a pod within the cluster:

1. Create a `sleep` pod in the cluster from which requests can be made:
```shell
curl -s https://raw.githubusercontent.com/iter8-tools/docs/v0.17.3/samples/modelmesh-serving/sleep.sh | sh -
```

2. Exec into the sleep pod:
```shell
kubectl exec --stdin --tty "$(kubectl get pod --sort-by={metadata.creationTimestamp} -l app=sleep -o jsonpath={.items..metadata.name} | rev | cut -d' ' -f 1 | rev)" -c sleep -- /bin/sh
```

3. Send requests:
```shell
cat grpc_input.json \
| grpcurl -vv -plaintext -proto kserve.proto -d @ \
  -authority wisdom.modelmesh-serving \
  modelmesh-serving.modelmesh-serving:8033 \
  inference.GRPCInferenceService.ModelInfer \
| grep -e app-version
```

The output includes the version of the application that responded (in the `app-version` response header). In this example:

```
app-version: wisdom-0
```

??? note "To send requests from outside the cluster"
    To configure the release for traffic from outside the cluster, a suitable Istio `Gateway` is required. For example, this [sample gateway](https://raw.githubusercontent.com/kalantar/docs/release/samples/iter8-sample-gateway.yaml). When using the Iter8 `release` chart, set the `gateway` field to the name of your `Gateway`. Finally, to send traffic:

    (a) In a separate terminal, port-forward the ingress gateway:
    ```shell
    kubectl -n istio-system port-forward svc/istio-ingressgateway 8080:80
    ```
    (b) Download the proto file and sample input:
    ```shell
    curl -sO https://raw.githubusercontent.com/iter8-tools/docs/v0.17.3/samples/modelmesh-serving/kserve.proto
    curl -sO https://raw.githubusercontent.com/iter8-tools/docs/v0.17.3/samples/modelmesh-serving/grpc_input.json
    ```
    \(c) Send requests using the `Host` header:
    ```shell
    cat grpc_input.json | \
    grpcurl -vv -plaintext -proto kserve.proto -d @ \
    -authority wisdom.modelmesh-serving \
    localhost:8080 inference.GRPCInferenceService.ModelInfer \
    | grep -e app-version
    ```

## Deploy candidate

A candidate version of the model can be deployed simply by adding a second version to the list of versions comprising the application:

```shell
cat <<EOF | helm upgrade --install wisdom --repo https://iter8-tools.github.io/iter8 release --version 0.18 -f -
environment: kserve-modelmesh-istio
application: 
  metadata:
    labels:
      app.kubernetes.io/name: wisdom
    annotations:
      serving.kserve.io/secretKey: localMinIO
  modelFormat: sklearn
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    storageUri: "s3://modelmesh-example-models/sklearn/mnist-svm.joblib"
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    storageUri: "s3://modelmesh-example-models/sklearn/mnist-svm.joblib"
  strategy: blue-green
EOF
```

??? note "About the candidate"
    In this tutorial, the model source (field `storageUri`) for the candidate version is the same as for the primary version of the model. In a real example, this would be different. The version label (`app.kubernetes.io/version`) can be used to distinguish between versions.

When the candidate version is ready, the Iter8 controller will Iter8 will automatically reconfigure the routing so that requests are sent to both versions.

### Verify Routing

You can verify the routing configuration by inspecting the `VirtualService` and/or by sending requests as described above. Requests will be handled equally by both versions. Output will be something like:

```
app-version: wisdom-0
...
app-version: wisdom-1
```

## Modify weights (optional)

To modify the request distribution between versions, add a `weight` to each version. The weights are relative to each other.

```shell
cat <<EOF | helm upgrade --install wisdom --repo https://iter8-tools.github.io/iter8 release --version 0.18 -f -
environment: kserve-modelmesh-istio
application: 
  metadata:
    labels:
      app.kubernetes.io/name: wisdom
    annotations:
      serving.kserve.io/secretKey: localMinIO
  modelFormat: sklearn
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    storageUri: "s3://modelmesh-example-models/sklearn/mnist-svm.joblib"
    weight: 30
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    storageUri: "s3://modelmesh-example-models/sklearn/mnist-svm.joblib"
    weight: 70
  strategy: blue-green
EOF
```

Iter8 automatically reconfigures the routing to distribute traffic between the versions based on the new weights.

### Verify Routing

You can verify the routing configuration by inspecting the `VirtualService` and/or by sending requests as described above. 70 percent of requests will now be handled by the candidate version (`wisdom-1`); the remaining 30 percent by the primary version (`wisdom-0`).

## Promote candidate

```shell
cat <<EOF | helm upgrade --install wisdom --repo https://iter8-tools.github.io/iter8 release --version 0.18 -f -
environment: kserve-modelmesh-istio
application: 
  metadata:
    labels:
      app.kubernetes.io/name: wisdom
    annotations:
      serving.kserve.io/secretKey: localMinIO
  modelFormat: sklearn
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
  strategy: blue-green
EOF
```

??? note "What is different?"
    The version label (`app.kubernetes.io/version`) of the primary version was updated. In a real world example, `storageUri` would also be updated (with that from the candidate version).

Once the (reconfigured) primary `InferenceService` ready, the Iter8 controller will automatically reconfigure the routing to send all requests to it.

### Verify Routing

You can verify the routing configuration by inspecting the `VirtualService` and/or by sending requests as described above. They will all be handled by the primary version. Output will be something like:

```
app-version: wisdom-0
```

## Cleanup

Delete the models are their routing:

```shell
helm delete wisdom
```

If you used the `sleep` pod to generate load, remove it:

```shell
kubectl delete deploy sleep
```

Uninstall Iter8 controller:

--8<-- "docs/getting-started/uninstall.md"
