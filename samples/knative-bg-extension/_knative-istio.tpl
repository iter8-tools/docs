{{- define "env.knative-istio" }}

{{- /* Prepare versions for simpler processing */}}
{{- $versions := include "normalize.versions.ksvc" . | mustFromJson }}

{{- range $i, $v := $versions }}
{{- /* KNative Service */}}
{{ include "env.knative.version.ksvc" $v }}
---
{{- end }} {{- /* range $i, $v := $versions */}}

{{- /* Service */}}
{{ include "env.knative-istio.service" . }}
---

{{- /* routemap (and other strategy specific objects) */}}
{{- if eq "blue-green" .Values.application.strategy }}
{{ include "env.knative-istio.blue-green" . }}
{{- end }} {{- /* if eq ... .Values.application.strategy */}}

{{- end }} {{- /* define "env.knative-istio" */}}