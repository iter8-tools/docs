---
template: main.html
---

# Your First Experiment

!!! tip "Benchmark an HTTP Service"
    Get started with your first [Iter8 experiment](concepts.md#what-is-an-iter8-experiment) by benchmarking an HTTP service. 
    
***

## 1. Install Iter8
=== "Brew"
    Install the latest stable release of the Iter8 CLI using `brew` as follows.

    ```shell
    brew tap iter8-tools/iter8
    brew install iter8
    ```
    
=== "Binaries"
    Replace `${TAG}` below with the [latest or any desired Iter8 release tag](https://github.com/iter8-tools/iter8/releases). For example,
    ```shell
    export TAG=v0.9.0
    ```

    === "darwin-amd64 (MacOS)"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-darwin-amd64.tar.gz | tar xvz -
        ```
        Move `darwin-amd64/iter8` to any directory in your `PATH`.

    === "linux-amd64"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-linux-amd64.tar.gz | tar xvz -
        ```
        Move `linux-amd64/iter8` to any directory in your `PATH`.

    === "linux-386"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-linux-386.tar.gz | tar xvz -
        ```
        Move `linux-386/iter8` to any directory in your `PATH`.

    === "windows-amd64"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-windows-amd64.tar.gz | tar xvz -
        ```
        Move `windows-amd64/iter8.exe` to any directory in your `PATH`.


=== "Source"
    Build Iter8 from source as follows. Go `1.17+` is a pre-requisite.
    ```shell
    # you can replace master with a specific tag, for example, v0.9.0
    export REF=master
    https://github.com/iter8-tools/iter8.git?ref=$REF
    cd iter8
    make install
    ```

=== "Go 1.17+"
    Install the latest stable release of the Iter8 CLI using `go 1.17+` as follows.

    ```shell
    go install github.com/iter8-tools/iter8@latest
    ```
    You can now run `iter8` (from your gopath bin/ directory)

## 2. Launch experiment
We will benchmark the HTTP service whose URL (`url`) is https://httpbin.org/get. 

The `iter8 launch` command downloads an [experiment chart](concepts.md#experiment-chart) from [Iter8 hub](concepts.md#iter8-hub), combines the chart with values to generate the `experiment.yaml` file, runs the experiment, and writes results into the `result.yaml` file. Launch the benchmarking experiment as follows.

```shell
iter8 launch -c load-test-http --set url=https://httpbin.org/get
```

## 3. View report
View a report containing the metrics collected during this experiment in HTML or text formats as follows.

=== "HTML"
    ```shell
    iter8 report -o html > report.html
    # open report.html with a browser. In MacOS, you can use the command:
    # open report.html
    ```

    ??? note "The HTML report looks like this"
        ![HTML report](images/report.html.png)

=== "Text"
    ```shell
    iter8 report
    ```

    ??? note "The text report looks like this"
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
