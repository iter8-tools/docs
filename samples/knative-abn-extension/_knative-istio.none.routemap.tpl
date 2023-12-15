{{- define "env.knative-istio.none.routemap" }}

{{- $versions := include "normalize.versions.ksvc" . | mustFromJson }}

apiVersion: v1
kind: ConfigMap
{{- template "routemap.metadata" . }}
data:
  strSpec: |
    versions: 
    {{- range $i, $v := $versions }}
    - resources:
      - gvrShort: ksvc
        name: {{ template "ksvc.name" $v }}
        namespace: {{ template "ksvc.namespace" $v }}
      weight: {{ $v.weight }}
    {{- end }} {{- /* range $i, $v := $versions */}}
{{- end }} {{- /* define "env.knative-istio.none.routemap" */}}
