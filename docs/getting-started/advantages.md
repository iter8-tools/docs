---
template: main.html
---

# Advantages

## Use any resources

Iter8 allows an application to be composed of resources of any type, including custom resources defined by custom resource definitions (CRDs). In Iter8, the set of resources that make up an application version is declarative making it easy to extend for new resource types. The same extension mechanism also allows Iter8 to be used with any service mesh or ingress.

## A/B/n testing of backend components

Using the Iter8 SDK enables frontend application components to reliably associate business metrics with the contributing versions of the backend. This addresses a key challenge of [A/B/n testing](../user-guide/abn/about.md) of backend application components/ML models.

## Simplified performance testing

Iter8 makes it easy to start running performance tests. Tests are composed with [easily configurable tasks](../user-guide/performance/parameters.md). Furthermore, there is no need to setup and configure an external metrics database -- Iter8 captures the metrics data and provides a REST API to access it, allowing it to be visualized and evaluated in Grafana.

# Comparison to other tools 

[Flagger](https://flagger.app/) and [Argo Rollouts](https://argo-rollouts.readthedocs.io/en/stable/) share similarities with Iter8. Both provide support for advanced application rollout on Kubernetes with blue-green and canary analysis. They work with many service meshes and ingress products to provide this support. Users specify the desired rollout using a Kubernetes custom resource.

Iter8 is inspired by both projects. However, Iter8 differs in several regards. For example, with Iter8:

- Applications can be composed of any resource type. For example, it works with machine learning applications built using KServe `InferenceService` resources out of the box. To do so, Iter8 allows the user to specify the resources being deployed as part of the specification of the rollout instead of assuming a particular pattern.

- Users can A/B/n test application backend components. Beyond providing HTTP header and cookie-based routing, Iter8 provides a client SDK with a simple API that allows users to write frontend components designed to focus A/B/n testing on the backend components.

- No custom resource is required to specify rollouts. Both Flagger and Argo Rollouts, requires the user to install and use a custom resource type to define rollouts. In Iter8, users specify rollouts using Helm configuration files.