---
template: main.html
---

# A/B Experiment

Given two versions of an application component, compare their metrics to determine a winner. 

<p align='center'>
  <img alt-text="abn experiment" src="../images/abn.png" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md). Understand the main [concepts](../../getting-started/concepts.md) behind Iter8 experiments.
    2. Deploy the sample two-tier application in your Kubernetes cluster.
    ```shell
    curl -qs https://raw.githubusercontent.com/kalantar/ab-example/main/frontend/deploy.yaml \
    | sed -e "s#FRONTEND_TAG#kalantar/frontend-node#" \
    | kubectl apply -f -

    curl -qs https://raw.githubusercontent.com/kalantar/ab-example/main/backend/deploy.yaml \
    | sed -e "s#\$NAME#backend#" \
    | sed -e "s#\$VERSION#v1#" \
    | sed -e "s#\$TRACK#default#" \
    | kubectl apply -f -
    kubectl label deployment backend iter8.tools/abn=true
    ```
    3. Generate load.  In separate shells, port-forward requests to the frontend service and generate load for multiple users.  For example:
    ```shell
    kubectl port-forward svc/frontend 8090:8090
    ```
    ```shell
    watch -x curl -s localhost:8090/buy -H "X-User: foo"
    ```
    ```shell
    watch -x curl -s localhost:8090/buy -H "X-User: foobar"
    ```
    4. Add (and/or) update the Iter8 helm repository:
    ```shell
    helm repo add iter8 https://iter8-tools.github.io/hub
    ```
    ```shell
    helm repo update
    ```
***

## Launch Iter8 A/B(/n) service

If not already deployed, deploy the Iter8 A/B(/n) service. Specify which Kubernetes resource(s) to watch in which namespaces.

```shell
helm install iter8-abn iter8/iter8-abn \
--set resources='{deployments,services}' \
--set namespaces='{default}'     --set image=iter8/iter8:0.11
```

??? warn "Currently supported resource types"
    Iter8 currently supports watching Kubernetes deployments and services as well as Knative services.
    The Helm chart used to deploy the service can be easily extended to support additional resource types.

## Deploy candidate version

```shell
curl -qs https://raw.githubusercontent.com/kalantar/ab-example/main/backend/deploy.yaml \
| sed -e "s#\$NAME#backend-candidate#" \
| sed -e "s#\$VERSION#v2#" \
| sed -e "s#\$TRACK#candidate#" \
| kubectl apply -f -
```

## Mark candidate ready

Once the candidate version is ready to receive user traffic, for example, when the pods are `Ready`, label the deployment object as a participant in an A/B(/n) experiment:.

```shell
kubectl label deployment backend-candidate iter8.tools/abn=true
```

Once labeled, subsequent requests for a backend will include the candidate version as a valid version to receive traffic.

To terminate traffic to the candidate version, remove the `iter8.tools/abn` label.

??? note "Details about labels"
    The Iter8 watches resources where the label `iter8-tools/abn` set to `true`. The following labels are expected to be present identifying the role of the resource in an A/B(/n) experiment. Note that an application _version_ might be composed of multiple resources. Iter8 expects only 1 of these resources to be labeled

    1. `app.kubernetes.io/name`: application name
    2. `app.kubernetes.io/version`: version name
    3. `iter8-tools/track`: track identifier (used for routing)

## Launch experiment

```shell
iter8 k launch \
--set abnmetrics.application=default/backend \
--set "tasks={abnmetrics}" \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

This experiment periodically (once a minute) reads the `abn` metrics associated with the `backend` application in the `default` namespace.

## Inspect experiment report

```shell
iter8 k report
```

??? note "Sample output from report"
    ```
    Experiment summary:
    *******************

    Experiment completed: false
    No task failures: true
    Total number of tasks: 1
    Number of completed tasks: 18

    Latest observed values for metrics:
    ***********************************

    Metric                   | candidate | default
    -------                  | -----     | -----
    abn/sample_metric/count  | 765.00    | 733.00
    abn/sample_metric/max    | 100.00    | 100.00
    abn/sample_metric/mean   | 50.11     | 49.64
    abn/sample_metric/min    | 0.00      | 0.00
    abn/sample_metric/stddev | 28.63     | 29.25
    ```
The output allows you to compare the versions against each other and select a winner. Since the experiment runs periodically, you should expect the values in the report to change over time.

Once a winner is identified, it can be promoted and the canidiate versions can be deleted.

## Promote the candidate version [to be deleted]

### Update the default version

Redeploy the `backend` deployment using the new image (`v2`) as the `default` track:

```shell
curl -qs https://raw.githubusercontent.com/kalantar/ab-example/main/backend/deploy.yaml \
| sed -e "s#\$NAME#backend#" \
| sed -e "s#\$VERSION#v2#" \
| sed -e "s#\$TRACK#default#" \
| kubectl apply -f -
```

### Remove the candidate version

```shell
kubectl delete deployment backend-candidate
```

## Cleanup

### Delete the experiment

```shell
iter8 k delete
```

### Delete the A/B(/n) service

```shell
helm delete iter8-abn
```

### Delete sample application

```shell
kubectl delete \
deploy/frontend deploy/backend deploy/backend-candidate \
svc/frontend svc/backend svc/backend-candidate \
secret/backend.iter8abnmetrics
```
