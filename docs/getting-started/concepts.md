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
Iter8 introduces the notion of an *experiment*. An experiment is simply a list of tasks that are executed in a specific sequence.

<p align='center'>
  <img alt-text="Iter8 experiment" src="../../images/authoring.png" width="40%" />
</p>

Iter8 provides a variety of configurable tasks. Authoring an experiment is as simple as specifying the names of the tasks and specifying their parameter values. The following are some examples of tasks provided by Iter8.

1.  Tasks for **generating load and collecting built-in metrics** for HTTP and gRPC services.
2.  A task for verifying **service-level objectives (SLOs)** for apps or app versions.
3.  A task for fetching **custom metrics** from any database(s) or REST API(s).
4.  A task for checking if an object **exists** in the Kubernetes cluster and is **ready**.

In addition to pre-defined tasks, Iter8 packs a number of powerful features that facilitate experimentation. They include the following.

1.  **HTML/text reports** that promote end-user understanding of experiment results through visual insights.
2.  **Assertions** that verify whether the target app satisfies the specified SLOs or not during/after an experiment.
3.  **Multi-loop experiments** that can be executed periodically instead of just once (single-loop).
4.  **Local experiments** that enable you to run experiments in your local environment.
5.  **Iter8 GitHub Action** that enables you to invoke the Iter8 CLI within a GitHub Actions workflow.

## Under the covers
In order to execute an experiment inside Kubernetes, Iter8 uses a Kubernetes [job](https://kubernetes.io/docs/concepts/workloads/controllers/job/) (single-loop) or a Kubernetes [cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) (multi-loop) workload, along with a Kubernetes secret. Iter8 instantiates all experiments using a Helm chart, that is also pre-defined and provided by Iter8.

<p align='center'>
  <img alt-text="Iter8 experiment" src="../../images/underthecovers.png" width="100%" />
</p>

## Implementation
Iter8 is written in `go` and builds on a few awesome open source projects including:

- [Helm](https://helm.sh)
- [Fortio](https://github.com/fortio/fortio)
- [ghz](https://ghz.sh)
- [plotly.js](https://github.com/plotly/plotly.js)
