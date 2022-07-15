---
template: main.html
hide:
- navigation
- toc
---

# Drop-dead Simple Performance Testing for Knative Services

> Launch performance tests for Knative apps (HTTP and gRPC services) and validate their service-level objectives (SLOs) **in seconds**.

Performance testing is a core building block in the robust delivery of Knative HTTP and gRPC apps (services). One way to accomplish this is by sending a stream of requests to the target service, and evaluating the responses for error and latency-related violations. From a developer’s perspective, this approach involves three main dimensions — namely, i) the load-related characteristics of the request stream, such as the request rate; ii) the shape of the requests, in particular, whether the service requires any payload/data to be sent as part of the requests; and iii) the service-level objectives (SLOs) used to validate the quality of the target service.

This article shows how you can use [Iter8](https://iter8.tools), the open source Kubernetes release optimizer, to flexibly launch performance tests for Knative apps **in seconds**, with precise control over all of the above. We begin with an [overview of Iter8](#overview-of-iter8), followed by [tutorials for performance testing of Knative HTTP and gRPC apps](#quick-start) that you can use to **get started quickly**, and conclude with some [useful variations and enhancements to these tutorials](#what-next) that you can try next.

## Overview of Iter8

[Iter8](https://iter8.tools) is the Kubernetes release optimizer built for DevOps, MLOps, SRE and data science teams. [Iter8](https://iter8.tools) makes it easy to ensure that Kubernetes apps and ML models perform well and maximize business value. Specifically, [Iter8](https://iter8.tools) introduces the notion of *experiments* that simplifies the collection of performance and business metrics for apps and ML models, assessment and comparison of one or more app/ML model versions, validation of service-level objectives (SLOs), promoting the winning version, and maximizing business value during each release.

<img src="https://iter8.tools/0.11/images/iter8-intro-dark.png" alt="Iter8 experiment" width="600"/>

### Why Iter8?
[Iter8](https://iter8.tools) is packed with powerful features that simplify performance testing of Knative apps. We highlight a few of them below.

1.  **Built-in latency and error-related metrics for HTTP and gRPC services.** Eliminates the need to set up and configure metrics databases for performance testing.
2.  **Well-defined notion of service-level objectives (SLOs)**. Makes it simple to define and verify SLOs in experiments.
3.  **Readiness checks.** Performance testing portion of the experiment begins only after the Knative service is ready.
4.  **HTML/text reports**. Ideal for promoting human understanding of experiment outcomes through visual insights. 
5. **Assertions**. Launch an experiment and assert whether the target app satisfies the specified SLOs or not. Ideal for automation: CI/CD/GitOps pipelines can branch off into different paths depending upon whether the assertions are true or false.

## Quick start

### Install Iter8 CLI

Install the latest stable release of the Iter8 CLI using `brew` as follows.

```shell
brew tap iter8-tools/iter8
brew install iter8@0.11
```

You can also install the Iter8 CLI using pre-built binaries, or use Iter8 inside a GitHub Actions pipeline as described [here](https://iter8.tools/0.11/getting-started/install/).

### HTTP performance test in seconds
<img src="https://iter8.tools/0.11/getting-started/images/http.png" alt="HTTP performance test" width="800"/>

### gRPC performance test in seconds
<img src="https://iter8.tools/0.11/tutorials/images/grpc.png" alt="gRPC performance test" width="800"/>

## What next?
