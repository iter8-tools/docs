---
template: main.html
hide:
- navigation
- toc
---

# Drop-dead Simple Performance Testing for Knative Services

> Launch performance tests for Knative apps (HTTP and gRPC services) and validate their service-level objectives (SLOs) **in seconds**.

Performance testing is a core building block in the robust delivery of Knative HTTP and gRPC apps (services). One way to accomplish this is by sending a stream of requests to the target service, and evaluating the responses for error and latency-related violations. From a developer’s perspective, this approach involves three main dimensions — namely, i) the load-related characteristics of the request stream, such as the request rate; ii) the shape of the requests, in particular, whether the service requires any payload/data to be sent as part of the requests; and iii) the service-level objectives (SLOs) used to validate the quality of the target service.

This article shows how you can use [Iter8](https://iter8.tools), the open source Kubernetes release optimizer, to flexibly launch performance tests for Knative apps **in seconds**, with precise control over all of the above. We begin with an [overview of Iter8](#overview-of-iter8), followed by [tutorials for performance testing of Knative HTTP and gRPC apps](#quick-start) that you can use to **get started quickly**, and conclude with some [useful variations and enhancements to these tutorials](#what-next) that you can try next.

## Overview of Iter8

[Iter8](https://iter8.tools) is the Kubernetes release optimizer built for DevOps, MLOps, SRE and data science teams. [Iter8](https://iter8.tools) makes it easy to ensure that Kubernetes apps and ML models perform well and maximize business value. Specifically, [Iter8](https://iter8.tools) introduces the notion of *experiments* that simplifies the collection of performance and business metrics for apps and ML models, assessment and comparison of one or more app/ML model versions, validation of service-level objectives (SLOs), promoting the winning version, and maximizing business value during each release.

<img src="https://iter8.tools/0.11/images/iter8-intro-dark.png" alt="Iter8 experiment" width="600"/>

### Why Iter8?
[Iter8](https://iter8.tools) is packed with powerful features that simplify performance testing of Knative apps. We highlight a few of them below.

1.  **Built-in latency and error-related metrics for HTTP and gRPC services.** Eliminates the need to set up and configure metrics databases during performance testing.
2.  **Well-defined notion of service-level objectives (SLOs)**. Makes it simple to define and verify SLOs in experiments.
3.  **Readiness checks.** The performance testing portion of the experiment begins only after the Knative service is ready.
4.  **HTML/text reports**. Promotes human understanding of experiment results through visual insights. 
5. **Assertions**. Launch an experiment and assert whether the target app satisfies the specified SLOs or not. Simplifies CI/CD/GitOps pipelines: they can branch off into different paths depending upon whether the assertions are true or false.

## Quick start

### Install Iter8 CLI

Install the latest stable release of the Iter8 CLI using `brew` as follows.

```shell
brew tap iter8-tools/iter8
brew install iter8@0.11
```

You can also install the Iter8 CLI using pre-built binaries, or use Iter8 inside a GitHub Actions pipeline as described [here](https://iter8.tools/0.11/getting-started/install/).

### Tutorial: HTTP performance test in seconds
In this tutorial, we will launch an Iter8 experiment that generates load for a Knative HTTP service and validates its service-level objectives (SLOs). The setup of this experiment is illustrated in the figure below.

<img src="https://iter8.tools/0.11/getting-started/images/http.png" alt="HTTP performance test" width="800"/>

Install Knative in your Kubernetes cluster, and deploy your Knative HTTP Service as described in [this Knative tutorial](https://knative.dev/docs/getting-started/first-service/). 

Launch the Iter8 experiment as follows.

```shell
iter8 k launch \
--set "tasks={ready,http,assess}" \
--set ready.ksvc=hello \
--set http.url=http://httpbin.default/get \
--set assess.SLOs.upper.http/latency-mean=200 \
--set assess.SLOs.upper.http/error-count=0 \
--set runner=job
```

???+ note "About this experiment"
    This experiment consists of three [tasks](concepts.md#tasks), namely, [ready](../user-guide/tasks/ready.md), [http](../user-guide/tasks/http.md), and [assess](../user-guide/tasks/assess.md). 
    
    The [ready](../user-guide/tasks/ready.md) task checks if the Knative service named `hello` exists and is ready.
    
    The [http](../user-guide/tasks/http.md) task sends requests to the cluster-local HTTP service whose URL is `http://hello.default`, and collects [Iter8's built-in HTTP load test metrics](../user-guide/tasks/http.md#metrics). 
    
    The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) the mean latency of the service does not exceed 200 msec, and ii) there are no errors (4xx or 5xx response codes) in the responses. 
    
    This is a [single-loop](concepts.md#loops) [Kubernetes experiment](concepts.md#execution-environments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](concepts.md#runners) value is set to `job`, which enables Iter8 to use a Kubernetes `job` to execute this experiment in the cluster.

## View experiment report
=== "Text"
    ```shell
    iter8 k report
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

=== "HTML"
    ```shell
    iter8 k report -o html > report.html # view in a browser
    ```

    ??? note "The HTML report looks like this"
        ![HTML report](https://iter8.tools/0.11/getting-started/images/report.html.png)

## Assert experiment outcomes
Assert that the experiment completed without failures, and all SLOs are satisfied. The timeout flag below specifies a period of 120 sec for assert conditions to be satisfied.

```shell
iter8 k assert -c completed -c nofailure -c slos --timeout 120s
```

If the assert conditions are satisfied, the above command exits with code 0; else, it exits with code 1. Assertions are especially useful inside CI/CD/GitOps pipelines. Depending on the exit code of the assert command, your pipeline can branch into different actions.

## View experiment logs
Logs are useful when debugging an experiment.

```shell
iter8 k log
```

??? note "Sample experiment logs"
    ```shell
    INFO[2022-06-27 11:50:39] inited Helm config                           
    INFO[2022-06-27 11:50:39] experiment logs from Kubernetes cluster       indented-trace=below ... 
      time=2022-06-27 15:48:59 level=info msg=task 1: ready : started
      time=2022-06-27 15:48:59 level=info msg=task 1: ready : completed
      time=2022-06-27 15:48:59 level=info msg=task 2: ready : started
      time=2022-06-27 15:48:59 level=info msg=task 2: ready : completed
      time=2022-06-27 15:48:59 level=info msg=task 3: http : started
      time=2022-06-27 15:49:11 level=info msg=task 3: http : completed
      time=2022-06-27 15:49:11 level=info msg=task 4: assess : started
      time=2022-06-27 15:49:11 level=info msg=task 4: assess : completed
    ```

## Cleanup
Remove the [Kubernetes objects](https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/) created by the `iter8 k launch` command.

```shell
iter8 k delete
```

### Tutorial: gRPC performance test in seconds
In this tutorial, we will launch an Iter8 experiment that generates load for a Knative gRPC service and validates its service-level objectives (SLOs). The setup of this experiment is illustrated in the figure below.

<img src="https://iter8.tools/0.11/tutorials/images/grpc.png" alt="gRPC performance test" width="800"/>

## What next?

1. Local experiments
2. Load characteristics ...
3. Expanded SLOs ... 
