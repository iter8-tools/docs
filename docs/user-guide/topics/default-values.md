---
template: main.html
---

# Default Chart Values

An Iter8 experiment chart is a Helm chart, and hence inherits the latter's [file structure](https://helm.sh/docs/topics/charts/#the-chart-file-structure). In particular, the chart contains the `values.yaml` which provides the default parameter values for the chart.

For example, to view the default values for the `load-test-grpc` chart, do as follows:

```shell
iter8 hub
cat charts/load-test-grpc/values.yaml
```

