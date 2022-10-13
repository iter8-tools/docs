# SLO Validation with Chaos Injection

Perform a joint Iter8 and [LitmusChaos](https://litmuschaos.io/) experiment. This joint experiment enables you to verify if an app continues to be resilient (satisfies SLOs) in the midst of chaos (pod kill).

In the tutorial, the app consists of a Kubernetes service and deployment. The chaos experiment kills the app's pods intermittently. At the same time, the Iter8 experiment performs a load test of the app and validates its [service-level objectives (SLOs)](../../getting-started/concepts.md#service-level-objectives). 


![Chaos with SLO Validation](images/slo-validation-chaos.png)

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md). Understand the main [concepts](../../getting-started/concepts.md) behind Iter8 experiments.
    2. Ensure that you have the [kubectl](https://kubernetes.io/docs/reference/kubectl/) CLI.
    3. Install [Litmus](https://litmuschaos.io/) in Kubernetes using [these steps](https://docs.litmuschaos.io/docs/getting-started/installation).
    4. Create the `httpbin` deployment file.
    ```shell
    cat <<EOF >>deploy.yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: httpbin
      labels:
        app: httpbin
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: httpbin
      template:
        metadata:
          labels:
            app: httpbin
        spec:
          containers:
          - name: httpbin
            image: kennethreitz/httpbin
            ports:
            - containerPort: 80
          initContainers:
          - name: init-myservice
            image: busybox:1.28
            command: ['sh', '-c', 'sleep 1']
    EOF
    ```
    5. Create the `httpbin` deployment.
    ```shell
    kubectl apply -f deploy.yaml
    ```
    6. Create the `httpbin` service.
    ```shell
    kubectl expose deploy httpbin --port=80
    ```

***

## Launch experiments
Launch the LitmusChaos and Iter8 experiments as described below.
=== "LitmusChaos"
    ```shell
    helm install httpbin litmuschaos \
    --repo https://iter8-tools.github.io/hub/ \
    --set applabel='app=httpbin' \
    --set totalChaosDuration=3600 \
    --set chaosInterval=5
    ```

    ??? note "About this LitmusChaos experiment"
        This experiment consists of three [tasks](../getting-started/concepts.md#iter8-experiment), namely, [ready](../user-guide/tasks/ready.md), [grpc](../user-guide/tasks/grpc.md), and [assess](../user-guide/tasks/assess.md). 
        
        The [ready](../user-guide/tasks/ready.md) task checks if the `hello` deployment exists and is available, and the `hello` service exists. 
        
        The [grpc](../user-guide/tasks/grpc.md) task sends call requests to the `helloworld.Greeter.SayHello` method of the cluster-local gRPC service with host address `hello.default:50051`, and collects [Iter8's built-in gRPC load test metrics](../user-guide/tasks/grpc.md#metrics). 
        
        The [assess](../user-guide/tasks/assess.md) task verifies if the app satisfies the specified SLOs: i) there are no errors, ii) the mean latency of the service does not exceed 50 msec, and iii) the `97.5`th percentile latency does not exceed 200 msec. 
        
        This is a [single-loop](../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../getting-started/concepts.md#kubernetes-experiments) where all the previously mentioned tasks will run once and the experiment will finish. Hence, its [runner](../getting-started/concepts.md#runners) value is set to `job`.

=== "Iter8" 
    ```shell
    iter8 k launch \
    --set "tasks={ready,http,assess}" \
    --set ready.deploy=httpbin \
    --set ready.service=httpbin \
    --set ready.chaosengine=litmuschaos-httpbin \
    --set ready.timeout=60s \
    --set http.url=http://httpbin.default/get \
    --set http.duration=30s \
    --set http.qps=20 \
    --set assess.SLOs.upper.http/latency-mean=50 \
    --set assess.SLOs.upper.http/latency-p99=100 \
    --set assess.SLOs.upper.http/error-count=0 \
    --set runner=job
    ```

    ??? note "About this Iter8 experiment"
        Please see [here](../../../getting-started/your-first-experiment/#launch-experiment).

*** 

## Observe experiments
Observe the LitmusChaos and Iter8 experiments as follows. The chaos and Iter8 experiments 

=== "LitmusChaos"
    Verify that the phase of the chaos experiment is `Running`.
    ```shell
    kubectl get chaosresults/litmuschaos-httpbin-pod-delete -n default \
    -ojsonpath='{.status.experimentStatus.phase}'
    ```

    ??? note "On completion of the LitmusChaos experiment"
        After the LitmusChaos experiment completes (in ~3600 sec), the phase of the experiment will change to `Completed`. At that point, you can verify that the chaos experiment returns a `Pass` verdict. The `Pass` verdict states that the application is still running after chaos has ended.
        ```shell
        kubectl get chaosresults/litmuschaos-httpbin-pod-delete -n default \
        -o=jsonpath='{.status.experimentStatus.verdict}'
        ```

=== "Iter8"
    Due to chaos injection, and the fact that the number of replicas of the app in the deployment manifest is set to 1, the SLOs are not expected to be satisfied within the Iter8 experiment. Verify this is the case.

    ```shell
    # the SLOs assertion is expected to fail
    iter8 k assert -c completed -c nofailure -c slos --timeout 30s
    ```

    For a more detailed report of the Iter8 experiment, run the `report` command.
    ```shell
    iter8 k report
    ```

***

## Cleanup experiments

Clean up the LitmusChaos and Iter8 experiments as described below.

=== "LitmusChaos"
    ```shell
    helm uninstall httpbin
    ```

=== "Iter8"
    ```shell
    iter8 k delete
    ```

***

## Scale app and retry
Scale up the app so that replica count is increased to 3. 
```shell
kubectl scale --replicas=3 -n default deploy/httpbin
```

The scaled app is now more resilient. Performing the same experiments as above will now result in SLOs being satisfied and a winner being found. Retry [this step](launch-experiments) and [this step](observe-experiments). You should now find that SLOs are satisfied.

***

## Cleanup

Cleanup the app as follows.

```shell
kubectl delete svc/httpbin
kubectl delete deploy/httpbin
```

Cleanup the experiments as described [here](#5-cleanup-experiments).

Uninstall LitmusChaos from your cluster as described [here](https://docs.litmuschaos.io/docs/user-guides/uninstall-litmus/).

***

??? note "Some variations and extensions of this experiment"
    1. Reuse the above experiment with *your* app by replacing the `httpbin` app with *your* app, and modifying the Helm values appropriately.
    2. gRPC. Variations of both HTTP and gRPC load tests.
    3. Litmus makes it possible to inject [over 51 types of Chaos](https://hub.litmuschaos.io/). Modify the Helm chart to use any of these other types of chaos experiments.
