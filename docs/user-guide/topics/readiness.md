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

## Extending to new resource kinds

1. Download chart(s):

    ```shell
    iter8 hub -c load-test-http
    ```

2. Add an entry for the new resource type in `charts/iter8lib/templates/_task.ready.tpl`.
Identify the group/version/resource and identify the condition that should be checked. For example, to add a readiness check for a Knative `Service`:

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

3. Add new `apiGroup` to rbac `Role` definition, `charts/iter8lib/templates/_task-ready-rbac.tpl`:

    ```yaml
    {{- if .Values.ready.ksrv }}
    - apiGroups: ["apps"]
      resourceNames: [{{ .Values.ready.ksrv | quote }}]
      resources: ["services"]
      verbs: ["get"]
    {{- end }}
    ```

4. Use the locally modified chart with the `--noDownload` option:

    ```shell
    iter8 k launch --noDownload -c load-test-http \
    --set ready.ksrv=myservice \
    ...
    ```
