---
template: main.html
---

# ready

Check if a Kubernetes object exists and is ready.

## Usage example
In the following example, the `ready` task checks if a deployment named `httpbin-prod` exists and its [availability condition](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) is set to true, and a service named `httpbin` exists.
```shell
iter8 k launch \
--set "tasks={ready,http}" \
--set ready.deploy=httpbin-prod \
--set ready.service=httpbin \
--set http.url=http://httpbin.default/get \
--set runner=job
```

## Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| deploy  | string | Name of a Kubernetes deployment. The task checks if the deployment exists and its `available` condition is set to true. |
| service | string | Name of a Kubernetes service. The task checks if the service exists. |
| timeout | string | Timeout for readiness check to succeed. Default value is `60s`. |
| namespace | string | The namespace under which to look for the Kubernetes objects. For experiments that run inside a Kubernetes cluster, the default value of this field is the [namespace of the Iter8 experiment](../topics/group.md); for experiments that run in the local environment, it is the [`default`](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) namespace. |


## Extensions

Iter8 can be easily extended to support readiness checks for any type of Kubernetes object (including objects with custom resource types). Please consider submitting a pull request for such extensions. Readiness checking in Iter8 involves two templates, namely, [`task.ready`](https://github.com/iter8-tools/iter8/blob/master/charts/iter8/templates/_task-ready.tpl) and [`k.role`](https://github.com/iter8-tools/iter8/blob/master/charts/iter8/templates/_k-role.tpl). Extending the readiness checks to new resource types involves modifying these templates.

### Example

Consider an extension that will enable experiment authors to define readiness check for [Knative services](https://knative.dev/docs/serving) as follows. The ready task should succeed if the Knative service named `httpbin` exists, and has its `Ready` condition set to true.

```shell
iter8 k launch \
--set "tasks={ready,http}" \
--set ready.ksvc=httpbin \
--set http.url=http://httpbin.default/get \
--set runner=job
```

The following changes to the `task.ready` and `k.role` templates will accomplish the above.

=== "task.ready"
    Define the group/version/resource (GVR) and the condition that should be checked for a Knative `Service`.

    ```yaml linenums="1"
    {{- if .Values.ready.ksvc }}
    # task: determine if Knative Service exists and is ready
    - task: ready
      with:
        name: {{ .Values.ready.ksvc | quote }}
        group: serving.knative.dev
        version: v1
        resource: services
        condition: Ready
    {{- if $namespace }}
        namespace: {{ $namespace }}
    {{- end }}
    {{- if .Values.ready.timeout }}
        timeout: {{ .Values.ready.timeout }}
    {{- end }}
    {{- end }}
    ```

=== "k.role"
    Add the Knative `apiGroup` to the role named `{{ .Release.Name }}-ready`.

    ```yaml linenums="1"
    {{- if .Values.ready.ksvc }}
    - apiGroups: ["apps"]
      resourceNames: [{{ .Values.ready.ksvc | quote }}]
      resources: ["services"]
      verbs: ["get"]
    {{- end }}
    ```