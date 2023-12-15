{{- define "env.knative-istio" }}

{{- /* Prepare versions for simpler processing */}}
{{- $versions := include "normalize.versions.ksvc" . | mustFromJson }}

{{- range $i, $v := $versions }}
{{- /* KNative Service */}}
{{ include "env.knative.version.ksvc" $v }}
---
{{- end }} {{- /* range $i, $v := $versions */}}

{{- /* routemap (and other strategy specific objects) */}}
{{- if not .Values.application.strategy }}
{{ include "env.knative-istio.none" . }}
{{- else if eq "none" .Values.application.strategy }}
{{ include "env.knative-istio.none" . }}
{{- end }} {{- /* if eq ... .Values.application.strategy */}}

{{- end }} {{- /* define "env.knative-istio" */}}