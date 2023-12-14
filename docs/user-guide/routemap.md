---
template: main.html
---

# Routemaps

A _routemap_ contains a description of each version of an application and may contain one or more routing templates. The description of versions is used by Iter8 to identify which versions of the application are available at any moment. Whenever versions become available or disappear, any defined routing templates will be applied. This results in the automatic reconfiguration of the routing.

## Version list

A version is a list of resources that, when available and ready, indicate that the version is available. For example, the following describes a application with two versions. In this case, each version has a `Service` and a `Deployment`. In one version, they are named `httpbin-0`, and in the other, `httpbin-1`:

```yaml
versions:
- resources:
  - gvrShort: svc
    name: httpbin-0
    namespace: default
  - gvrShort: deploy
    name: httpbin-0
    namespace: default
- resources:
  - gvrShort: svc
    name: httpbin-1
    namespace: default
  - gvrShort: deploy
    name: httpbin-1
    namespace: default
```

Note that the resources types are specified using a short name. In this example, `svc` and `deploy`. A short name is used to simplify the specification of the resources in a version. A mapping of short name to Kubernetes group, version, resource is captured in the configuration of the Iter8 controller. This set can be [extended](controller/extensions.md) to include any types including custom resources; that is, those defined by a CRD.

A version may optionally specify an integer `weight` indicating the proportion of traffic that should be sent to it relative to other versions.

## Routing templates

A routing template is a Go template that is applied each time a new version becomes available or an old one goes away. Multiple templates can be defined/applied.

THe application of the templates allows Iter8 to automatically reconfigure the routing when versions come and go. For example, the template created in the [blue-green release tutorial](../getting-started/first-release.md) contains an Istio `VirtualService`. Applying the template to the available versions yields the necessary `VirtualSerivce` definition. In this tutorial, the template definition is:

```yaml
blue-green:
  gvrShort: vs
  template: |
    apiVersion: networking.istio.io/v1beta1
    kind: VirtualService
    metadata:
      name: httpbin
    spec:
      gateways:
      - mm-gateway
      - mesh
      hosts:
      - httpbin.default
      - httpbin.default.svc
      - httpbin.default.svc.cluster.local
      http:
      - route:
        # primary version
        - destination:
            host: httpbin-0.default.svc.cluster.local
          {{- if gt (index .Weights 1) 0 }}
          weight: {{ index .Weights 0 }}
          {{- end }}
          headers:
            response:
              add:
                app-version: httpbin-0
        # other versions
        {{- if gt (index .Weights 1) 0 }}
        - destination:
            host: httpbin-1.default.svc.cluster.local
          weight: {{ index .Weights 1 }}
          headers:
            response:
              add:
                app-version: httpbin-1
        {{- end }}
```

## Implementation

A routemap is implemented as a `ConfigMap` with the following labels:
- `iter8.tools/kind` with value `routemap`
- `iter8.tools/version` with value corresponding the version of the controller being used

The list of versions and routing templates are stored as stringified YAML.