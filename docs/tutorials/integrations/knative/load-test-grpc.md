---
template: main.html
---

# Benchmark and Validate a Knative gRPC service

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
    3. Update the Knative service deployed above to a gRPC service as follows.
    ```shell
    kn service update hello \
    --image docker.io/grpc/java-example-hostname:latest \
    --port 50051 \
    --revision-name=grpc
    ```

***

Benchmark and validate SLOs for the Knative gRPC service by launching an Iter8 experiment.

```shell
iter8 launch -c load-test-grpc \
--set-string host="hello.default.127.0.0.1.sslip.io:50051" \
--set-string call="helloworld.Greeter.SayHello" \
--set-string protoURL="https://raw.githubusercontent.com/grpc/grpc-java/master/examples/example-hostname/src/main/proto/helloworld/helloworld.proto" \
--set data.name="frodo" \
--set SLOs.grpc/error-rate=0 \
--set SLOs.grpc/latency/mean=400 \
--set SLOs.grpc/latency/p90=500 \
--set SLOs.grpc/latency/p'97\.5'=600
```

Please refer to [the usage documentation for the `load-test-grpc` experiment chart](../../load-test-grpc/basicusage.md) that describes how to parameterize this experiment, assert SLOs, and view experiment reports.    
