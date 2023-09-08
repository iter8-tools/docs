=== "Helm"
    ```shell
    helm delete iter8
    ```
    
=== "Kustomize"
    === "Namespace scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/namespaceScoped?ref=v0.17.1'
        ```

    === "Cluster scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/clusterScoped?ref=v0.17.1'
        ```
