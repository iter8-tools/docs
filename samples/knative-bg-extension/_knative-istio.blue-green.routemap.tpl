{{- define "env.knative-istio.blue-green.routemap" }}

{{- $APP_NAME := (include "application.name" .) }}
{{- $APP_NAMESPACE := (include "application.namespace" .) }}
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
      - gvrShort: cm
        name: {{ $v.VERSION_NAME }}-weight-config
        namespace: {{ $v.VERSION_NAMESPACE }}
      weight: {{ $v.weight }}
    {{- end }} {{- /* range $i, $v := $versions */}}
    routingTemplates:
      {{ .Values.application.strategy }}:
        gvrShort: vs
        template: |
          apiVersion: networking.istio.io/v1beta1
          kind: VirtualService
          metadata:
            name: {{ $APP_NAME }}
            namespace: {{ $APP_NAMESPACE }}
          spec:
            gateways:
            {{- if .Values.gateway }}
            - {{ .Values.gateway }}
            {{- end }}
            - mesh
            hosts:
            - {{ $APP_NAME }}.{{ $APP_NAMESPACE }}
            - {{ $APP_NAME }}.{{ $APP_NAMESPACE }}.svc
            - {{ $APP_NAME }}.{{ $APP_NAMESPACE }}.svc.cluster.local
            http:
            - name: {{ $APP_NAME }}
              route:
              # primary version
              {{- $v := (index $versions 0) }}
              - destination:
                  host: {{ template "ksvc.name" $v }}.{{ $APP_NAMESPACE }}.svc.cluster.local
                  port:
                    number: {{ $v.port }}
                {{- if gt (len $versions) 1 }}
                {{ `{{- if gt (index .Weights 1) 0 }}` }}
                weight: {{ `{{ index .Weights 0 }}` }}
                {{ `{{- end }}` }}
                {{- end  }}
                headers: 
                  request:
                    add:
                      Host: {{ template "ksvc.name" $v }}.{{ $APP_NAMESPACE }}.svc.cluster.local
                  response:
                    add:
                      app-version: {{ template "ksvc.name" $v }}
              # other versions
              {{- range $i, $v := (rest $versions) }}
              {{ `{{- if gt (index .Weights ` }}{{ print (add1 $i) }}{{ `) 0 }}` }}
              - destination:
                  host: {{ template "ksvc.name" $v }}.{{ $APP_NAMESPACE }}.svc.cluster.local
                  port:
                    number: {{ $v.port }}
                weight: {{ `{{ index .Weights ` }}{{ print (add1 $i) }}{{ ` }}` }}
                headers:
                  request:
                    add:
                      Host: {{ template "ksvc.name" $v }}.{{ $APP_NAMESPACE }}.svc.cluster.local
                  response:
                    add:
                      app-version: {{ template "ksvc.name" $v }}
              {{ `{{- end }}` }}
              {{- end }}
{{- end }} {{- /* define "env.knative-istio.blue-green.routemap" */}}
