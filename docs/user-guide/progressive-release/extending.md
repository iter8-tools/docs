---
template: main.html
---

# Extending progressive release to other deployment environments

The `release` chart can be easily extended to include other deployment environments. Please consider contributing any extensions to Iter8. We briefly descrive how to extend the chart for an [Knative](https://knative.dev/docs/) application. 

## Approach

Modify the release chart after forking the [Iter8 project](https://github.com/iter8-tools/iter8). The chart to be extended is in the `charts/release` subfolder. The file `release.yaml` is the starting point. For each valid environment, the chart contains a set of files defining the resources that should be created.  These may include:

- the application object(s)
- [routemaps](../routemap.md) for different traffic patterns
- configmaps used to specify request distribution (blue-green pattern only)
- a service defining a common entry for requests (if needed)

Note that the file naming helps identify related template files.

## Example (KNative Service)

For example, to implement a blue-green release for Knative servives, the following files could be added.

- `_knative-istio.tpl` - describing the application objects should be deployed
- `_knative-istio.version.ksvc.tpl` - describe the Knative service object that should be deployed for a version
- `_knative-istio.blue-green.tpl` - identifies any objects that should be deployed to support the blue-green traffic pattern
- `_knative-istio.blue-green.routemap.tpl` - the routemap definition

Finally, update `release.yaml` to include `knative-istio` as a valid option.

## Extend the controller

The Iter8 controller will need to be extended to give permission to Iter8 to watch Knative service objects. Configure the deployment of the controller to enable this, (re-)install the controller using the following additional options:

```shell
--set resourceTypes.ksvc.Group=serving.knative.dev \
--set resourceTypes.ksvc.Version=v1 \
--set resourceTypes.ksvc.Resource=services \
--set "resourceTypes.ksvc.conditions[0]=Ready"
```

## Using a modified chart

Reference the location of the local copy of the chart instead of using the `--repo` and `--version` options. For example assuming the location is `$CHART`, a deployment of 2 versions of the Knative `hello` service with a 30-70 traffic split would be:

```shell
cat <<EOF | helm upgrade --install hello $CHART -f -
environment: knative-istio
application:
  versions:
  - ksvcSpecification:
      spec:
        template:
          spec:
            containers:
            - image: ghcr.io/knative/helloworld-go:latest
              ports:
              - containerPort: 80
              env:
              - name: TARGET
                value: "v1"
    weight: 30
  - ksvcSpecification:
      spec:
        template:
          spec:
            containers:
            - image: ghcr.io/knative/helloworld-go:latest
              ports:
              - containerPort: 80
              env:
              - name: TARGET
                value: "v2"
    weight: 70
  strategy: blue-green
EOF
```

<!-- 
At the time of writing, this was tested locally as follows. These may not be minimal requirements.
(1) Created a rootful podman machine with 6 CPU and 24 GB memory. Set it run docker API. (used podman desktop)
alias docker=podman
(3) Created kind cluster (slightly modified from https://knative.dev/blog/articles/set-up-a-local-knative-environment-with-kind/)
export KIND_EXPERIMENTAL_PROVIDER=podman
cat > clusterconfig.yaml <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
    ## expose port 31080 of the node to port 80 on the host
  - containerPort: 31080
    hostPort: 80
    ## expose port 31443 of the node to port 443 on the host
  - containerPort: 31443
    hostPort: 443
EOF
kind create cluster --name knative --config clusterconfig.yaml
(4) Install Knative Serving (https://knative.dev/docs/install/yaml-install/serving/install-serving-with-yaml) inclusing Istio as the networking layer
(5) Run Iter8
helm upgrade --install --repo https://iter8-tools.github.io/iter8 --version 0.18 iter8 controller \
--set clusterScoped=true --set resourceTypes.ksvc.Group=serving.knative.dev \
--set resourceTypes.ksvc.Version=v1 \
--set resourceTypes.ksvc.Resource=services \
--set "resourceTypes.ksvc.conditions[0]=Ready"
(6) Deploy 2 versions of a Knative service with a 30-70 request distribution
cat <<EOF | helm upgrade --install hello $CHART -f -
environment: knative-istio
application:
  versions:
  - ksvcSpecification:
      spec:
        template:
          spec:
            containers:
            - image: ghcr.io/knative/helloworld-go:latest
              ports:
              - containerPort: 80
              env:
              - name: TARGET
                value: "v1"
    weight: 30
  - ksvcSpecification:
      spec:
        template:
          spec:
            containers:
            - image: ghcr.io/knative/helloworld-go:latest
              ports:
              - containerPort: 80
              env:
              - name: TARGET
                value: "v2"
    weight: 70
  strategy: blue-green
EOF
(7) Create sleep pod in cluster for testing and exec into it
curl -s https://raw.githubusercontent.com/iter8-tools/docs/v0.18.4/samples/kserve-serving/sleep.sh | sh -
kubectl exec --stdin --tty "$(kubectl get pod --sort-by={metadata.creationTimestamp} -l app=sleep -o jsonpath={.items..metadata.name} | rev | cut -d' ' -f 1 | rev)" -c sleep -- /bin/sh
(8) Send test requests
curl hello.default -s -D - | grep -e Hello -e app-version
-->
