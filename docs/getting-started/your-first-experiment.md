---
template: main.html
---

# Your First Experiment

!!! tip "Benchmark an HTTP Service"
    Get started with your first [Iter8 experiment](concepts.md#what-is-an-iter8-experiment) by benchmarking an HTTP service. 
    
***

## 1. Install Iter8 CLI
--8<-- "docs/getting-started/installiter8cli.md"

## 2. Launch experiment
The `iter8 launch` subcommand downloads an [experiment chart](concepts.md#experiment-chart) from [Iter8 hub](concepts.md#iter8-hub), combines the chart with values to generate the `experiment.yaml` file, runs the experiment, and writes results into the `result.yaml` file. 

Use `iter8 launch` to benchmark the HTTP service whose URL is https://httpbin.org/get.

```shell
iter8 launch -c load-test-http --set url=https://httpbin.org/get
```

## 3. View experiment report
=== "HTML"
    ```shell
    iter8 report -o html > report.html
    # open report.html with a browser. In MacOS, you can use the command:
    # open report.html
    ```

    ???+ note "The HTML report looks like this"
        ![HTML report](images/report.html.png)

=== "Text"
    ```shell
    iter8 report
    ```

    ???+ note "The text report looks like this"
        ```shell
        Experiment summary:
        *******************

          Experiment completed: true
          No task failures: true
          Total number of tasks: 1
          Number of completed tasks: 1

        Latest observed values for metrics:
        ***********************************

          Metric                              |value
          -------                             |-----
          built-in/http-error-count           |0.00
          built-in/http-error-rate            |0.00
          built-in/http-latency-max (msec)    |203.78
          built-in/http-latency-mean (msec)   |17.00
          built-in/http-latency-min (msec)    |4.20
          built-in/http-latency-p50 (msec)    |10.67
          built-in/http-latency-p75 (msec)    |12.33
          built-in/http-latency-p90 (msec)    |14.00
          built-in/http-latency-p95 (msec)    |15.67
          built-in/http-latency-p99 (msec)    |202.84
          built-in/http-latency-p99.9 (msec)  |203.69
          built-in/http-latency-stddev (msec) |37.94
          built-in/http-request-count         |100.00
        ```

Congratulations! :tada: You completed your first Iter8 experiment.

???+ tip "Next steps"
    1. Learn more about [benchmarking and validating HTTP services with service-level objectives (SLOs)](../tutorials/load-test-http/usage.md).
    2. Learn more about [benchmarking and validating gRPC services with service-level objectives (SLOs)](../tutorials/load-test-grpc/usage.md).
