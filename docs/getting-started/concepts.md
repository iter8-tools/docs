---
template: main.html
---

# Iter8
Iter8 is the Kubernetes release optimizer built for DevOps, MLOps, SRE and data science teams. Iter8 makes it easy to ensure that Kubernetes apps and ML models perform well and maximize business value.

## Iter8 experiment
<p align='center'>
  <img alt-text="load-test-http" src="../../images/iter8-intro-dark.png" width="70%" />
</p>

Iter8 experiments make it simple to collect performance and business metrics for apps and ML models, assess and compare one or more app/ML model versions, validate [service-level objectives (SLOs)](#service-level-objectives), promote the winning version, and maximize business value during each release.

### Tasks
An experiment is a set of tasks that are executed in a specific sequence. Iter8 provides pre-defined tasks for various functions such as the following:

1. Generating load and collecting Iter8's built-in metrics for HTTP services.
2. Generating load and collecting Iter8's built-in metrics for gRPC services.
3. Fetching metrics from databases and other REST APIs.
4. Assessing [SLOs](#service-level-objectives) for app versions.
5. Checking if app is ready.

### Loops
A single loop of an experiment involves the entire sequence of [tasks](#tasks) in the experiment executing once. You can configure experiments to loop once and finish (**single-loop**), or loop repeatedly over time (**multi-loop**).

### Execution environments
Experiments can run inside Kubernetes clusters, in the local environment, and within CI/CD/GitOps workflows including GitHub Actions workflows. Experiments that are executed inside Kubernetes clusters are referred to as **Kubernetes experiments**. All other experiments are referred to as **local experiments**.

### Runner
A [single-loop](#loops) [Kubernetes experiment](#execution-environments) uses the Kubernetes [job](https://kubernetes.io/docs/concepts/workloads/controllers/job/) workload as its runner. A [multi-loop](#loops) [Kubernetes experiment](#execution-environments) uses the Kubernetes [cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) workload as its runner.

### Specifying an experiment
Specifying an Iter8 experiment involves specifying the list of [tasks](#tasks) executed during the experiment and their [parameters](../user-guide/topics/parameters.md). Additionally, Kubernetes experiments involve specifying the [runner](#runner).

## Service-level objectives

Service-level objectives (SLOs) are acceptable limits for an app's metric values. Both upper and lower limits on metric values can be specified as SLOs in Iter8 experiments.


<!-- ## Features at a glance

- **Benchmark and validate HTTP and gRPC services** 
    
    Iter8 experiments can generate requests for HTTP and gRPC services, collect built-in latency and error-related metrics, and validate SLOs.

- **A/B(/n) testing** 
      
    Grow your business with every release. Iter8 experiments can compare multiple versions based on business value and identify a winner.

- **Simple to use** 
      
    Get started with Iter8 in seconds using pre-packaged experiment charts. Run Iter8 experiments locally, inside Kubernetes, or inside your CI/CD/GitOps pipelines.

- **App frameworks** 
      
    Use with any app, serverless, or ML framework. Iter8 works with Kubernetes deployments, statefulsets, Knative services, KServe/Seldon ML deployments, or other custom Kubernetes resource types. -->

## Implementation
Iter8 is written in `go` and builds on a few awesome open source projects including:

- [Helm](https://helm.sh)
- [Fortio](https://github.com/fortio/fortio)
- [ghz](https://ghz.sh)
- [plotly.js](https://github.com/plotly/plotly.js)
