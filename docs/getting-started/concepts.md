---
template: main.html
---

# Iter8
Iter8 is the Kubernetes testing platform built for DevSecOps, MLOps, SRE and data science teams. Iter8 makes it easy to ensure that new versions of apps and ML models perform well and maximize business value.

## Iter8 Experiment
<p align='center'>
  <img alt-text="load-test-http" src="../../images/iter8-intro-dark.png" width="70%" />
</p>

Iter8 experiments makes it simple to collect performance and business metrics for apps and ML models, assess and compare one or more app/ML model versions, validate [service-level objectives (SLOs)](#service-level-objectives), promote the winning version, and maximize business value during each release.

Specifying an Iter8 experiment involves specifying the list of [tasks](#tasks) executed during the experiment, their parameter values, and the experiment [runner](#loops-and-runner).

### Tasks
An experiment is a sequence of tasks. Iter8 provides pre-defined tasks for various functions such as generating load and collecting Iter8's built-in metrics for HTTP and gRPC services, collecting custom metrics for one or more versions of an app from databases, assessing [SLOs](#service-level-objectives), and checking if the application is ready.

### Loops
Experiments have a concept of loops. When the sequence of tasks in an experiment are executed once from start to finish, a single loop of the experiment is said to have completed. An experiment can have a single loop or multiple loops.

### Execution Environments
Experiments can run inside Kubernetes clusters, in the local environment, and within CI/CD/GitOps workflows including GitHub Actions workflows.

### Runners
Runner refers to the type of Kubernetes workload used to run the experiment inside a Kubernetes cluster. Experiments use either a Kubernetes [job](https://kubernetes.io/docs/concepts/workloads/controllers/job/) or a Kubernetes [cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) as its runner. The former is used for single-loop experiments, while the latter is used for multi-loop experiments, where each loop is executed repeatedly at periodic intervals.

## Service-level Objectives

Service-level objectives (SLOs) represent the range of metric values for an application that is considered acceptable. Both upper and lower limits on metric values can be specified as SLOs in Iter8 experiments.


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
