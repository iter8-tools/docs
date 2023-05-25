=== "Helm"
    Install the Iter8 controller using `helm` as follows.

    ```shell
    helm install --repo https://iter8-tools.github.io/iter8 iter8-traffic traffic
    ```
    
=== "Kustomize"
    Install the Iter8 controller using `kustomize` as follows.

    === "namespace scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/traffic/namespaceScoped?ref=v0.14.8'
        ```

    === "cluster scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/traffic/clusterScoped?ref=v0.14.8'
        ```
