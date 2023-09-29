---
template: main.html
---

# Iter8
Iter8 is the Kubernetes release optimizer built for DevOps, MLOps, SRE and data science teams. Iter8 automates traffic control for new versions of apps/ML models in the cluster and visualizes their performance metrics.

## Use-cases

Iter8 simplifies a variety of traffic engineering and metrics-driven validation use-cases. They are illustrated below.

![Iter8 use-cases](images/iter8usecases.png)

## Design

Iter8 provides three inter-related components to support the above use-cases.

1. Iter8 controller
2. Client SDK
3. Composable tasks

=== "Iter8 controller"

    Iter8 provides a controller that automatically and dynamically reconfigures routing resources based on the state of Kubernetes apps/ML models. 
    
    The following picture illustrates a blue-green rollout scenario that is orchestrated by this controller.

    ![Blue-green](../tutorials/integrations/kserve-mm/images/blue-green.png)
    
    As part of the dynamic reconfiguration of route resources, the Iter8 controller also checks for readiness (for e.g., in KServe ModelMesh), availability (for e.g., in Kubernetes deployments) and other relevant status conditions before configuring traffic splits to candidate versions. Similarly, before candidate versions are deleted, the Iter8 controller uses finalizers to first ensure that all traffic flows to the primary version of the ML model. This makes for a very high-degree of reliability and zero-downtime/loss-less rollouts of new app/ML model versions. Users do not get this level of reliability out-of-the-box with a vanilla service mesh.

    With Iter8, the barrier to entry for end-users is significantly reduced. In particular, by just providing names of their ML serving resources, and (optional) traffic weights/labels, end users can get started with their release optimization use cases rapidly. Further, Iter8 does not limit the capabilities of the underlying service mesh in any way. This means more advanced teams still get to use all the power of the service-mesh alongside the reliability and ease-of-use that Iter8 brings.

    In addition, Iter8 provides a simple metrics store, eliminating the need for an external database.

=== "Client SDK"

    Iter8 provides a client-side SDK to facilitate routing as well as metrics collection task associated with distributed (i.e., client-server architecture-based) A/B/n testing in Kubernetes. 
    
    The following picture illustrates the use of the SDK for A/B testing.

    ![A/B testing](images/abn.png)

    Iter8's SDK is designed to handle user stickiness, collection of business metrics, and decoupling of front-end and back-end releases processes during A/B/n testing.

=== "Composable tasks"

    Iter8 introduces a set of tasks which which can be composed in order to conduct a variety of performance tests.
    
    The following picture illustrates a performance test for an HTTP application, and this test consists of two tasks.

    ![Iter8 performance test](images/kubernetesusage.png)

    In addition to load testing HTTP and gRPC services, Iter8 tasks can perform other actions such as sending notifications to Slack or GitHub.

## Advantages
Iter8 has several advantages compared to other tooling in this space. 

First, Iter8 has no restrictions on the types of resources that make up a version of an application. This includes custom resources; that is, those defined by a custom resource definition (CRD). Because the set of resources that comprise a version is declarative, it is easy to [extend](../user-guide/topics/extensions.md). Note that this same extension mechanism also allows Iter8 to be used with any service mesh.

Second, the Iter8 client SDK addresses a key challenge to [A/B/n testing](../user-guide/topics/ab_testing.md): the decoupling of the front-end release process from that of the back-end. It allows the front-end to reliably associate business metrics with the contributing version of the back-end.

Finally, Iter8 simplifies performance testing by reducing the set up time needed to start testing. Tests can be easily specified as a sequence of [easily configured tasks](../user-guide/topics/parameters.md). Further, there is no need to setup and configure an external metrics database -- Iter8 captures the metrics data and provides a REST API allowing it to be visualized and evaluated in Grafana.

## Implementation
Iter8 is written in `go` and builds on a few awesome open source projects including:

- [Helm](https://helm.sh)
- [Istio](https://istio.io)
- [plotly.js](https://github.com/plotly/plotly.js)
- [Fortio](https://github.com/fortio/fortio)
- [ghz](https://ghz.sh)
- [Grafana](https://grafana.com/)
