---
template: main.html
---

# Iter8 controller extensions

Iter8 can be easily extended to watch any type of Kubernetes object (including objects with custom resource definitions) as part of a application version. 

For example, to include a Knative service as part of an version definition, add the following to the map of `resourceTypes` in the [`values.yaml`](https://raw.githubusercontent.com/iter8-tools/iter8/v0.17.1/charts/controller/values.yaml) file used to configure the controller. The addition identifies the Kubernetes group, version, and resource (GVR) and the status condition that should be checked for readiness.

```yaml
ksvc:
    Group: serving.knative.dev
    Version: v1
    Resource: services
    conditions:
    - name: Ready
      status: "True"
```
