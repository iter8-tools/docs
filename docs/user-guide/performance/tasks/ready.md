---
template: main.html
---

# ready

Check if a Kubernetes object exists and is ready.

## Usage example

In the following example, the `ready` task checks if a deployment named `httpbin-prod` exists and its [availability condition](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) is set to true, and a service named `httpbin` exists.
```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 1.1 httpbin-test iter8 \
--set "tasks={ready,http}" \
--set ready.deploy=httpbin-prod \
--set ready.service=httpbin \
--set http.url=http://httpbin.default/get
```

## Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| deploy  | string | Name of a Kubernetes deployment. The task checks if the deployment exists and its `Available` condition is set to true. |
| service | string | Name of a Kubernetes service. The task checks if the service exists. |
| ksvc | string | Name of a Knative service. The task checks if the service exists and its `Ready` condition is set to true. |
| timeout | string | Timeout for readiness check to succeed. Default value is `60s`. |
| namespace | string | The namespace under which to look for the Kubernetes objects. |


## Extensions

Iter8 can be easily extended to support readiness checks for any type of Kubernetes object (including objects with custom resource types). To do so, add the new resource type to the list of known types defined in the default [`values.yaml` file](https://github.com/iter8-tools/iter8/blob/v1.1.1/charts/iter8/values.yaml) for the chart.

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
