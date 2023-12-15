{{- define "env.knative-istio.none" }}

{{- /* prepare versions for simpler processing */}}
{{- $versions := include "normalize.versions.ksvc" . | mustFromJson }}

{{- /* routemap */}}
{{ include "env.knative-istio.none.routemap" . }}

{{- end }} {{- /* define "env.knative-istio.none" */}}