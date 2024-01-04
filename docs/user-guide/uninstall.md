---
template: main.html
---

## With `helm`

```shell
helm delete iter8
```

## With `kustomize`

Uninstall with one of the following, depending on whether Iter8 was installed and configured for a single namespace (namespace-scoped) or for the whole cluster (cluster-scoped).

=== "Namespace-scoped"
    ```shell
    kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/namespaceScoped?ref=v1.1.1'
    ```

=== "Cluster-scoped"
    ```shell
    kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/clusterScoped?ref=v1.1.1'
    ```
