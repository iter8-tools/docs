=== "Helm"
    ```shell
    helm delete iter8
    ```
    
=== "Kustomize"
    === "namespace scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/namespaceScoped?ref=v0.16.0'
        ```

    === "cluster scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/clusterScoped?ref=v0.16.0'
        ```
