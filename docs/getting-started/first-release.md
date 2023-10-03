# Your first blue-green release

This tutorial shows how Iter8 can be used to release a basic Kubernetes application using a blue-green rollout strategy. The user declaratively describes the desired application state at a given moment. An Iter8 `release` chart ensures that Iter8 can automatically respond to automatically deploy the application components and configure the necessary routing.

??? note "Motivation"
    Today, the `routing-template` helm chart provides a set of templates that enables a user to easily implement common routing patterns for some kinds of applications. These templates rely on the specification of an `action` making the approach imperative rather than declarative. Furthermore, this approach works only after the user has deployed application versions.

    We propose extending the concepts in the `routing-template` chart to a `release` chart. The configutation (`values.yaml`) would now include an application description allowing the chart to deploy the application and be fullly declarative.

???+ warning "Before you begin"
    1. Ensure that you have a Kubernetes cluster and the [`kubectl`](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs. You can create a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
    2. Install [Istio](https://istio.io). It suffices to install the [demo profile](https://istio.io/latest/docs/setup/getting-started/), for example by using: `istioctl install --set profile=demo -y`

## Install the Iter8 controller

--8<-- "docs/getting-started/install.md"

## Deploy initial version

```shell
cat <<EOF | helm upgrade --install httpbin charts/release -f -
environment: deployment-istio
# gateway: my-gateway
application: 
  # metadata:
  #   name: httpbin   # default is .Release.Name
  #   namespace: default  # default is .Release.Namespace
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    # image: kennethreitz/httpbin:v0
    image: kennethreitz/httpbin
    port: 80
    # deploymentSpecification:
    # serviceSpecification:
  strategy: blue-green
EOF
```

???+ note "What happens?"
    _Application components_

    - Because `environment` is set to `deployment-istio`, the following application components are created:
        - `Deployment` `default/httpbin-0` using image `kennethreitz/httpbin` listening on port `80` will be deployed. It will have label `iter8.tools/watch=true`.
        - `Service` `default/httpbin-0` on port `80` will be deployed. It has label `iter8.tools/watch=true`.
    - The namespace `default` is inherited from `application.metadata.namespace` since it is not specified in the version
    - The name `httpbin-0` is derived from `application.metadata.name` by adding `-0` (index of the version in `versions`) since it is not specified in the version.
    - Alternatively, a `deploymentSpecification` and/or a `serviceSpecification` could have been specified.

    _Routing components_

    - `Service` of type `ExternalName` named `default/httpbin` pointing at `istio-ingressgateway.istio-system` is deployed.
    - **Assumes**`Gateway` named `gateway` exists and that it is configured for traffic to  `httpbin.default.svc.cluster.local`.

    _Iter8 components_

    - The routemap (`ConfigMap` `httpbin-routemap`) is created with 1 version and a single routing template.
    - `ConfigMap` `httpbin-0` (used to manage the proportion of traffic sent to the first version) is created with label `iter8.tools/weight: 100`. It has label `iter8.tools/watch=true`.

    _What else happens?_

    Once the application components are ready, the Iter8 controller will trigger the routing template defined in the routemap. As a consequence, a `VirtualService` named `default/httpbin` will be created. It will send all traffic sent to the service `httpbin` to the deployed version `httpbin-0`.

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
    curl httpbin.default -s -D - | grep -e '^HTTP' -e app-version
    ```

=== "From outside the cluster"
    1. In a separate terminal, port-forward the ingress gateway:
    ```shell
    kubectl -n istio-system port-forward svc/istio-ingressgateway 8080:80
    ```

    2. Send requests:
    ```shell
    curl -H 'Host: httpbin.default' localhost:8080 -s -D - | grep -e '^HTTP' -e app-version
    ```

??? note "Sample output"
    The output identifies the success of the request and the version of the application that responds. For example:

    ```
    HTTP/1.1 200 OK
    app-version: httpbin-0
    ```

## Deploy candidate

```shell
cat <<EOF | helm upgrade --install httpbin charts/release -f -
environment: deployment-istio
application: 
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    # image: kennethreitz/httpbin:v0
    image: kennethreitz/httpbin
    port: 80
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    # image: kennethreitz/httpbin:v1
    image: kennethreitz/httpbin
    port: 80
  strategy: blue-green
EOF
```

???+ note "What happens?"
    _Application components_

    - Since the definition for the first version does not change, no changes to the `Deployment` or `Service` occur.
    - `Deployment` `default/httpbin-1` using image `kennethreitz/httpbin` listening on port `80` is deployed. Has label `iter8.tools/watch=true`.
    - `Service` `default/httpbin-1` on port `80` is deployed. Has label `iter8.tools/watch=true`.

    _Routing components_

    - no changes

    _Iter8 components_

    - The routemap (`ConfigMap` `httpbin-routemap`) is updated with 2 versions and an updated `routingTemplate`.
    - `ConfigMap` `httpbin-0` (used to manage the proportion of traffic sent to the first version) is updated (label `iter8.tools/weight` is changed to 50),
    - `ConfigMap` `httpbin-1` (used to manage the proportion of traffic sent to the second version) is created with label `iter8.tools/weight=50`. It has label `iter8.tools/watch=true`.

    _What else happens?_

    Once the application components are ready, the Iter8 controller will trigger the routing template defined in the routemap. As a consequence, the `VirtualService` `httpbin` will be updated to distribute traffic between versions based on the weights.

## Modify weights (optional)

```shell
cat <<EOF | helm upgrade --install httpbin charts/release -f -
environment: deployment-istio
application: 
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    # image: kennethreitz/httpbin:v0
    image: kennethreitz/httpbin
    port: 80
    weight: 30
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    # image: kennethreitz/httpbin:v1
    image: kennethreitz/httpbin
    port: 80
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

    - `ConfigMap` `httpbin-0` (used to manage the proportion of traffic sent to the first version) is updated (label `iter8.tools/weight` changes).
    - `ConfigMap` `httpbin-1` (used to manage the proportion of traffic sent to the second version) is updated (label `iter8.tools/weight` changes).

    _What else happens?_

    Since the configmaps used to manage traffic distribution are modified, the Iter8 controller will trigger the routing template defined in the routemap. As a consequence, the `VirtualService` `httpbin` will be updated to distribute traffic between versions based on the new weights.

## Promote candidate

```shell
cat <<EOF | helm upgrade --install httpbin charts/release -f -
environment: deployment-istio
application: 
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    # image: kennethreitz/httpbin:v1
    image: kennethreitz/httpbin
    port: 80
  strategy: blue-green
EOF
```

???+ note "What happens?"
    _Application components_

    - Since the definition for the first version has changed (image and label), the `Deployment` object is updated. In this case, there are no changes to the `Service`.
    - `Deployment` and `Service` named `default/httpbin-1` are deleted because the second version has been removed.

    _Routing components_

    - no changes

    _Iter8 components_

    - The routemap (`ConfigMap` `httpbin-routemap`) is updated with 1 version and an updated `routingTemplate`.
    - `ConfigMap` `httpbin-0` (used to manage the proportion of traffic sent to the first version) is updated (label `iter8.tools/weight` is set to 100).
    - `ConfigMap` `httpbin-1` (used to manage the proportion of traffic sent to the second version) is deleted.

    _What else happens?_

    Once the application components are ready, the Iter8 controller will trigger the routing template defined in the routemap. As a consequence, the `VirtualService` `httpbin` will be updated to send all traffic to the single version.

## Cleanup

```shell
helm delete httpbin
```