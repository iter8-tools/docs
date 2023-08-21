=== "Helm"
    === "namespace scoped"
        ```shell
        helm install --repo https://iter8-tools.github.io/iter8 iter8 controller
        ```

    === "cluster scoped"
        ```shell
        helm install --repo https://iter8-tools.github.io/iter8 iter8 controller \
        --set clusterScoped=true
        ```
    
=== "Kustomize"
    === "namespace scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/namespaceScoped?ref=v0.15.3'
        ```

    === "cluster scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/iter8/clusterScoped?ref=v0.15.3'
        ```
