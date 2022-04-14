---
template: main.html
---

# Chart Values

An Iter8 experiment chart is a Helm chart, and inherits the latter's [file structure](https://helm.sh/docs/topics/charts/#the-chart-file-structure). In particular, the chart contains the `values.yaml` file which documents all the parameters supported by the chart, and their (non-nil) default values (if any).

***

For example, to view `values.yaml` file for the `load-test-http` chart, do as follows:

```shell
iter8 hub
cat charts/load-test-http/values.yaml
```

