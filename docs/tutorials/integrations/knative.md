---
template: main.html
hide:
- navigation
- toc
---

# Drop-dead Simple Performance Testing for Knative Apps

Performance testing is an essential building block that enables robust continuous delivery of Kubernetes and Knative apps. A standard approach to accomplishing this is to send a stream of requests to the target service, and evaluate the responses for error and latency-related violations. From a developer’s perspective, setting up a performance test involves three key dimensions — namely, i) the load-related characteristics of the request stream, such as the request rate; ii) the shape of the requests, in particular, whether the service requires any payload/data to be sent as part of its requests; and iii) the service-level objectives (SLOs) used to validate the quality of the target service.

This article shows how you can use [Iter8](https://iter8.tools), the open source Kubernetes release optimizer, to flexibly set up a performance test with precise control over all of the above, within seconds.

## Overview of Iter8

One or two line descriptions of Iter8... 

: Iter8 experiment ... 

The power of the Iter8 solution (HTTP, gRPC, load profile, request payload/data, SLO validation and reporting)... set up the Knative performance testing illustrations.

: Performance testing Knative HTTP apps
: Performance testing Knative gRPC apps

## Quick start

### Install Iter8 CLI and setup Knative apps

### Sample HTTP and gRPC apps

### HTTP performance testing in seconds

### gRPC performance testing in seconds

## What next?
