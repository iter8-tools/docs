---
template: main.html
---

# About progressive release

Progressive release uses the Iter8 `release` Helm chart. Options are described in the default [values.yaml](https://github.com/iter8-tools/iter8/blob/v0.18.6/charts/release/values.yaml) file. The progressive release tutorials show how it can be used to progressively release versions of an application or ML model.

The chart provided by Iter8 supports many common deployment scenarios including:

- Applications composed of a `Deployment` and a `Service` object using Istio as a service mesh
- Applications composed of a `Deployment` and a `Service` object using the Kubernetes Gateway API
- ML models deployed to KServe (using KNative)
- ML models deploy to KServe ModelMesh using Istio as a service mesh

## Other deployment environments

The `release` chart can be easily [extended](extending.md) to include other deployment environments. Please consider contributing any extensions to Iter8.

