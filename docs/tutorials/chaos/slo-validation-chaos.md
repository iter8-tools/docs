# SLO Validation with Chaos Injection

Perform a joint Iter8 and [LitmusChaos](https://litmuschaos.io/) experiment. This joint experiment enables you to verify if an app continues to be resilient (satisfies SLOs) in the midst of chaos (pod kill).

In the tutorial, the app consists of a Kubernetes service and deployment. The chaos experiment kills the app's pods intermittently. At the same time, the Iter8 experiment performs a load test of the app and validates its [service-level objectives (SLOs)](../../getting-started/concepts.md#service-level-objectives). 


![Chaos with SLO Validation](images/slo-validation-chaos.png)

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md). Understand the main [concepts](../../getting-started/concepts.md) behind Iter8 experiments.
    2. Ensure that you have the [Helm v3.x](https://helm.sh/docs/intro/install/) and [kubectl](https://kubernetes.io/docs/reference/kubectl/) CLIs.
    3. Install [Litmus](https://litmuschaos.io/) in Kubernetes cluster. You can perform a quick install as follows.
    ```shell
    helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/
    kubectl create ns litmus
    helm install chaos litmuschaos/litmus --namespace=litmus --set portal.frontend.service.type=NodePort
    ```
    Verify that Litmus is install correctly as described [here](https://docs.litmuschaos.io/docs/getting-started/installation/#verify-your-installation).

***

## 1. Create app
Deploy the `httpbin` app in the Kubernetes cluster.

```shell
kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
kubectl expose deploy httpbin --port=80
```

***

## 2. Launch experiments
Launch the LitmusChaos and Iter8 experiments as described below.
=== "LitmusChaos"
    ```shell
    helm install httpbinchaos litmuschaos \
    --repo https://iter8-tools.github.io/hub/ \
    --set applabel='app.kubernetes.io/name=hello'
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
    --set ready.chaosengine=httpbinchaos \
    --set ready.timeout=60s \
    --set http.url=http://httpbin.default/get \
    --set assess.SLOs.upper.http/latency-mean=50 \
    --set assess.SLOs.upper.http/error-count=0 \
    --set runner=job
    ```

    ??? note "About this Iter8 experiment"
        Please see [here](../../../getting-started/your-first-experiment/#launch-experiment).

## 3. Observe experiments
Observe the LitmusChaos and Iter8 experiments as follows. The chaos and Iter8 experiments 

=== "LitmusChaos"
    Verify that the phase of the chaos experiment is `Completed`.
    ```shell
    kubectl get chaosresults/httpbinchaos -n default \
    -ojsonpath='{.status.experimentStatus.phase}'
    ```

    Verify that the chaos experiment returned a `Pass` verdict. The `Pass` verdict states that the application is still running after chaos has ended.
    ```shell
    kubectl get chaosresults/httpbinchaos -n default \
    -o=jsonpath='{.status.experimentStatus.verdict}'
    ```

=== "Iter8"
    Due to chaos injection, and the fact that the number of replicas of the app in the deployment manifest is set to 1, the SLOs are not expected to be satisfied within the Iter8 experiment. Verify this is the case.

    ```shell
    # the SLOs assertion is expected to fail
    iter8 k assert -c completed -c nofailure -c slos --timeout 120s
    ```

    For a more detailed report of the Iter8 experiment, run the `report` command.
    ```shell
    iter8 k report
    ```

## 5. Cleanup experiments

Clean up the LitmusChaos and Iter8 experiments as described below.

=== "LitmusChaos"
    ```shell
    helm uninstall httpbinchaos
    ```

=== "Iter8"
    ```shell
    iter8 k delete
    ```


## 6. Scale app and retry
Scale up the app so that replica count is increased to 2. 
```shell
kubectl scale --replicas=2 -n default deploy/httpbin
```

The scaled app is now more resilient. Performing the same experiments as above will now result in SLOs being satisfied and a winner being found. Retry [this step](#2-launch-experiments) and [this step](#3-observe-experiments). You should now find that SLOs are satisfied.

## 7. Cleanup

Cleanup the app as follows.

```shell
kubectl delete svc httpbin
kubectl delete deploy httpbin
```

Cleanup the experiments as described [here](#5-cleanup-experiments).

Delete LitmusChaos from your cluster as follows.
```shell
helm uninstall chaos -n litmus
kubectl delete ns litmus
```

***

??? note "Some variations and extensions of this experiment"
    1. Reuse the above experiment with *your* app by replacing the `httpbin` app with *your* app, and modifying the Helm values appropriately.
    2. gRPC. Variations of both HTTP and gRPC load tests.
    3. Litmus makes it possible to inject [over 51 types of Chaos](https://hub.litmuschaos.io/). Modify the Helm chart to use any of these other types of chaos experiments.
