---
template: main.html
---

# Benchmark and Validate a Knative HTTP service

???+ note "Before you begin"
    1. [Install Iter8 CLI](../../../getting-started/install.md).
    2. [Install Knative and deploy your first Knative Service](https://knative.dev/docs/getting-started/first-service/). As noted at the end of the Knative tutorial, when you curl the Knative service,
    ```shell
    curl http://hello.default.127.0.0.1.sslip.io
    ```
    you should see the expected output as follows.
    ```
    Hello World!
    ```

***

Benchmark and validate SLOs for the Knative HTTP service by launching an Iter8 experiment.

```shell
iter8 launch -c load-test-http \
--set url=http://hello.default.127.0.0.1.sslip.io \
--set SLOs.http/error-rate=0 \
--set SLOs.http/latency-mean=50 \
--set SLOs.http/latency-p90=100 
```

Please refer to [the usage documentation for the `load-test-http` experiment chart](../../load-test-http/basicusage.md) that describes how to parameterize this experiment, assert SLOs, and view experiment reports.