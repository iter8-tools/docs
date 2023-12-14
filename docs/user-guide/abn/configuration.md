---
template: main.html
---

# Configuration

The Iter8 controller watches for changes to the application enabling it to maintain a list of available versions. To configure the controller to be aware of the composition of application versions, create a [routemap](../routemap.md) for the application (no `routingTemplate` is needed). This is a Kubernetes `ConfigMap` such as:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: backend
  labels:
    app.kubernetes.io/managed-by: iter8
    iter8.tools/kind: routemap
    iter8.tools/version: "v0.18"
immutable: true
data:
  strSpec: |
    versions:
    - resources:
      - gvrShort: svc
        name: backend
        namespace: default
      - gvrShort: deploy
        name: backend
        namespace: default
    - resources:
      - gvrShort: svc
        name: backend-candidate-1
        namespace: default
      - gvrShort: deploy
        name: backend-candidate-1
        namespace: default
```

This `ConfigMap` describes an application `backend`. It identifies two versions of the application. The first is comprised of a Kubernetes `Deployment` and a `Service` object both named `backend` in the `default` namespace.  The second is comprised of the same resource types named `backend-candidate-1` in the same namespace.
