---
template: main.html
---

# Chart Documentation

An Iter8 experiment chart is a Helm chart, and hence inherits the latter's [file structure](https://helm.sh/docs/topics/charts/#the-chart-file-structure). In particular, the chart contains the following files.

* Chart.yaml: A YAML file containing information about the chart
* README.md: A human-readable README file
* values.yaml: The default configuration values for this chart

For example, to download the `load-test-grpc` chart and view its `values.yaml` , do as follows:

```shell
iter8 hub -c load-test-grpc
cd load-test-grpc
cat values.yaml
```

