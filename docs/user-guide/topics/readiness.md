---
template: main.html
---

# Checking readiness of Kubernetes resources

When Iter8 experiment rely on Kubernetes resources, the readiness of those resources can be checked first.  

For example, to check the readiness of a `Deployment` resource:

```shell
iter8 k -c load-test-http \
--set url=http://httpbin.default \
--set ready.deploy=httpbin
```

Currently, Iter8 provides native support for the following resource types:

- `Deployment`: resource is ready if it exists and has the status condition `Available` set to `True`
- `Service`: resource is ready if it exists

The option `ready.timeout` can be used to specify the maximum time that the experiment should wait for the readiness condition to be satisfied.

## Readiness checking other resources

Iter8 can be extended to support readiness checking of additional resource types beyond deployments and services. Please consider submitting a pull request with new resource types.

In brief, the current implementation uses two templates [`task.ready`](https://github.com/iter8-tools/iter8/blob/master/charts/iter8lib/templates/_task.ready.tpl#L24) and [`task.ready.rbac`](https://github.com/iter8-tools/iter8/blob/master/charts/iter8lib/templates/_task-ready-rbac.tpl#L1). To support a new type:

- Modify [`task.ready`](https://github.com/iter8-tools/iter8/blob/master/charts/iter8lib/templates/_task.ready.tpl#L24) to define the group/version/resource and, optionally, a condition that should be checked. For example, to add a readiness check for a Knative `Service`, the following might be added:

    ```yaml
    {{- if .Values.ready.ksrv }}
    # task: determine if Knative Service is Ready
    - task: k8s-object-ready
      with:
        name: {{ .Values.ready.ksrv | quote }}
        group: serving.knative.dev
        version: v1
        resource: services
        condition: Ready
    {{- include "task.ready.tn" . | indent 4 }}
    {{ end }}
    ```
    <!-- https://github.com/knative/specs/blob/main/specs/serving/knative-api-specification-1.0.md#service-1 -->

- Add a new `apiGroup` to [`task.ready.rbac`](https://github.com/iter8-tools/iter8/blob/master/charts/iter8lib/templates/_task-ready-rbac.tpl#L1) template for `Role`. For example:

    ```yaml
    {{- if .Values.ready.ksrv }}
    - apiGroups: ["apps"]
      resourceNames: [{{ .Values.ready.ksrv | quote }}]
      resources: ["services"]
      verbs: ["get"]
    {{- end }}
    ```
