=== "Helm"
    Delete the Iter8 controller using `helm` as follows.

    ```shell
    helm delete iter8
    ```
    
=== "Kustomize"
    Delete the Iter8 controller using `kustomize` as follows.

    === "namespace scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/namespaceScoped?ref=v0.15.3'
        ```

    === "cluster scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/clusterScoped?ref=v0.15.3'
        ```
