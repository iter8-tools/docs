=== "Helm"
    ```shell
    helm delete iter8
    ```
    
=== "Kustomize"
    === "namespace scoped"
    <!-- TODO: should these be bumped to v0.16? -->
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/namespaceScoped?ref=v0.15.3'
        ```

    === "cluster scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/clusterScoped?ref=v0.15.3'
        ```
