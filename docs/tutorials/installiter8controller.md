=== "Helm"
    === "namespace scoped"
        ```shell
        helm install --repo https://iter8-tools.github.io/iter8 --version 0.1.11 iter8 controller
        ```

    === "cluster scoped"
        ```shell
        helm install --repo https://iter8-tools.github.io/iter8 --version 0.1.11 iter8 controller \
        --set clusterScoped=true
        ```
    
=== "Kustomize"
    === "namespace scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/namespaceScoped?ref=v0.16.1'
        ```

    === "cluster scoped"
        ```shell
        kubectl apply -k 'https://github.com/iter8-tools/iter8.git/kustomize/controller/clusterScoped?ref=v0.16.1'
        ```
