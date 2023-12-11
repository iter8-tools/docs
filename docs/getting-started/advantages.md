---
template: main.html
---

# Advantages
First, Iter8 has no restrictions on the types of resources that make up a version of an application. This includes custom resources; that is, those defined by a custom resource definition (CRD). Because the set of resources that comprise a version is declarative, it is easy to [extend](../user-guide/topics/extensions.md). Note that this same extension mechanism also allows Iter8 to be used with any service mesh.

Second, the Iter8 client SDK addresses a key challenge to [A/B/n testing](../user-guide/topics/ab_testing.md): the decoupling of the front-end release process from that of the back-end. It allows the front-end to reliably associate business metrics with the contributing version of the back-end.

Finally, Iter8 simplifies performance testing by reducing the set up time needed to start testing. Tests can be easily specified as a sequence of [easily configured tasks](../user-guide/topics/parameters.md). Further, there is no need to setup and configure an external metrics database -- Iter8 captures the metrics data and provides a REST API allowing it to be visualized and evaluated in Grafana.

# Comparison to other tools 

Other tools with similarities to Iter8 are [Flagger](https://flagger.app/) and [Argo Rollouts](https://argo-rollouts.readthedocs.io/en/stable/). 
Both provide support advanced application rollout on Kubernetes with blue-green and canary analysis. They work with a many service mesh and ingress products to provide this support.
Users specify the desired rollout using a Kubernetes custom resource.

Iter8 was heavily inspired by both projects. However, Iter8 differs in several regards. For example, with Iter8:

- Applications can be composed of any resource types. For example, it works with machine learning applications built using KServe `InferenceService` resources out of the box. To do so, Iter8 allows the user to specify the resources being deployed as part of the specification of the rollout instead of assuming a particular pattern.

- Users can A/B/n test application backend components. Beyond providing HTTP header and cookie based routing, Iter8 provides a client SDK with a very simple API that allows users to write frontend components designed to decouple A/B/n testing of the backend components.

- No custom resource is required. Both Flagger and Argo Rollouts, requires the user to install and use a custom resource type. Rollouts are specified using instances of this custom type. Instead, Iter8 provides allows users to specify an application rollout using a Helm configuration file.