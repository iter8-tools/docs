---
template: main.html
---

# A/B Testing with the Iter8 SDK

This tutorial describes how to do A/B testing of a backend component using the [Iter8 SDK](../user-guide/abn/about.md). 

![A/B/n testing](../tutorials/images/abn.png)

***

???+ warning "Before you begin"
    1. Ensure that you have a Kubernetes cluster and the [`kubectl`](https://kubernetes.io/docs/reference/kubectl/) and [`helm`](https://helm.sh/) CLIs. If using a local cluster (for example, [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/)), we recommend providing the cluster with at least 16GB of memory.
    2. Have Grafana available. For example, Grafana can be installed on your cluster as follows:
    ```shell
    kubectl create deploy grafana --image=grafana/grafana
    kubectl expose deploy grafana --port=3000
    ```
 
## Install the Iter8 controller

--8<-- "docs/getting-started/install.md"

## Deploy the sample application

A simple sample two-tier application using the Iter8 SDK is provided. Note that only the frontend component uses the Iter8 SDK. Deploy both the frontend and backend components of this application as described in each tab:

=== "Frontend"
    Install the frontend component using an implementation in the language of your choice:

    === "Node"
        ```shell
        kubectl create deployment frontend --image=iter8/abn-sample-frontend-node:0.17.3
        kubectl expose deployment frontend --name=frontend --port=8090
        ```

    === "Go"
        ```shell
        kubectl create deployment frontend --image=iter8/abn-sample-frontend-go:0.17.3
        kubectl expose deployment frontend --name=frontend --port=8090
        ```
    
    The frontend component is implemented to call `Lookup()` before each call to the backend component. The frontend component uses the returned version number to route the request to the recommended version of the backend component.

=== "Backend"
    Release an initial version of the backend named `backend`:

    ```shell
    cat <<EOF | helm upgrade --install backend --repo https://iter8-tools.github.io/iter8 release --version 1.1 -f -
    environment: deployment
    application: 
      port: 8091
      versions:
      - metadata:
          name: backend
        image: iter8/abn-sample-backend:0.17-v1
    EOF
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
Recommendation: {"Id":19,"Name":"sample","Source":"backend-74ff88c76d-nb87j"}
```

## Deploy candidate

A candidate version of the *backend* component can be deployed simply by adding a second version to the list of versions:

```shell
cat <<EOF | helm upgrade --install backend --repo https://iter8-tools.github.io/iter8 release --version 1.1 -f -
environment: deployment
application: 
  port: 8091
  versions:
  - metadata:
      name: backend
    image: iter8/abn-sample-backend:0.17-v1
  - metadata:
      name: backend-candidate-1
    image: iter8/abn-sample-backend:0.17-v2
EOF
```

While the candidate version is deploying, `Lookup()` will return only the version index number `0`; that is, the first, or primary, version of the model.
Once the candidate version is ready, `Lookup()` will return both `0` and `1`, the indices of both versions, so that requests can be distributed across both versions.

Once both backend versions are responding to requests, the output of the load generator will include recommendations from the candidate version. In this example, you should see something like:

```
Recommendation: {"Id":19,"Name":"sample","Source":"backend-candidate-1-56cb7cd5cf-bkrjv"}
```

## Compare versions using Grafana

Inspect the metrics using Grafana. If Grafana is deployed to your cluster, port-forward requests as follows:

```shell
kubectl port-forward service/grafana 3000:3000
```

Open Grafana in a browser by going to [http://localhost:3000](http://localhost:3000) and login. The default username/password are `admin`/`admin`.

[Add a JSON API data source](http://localhost:3000/connections/datasources/marcusolsson-json-datasource) `default/backend` with the following parameters:

* URL: `http://iter8.default:8080/abnDashboard`
* Query string: `namespace=default&application=backend`

[Create a new dashboard](http://localhost:3000/dashboards) by *import*. Copy and paste the contents of the [`abn` Grafana dashboard](https://raw.githubusercontent.com/iter8-tools/iter8/v1.1.1/grafana/abn.json) into the text box and *load* it. Associate it with the JSON API data source above.

The Iter8 dashboard allows you to compare the behavior of the two versions of the backend component against each other and select a winner. Since user requests are being sent by the load generation script, the values in the report may change over time. The Iter8 dashboard will look like the following:

![A/B dashboard](../tutorials/images/abnDashboard.png)

Once you identify a winner, it can be promoted, and the candidate version deleted.

## Promote candidate

To promote the candidate version (`backend-candidate-1`), re-release the application, updating the image of the primary (the first) version to use the image of the candidate version and remove the candidate version:

```shell
cat <<EOF | helm upgrade --install backend --repo https://iter8-tools.github.io/iter8 release --version 1.1 -f -
environment: deployment
application: 
  port: 8091
  versions:
  - metadata:
      name: backend
    image: iter8/abn-sample-backend:0.17-v2
EOF
```

Calls to `Lookup()` will now recommend that all traffic be sent to the new primary version `backend` (currently serving the promoted version of the code).

The output of the load generator will again show just `backend_0`:

```
Recommendation: {"Id":19,"Name":"sample","Source":"backend-74ff88c76d-nb87j"}
```

## Cleanup

Delete the sample application:

```shell
kubectl delete svc/frontend deploy/frontend
helm delete backend
```

Uninstall the Iter8 controller:

--8<-- "docs/getting-started/uninstall.md"

If you installed Grafana, you can delete it as follows:

```shell
kubectl delete svc/grafana deploy/grafana
```

***

Congratulations! :tada: You completed your first A/B test with Iter8.

***
