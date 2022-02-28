---
template: main.html
---

# Your First Experiment

Get started with your first [Iter8 experiment](concepts.md#what-is-an-iter8-experiment) by benchmarking an HTTP service. 
    
***

## 1. Install Iter8 CLI
--8<-- "docs/getting-started/installiter8cli.md"

## 2. Launch experiment
Use `iter8 launch` to benchmark the HTTP service whose URL is https://httpbin.org/get.

```shell
iter8 launch -c load-test-http --set url=https://httpbin.org/get
```

The `iter8 launch` subcommand downloads an [experiment chart](concepts.md#experiment-chart) from [Iter8 hub](concepts.md#iter8-hub), combines the chart with values that are set in order to generate the `experiment.yaml` file, runs the experiment, and writes results into the `result.yaml` file. 

## 3. View experiment report
--8<-- "docs/getting-started/expreport.md"


Congratulations! :tada: You completed your first Iter8 experiment.

???+ tip "Next steps"
    1. Learn more about [benchmarking and validating HTTP services with service-level objectives (SLOs)](../tutorials/load-test-http/basicusage.md).
    2. Learn more about [benchmarking and validating gRPC services with service-level objectives (SLOs)](../tutorials/load-test-grpc/basicusage.md).
