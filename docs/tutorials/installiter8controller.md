=== "Helm"
    === "Namespace scoped"
        ```shell
        helm install --repo https://iter8-tools.github.io/iter8 --version 0.1.12 iter8 controller
        ```

    === "Cluster scoped"
        ```shell
        helm install --repo https://iter8-tools.github.io/iter8 --version 0.1.12 iter8 controller \
        --set clusterScoped=true
        ```
    
=== "Kustomize"
    === "Namespace scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/namespaceScoped?ref=v0.17.1'
        ```

    === "Cluster scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/clusterScoped?ref=v0.17.1'
        ```
