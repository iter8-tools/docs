---
template: main.html
---

## Uninstall with `helm`

If installed with `helm`, uninstall with:

```shell
helm delete iter8
```

## Uninstall with `kustomize`

If installed with `kustomize`, uninstall with one of the following:

=== "Namespace scoped"
    ```shell
    kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/namespaceScoped?ref=v0.18.3'
    ```

=== "Cluster scoped"
    ```shell
    kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/clusterScoped?ref=v0.18.3'
    ```
