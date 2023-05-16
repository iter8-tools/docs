=== "Helm"
    Install the Iter8 controller using `helm` as follows.

    ```shell
    helm install --repo https://iter8-tools.github.io/iter8 iter8-traffic traffic
    ```
    
=== "Kustomize"
    Install the Iter8 controller using `kustomize` as follows.

    === "cluster scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/hub.git/kustomize/traffic/clusterScoped?ref=traffic-templates-0.1.1'
        ```

    === "namespace scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/hub.git/kustomize/traffic/namespaceScoped?ref=traffic-templates-0.1.1'
        ```
