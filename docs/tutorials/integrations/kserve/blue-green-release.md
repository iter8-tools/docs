---
template: main.html
---

# Your first blue-green release

This tutorial shows how Iter8 can be used to release ML models hosted in a KServe environment using a blue-green rollout strategy. The user declaratively describes the desired application state at any given moment. An Iter8 `release` chart ensures that Iter8 can automatically respond to automatically deploy the application components and configure the necessary routing.

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
  versions:
  - metadata:
      labels:
        app.kubernetes.io/name: wisdom
        app.kubernetes.io/version: v0
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
    # inferenceServiceSpecification:
  strategy: blue-green
EOF
```

???+ note "What happens?"
    _Application components_

    - Because `environment` is set to `kserve`, the following application components are created:
        - `InferenceService` `default/wisdom-0`. It will have label `iter8.tools/watch=true`.
    - The namespace `default` is inherited from the helm release namespace since it is not specified in either the version or in `application.metadata.namespace`.
    - The name `wisdom-0` is derived from the helm release name since it is not specified in either the version or in `application.metadata.name`. `-0` (the index of the version in `versions`) is appended to the base name.
    - Alternatively, an `inferenceServiceSpecification` could have been specified.

    _Routing components_

    - `Service` of type `ExternalName` named `default/wisdom` pointing at `knative-local-gateway.istio-system` is deployed.

    _Iter8 components_

    - The routemap (`ConfigMap` `wisdom-routemap`) is created with 1 version and a single routing template.
    - `ConfigMap` `wisdom-0-weight-config` (used to manage the proportion of traffic sent to the first version) is created with annotation `iter8.tools/weight`. It has label `iter8.tools/watch=true`.

    _What else happens?_

    Once the application components are ready, the Iter8 controller will trigger the routing template defined in the routemap. As a consequence, a `VirtualService` named `default/wisdom` will be created. It will send all traffic sent to the service `wisdom` to the deployed version `wisdom-0`.

## Sending requests

To send requests to the application:

=== "From within the cluster"
    1. Create a `sleep` pod in the cluster from which requests can be made:
    ```shell
    curl -s https://raw.githubusercontent.com/iter8-tools/docs/v0.17.3/samples/kserve-serving/sleep.sh | sh -
    ```

    2. Exec into the sleep pod:
    ```shell
    kubectl exec --stdin --tty "$(kubectl get pod --sort-by={metadata.creationTimestamp} -l app=sleep -o jsonpath={.items..metadata.name} | rev | cut -d' ' -f 1 | rev)" -c sleep -- /bin/sh
    ```

    3. Send requests:
    ```shell
    curl -H 'Content-Type: application/json' \
    http://wisdom.default -d @input.json -s -D -  \
    | grep -e HTTP -e app-version
    ```

=== "From outside the cluster"
    1. In a separate terminal, port-forward the ingress gateway:
    ```shell
    kubectl -n istio-system port-forward svc/istio-ingressgateway 8080:80
    ```

    2. Send requests:
    ```shell
    curl -H 'Content-Type: application/json' \
    -H 'Host: wisdom.default' \
    localhost:8080 -d @input.json -s -D - \
    | grep -e '^HTTP' -e app-version
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
  versions:
  - metadata:
      labels:
        app.kubernetes.io/name: wisdom
        app.kubernetes.io/version: v0
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
  - metadata:
      labels:
        app.kubernetes.io/name: wisdom
        app.kubernetes.io/version: v1
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
  strategy: blue-green
EOF
```

???+ note "What happens?"
    _Application components_

    - Since the definition for the first version does not change, there is no change to the `InferenceService` named `default/wisdom-0`.
    - `InferenceService` named `default/wisdom-1` is deployed. It has label `iter8.tools/watch=true`.

    _Routing components_

    - no changes

    _Iter8 components_

    - The routemap (`ConfigMap` `wisdom-routemap`) is updated with 2 versions and an updated `routingTemplate`.
    - `ConfigMap` `wisdom-0-weight-config` (used to manage the proportion of traffic sent to the first version) is updated (annotation `iter8.tools/weight` is updated),
    - `ConfigMap` `wisdom-1-weight-config` (used to manage the proportion of traffic sent to the second version) is created with annnotation `iter8.tools/weight`. It has label `iter8.tools/watch=true`.

    _What else happens?_

    Once the application components are ready, the Iter8 controller will trigger the routing template defined in the routemap. As a consequence, the `VirtualService` `wisdom` will be updated to distribute traffic between versions based on the weights.

## Modify weights (optional)

```shell
cat <<EOF | helm upgrade --install wisdom $CHARTS/release -f -
environment: kserve
application: 
  versions:
  - metadata:
      labels:
        app.kubernetes.io/name: wisdom
        app.kubernetes.io/version: v0
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
    weight: 30
  - metadata:
      labels:
        app.kubernetes.io/name: wisdom
        app.kubernetes.io/version: v1
    modelFormat: sklearn
    runtime: kserve-mlserver
    storageUri: "gs://seldon-models/sklearn/mms/lr_model"
    weight: 70
  strategy: blue-green
EOF
```

???+ note "What happens?"
    _Application components_

    - no changes

    _Routing components_

    - no changes

    _Iter8 components_

    - `ConfigMap` `wisdom-0-weight-config` (used to manage the proportion of traffic sent to the first version) is updated (annotation `iter8.tools/weight` changes).
    - `ConfigMap` `wisdom-1-weight-config` (used to manage the proportion of traffic sent to the second version) is updated (annotation `iter8.tools/weight` changes).

    _What else happens?_

    Since the configmaps used to manage traffic distribution are modified, the Iter8 controller will trigger the routing template defined in the routemap. As a consequence, the `VirtualService` `wisdom` will be updated to distribute traffic between versions based on the new weights.

## Promote candidate

ç
???+ note "What happens?"
    _Application components_

    - Since the definition for the first version has changed (label and `storageUri`), the `InferenceService` object is updated.
    - The `InferenceService` named `default/wisdom-1` is deleted because the second version has been removed.

    _Routing components_

    - no changes

    _Iter8 components_

    - The routemap (`ConfigMap` `wisdom-routemap`) is updated with 1 version and an updated `routingTemplate`.
    - `ConfigMap` `wisdom-0-weight-config` (used to manage the proportion of traffic sent to the first version) is updated (annotation `iter8.tools/weight` is set to 100).
    - `ConfigMap` `wisdom-1-weight` (used to manage the proportion of traffic sent to the second version) is deleted.

    _What else happens?_

    Once the application components are ready, the Iter8 controller will trigger the routing template defined in the routemap. As a consequence, the `VirtualService` `wisdom` will be updated to send all traffic to the single (promoted) version.

## Cleanup

```shell
helm delete wisdom
```

Uninstall Iter8 controller:

--8<-- "docs/getting-started/uninstall.md"
