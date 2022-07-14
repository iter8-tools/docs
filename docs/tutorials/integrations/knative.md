---
template: main.html
hide:
- navigation
- toc
---

# Drop-dead Simple Performance Testing for Knative Services

Performance testing is a core building block in the robust delivery of Knative services. One approach to accomplishing this is to send a stream of requests to the target service, and evaluate the responses for error and latency-related violations. From a developer’s perspective, setting up a performance test involves three key dimensions — namely, i) the load-related characteristics of the request stream, such as the request rate; ii) the shape of the requests, in particular, whether the service requires any payload/data to be sent as part of the requests; and iii) the service-level objectives (SLOs) used to validate the quality of the target service.

This article shows how you can use [Iter8](https://iter8.tools), the open source Kubernetes release optimizer, to flexibly set up performance tests for Knative apps **in seconds**, with precise control over all of the above.

## Overview of Iter8

Iter8 is the Kubernetes release optimizer built for DevOps, MLOps, SRE and data science teams. Iter8 makes it easy to ensure that Kubernetes apps and ML models perform well and maximize business value. Specifically, Iter8 introduces the notion of *experiments* that simplifies the collection of performance and business metrics for apps and ML models, assessment and comparison of one or more app/ML model versions, validation of service-level objectives (SLOs), promoting the winning version, and maximizing business value during each release.

<img src="https://iter8.tools/0.11/images/iter8-intro-dark.png" alt="Iter8 experiment" width="500"/>

### Why Iter8?
Iter8 is packed with powerful features that simultaneously enhance and simplify the release engineering/CI/CD workflows for Knative apps. We highlight three of them below.

1.  Iter8 provides a unified way to benchmark and validate both HTTP and gRPC services. This enables developers to evolve a common set of CI/CD best practices and design patterns that are applicable to both these types of services.
2.  Built-in metrics and SLO validation
3.  Knative readiness checks
4.  Iter8’s CLI generates intuitive reports that describe the results of an experiment. Reports are ideal for promoting human understanding of experiment outcomes through visual insights. Iter8 CLI also provides a simple way to assert conditions, such as whether the target app satisfies the SLOs specified in the experiment. Assertions are ideal from an automation perspective, since they allow CI/CD pipelines and/or scripts to branch off into different paths depending on whether the assertions are true.
5.  Iter8 CLI makes it possible to launch experiments in seconds by supplying a few relevant parameter values, which conceptually and operationally simplify even advanced benchmarking and validation tasks.

## Quick start

### Install Iter8 CLI

### HTTP performance test in seconds
<img src="https://iter8.tools/0.11/getting-started/images/http.png" alt="HTTP performance test" width="700"/>

### gRPC performance test in seconds
<img src="https://iter8.tools/0.11/tutorials/images/grpc.png" alt="gRPC performance test" width="700"/>

## What next?
