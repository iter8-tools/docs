=== "Helm"
    Delete the Iter8 controller using `helm` as follows.

    ```shell
    helm delete iter8-traffic
    ```
    
=== "Kustomize"
    Delete the Iter8 controller using `kustomize` as follows.

    === "cluster scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/hub.git/kustomize/traffic/clusterScoped?ref=traffic-0.1.3'
        ```

    === "namespace scoped"
        ```shell
        kubectl delete -k 'https://github.com/iter8-tools/hub.git/kustomize/traffic/namespaceScoped?ref=traffic-0.1.3'
        ```
