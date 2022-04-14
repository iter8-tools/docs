---
template: main.html
---

# Iter8

Kubernetes release optimizer built for DevOps, MLOps, SRE, and data science teams.

## Iter8 experiment
<p align='center'>
  <img alt-text="load-test-http" src="../../images/iter8-intro-dark.png" width="70%" />
</p>

Iter8 experiments make it simple to collect performance and business metrics for apps and ML models, assess, compare and validate one or more app/ML model versions, promote the winning version, and maximize business value in each release.

## Experiment chart
Iter8 Experiment charts are [Helm charts](https://helm.sh/docs/topics/charts/) under the covers, and enable simple, declarative, and reusable metrics-driven experiments. Iter8 combines experiment charts with user supplied values to generate runnable `experiment.yaml` files or Kubernetes experiment manifests. The former is used for running experiments in the local environment, while the latter is used for running experiments inside Kubernetes.

Iter8 uses experiment charts located in the [Iter8 GitHub repo](https://github.com/iter8-tools/iter8) by default. You can create, package and host Iter8 experiment charts in any GitHub repo and use them with Iter8 CLI, and other Iter8 components.

## Features at a glance

- **Benchmark and validate HTTP and gRPC services** 
    
    Iter8 experiments can generate requests for HTTP and gRPC services, collect built-in latency and error-related metrics, and validate SLOs.

- **A/B(/n) testing** 
      
    Grow your business with every release. Iter8 experiments can compare multiple versions based on business value and identify a winner.

- **Simple to use** 
      
    Get started with Iter8 in seconds using pre-packaged experiment charts. Run Iter8 experiments locally, inside Kubernetes, or inside your CI/CD/GitOps pipelines.

- **App frameworks** 
      
    Use with any app, serverless, or ML framework. Iter8 works with Kubernetes deployments, statefulsets, Knative services, KServe/Seldon ML deployments, or other custom Kubernetes resource types.

## Implementation
Iter8 is written in `go` and builds on a few awesome open source projects including:

- [Helm](https://helm.sh)
- [Fortio](https://github.com/fortio/fortio)
- [ghz](https://ghz.sh)
- [plotly.js](https://github.com/plotly/plotly.js)
