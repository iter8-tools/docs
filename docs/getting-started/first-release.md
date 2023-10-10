---
template: main.html
---

# Your first blue-green release

This tutorial shows how Iter8 can be used to release a basic Kubernetes application using a blue-green rollout strategy. The user declaratively describes the desired application state at any given moment. An Iter8 `release` chart deploys the application components and configures Iter8 to automatically configure routing to implement a the blue-green rollout.

In a blue-green rollout, a percentage of requests are directed to a candidate version of the model. This percentage can be changed over time.

![Blue-green rollout](../tutorials/images/blue-green.png)

In this tutorial, we use the Istio service mesh to distribute inference requests between different versions of a model.

???+ warning "Before you begin"
    1. Ensure that you have a Kubernetes cluster and the [`kubectl`](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs. You can create a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
    2. Install [Istio](https://istio.io). It suffices to install the [demo profile](https://istio.io/latest/docs/setup/getting-started/), for example by using: 
    ```shell
    istioctl install --set profile=demo -y
    ```

## Install the Iter8 controller

--8<-- "docs/getting-started/install.md"

```shell
export IMG=kalantar/iter8:20241010-0930
export CHARTS=/Users/kalantar/projects/go.workspace/src/github.com/iter8-tools/iter8/charts
helm upgrade --install iter8 $CHARTS/controller \
--set image=$IMG --set logLevel=trace \
--set clusterScoped=true
```

## Deploy initial version

Deploy the initial version of the application using the Iter8 `release` chart by identifying the environment into which it should be deployed, a list of the versions to be deployed (just one), and the rollout strategy to be used:

```shell
cat <<EOF | helm upgrade --install httpbin $CHARTS/release -f -
environment: deployment-istio
application: 
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    image: kennethreitz/httpbin
    port: 80
  strategy: blue-green
EOF
```

??? note "What happens?"
    Because `environment` is set to `deployment-istio`, a `Deployment` and a `Service` object are created.
        - The namespace `default` is inherited from the helm release namespace since it is not specified in the version or in `application.metadata`.
        - The name `httpbin-0` is derived from the helm release name since it is not specified in the version or in `application.metadata`. The names is derived by appending the index of the version in the list of versions; `-0` in this case.
        - Alternatively, a `deploymentSpecification` and/or a `serviceSpecification` could have been specified.

    To support routing, a `Service` (of type `ExternalName`) named `default/httpbin` pointing at `istio-ingressgateway.istio-system` is deployed. The name is the helm release name since it not specified in `application.metadata`. Further, an Iter8 [routemap](../user-guide/topics/routemap.md) is created. Finally, to support the blue-green rollout, a `ConfigMap` (`httpbin-0-weight-config`) is created to be used to manage the proportion of traffic sent to this version.
     - **Assumes**`Gateway` named `gateway` exists and that it is configured for traffic to  `httpbin.default.svc.cluster.local`.

Once the application components are ready, the Iter8 controller atuomatically configures the routing by creating an Istio `VirtualService`. It is configured to route all traffic to the the only deployed version, `httpbin-0`.

### Verify routing

You can verify the routing configuration by inspecting the `VirtualService`:

```shell
kubectl get virtualservice httpbin -o yaml
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
    The output identifies the success of the request (the HTTP return code) and the version of the application that responded (the `app-version` header). For example:

    ```
    HTTP/1.1 200 OK
    app-version: httpbin-0
    ```

## Deploy candidate

A candidate can deployed by simply adding a second version to the list of versions composing the application:

```shell
cat <<EOF | helm upgrade --install httpbin $CHARTS/release -f -
environment: deployment-istio
application: 
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    image: kennethreitz/httpbin
    port: 80
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    image: kennethreitz/httpbin
    port: 80
  strategy: blue-green
EOF
```

??? note "About the candidate version"
    In this tutorial, the candidate image is the same as the one for the primary version. In a real world example, it would be different. The version label (`app.kubernetes.io/version`) can be used to distinguish between versions.

When the second version is deployed and ready, the Iter8 controller automatically reconfigures the routing; the `VirtualService` is updated to distribute traffic between versions based on the weights.

### Verify routing

You can verify the routing configuration by inspecting the `VirtualService` and/or by sending requests as described above. Requests will now be handled equally by both versions.

## Modify weights (optional)

To modify the request distribution between the versions, add a `weight` to each version:

```shell
cat <<EOF | helm upgrade --install httpbin $CHARTS/release -f -
environment: deployment-istio
application: 
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    image: kennethreitz/httpbin
    port: 80
    weight: 30
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    image: kennethreitz/httpbin
    port: 80
    weight: 70
  strategy: blue-green
EOF
```

Iter8 automatically reconfigures the routing (modifies the `VirtualService`) to distribute traffic between the versions based on the new weights.

### Verify routing

You can verify the routing configuration by inspecting the `VirtualService` and/or by sending requests as described above. Seventy percent of requests will now be handled by the candidate version; the remaining thirty percent by the the primary version.

## Promote candidate

The candidate can be promoted by redefinig the primary version and removing the candidate:

```shell
cat <<EOF | helm upgrade --install httpbin $CHARTS/release -f -
environment: deployment-istio
application: 
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    image: kennethreitz/httpbin
    port: 80
  strategy: blue-green
EOF
```
??? note "What is different?"
    The version label (`app.kubernetes.io/version`) was updated. In a real world example, the image would also have been updated.

Once the application components are ready, the Iter8 controller will automatically reconfigure the routing to send all traffic to the single version.

### Verify routing

You can verify the routing configuration by inspecting the `VirtualService` and/or by sending requests as described above. They will all be handled by the primary version.

## Cleanup

Delete the application and its routing configuration:

```shell
helm delete httpbin
```

Uninstall Iter8 controller:

--8<-- "docs/getting-started/uninstall.md"
