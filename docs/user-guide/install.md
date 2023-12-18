---
template: main.html
---

# Install options

By default, Iter8 uses [BadgerDB](https://dgraph.io/docs/badger/) to store metrics from A/B/n and performance tests. BadgerDB is not suitable for production use. To install for production, use another database. Database options are listed [here](metrics_store.md).

## Install with `helm`

Iter8 can be installed and configured to watch resources either in a single namespace (namespace-scoped) or in the whole cluster (cluster-scoped). 

=== "Namespace-scoped"
    ```shell
    helm install --repo https://iter8-tools.github.io/iter8 --version 0.18 iter8 controller
    ```

=== "Cluster-scoped"
    ```shell
    helm install --repo https://iter8-tools.github.io/iter8 --version 0.18 iter8 controller \
    --set clusterScoped=true
    ```

To install Iter8 in a non-default namespace, use the `-n` option.

## Install with `kustomize`

Iter8 can be installed and configured to watch resources either in a single namespace (namespace-scoped) or in the whole cluster (cluster-scoped). 

=== "Namespace-scoped"
    ```shell
    kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/namespaceScoped?ref=v0.18.3'
    ```

=== "Cluster-scoped"
    ```shell
    kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/clusterScoped?ref=v0.18.3'
    ```

To install Iter8 in a non-default namespace, download the `kustomize` folder and modify the `namespace` field in the `kustomization.yaml` file.

## Install on OpenDataHub

See [https://github.com/opendatahub-io-contrib/odh-contrib-manifests/tree/main/iter8](https://github.com/opendatahub-io-contrib/odh-contrib-manifests/tree/main/iter8)