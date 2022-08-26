---
template: main.html
---

# Use Iter8 in a GitHub Actions workflow

Install the latest version of the Iter8 CLI using `iter8-tools/iter8@v0.11`. Once installed, the Iter8 CLI can be used as documented in various tutorials. For example:

```yaml linenums="1"
# install Iter8 CLI
- uses: iter8-tools/iter8@v0.11
# launch a local experiment
- run: |
    iter8 launch --set "tasks={http}" --set http.url=http://httpbin.org/get
# launch an experiment inside Kubernetes;
# this assumes that your Kubernetes cluster is accessible 
# from the GitHub Actions pipeline
- run: |
    iter8 k launch --set "tasks={http}" \
    --set http.url=http://httpbin.org/get \
    --set runner=job
```

# Use Iter8 to trigger a GitHub Actions workflow

Iter8 provides a `github` task that sends a [repository_dispatch](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#repository_dispatch) which can trigger the workflows in the default branch of a GitHub repository.

The `github` task has the following parameters:

| Name | Type | Required | Default value | Description |
| ---- | ---- | -------- | ------------- | ----------- |
| owner | string | Yes | N/A | Owner of the GitHub repository |
| repo | string | Yes | N/A | GitHub repository |
| token | string | Yes | N/A | Authorization token |
| payloadTemplateURL | string | No | [https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-github.tpl](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-github.tpl) | URL to a payload template |
| softFailure | bool | No | true | Indicates the task and experiment should not fail if the task cannot successfully send the request |

## Example

1. Create a new repository on GitHub.
2. Add the following workflow.
```yaml
name: iter8 notify test
on:
  repository_dispatch:
    types: iter8
jobs:
  my-job:
    runs-on: ubuntu-latest
    steps:
      - run: 'echo "payload: ${{ toJson(github.event.client_payload) }}"'
```

    Note that this workflow has one job that will print out the `client_payload`. The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-github.tpl) is configured with `client_payload` set to experiment report. This means that this job will simply print out the entire experiment report.

    Also note that the `on.repository_dispatch.types` is set to `iter8`. The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-github.tpl) is configured with `event_type` set to `iter8`. This indicates that once the `repository_dispatch` has been sent, only workflows on the default branch with `on.repository_dispatch.types` set to `iter8` will be triggered.

3. Create a GitHub [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) for the `token` parameter.
4. Ensure that you have a Kubernetes cluster and the [`kubectl` CLI](https://kubernetes.io/docs/reference/kubectl/). You can create a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
5. Deploy the sample HTTP service in the Kubernetes cluster.
```shell
kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
kubectl expose deploy httpbin --port=80
```
6. Launch the experiment with the `github` task with the appropriate values.
```shell
iter8 launch \
--noDownload \
--set "tasks={http,assess,github}" \
--set http.url=http://127.0.0.1/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0 \
--set github.owner=<GitHub owner> \
--set github.repo=<GitHub repository> \
--set github.token=<GitHub token> \
--set runner=job
```
7. Verify that the workflow has been triggered after the experiment has completed.

??? note "Some variations and extensions of the `github` task"
    1. The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-github.tpl) sends the entirety of the experiment report. In your workflow, you can read from the report and use that data for control flow or use snippets of that data in different actions. For example, you can check to see if there have been any task failures during the experiment and perform different actions.
    2. You do not need to use the default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-github.tpl). You can provide your own payload by overriding the default of the `payloadTemplateURL`. For example, instead of sending the entirety of the experiment report, you can create a payload template that only sends a subset.