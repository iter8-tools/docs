---
template: main.html
---

# Using new resource types

The Iter8 performance test task [`ready`](tasks/ready.md) ensures that an object exists and is ready. To use this task with new resource types, including CRDs, add the new resource type to the list of known types defined in the default [`values.yaml` file](https://github.com/iter8-tools/iter8/blob/v1.1.1/charts/iter8/values.yaml) for the chart.  Alternatively, the new type can be specified at run time with the `--set` option.

### Example

To include a Knative service as part of a version definition, add the following to the map of `resourceTypes` in the [`values.yaml`](https://github.com/iter8-tools/iter8/blob/v1.1.1/charts/iter8/values.yaml) file used to configure the controller. The addition identifies the Kubernetes group, version, and resource (GVR) and the status condition that should be checked for readiness.

```yaml
ksvc:
    Group: serving.knative.dev
    Version: v1
    Resource: services
    conditions:
    - Ready
```

Alternatively, to set the values at run time:

```shell
--set resourceTypes.ksvc.Group=serving.knative.dev \
--set resourceTypes.ksvc.Version=v1 \
--set resourceTypes.Resource=services \
--set "resourceTypes.conditions[0]=Ready"
```