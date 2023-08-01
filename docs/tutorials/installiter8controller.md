=== "Helm"
    Install the Iter8 controller using `helm` as follows.

    === "namespace scoped"
        ```shell
        helm install --repo https://iter8-tools.github.io/iter8 iter8 traffic
        ```

    === "cluster scoped"
        ```shell
        helm install --repo https://iter8-tools.github.io/iter8 iter8 traffic \
        --set clusterScoped=true
        ```
    
=== "Kustomize"
    Install the Iter8 controller using `kustomize` as follows.

    === "namespace scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/namespaceScoped?ref=v0.15.3'
        ```

    === "cluster scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/clusterScoped?ref=v0.15.3'
        ```
