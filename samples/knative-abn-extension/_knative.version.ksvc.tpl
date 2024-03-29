{{- define "env.knative.version.ksvc" }}

{{- /* compute basic metadata */}}
{{- $metadata := include "application.version.metadata" . | mustFromJson }}

{{- /* define Service */}}
apiVersion: serving.knative.dev/v1
kind: Service
{{- if .ksvcSpecification }}
metadata:
{{- if .ksvcSpecification.metadata }}
  {{ toYaml (merge .ksvcSpecification.metadata $metadata) | nindent 2 | trim }}
{{- else }}
  {{ toYaml $metadata | nindent 2 | trim }}
{{- end }} {{- /* if .ksvcSpecification.metadata */}}
spec:
  {{ toYaml .ksvcSpecification.spec | nindent 2  | trim }}
{{- else }} {{- /* if .ksvcSpecification */}}
{{- if not .image }} {{- /* require .image */}}
{{- print "missing field: image required when deploymentSpecification absent" | fail }}
{{- end }} {{- /* if not .image */}}
{{- if not .port }} {{- /* require .port */}}
{{- print "missing field: port required when deploymentSpecification absent" | fail }}
{{- end }} {{- /* if not .port */}}
metadata:
  {{ toYaml $metadata | nindent 2 | trim }}
spec:
  template:
    spec:
      containers:
        - image: {{ .image }}
          ports:
            - containerPort: {{ .port }}
{{- end }} {{- /* if .ksvcSpecification */}}

{{- end }} {{- /* define "env.knative.version.ksvc" */}}
