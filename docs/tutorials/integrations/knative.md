---
template: main.html
hide:
- navigation
- toc
---

# Drop-dead Simple Performance Testing for Knative Apps

Performance testing is a core building block in the robust delivery of Knative apps. A standard approach to accomplishing this is to send a stream of requests to the target service, and evaluate the responses for error and latency-related violations. From a developer’s perspective, setting up a performance test involves three key dimensions — namely, i) the load-related characteristics of the request stream, such as the request rate; ii) the shape of the requests, in particular, whether the service requires any payload/data to be sent as part of its requests; and iii) the service-level objectives (SLOs) used to validate the quality of the target service.

This article shows how you can use [Iter8](https://iter8.tools), the open source Kubernetes release optimizer, to flexibly set up performance tests for Knative apps **in seconds**, with precise control over all of the above.

## Overview of Iter8

One or two line descriptions of Iter8... (from iter8.tools concepts.md)

<img src="https://iter8.tools/0.11/images/iter8-intro-dark.png" alt="Iter8 experiment" width="500"/>

Describe the power of the Iter8 solution (HTTP, gRPC, load profile, request payload/data, SLO validation and reporting).

## Quick start

### Install Iter8 CLI and setup Knative apps

### Setup sample HTTP and gRPC apps

### HTTP performance testing in seconds
<img src="https://iter8.tools/0.11/getting-started/images/http.png" alt="HTTP performance test" width="700"/>

### gRPC performance testing in seconds
<img src="https://iter8.tools/0.11/tutorials/images/grpc.png" alt="gRPC performance test" width="700"/>

## What next?
