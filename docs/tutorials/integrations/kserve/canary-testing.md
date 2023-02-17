---
template: main.html
---

# Canary Testing

This tutorial shows how easy it is validate SLOs for multiple versions of a model in [KServe](https://kserve.github.io/website/0.10/) when fetching metrics from a metrics database like Prometheus. We show this using the `sklearn-iris` model used to describe [canary rollouts](https://kserve.github.io/website/0.10/modelserving/v1beta1/rollout/canary-example/) in KServe. 

???+ "Before you begin"
    1. Try [your first experiment](../../../getting-started/your-first-experiment.md). Understand the main [concepts](../../../getting-started/concepts.md) behind Iter8 experiments.
    2. Ensure that you have the [kubectl](https://kubernetes.io/docs/reference/kubectl/) CLI.
    3. Have access to a cluster running [KServe](https://kserve.github.io/website). You can create a [KServe Quickstart](https://kserve.github.io/website/0.10/get_started/#before-you-begin) environment as follows:
    ```shell
    curl -s "https://raw.githubusercontent.com/kserve/kserve/release-0.10/hack/quick_install.sh" | bash
    ```
    4. Install Prometheus monitoring for KServe [using these instructions](https://github.com/kserve/kserve/tree/master/docs/samples/metrics-and-monitoring#install-prometheus).

## Experiment Setup

Deploy two models to compare and generate load against them. We follow the instructions for the [KServe canary rollout example](https://kserve.github.io/website/0.10/modelserving/v1beta1/rollout/canary-example/) to deploy the models.

### Create InferenceService for Initial Model

```shell
kubectl apply -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris"
spec:
  predictor:
    model:
      modelFormat:
        name: sklearn
      storageUri: "gs://kfserving-examples/models/sklearn/1.0/model"
EOF
```

### Update InferenceService with a Canary Model

```shell
kubectl apply -f - <<EOF
apiVersion: "serving.kserve.io/v1beta1"
kind: "InferenceService"
metadata:
  name: "sklearn-iris"
spec:
  predictor:
    canaryTrafficPercent: 10
    model:
      modelFormat:
        name: sklearn
      storageUri: "gs://kfserving-examples/models/sklearn/1.0/model-2"
EOF
```

### Generate Load

Port forward requests to the cluster:

```shell
INGRESS_GATEWAY=$(kubectl get svc --namespace istio-system --selector="app=istio-ingressgateway" --output jsonpath='{.items[0].metadata.name}')
kubectl port-forward --namespace istio-system svc/$INGRESS_GATEWAY 8080:80
```

Send prediction requests to the inference service. The following script generates about one request a second. In a production cluster, this step is not required since your inference service will receive requests from real users.

```shell
SERVICE_HOSTNAME="sklearn-iris.default.example.com"
# kubectl get inferenceservice sklearn-iris -o jsonpath='{.status.url}' | cut -d "/" -f 3

cat <<EOF > "./iris-input.json"
{
  "instances": [
    [6.8,  2.8,  4.8,  1.4],
    [6.0,  3.4,  4.5,  1.6]
  ]
}
EOF

while true; do 
  curl -H "Host: ${SERVICE_HOSTNAME}" \
    http://localhost:8080/v1/models/sklearn-iris:predict \
    -d @./iris-input.json
  sleep 1
done
```

## Launch Iter8 Experiment

Launch an Iter8 experiment inside the Kubernetes cluster:

```shell
iter8 k launch \
--set "tasks={ready,custommetrics,assess}" \
--set ready.isvc=sklearn-iris \
--set ready.timeout=180s \
--set custommetrics.templates.kserve-prometheus="https://gist.githubusercontent.com/kalantar/adc6c9b0efe483c00b8f0c20605ac36c/raw/c4562e87b7ac0652b0e46f8f494d024307bff7a1/kserve-prometheus.tpl" \
--set custommetrics.values.labels.service_name=sklearn-iris-predictor-default \
--set 'custommetrics.versionValues[0].labels.revision_name=sklearn-iris-predictor-default-00002' \
--set 'custommetrics.versionValues[1].labels.revision_name=sklearn-iris-predictor-default-00001' \
--set "custommetrics.values.latencyPercentiles={50,75,90,95}" \
--set assess.SLOs.upper.kserve-prometheus/error-count=0 \
--set assess.SLOs.upper.kserve-prometheus/latency-mean=25 \
--set assess.SLOs.upper.kserve-prometheus/latency-p'90'=40 \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

??? note "About this experiment"
    This experiment consists of three [tasks](../../../getting-started/concepts.md#iter8-experiment), namely, [ready](../../../user-guide/tasks/ready.md), [custommetrics](../../../user-guide/tasks/custommetrics.md) and [assess](../../../user-guide/tasks/assess.md). 

    The [ready](../../../user-guide/tasks/ready.md) task checks if the `sklearn-iris` InferenceService exists and is `Ready`. 
    
    The [custommetrics](../../../user-guide/tasks/custommetrics.md) task reads metrics from a Prometheus service as defined by the [template](https://gist.githubusercontent.com/kalantar/adc6c9b0efe483c00b8f0c20605ac36c/raw/c4562e87b7ac0652b0e46f8f494d024307bff7a1/kserve-prometheus.tpl). The template is parameterised using labels for service and revision name. You can identify the revision names from the `InferenceService`:

    ```shell
    kubectl get isvc sklearn-iris -o json \
    | jq -r '.status.components.predictor.traffic | .[] | .revisionName'
    ```
    The service name is the prefix (remove the trailing `-ddddd`).


    The [assess](../../../user-guide/tasks/assess.md) task verifies if the model satisfies the specified SLOs:

    - there are no errors
    - the mean latency of the prediction does not exceed 25 msec, and
    - the 90th percentile latency for prediction does not exceed 40 msec. 
    
    This is a [multi-loop](../../../getting-started/concepts.md#iter8-experiment) [Kubernetes experiment](../../../getting-started/concepts.md#kubernetes-experiments); its [runner](../../../getting-started/concepts.md#runners) is `cronjob`. The `cronjobSchedule` expression specifies the frequency of the experiment execution -- periodically refreshing the metric values and performing SLO validation using the updated values.

***

You can assert experiment outcomes, view an experiment report, and view experiment logs as described in [your first experiment](../../../getting-started/your-first-experiment.md).

## Clean up

To clean up, delete the Iter8 experiment:

```shell
iter8 k delete
```

Remove the `InferenceService` and the request data:
```shell
kubectl delete inferenceservice sklearn-iris
rm ./iris-input.json
```

You can remove Prometheus using [these instructions](https://github.com/kserve/kserve/tree/master/docs/samples/metrics-and-monitoring#removal).
