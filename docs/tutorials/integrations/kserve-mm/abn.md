---
template: main.html
---

# A/B Testing a backend ML model

This tutorial describes how to do A/B testing as part of the release of a backend ML model hosted on [KServe ModelMesh](https://github.com/kserve/modelmesh) using the [Iter8 SDK](../../../user-guide/abn/about.md). 

![A/B/n testing](../../images/abn.png)

***

???+ warning "Before you begin"
    1. Ensure that you have the [`kubectl`](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs installed.
    2. Have access to a cluster running [KServe ModelMesh Serving](https://github.com/kserve/modelmesh-serving). For example, you can create a modelmesh-serving [Quickstart](https://github.com/kserve/modelmesh-serving/blob/release-0.11/docs/quickstart.md) environment.  If using the Quickstart environment, your default namespace will be changed to `modelmesh-serving`. If using a local cluster (for example, [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/)), we recommend providing the cluster with at least 16GB of memory.
    3. Have Grafana available. For example, Grafana can be installed on your cluster as follows:
    ```shell
    kubectl create deploy grafana --image=grafana/grafana
    kubectl expose deploy grafana --port=3000
    ```
 
## Install the Iter8 controller

--8<-- "docs/getting-started/install.md"

## Deploy the sample application

A simple sample two-tier application using the Iter8 SDK is provided. Note that only the frontend component uses the Iter8 SDK. Deploy both the frontend and backend components:

### Frontend

The frontend component uses the Iter8 SDK method `Lookup()` before each call to the backend (ML model). The frontend uses the returned version number to route the request to the recommended version of backend.

```shell
kubectl create deployment frontend --image=iter8/abn-sample-mm-frontend-go:0.17.3
kubectl expose deployment frontend --name=frontend --port=8090
```

### Backend

The backend application component is an ML model. Release it using the Iter8 `release` chart:

```shell
cat <<EOF | helm upgrade --install backend --repo https://iter8-tools.github.io/iter8 release --version 1.1 -f -
environment: kserve-modelmesh-istio
application: 
  metadata:
    labels:
      app.kubernetes.io/name: backend
    annotations:
      serving.kserve.io/secretKey: localMinIO
  modelFormat: sklearn
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
EOF
```

Wait for the backend model to be ready:

```shell
kubectl wait --for condition=ready isvc/backend-0 --timeout=600s
```

## Generate load

In one shell, port-forward requests to the frontend component:
    ```shell
    kubectl port-forward service/frontend 8090:8090
    ```

In another shell, run a script to generate load from multiple users:
    ```shell
    curl -s https://raw.githubusercontent.com/iter8-tools/docs/v0.18.3/samples/abn-sample/generate_load.sh | sh -s --
    ```
 
The load generator and sample frontend application outputs the backend that handled each recommendation. With just one version is deployed, all requests are handled by `backend-0`. In the output you will see something like:

```
Recommendation: backend-0__isvc-3642375d03
```

## Deploy candidate

A candidate version of the model can be deployed simply by adding a second version to the list of versions:

```shell
cat <<EOF | helm upgrade --install backend --repo https://iter8-tools.github.io/iter8 release --version 1.1 -f -
environment: kserve-modelmesh-istio
application: 
  metadata:
    labels:
      app.kubernetes.io/name: backend
    annotations:
      serving.kserve.io/secretKey: localMinIO
  modelFormat: sklearn
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v0
    storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
EOF
```

??? note "About the candidate"
    In this tutorial, the model source (field `storageUri`) for the candidate version is the same as for the primary version of the model. In a real example, this would be different. The version label (`app.kubernetes.io/version`) can be used to distinguish between versions.

Until the candidate version is ready, calls to `Lookup()` will return only the version index number `0`; that is, the first, or primary, version of the model.
Once the candidate version is ready, `Lookup()` will return both `0` and `1`, the indices of both versions, so that requests can be distributed across both versions.

Once both backend versions are responding to requests, the output of the load generator will include recommendations from the candidate version. In this example, you should see something like:

```
Recommendation: backend-1__isvc-3642375d03
```

## Compare versions using Grafana

Inspect the metrics using Grafana. If Grafana is deployed to your cluster, port-forward requests as follows:

```shell
kubectl port-forward service/grafana 3000:3000
```

Open Grafana in a browser by going to [http://localhost:3000](http://localhost:3000) and login. The default username/password are `admin`/`admin`.

[Add a JSON API data source](http://localhost:3000/connections/datasources/marcusolsson-json-datasource) `modelmesh-serving/backend` with the following parameters:

* URL: `http://iter8.modelmesh-serving:8080/abnDashboard`
* Query string: `namespace=modelmesh-serving&application=backend`

[Create a new dashboard](http://localhost:3000/dashboards) by *import*. Copy and paste the contents of the [`abn` Grafana dashboard](https://raw.githubusercontent.com/iter8-tools/iter8/v1.1.1/grafana/abn.json) into the text box and *load* it. Associate it with the JSON API data source above.

The Iter8 dashboard allows you to compare the behavior of the two versions of the backend component against each other and select a winner. Since user requests are being sent by the load generation script, the values in the report may change over time. The Iter8 dashboard will look like the following:

![A/B dashboard](../../images/abnDashboard.png)

Once you identify a winner, it can be promoted, and the candidate version deleted.

## Promote candidate

The candidate can be promoted by redefining the primary version and removing the candidate:

```shell
cat <<EOF | helm upgrade --install backend --repo https://iter8-tools.github.io/iter8 release --version 1.1 -f -
environment: kserve-modelmesh-istio
application: 
  metadata:
    labels:
      app.kubernetes.io/name: backend
    annotations:
      serving.kserve.io/secretKey: localMinIO
  modelFormat: sklearn
  versions:
  - metadata:
      labels:
        app.kubernetes.io/version: v1
    storageUri: s3://modelmesh-example-models/sklearn/mnist-svm.joblib
EOF
```

??? note "What is different?"
    The version label (`app.kubernetes.io/version`) of the primary version was updated. In a real world example, `storageUri` would also be updated (with that from the candidate version).

Calls to `Lookup()` will now recommend that all traffic be sent to the new primary version `backend-0` (currently serving the promoted version of the code).

The output of the load generator will again show just `backend_0`:

```
Recommendation: backend-0__isvc-3642375d03
```

## Cleanup

Delete the backend:

```shell
helm delete backend
```

Delete the frontend:

```shell
kubectl delete deploy/frontend svc/frontend
```

Uninstall Iter8 controller:

--8<-- "docs/getting-started/uninstall.md"

If you installed Grafana, you can delete it as follows:

```shell
kubectl delete svc/grafana deploy/grafana
```
