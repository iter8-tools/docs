---
template: main.html
---

# Iter8 GitHub Action

The [Iter8 GitHub Action](https://github.com/marketplace/actions/run-iter8-experiment) enables you to load test, benchmark, and validate gRPC services with service-level objectives (SLOs) inside GitHub Action workflows. This tutorial shows how.

## Basic example

```yaml
- uses: iter8-tools/iter8-action@v1
  with:
    chart: load-test-grpc
    valuesFile: experiment-config.yaml
```

A sample `experiment-config.yaml` is as follows.

```yaml
# An earlier step in the workflow is assumed to have started the gRPC service
host: 127.0.0.1:50051
call: helloworld.Greeter.SayHello
protoURL: https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto
```

Details of the configuration parameters that can be set are [here](basicusage.md). Default values are [here](https://github.com/iter8-tools/hub/blob/main/charts/load-test-grpc/values.yaml).

## Complete example

A complete GitHub Actions workflow which exercises the Iter8 Action using the `load-test-grpc` experiment is available as part of [the Iter8 docs repo](https://github.com/iter8-tools/docs). Run this example as follows.

1. Fork the Iter8 docs repo: <https://github.com/iter8-tools/docs> to your organization, *myorg*.

2. If necessary, enable GitHub Actions workflows on the your forked repo by navigating to <https://github.com/myorg/docs/actions>.

3. Run the workflow:

    - Navigate to the **Actions** tab of your forked repository: <https://github.com/myorg/docs/actions>.
    - Select the workflow **end-to-end tests**.
    - Click the **Run workflow** button.

4. When the workflow has completed, there will be a new entry for the execution. Select the latest, then the entry for the **local grpc tests** job. View the logs generated by the Iter8 GitHub Action by clicking on the steps labeled *Run iter8-tools/iter8-action@v1*.

## Iter8 Action inputs

The list of inputs that can be configured for the Iter8 GitHub Action is documented [here](https://github.com/iter8-tools/iter8-action#action-inputs).