{{- define "env.knative-istio.blue-green" }}

{{- /* prepare versions for simpler processing */}}
{{- $versions := include "normalize.versions.ksvc" . | mustFromJson }}

{{- /* weight-config ConfigMaps */}}
{{- range $i, $v := $versions }}
{{ include "configmap.weight-config" $v }}
---
{{- end }} {{- /* range $i, $v := $versions */}}

{{- /* routemap */}}
{{ include "env.knative-istio.blue-green.routemap" . }}

{{- end }} {{- /* define "env.knative-istio.blue-green" */}}