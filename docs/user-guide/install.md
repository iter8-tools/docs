---
template: main.html
---

# Install options

## With `helm`

Iter8 can be installed and configured to watch resources either in a single namespace (namespace-scoped) or in the whole cluster (cluster-scoped). 

=== "Namespace-scoped"
    ```shell
    helm install --repo https://iter8-tools.github.io/iter8 --version 1.1 iter8 controller
    ```

=== "Cluster-scoped"
    ```shell
    helm install --repo https://iter8-tools.github.io/iter8 --version 1.1 iter8 controller \
    --set clusterScoped=true
    ```

To install Iter8 in a non-default namespace, use the `-n` option.

## With `kustomize`

Iter8 can be installed and configured to watch resources either in a single namespace (namespace-scoped) or in the whole cluster (cluster-scoped). 

=== "Namespace-scoped"
    ```shell
    kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/namespaceScoped?ref=v1.1.1'
    ```

=== "Cluster-scoped"
    ```shell
    kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/clusterScoped?ref=v1.1.1'
    ```

To install Iter8 in a non-default namespace, download the `kustomize` folder and modify the `namespace` field in the `kustomization.yaml` file.

## Install for production use

By default, Iter8 uses [BadgerDB](https://dgraph.io/docs/badger/) to store metrics from A/B/n and performance tests. BadgerDB is not suitable for production use. To install for production, use [Redis](metrics_store.md).

## Install using Rancher Desktop

Rancher does not support a `standard` storage class by default. Install the [local path provisioner](https://github.com/rancher/local-path-provisioner/):

```shell
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.26/deploy/local-path-storage.yaml
```

And set `metrics.badgerdb.storageClassName` when starting the Iter8 controller:

```shell
--set metrics.badgerdb.storageClassName=local-path
```

## Install on OpenDataHub

See [here](https://github.com/opendatahub-io-contrib/odh-contrib-manifests/tree/main/iter8).