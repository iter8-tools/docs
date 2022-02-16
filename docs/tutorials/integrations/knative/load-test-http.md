---
template: main.html
---

# Benchmark and Validate a Knative HTTP service

???+ note "Before you begin"
    1. [Install Iter8](../../../getting-started/install.md).
    2. [Install Knative and deploy your first Knative Service](https://knative.dev/docs/getting-started/first-service/). As noted at the end of the Knative tutorial, when you curl the Knative service,
    ```shell
    curl http://hello.default.127.0.0.1.sslip.io
    ```
    you should see the expected output as follows.
    ```
    Hello World!
    ```

### 1. Launch experiment
We will benchmark and validate SLOs for the Knative HTTP service by launching an Iter8 experiment.

```shell
iter8 launch -c load-test-http \
          --set url=http://hello.default.127.0.0.1.sslip.io \
          --set SLOs.error-rate=0 \
          --set SLOs.latency-mean=50 \
          --set SLOs.latency-p90=100 \
          --set SLOs.latency-p'97\.5'=200
```

In the above experiment, the following SLOs are validated for the Knative service.
- error rate is 0
- mean latency is under 50 msec
- 90th percentile latency is under 100 msec
- 97.5th percentile latency is under 200 msec

### 2. Assert outcomes
Assert that the experiment completed without any failures and SLOs are satisfied.

```shell
iter8 assert -c completed -c nofailure -c slos
```

### 3. View report
View a report of the experiment in HTML or text formats as follows.

=== "HTML"
    ```shell
    iter8 report -o html > report.html
    # open report.html with a browser. In MacOS, you can use the command:
    # open report.html
    ```

=== "Text"
    ```shell
    iter8 report -o text
    ```
