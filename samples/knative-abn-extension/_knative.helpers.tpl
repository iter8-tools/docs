{{- /* deployment specific wrapping for normalize.versions */}}
{{- define "normalize.versions.ksvc" }}
{{- $versions := include "normalize.versions" . | mustFromJson }}
  {{- $normalizedVersions := list }}
  {{- range $i, $v := $versions -}}
    {{- $version := merge $v }}

    {{- $version = set $version "port" (pluck "port" (dict "port" 80) $.Values.application $v | last) }}
    {{- if (and $v.ksvcSpecification $v.ksvcSpecification.ports) }}
      {{- $version = set $version "port" (pluck "port" $v (index $v.ksvcSpecification.ports 0) | last) }}
    {{- end }}

    {{- $normalizedVersions = append $normalizedVersions $version }}
  {{- end }} {{- /* range $i, $v := $versions */}}
  {{- mustToJson $normalizedVersions }}
{{- end }} {{- /* define "normalize.versions.ksvc" */}}

{{- /* Identify the name of a Knative Service object */ -}}
{{- define "ksvc.name" -}}
{{- if (and .ksvcSpecification .ksvcSpecification.metadata .ksvcSpecification.metadata.name) -}}
{{ .ksvcSpecification.metadata.name }}
{{- else -}}
{{ .VERSION_NAME }}
{{- end -}}
{{- end }} {{- /* define "ksvc.name" */ -}}

{{- /* Identify the namespace of a Knative Service object */ -}}
{{- define "ksvc.namespace" -}}
{{- if (and .ksvcSpecification .ksvcSpecification.metadata .ksvcSpecification.metadata.namespace) -}}
{{ .ksvcSpecification.metadata.namespace }}
{{- else -}}
{{ .VERSION_NAMESPACE }}
{{- end -}}
{{- end }} {{- /* define "ksvc.namespace" */ -}}
