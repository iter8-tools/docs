---
template: main.html
---

# A/B Experiments

This tutorial describes an [A/B testing](../../user-guide/topics/ab_testing.md) experiment for a backend component.

![A/B/n experiment](images/abn.png)

***

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md). Understand the main [concepts](../../getting-started/concepts.md) behind Iter8 experiments.
 
## Launch Iter8 A/B/n service

Deploy the Iter8 A/B/n service. When deploying the service, specify which Kubernetes resource types to watch for each application. To watch for versions of the *backend* application in the *default* namespace, configure the service to watch for Kubernetes service and deployment resources:

```shell
helm install --repo https://iter8-tools.github.io/hub iter8-abn abn \
--set "apps.default.backend.resources={service,deployment}"
```

??? warn "Assumptions"
    To simplify specification, Iter8 assumes certain conventions:

    - The baseline track identifier is the application name
    - Track identifiers associated with candidate versions are of the form `<application_name>-candidate-<index>`
    - All resource objects for all versions are deployed in the same namespace
    - There is only 1 resource object of a given type in each version
    - The name of each object in the version associated with the baseline track is the application name
    - The name of each object in the version associate with a candidate track is of the form  `<application_name>-candidate-<index>` where index is 1, 2, etc.


## Deploy the sample application

Deploy both the frontend and backend components of the application as described in each tab:

=== "frontend"
    Install the frontend component using an implementation in the language of your choice:

    === "node"
        ```shell
        kubectl create deployment frontend --image=iter8/abn-sample-frontend-node:0.13
        kubectl expose deployment frontend --name=frontend --port=8090
        ```

    === "Go"
        ```shell
        kubectl create deployment frontend --image=iter8/abn-sample-frontend-go:0.13
        kubectl expose deployment frontend --name=frontend --port=8090
        ```
    
    The frontend component is implemented to call *Lookup()* before each call to the backend component. The frontend componet uses the returned track identifier to route the request to a version of the backend component.

=== "backend"
    Deploy version *v1* of the *backend* component, associating it with the track identifier *backend*.

    ```shell
    kubectl create deployment backend --image=iter8/abn-sample-backend:0.13-v1
    kubectl label deployment backend app.kubernetes.io/version=v1

    kubectl expose deployment backend --name=backend --port=8091
    ```

## Generate load

Generate load. In separate shells, port-forward requests to the frontend component and generate load for multiple users.  A [script](https://raw.githubusercontent.com/iter8-tools/docs/main/samples/abn-sample/generate_load.sh) is provided to do this. To use it:
    ```shell
    kubectl port-forward service/frontend 8090:8090
    ```
    ```shell
    curl -s https://raw.githubusercontent.com/iter8-tools/docs/main/samples/abn-sample/generate_load.sh | sh -s --
    ```

## Deploy a candidate version

Deploy version *v2* of the *backend* component, associating it with the track identifier *backend-candidate-1*.

```shell
kubectl create deployment backend-candidate-1 --image=iter8/abn-sample-backend:0.13-v2
kubectl label deployment backend-candidate-1 app.kubernetes.io/version=v2

kubectl expose deployment backend-candidate-1 --name=backend-candidate-1 --port=8091
```

Until the candidate version is ready; that is, until all expected resources are deployed and available, calls to *Lookup()* will return only the *backend* track identifier.
Once the candidate version is ready, *Lookup()* will return both track identifiers so that requests will be distributed between versions.

## Launch experiment

```shell
iter8 k launch \
--set abnmetrics.application=default/backend \
--set "tasks={abnmetrics}" \
--set runner=cronjob \
--set cronjobSchedule="*/1 * * * *"
```

??? note "About this experiment"
    This experiment periodically (in this case, once a minute) reads the `abn` metrics associated with the *backend* application component in the *default* namespace. These metrics are written by the frontend component using the *WriteMetric()* interface as a part of processing user requests.

## Inspect experiment report

Inspect the metrics:

```shell
iter8 k report
```

??? note "Sample output from report"
    ```
    Experiment summary:
    *******************

    Experiment completed: true
    No task failures: true
    Total number of tasks: 1
    Number of completed tasks: 1
    Number of completed loops: 3

    Latest observed values for metrics:
    ***********************************

    Metric                   | backend (v1) | backend-candidate-1 (v2)
    -------                  | -----        | -----
    abn/sample_metric/count  | 35.00        | 28.00
    abn/sample_metric/max    | 99.00        | 100.00
    abn/sample_metric/mean   | 56.31        | 52.79
    abn/sample_metric/min    | 0.00         | 1.00
    abn/sample_metric/stddev | 28.52        | 31.91
    ```
The output allows you to compare the versions against each other and select a winner. Since the experiment runs periodically, the values in the report will change over time.

Once a winner is identified, the experiment can be terminated, the winner can be promoted, and the candidate version(s) can be deleted.

To delete the experiment:

```shell
iter8 k delete
```

## Promote candidate version

Delete the candidate version:

```shell
kubectl delete deployment backend-candidate-1 
kubectl delete service backend-candidate-1
```

Update the version associated with the baseline track identifier *backend*:

```shell
kubectl set image deployment/backend abn-sample-backend=iter8/abn-sample-backend:0.13-v2
kubectl label --overwrite deployment/backend app.kubernetes.io/version=v2
```

## Cleanup

### Delete sample application

```shell
kubectl delete \
deploy/frontend deploy/backend deploy/backend-candidate-1 \
service/frontend service/backend service/backend-candidate-1
```

### Uninstall the A/B/n service

```shell
helm delete iter8-abn
```
