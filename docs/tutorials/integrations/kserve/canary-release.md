---
template: main.html
---

# Canary release of a KServe ML model

This tutorial shows how Iter8 can be used to release ML models hosted in a KServe environment using a canary rollout strategy. The user declaratively describes the desired application state at any given moment. An Iter8 `release` chart ensures that Iter8 can automatically respond to automatically deploy the application components and configure the necessary routing.

In a canary rollout, inference requests that match a particular pattern, for example those that have a particular header, are directed to the candidate version of the model. The remaining requests go to the primary, or initial, version of the model.

![Canary rollout](images/canary.png)

???+ warning "Before you begin"
    1. Ensure that you have the [`kubectl`](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs.
    2. Have access to a cluster running [KServe](https://kserve.github.io/website). You can create a [KServe Quickstart](https://kserve.github.io/website/0.11/get_started/#before-you-begin) environment as follows:
    ```shell
    curl -s "https://raw.githubusercontent.com/kserve/kserve/release-0.11/hack/quick_install.sh" | bash
    ```
<!-- Istio is installed as part of kserve install -->

## Install the Iter8 controller

--8<-- "docs/getting-started/install.md"

```shell
export IMG=kalantar/iter8:20231004-1030
export CHARTS=/Users/kalantar/projects/go.workspace/src/github.com/iter8-tools/iter8/charts
helm upgrade --install iter8 $CHARTS/controller \
--set image=$IMG --set logLevel=trace \
--set clusterScoped=true
```

## Deploy initial version

```shell
cat <<EOF | helm upgrade --install wisdom $CHARTS/release -f -
environment: kserve
application: 
  metadata:
    labels:
      app.kubernetes.io/name: wisdom
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
  strategy: canary
EOF
```

??? note "What happens?"
    - Because `environment` is set to `kserve`, the `InferenceService` `default/wisdom-0` is created. It will have label `iter8.tools/watch=true`.
        - The namespace `default` is inherited from the helm release namespace since it is not specified in either the version or in `application.metadata.namespace`.
        - The name `wisdom-0` is derived from the helm release name since it is not specified in either the version or in `application.metadata.name`. `-0` (the index of the version in `versions`) is appended to the base name.
        - Alternatively, an `inferenceServiceSpecification` could have been specified.
    - A `Service` of type `ExternalName` named `default/wisdom` pointing at `knative-local-gateway.istio-system` is created.
    - The routemap (`ConfigMap` `wisdom-routemap`) is created with 1 version and a single routing template.

Once the application components are ready, the Iter8 controller will trigger the creation of a `VirtualService` named `default/wisdom`. It will send all traffic sent to the service `wisdom` to the deployed version `wisdom-0`.

### Verify routing

You can send verify the routing configuration by inspecting the `VirtualService`:

```shell
kubectl get virtualservice wisdom -o yaml
```

You can also send requests:

=== "From within the cluster"
    1. Create a `sleep` pod in the cluster from which requests can be made:
    ```shell
    curl -s https://raw.githubusercontent.com/iter8-tools/docs/v0.17.3/samples/kserve-serving/sleep.sh | sh -
    ```

    2. Exec into the sleep pod:
    ```shell
    kubectl exec --stdin --tty "$(kubectl get pod --sort-by={metadata.creationTimestamp} -l app=sleep -o jsonpath={.items..metadata.name} | rev | cut -d' ' -f 1 | rev)" -c sleep -- /bin/sh
    ```

    3. To send requests without the header `traffic`:
    ```shell
    curl -H 'Content-Type: application/json' \
    http://wisdom.default -d @input.json -s -D -  \
    | grep -e HTTP -e app-version
    ```

    4. To send requests with the header `traffic: test`:
    ```shell
    curl -H 'Content-Type: application/json' \
    -H 'traffic: test' \
    http://wisdom.default -d @input.json -s -D -  \
    | grep -e HTTP -e app-version
    ```

=== "From outside the cluster"
    1. In a separate terminal, port-forward the ingress gateway:
    ```shell
    kubectl -n istio-system port-forward svc/istio-ingressgateway 8080:80
    ```

    2. Download the sample input:
    ```shell
    curl -sO https://raw.githubusercontent.com/iter8-tools/docs/v0.17.3/samples/kserve-serving/input.json
    ```

    3. To send requests without the header `traffic`:
    ```shell
    curl -H 'Content-Type: application/json' \
    -H 'Host: wisdom.default' \
    localhost:8080 -d @input.json -s -D - \
    | grep -e '^HTTP' -e app-version
    ```
    
    4. To send requests with the header `traffic: test`:
    ```shell
    curl -H 'Content-Type: application/json' \
    -H 'traffic: test' \
    localhost:8080 -d @input.json -s -D -  \
    | grep -e HTTP -e app-version
    ```

??? note "Sample output"
    The output identifies the success of the request and the version of the application that responds. For example:

    ```
    HTTP/1.1 200 OK
    app-version: wisdom-0
    ```

## Deploy candidate

```shell
cat <<EOF | helm upgrade --install wisdom $CHARTS/release -f -
environment: kserve
application: 
  metadata:
    labels:
      app.kubernetes.io/name: wisdom
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
  strategy: canary
EOF
```

??? note "What happens?"
    - Since the definition for the first version does not change, there is no change to the `InferenceService` named `default/wisdom-0`.
    - `InferenceService` named `default/wisdom-1` is deployed. It has label `iter8.tools/watch=true`.
    - The routemap (`ConfigMap` `wisdom-routemap`) is updated with 2 versions and an updated `routingTemplate`.

When the candidate version is ready, the Iter8 controller will trigger the reconfiguration of the `VirtualService`. Requests with the header `traffic` set to `true` will be sent to the candidate model. Still other requests will be sent to the primary model.

### Verify routing

You can send additional inference requests as described above. Those with header `traffic` set to `true` will be handled by the candidate model (`wisdom-1`) while all others will be handled by the primary version.

## Promote candidate

```shell
cat <<EOF | helm upgrade --install wisdom $CHARTS/release -f -
environment: kserve
application: 
  metadata:
    labels:
      app.kubernetes.io/name: wisdom
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
  strategy: canary
EOF
```

??? note "What happens?"
    - Since the definition for the first version has changed (label and `storageUri`), the `InferenceService` object is updated.
    - The `InferenceService` named `default/wisdom-1` is deleted because the second version has been removed.
    - The routemap (`ConfigMap` `wisdom-routemap`) is updated with 1 version and an updated `routingTemplate`.

In response to the changes in the application, the Iter8 controller will trigger the update of the `VirtualService` `wisdom` to send all traffic to the single (promoted) version.

### Verify Routing

You can verify the routing configuration by inspecting the `VirtualService` and/or by sending requests as described above. They will be handled by the primary version.

## Cleanup

Delete the application:

```shell
helm delete wisdom
```

Uninstall Iter8 controller:

--8<-- "docs/getting-started/uninstall.md"
