=== "Helm"
    Delete the Iter8 controller using `helm` as follows.

    ```shell
    helm delete iter8-traffic
    ```
    
=== "Kustomize"
    Delete the Iter8 controller using `kustomize` as follows.

    === "cluster scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/traffic/clusterScoped?ref=iter8-0.13.7'
        ```

    === "namespace scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/iter8.git/kustomize/traffic/namespaceScoped?ref=iter8-0.13.7'
        ```
