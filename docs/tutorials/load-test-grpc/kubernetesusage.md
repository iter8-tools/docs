---
template: main.html
---

# Benchmark and Validate Kubernetes gRPC Services

Benchmark, and validate a gRPC service inside a Kubernetes cluster using the  [`load-test-grpc` experiment](basicusage.md). The experiment is run inside the cluster. The gRPC service may be externally exposed or local to the cluster.

<p align='center'>
  <img alt-text="load-test-grpc" src="../images/kubernetesusage.png" width="90%" />
</p>

***

???+ warning "Before you begin"
    1. Try the [basic `load-test-grpc` tutorial](basicusage.md).
    2. Try the [Kubernetes usage tutorial for the `load-test-http` experiment](../load-test-http/kubernetesusage.md).
    3. Ensure that you have a Kubernetes cluster and the [`kubectl` CLI](https://kubernetes.io/docs/reference/kubectl/). You may run a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
    4. Deploy the sample gRPC service in the Kubernetes cluster.
    ```shell
    kubectl create deploy hello --image=docker.io/grpc/java-example-hostname:latest --port=50051
    kubectl expose deploy hello --port=50051
    ```

***

## Launch experiment

Launch a `load-test-grpc` experiment inside the Kubernetes cluster. Note that the gRPC host in this experiment is `hello.default`, which refers to a hostname inside the Kubernetes cluster, specifically, the `hello` service in the `default` namespace.

```shell
iter8 k launch -c load-test-grpc \
--set host="hello.default:50051" \
--set call="helloworld.Greeter.SayHello" \
--set protoURL="https://raw.githubusercontent.com/grpc/grpc-go/master/examples/helloworld/helloworld/helloworld.proto" \
--set ready.deploy=hello \
--set ready.service=hello \
--set ready.timeout=60s 
```

***

## Similarities to `load-test-http`

You can configure this experiment with various parameter values (in particular, you can specify SLOs) as described in the [basic `load-test-grpc` tutorial](basicusage.md). You can assert experiment conditions, view experiment reports, view experiment logs, and cleanup the experiment as described in the [Kubernetes usage tutorial for `load-test-http`](../load-test-http/kubernetesusage.md).