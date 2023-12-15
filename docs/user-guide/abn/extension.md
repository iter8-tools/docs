---
template: main.html
---

# Using new resource types

Like progressive release, A/B/n tests use the Iter8 `release` chart. Extending this chart for a new resource type is similar to [extending it for progressive release](../progressive-release/extension.md). The key difference in in the definition of the [routemap](../routemap.md). We briefly describe how to extend the chart for an [Knative](https://knative.dev/docs/) application.

## Example (Knative Service)

For example, to extend the chart to support A/B/n testing for Knative services, the following files could be added:

- `_knative-istio.tpl` - wrapper for all objects that should be deployed
- `_knative-istio.version.ksvc.tpl` - the Knative service object that should be deployed for a version
- `_knative-istio.none.tpl` - wrapper for all objects that should be deployed to support the no automated traffic pattern
- `_knative-istio.none.routemap.tpl` - the routemap definition
- `_knative.helpers.tpl` - supporting functions

An implementation of these is [here](https://github.com/iter8-tools/docs/tree/v0.18.13/samples/knative-abn-extension).

Note that many of these additions are the same as in the example for progressive release. Of the above files, only `_knative-istio.tpl` is different. 

Finally, update `release.yaml` to include `knative-istio` as a valid option:

```tpl
{{- else if eq "knative-istio" .Values.environment }}
  {{- include "env.knative-istio" . }}
```

## Extend the controller

The Iter8 controller will need to be restarted with permission to watch Knative service objects. Re-install the controller using the following additional options:

```shell
--set resourceTypes.ksvc.Group=serving.knative.dev \
--set resourceTypes.ksvc.Version=v1 \
--set resourceTypes.ksvc.Resource=services \
--set "resourceTypes.ksvc.conditions[0]=Ready"
```

## Using the modified chart

To use the modified chart to run A/B tests for Knative services, see the example for [progressive release](../progressive-release/extension.md#using-the-modified-chart). Unset the `strategy` and `weight` fields.