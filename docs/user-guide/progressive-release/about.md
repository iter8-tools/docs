---
template: main.html
---

# About progressive release

Progressive release is an approach to application/ML model release that involves deploying a candidate version at the same time as the current, or primary, version. A portion of user traffic is directed to the candidate version allowing it to be evaluated for behavior and performance. Once approved, the candidate can be promoted -- becoming the new primary version.

There are different strategies for distributing traffic between the primary and candidate version:

- **blue-green** - A percentage of requests are directed to a candidate version of the model. This percentage can be changed over time. The remaining requests go to the primary version.
- **canary** - Requests that match a particular pattern, for example those that have a particular header, are directed to the candidate version of the model. The remaining requests go to the primary version.
- **mirrored** - All requests are sent to the primary version. A percentage of requests are replicated and sent to a candidate version of the model. This percentage can be changed over time. Only responses from the primary version are returned to the user.

Progressive release uses the Iter8 `release` Helm chart. Options are described in the default [values.yaml](https://github.com/iter8-tools/iter8/blob/v0.18.6/charts/release/values.yaml) file. The progressive release tutorials show how it can be used to progressively release versions of an application or ML model.

The chart provided by Iter8 supports many common deployment scenarios including:

- Applications composed of a `Deployment` and a `Service` object using Istio as a service mesh
- Applications composed of a `Deployment` and a `Service` object using the Kubernetes Gateway API
- ML models deployed to KServe (using KNative)
- ML models deploy to KServe ModelMesh using Istio as a service mesh

## Other deployment environments

The `release` chart can be easily [extended](extending.md) to include other deployment environments. Please consider contributing any extensions to Iter8.

