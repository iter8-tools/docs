---
template: main.html
---

# Canary rollout of a KServe ML model

This tutorial shows how Iter8 can be used to implement a canary rollout of ML models hosted in a KServe environment. In a canary rollout, inference requests that match a particular pattern, for example those that have a particular header, are directed to the candidate version of the model. The remaining requests go to the primary, or initial, version of the model. Iter8 enables a canary rollout by automatically configuring the routing resources to distribute inference requests.

After a one-time initialization step, the end user merely deploys candidate models, evaluates them, and either promotes or deletes them. Iter8 automatically handles the underlying routing configuration.

![Canary rollout](images/canary.png)

???+ warning "Before you begin"
    1. Ensure that you have the [kubectl](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs.
    2. Have access to a cluster running [KServe](https://kserve.github.io/website). You can create a [KServe Quickstart](https://kserve.github.io/website/0.10/get_started/#before-you-begin) environment as follows:
    ```shell
    curl -s "https://raw.githubusercontent.com/kserve/kserve/release-0.11/hack/quick_install.sh" | bash
    ```
<!-- Istio 1.15 is installed as part of kserve install -->

## Install the Iter8 controller

--8<-- "docs/tutorials/installiter8controller.md"

## Initialize primary

### Application

Deploy the primary version of the application. In this tutorial, the application is a KServe model. Initialize the resources for the primary version of the model (`v0`) by deploying an `InferenceService` as follows:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-0
  labels:
    app.kubernetes.io/name: wisdom
    app.kubernetes.io/version: v0
    iter8.tools/watch: "true"
spec:
  predictor:
    minReplicas: 1
    model:
      modelFormat:
        name: sklearn
      runtime: kserve-mlserver
      storageUri: "gs://seldon-models/sklearn/mms/lr_model"
EOF
```

??? note "About the primary `InferenceService`"
    The base name (`wisdom`) and version (`v0`) are identified using the labels `app.kubernets.io/name` and `app.kubernets.io/version`, respectively. These labels are not required.

    Naming the instance with the suffix `-0` (and the candidate with the suffix `-1`) simplifies the routing initialization (see below). However, any name can be specified.
    
    The label `iter8.tools/watch: "true"` is required. It lets Iter8 know that it should pay attention to changes to this application resource.

You can inspect the deployed `InferenceService`. When the `READY` field becomes `True`, the model is fully deployed.

```shell
kubectl get inferenceservice wisdom-0
```
    
### Routing

Initialize the routing resources for the application to use a canary rollout strategy:

```shell
cat <<EOF | helm template routing --repo https://iter8-tools.github.io/iter8 routing-actions --version 0.1.5 -f - | kubectl apply -f -
appType: kserve
appName: wisdom
action: initialize
strategy: canary
EOF
```

The `initialize` action (with strategy `canary`) configures the (Istio) service mesh to route all requests to the primary version of the application (`wisdom-0`). It further defines the routing policy that will be used when changes are observed in the application resources. By default, this routing policy sends requests with the header `traffic` set to the value `test` to the candidate version and all remaining requests to the primary version. For detailed configuration options, see the [Helm chart](https://github.com/iter8-tools/iter8/blob/v0.16.0/charts/routing-actions/values.yaml).

## Verify routing

To verify the routing configuration, you can inspect the `VirtualService`:

```shell
kubectl get virtualservice -o yaml wisdom
```

To send inference requests to the model:

=== "From within the cluster"
    1. Create a `sleep` pod in the cluster from which requests can be made:
    ```shell
    curl -s https://raw.githubusercontent.com/iter8-tools/docs/v0.15.2/samples/kserve-serving/sleep.sh | sh -
    ```

    2. Exec into the sleep pod:
    ```shell
    kubectl exec --stdin --tty "$(kubectl get pod --sort-by={metadata.creationTimestamp} -l app=sleep -o jsonpath={.items..metadata.name} | rev | cut -d' ' -f 1 | rev)" -c sleep -- /bin/sh
    ```

    3. Make inference requests:
    ```shell
    cat wisdom.sh
    . wisdom.sh
    ```
    or, to send a request with header `traffic: test`:
    ```shell
    cat wisdom-test.sh
    . wisdom-test.sh
    ```

=== "From outside the cluster"
    1. In a separate terminal, port-forward the ingress gateway:
    ```shell
    kubectl -n istio-system port-forward svc/knative-local-gateway 8080:80
    ```

    2. Download the sample input:
    ```shell
    curl -sO https://raw.githubusercontent.com/iter8-tools/docs/v0.15.2/samples/kserve-serving/input.json
    ```

    3. Send inference requests:
    ```shell
    curl -H 'Content-Type: application/json' -H 'Host: wisdom.default' localhost:8080 -d @input.json -s -D - \
    | grep -e HTTP -e app-version
    ```
    Or, to send a request with header `traffic: test`:
    ```shell
      curl -H 'Content-Type: application/json' -H 'Host: wisdom.default' localhost:8080 -d @input.json -s -D - \
      -H 'traffic: test' \
      | grep -e HTTP -e app-version
    ```

Note that the model version responding to each inference request is noted in the response header `app-version`. In the requests above, we display only the response code and this header.

## Deploy candidate

Deploy a candidate model using a second `InferenceService`:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-1
  labels:
    app.kubernetes.io/name: wisdom
    app.kubernetes.io/version: v1
    iter8.tools/watch: "true"
spec:
  predictor:
    minReplicas: 1
    model:
      modelFormat:
        name: sklearn
      runtime: kserve-mlserver
      storageUri: "gs://seldon-models/sklearn/mms/lr_model"
EOF
```

??? note "About the candidate `InferenceService`"
    In this tutorial, the model source (field `spec.predictor.model.storageUri`) is the same as for the primary version of the model. In a real world example, this would be different.

## Verify routing changes

The deployment of the candidate triggers an automatic routing reconfiguration by Iter8. Inspect the `VirtualService` to see that the routing has been changed. Requests are now distributed between the primary model and the secondary model:

```shell
kubectl get virtualservice wisdom -o yaml
```

You can send additional inference requests as described above. They will be handled by both versions of the model.

## Promote candidate

Promoting the candidate involves redefining the primary version of the application and deleting the candidate version.

### Redefine primary

```shell
cat <<EOF | kubectl replace -f -
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: wisdom-0
  labels:
    app.kubernetes.io/name: wisdom
    app.kubernetes.io/version: v1
    iter8.tools/watch: "true"
spec:
  predictor:
    minReplicas: 1
    model:
      modelFormat:
        name: sklearn
      runtime: kserve-mlserver
      storageUri: "gs://seldon-models/sklearn/mms/lr_model"
EOF
```

??? note "What is different?"
    The version label (`app.kubernets.io/version`) was updated. In a real world example, `spec.predictor.model.storageUri` would also be updated.

### Delete candidate

Once the primary `InferenceService` has been redeployed, delete the candidate:

```shell
kubectl delete inferenceservice wisdom-1
```

### Verify routing changes

Inspect the `VirtualService` to see that the it has been automatically reconfigured to send requests only to the primary.

## Cleanup

If not already deleted, delete the candidate model:

```shell
kubectl delete isvc/wisdom-1
```

Delete routing:

```shell
cat <<EOF | helm template routing --repo https://iter8-tools.github.io/iter8 routing-actions --version 0.1.5 -f - | kubectl delete -f -
appType: kserve
appName: wisdom
action: initialize
strategy: canary
EOF
```

Delete primary:

```shell
kubectl delete isvc/wisdom-0
```

Uninstall Iter8 controller:

--8<-- "docs/tutorials/deleteiter8controller.md"
