---
template: main.html
---

# Iter8 controller extensions

To support automated traffic routing (re)configuration and A/B/n testing, the Iter8  controller watches resources specified to be part of the application and checks their readiness. In this way, Iter8 is able to identify when new versions of an application become available.

By default, Iter8 supports watching a limited set of resource types. These resource types are specified in the default controller chart [values.yaml](https://github.com/iter8-tools/iter8/blob/v0.16.6/charts/controller/values.yaml) file. They are:

- Kubernetes `Service`, `ConfigMap`, and `Deployment` resource types
- KServe `InferenceService` resource types
- Istio `VirtualService` resource types

The set of watched resource types can be extended by extending the list in the `valuses.yaml` file used when deploying the controller. For each new resource type, specify the Kubernetes group, version and resource. Furthermore, a list of status conditions and their expected values can be specified. To check for readiness, Iter8 ensures that each of the specified conditions matches the specified value. For example, the configuration for the `Deploymet` resouce type is:

```yaml
deploy:
    Group: apps
    Version: v1
    Resource: deployments
    conditions:
    - name: Available
      status: "True"
```