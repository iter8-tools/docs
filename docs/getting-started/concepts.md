---
template: main.html
---

# Iter8

Kubernetes release optimizer.

## Iter8 experiment
Iter8 experiments make it simple benchmark and validate HTTP and gRPC services with SLOs, and maximize business value with each release. Experiments can be run locally, inside Kubernetes, or inside CI/CD/GitOps pipelines.

<p align='center'>
  <img alt-text="load-test-http" src="../../images/iter8-intro-dark.png" width="70%" />
</p>

### Experiment chart
Experiment charts are [Helm charts](https://helm.sh/docs/topics/charts/) with a special structure that contain reusable experiment templates. Iter8 combines experiment charts with user supplied values to generate runnable `experiment.yaml` files. Iter8 experiment charts enable you to launch powerful release optimization experiments in seconds. Their usage is described in depth in various [Iter8 tutorials](../tutorials/load-test-http/usage.md).

#### Iter8 hub
Iter8 hub refers to the URL of a Helm repo that hosts Iter8 experiment charts. Unless stated otherwise, Iter8 hub refers to the default URL https://iter8-tools.github.io/hub/. You can host Iter8 experiment charts in any Helm repo, and use them with Iter8.

## Implementation
Iter8 is written in `go` and builds on a few awesome open source projects including:

- [Helm](https://helm.sh)
- [Fortio](https://github.com/fortio/fortio)
- [ghz](https://ghz.sh)
- [plotly.js](https://github.com/plotly/plotly.js)
