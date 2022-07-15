---
template: main.html
hide:
- navigation
- toc
---

# Dead Simple Performance Testing with SLOs for Knative Services

> Launch performance tests for Knative apps (HTTP and gRPC services) and validate their service-level objectives (SLOs) **in seconds**.

Performance testing is a core building block in the robust delivery of Knative HTTP and gRPC apps (services). One way to accomplish this is by sending a stream of requests to the target service, and evaluating the responses for error and latency-related violations. From a developer’s perspective, this approach involves three main dimensions — namely, i) the load-related characteristics of the request stream, such as the request rate; ii) the shape of the requests, in particular, whether the service requires any payload/data to be sent as part of the requests; and iii) the service-level objectives (SLOs) used to validate the quality of the target service.

You can use [Iter8](https://iter8.tools), the open source Kubernetes release optimizer, to flexibly launch performance tests for Knative apps **in seconds**, with precise control over all of the above. This article shows how.

## Overview of Iter8

[Iter8](https://iter8.tools) is the Kubernetes release optimizer built for DevOps, MLOps, SRE and data science teams. [Iter8](https://iter8.tools) makes it easy to ensure that Kubernetes apps and ML models perform well and maximize business value. Specifically, [Iter8](https://iter8.tools) introduces the notion of *experiments* that simplifies the collection of performance and business metrics for apps and ML models, assessment and comparison of one or more app/ML model versions, validation of service-level objectives (SLOs), promoting the winning version, and maximizing business value during each release.

<img src="https://iter8.tools/0.11/images/iter8-intro-dark.png" alt="Iter8 experiment" width="600"/>

### Why Iter8?
We highlight a few powerful features of [Iter8](https://iter8.tools) that simplify performance testing of Knative apps.

1.  **Built-in latency and error-related metrics for HTTP and gRPC services.** Eliminates the need to set up and configure metrics databases during performance testing.
2.  **Well-defined notion of service-level objectives (SLOs)**. Makes it simple to define and verify SLOs in experiments.
3.  **Readiness checks.** The performance testing portion of the experiment begins only after the Knative service is ready.
4.  **HTML/text reports**. Promotes human understanding of experiment results through visual insights. 
5. **Assertions**. Verify whether the target app satisfies the specified SLOs or not after an experiment. Simplified automation in CI/CD/GitOps pipelines. Pipelines can branch off into different paths depending upon whether the assertions are true or false.

## Quick start

Install Iter8 CLI using `brew` as follows. You can also install using pre-built binaries, or use Iter8 inside a GitHub Actions pipeline as described [here](https://iter8.tools/0.11/getting-started/install/).

```shell
brew tap iter8-tools/iter8
brew install iter8@0.11
```

### Tutorial: Performance test for Knative HTTP service
Launch an Iter8 experiment that generates load for a Knative HTTP service and validates its service-level objectives (SLOs). The setup of this experiment is illustrated in the figure below.

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

#### View experiment report
Once the experiment completes (~ 20 secs), view the experiment report as follows.

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

### Tutorial: Performance test for Knative gRPC service
In this tutorial, we will launch an Iter8 experiment that generates load for a Knative gRPC service and validates its service-level objectives (SLOs). The setup of this experiment is illustrated in the figure below.

<img src="https://iter8.tools/0.11/tutorials/images/grpc.png" alt="gRPC performance test" width="800"/>

Use the [Knative (`kn`) CLI]() to update the Knative service deployed in the [above tutorial](#tutorial-performance-test-for-knative-http-service) to a gRPC service as follows.

```shell
kn service update hello \
--image docker.io/grpc/java-example-hostname:latest \
--port 50051 \
--revision-name=grpc
```

Launch the Iter8 experiment as follows.

```shell
iter8 k launch \
--set "tasks={ready,grpc,assess}" \
--set ready.ksvc=hello \
--set grpc.host="hello.default:50051" \
--set grpc.call="helloworld.Greeter.SayHello" \
--set grpc.protoURL="https://raw.githubusercontent.com/grpc/grpc-java/master/examples/example-hostname/src/main/proto/helloworld/helloworld.proto" \
--set grpc.data.name="frodo" \
--set assess.SLOs.upper.grpc/error-rate=0 \
--set assess.SLOs.upper.grpc/latency/mean=400 \
--set assess.SLOs.upper.grpc/latency/p90=500 \
--set runner=job
```

???+ note "About this experiment"
    This experiment consists of three [tasks](../getting-started/concepts.md#tasks), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md).
    
    The [ready](../user-guide/tasks/ready.md) task checks if the `hello` deployment exists and is available, and the `hello` service exists. 
    
    The [grpc](../user-guide/tasks/grpc.md) task sends call requests to the `helloworld.Greeter.SayHello` method of the cluster-local gRPC service with host address `hello.default:50051`, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). 
    
    The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. 
    
    This is a [single-loop](../getting-started/concepts.md#loops) [Kubernetes experiment](../getting-started/concepts.md#execution-environments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

Once the experiment completes (~ 20 secs), view the experiment report as described in the [earlier tutorial](#tutorial-performance-test-for-knative-http-service).

## What next?

Try the following enhancements and variations of the above tutorials.

1. Assert experiment outcomes, view experiment logs, and clean up the experiment as documented in [this example](https://iter8.tools/0.11/getting-started/your-first-experiment/).
2. Run experiments in your local environment instead of a Kubernetes cluster. Local experiment is documented in [this example](https://iter8.tools/0.11/tutorials/integrations/local/).
3. Configure the [`http` task](https://iter8.tools/0.11/user-guide/tasks/http/) in the [HTTP tutorial](#tutorial-performance-test-for-knative-http-service), and the [`grpc` task](https://iter8.tools/0.11/user-guide/tasks/grpc/) in the [gRPC tutorial](#tutorial-performance-test-for-knative-grpc-service) with load related parameters such as the number of requests, queries per second, duration, number of parallel connections, and various types of payloads. Use the [`grpc` task](https://iter8.tools/0.11/user-guide/tasks/grpc/) for performance testing of streaming gRPC.
4. Configure the [assess task](https://iter8.tools/0.11/user-guide/tasks/assess/) in the [HTTP tutorial](#tutorial-performance-test-for-knative-http-service) with additional SLOs based on [Iter8's built-in HTTP metrics](https://iter8.tools/0.11/user-guide/tasks/http/#metrics). Similarly, configure the [assess task](https://iter8.tools/0.11/user-guide/tasks/assess/) in the [gRPC tutorial](#tutorial-performance-test-for-knative-grpc-service) with additional SLOs based on [Iter8's built-in gRPC metrics](https://iter8.tools/0.11/user-guide/tasks/grpc/#metrics).
