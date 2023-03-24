---
template: main.html
---

# Iter8
Iter8 is the Kubernetes release optimizer built for DevOps, MLOps, SRE and data science teams. Iter8 makes it easy to ensure that Kubernetes apps and ML models perform well and maximize business value.

Iter8 supports the following use-cases.

1.  Performance testing and SLO validation of HTTP services.
2.  Performance testing and SLO validation of gRPC services.
3.  SLO validation using custom metrics from any database(s) or REST API(s).
4.  A/B/n experiments.


## Iter8 experiment
Iter8 introduces the notion of an *experiment*. An experiment is simply a list of tasks that are executed in a specific sequence.

![Iter8 experiment](../../docs/images/authoring.png)

Iter8 provides a variety of configurable tasks. Authoring an experiment is as simple as specifying the names of the tasks and specifying their parameter values. The following are some examples of tasks provided by Iter8.

1.  Tasks for **generating load and collecting built-in metrics** for HTTP and gRPC services.
2.  A task for verifying **service-level objectives (SLOs)** for apps or app versions.
3.  A task for fetching **custom metrics** from any database(s) or REST API(s).
4.  A task for checking if an object **exists** in the Kubernetes cluster and is **ready**.

In addition to predefined tasks, Iter8 packs a number of powerful features that facilitate experimentation. They include the following.

1.  **HTML/text reports** that promote end-user understanding of experiment results through visual insights.
2.  **Assertions** that verify whether the target app satisfies the specified SLOs or not during/after an experiment.
3.  **Multi-loop experiments** that can be executed periodically instead of just once (single-loop).
4.  **Iter8 GitHub Action** that enables you to invoke the Iter8 CLI within a GitHub Actions workflow.

## Imperative and declarative experiments

You can use the Iter8 CLI to launch and manage experiments through the command line. This is the imperative style of experimentation. You can also use the Iter8 Autox controller to launch and manage experiments declaratively. [AutoX](../user-guide/topics/autox.md), short for “automated experiments”, allows Iter8 to detect changes to your Kubernetes resources objects and automatically start new experiments, allowing you to test your applications as soon as you release a new version.

## Under the covers
In order to execute an experiment inside Kubernetes, Iter8 uses a Kubernetes [job](https://kubernetes.io/docs/concepts/workloads/controllers/job/) (single-loop) or a Kubernetes [cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) (multi-loop) workload, along with a Kubernetes secret. Iter8 instantiates all experiments using a Helm chart, that is also predefined and provided by Iter8.

![Iter8 experiment](./../images/underthecovers.png" width="100%")

## Implementation
Iter8 is written in `go` and builds on a few awesome open source projects including:

- [Helm](https://helm.sh)
- [Fortio](https://github.com/fortio/fortio)
- [ghz](https://ghz.sh)
- [plotly.js](https://github.com/plotly/plotly.js)