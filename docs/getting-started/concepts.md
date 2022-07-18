---
template: main.html
---

# Iter8
Iter8 is the Kubernetes release optimizer built for DevOps, MLOps, SRE and data science teams. Iter8 makes it easy to ensure that Kubernetes apps and ML models perform well and maximize business value.

Iter8 supports the following use-cases.
1.  Performance testing and SLO validation of HTTP services.
2.  Performance testing and SLO validation of gRPC services.
3.  SLO validation using custom metrics from any database(s) or REST API(s).


## Iter8 experiment
Iter8 introduces the notion of an experiment, which is a set of configurable tasks that are executed in a specific sequence.

<p align='center'>
  <img alt-text="load-test-http" src="../../images/iter8-intro-dark.png" width="70%" />
</p>

Iter8 packs a number of powerful features that facilitate Kubernetes app testing and experimentation. They include the following.

1.  **Generating load and collecting built-in metrics for HTTP and gRPC services.** Simplifies performance testing by eliminating the need to setup and use metrics databases.
2.  **Well-defined notion of service-level objectives (SLOs).** Makes it simple to define and verify SLOs in experiments.
3.  **Custom metrics.** Enables the use of custom metrics from any database(s) or REST API(s) in experiments.
4.  **Readiness check.** The performance testing portion of the experiment begins only after the service is ready.
5.  **HTML/text reports.** Promotes human understanding of experiment results through visual insights.
6.  **Assertions.** Verifies whether the target app satisfies the specified SLOs or not after an experiment. Simplifies automation in CI/CD/GitOps pipelines: branch off into different paths depending upon whether the assertions are true or false.
7.  **Multi-loop experiments.** Experiment tasks can be executed periodically (multi-loop) instead of just once (single-loop). This enables Iter8 to refresh metric values and perform SLO validation using the latest metric values during each loop.
8.  **Experiment anywhere.** Iter8 experiments can be launched inside a Kubernetes cluster, in local environments, or inside a GitHub Actions pipeline.

### Kubernetes experiments
Experiments that are executed inside Kubernetes clusters are referred to as **Kubernetes experiments**. All other experiments are referred to as **local experiments**.

### Runner
A [single-loop](#iter8-experiment) [Kubernetes experiment](#kubernetes-experiments) uses the Kubernetes [job](https://kubernetes.io/docs/concepts/workloads/controllers/job/) workload as its runner. A [multi-loop](#iter8-experiment) [Kubernetes experiment](#kubernetes-experiments) uses the Kubernetes [cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) workload as its runner.

### Specifying an experiment
Specifying an Iter8 experiment involves specifying the [list of tasks executed during the experiment](#iter8-experiment) and their [parameters](../user-guide/topics/parameters.md). Additionally, [Kubernetes experiments](#kubernetes-experiments) involve specifying the [runner](#runner).

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
