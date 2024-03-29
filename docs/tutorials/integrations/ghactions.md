---
template: main.html
---

# Trigger a GitHub Actions workflow

Iter8 provides a [`github` task](../../user-guide/performance/tasks/github.md) that can be used in a performance test to send a [`repository_dispatch`](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#repository_dispatch) event to trigger workflows in the default branch of a GitHub repository.

## Example

In this example, you will run the [Your first performance test](../../getting-started/first-performance.md) but at the end of the performance test, Iter8 will trigger a workflow on GitHub.

In this simple example, the workflow will simply print out a test summary that it will receive with the `repository_dispatch`. In a more sophisticated scenario, the workflow could, for example, read from the test summary and determine what to do next.

To summarize what will happen, you will create a new GitHub repository, add a workflow that will respond to the `github` task, set up and run a performance test, and check if the workflow was triggered.

The `github` task requires the name of a repository, the name of the owner, as well as an authentication token in order to send the `repository_dispatch`. To see a full list of the `github` task parameters, see [here](../../user-guide/performance/tasks/github.md#parameters).

1. Create a new repository on GitHub.
2. Add the following workflow.

```yaml
name: iter8 `github` task test
on:
  repository_dispatch:
    types: iter8
jobs:
  my-job:
    runs-on: ubuntu-latest
    steps:
      - run: 'echo "payload: ${{ toJson(github.event.client_payload) }}"'
```

Note that this workflow has one job that will print out the `client_payload`. The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/iter8/v1.1.1/templates/notify/_payload-github.tpl) is configured with `client_payload` set to `.Report`, a summary of the performance test.

Also note that the `on.repository_dispatch.types` is set to `iter8`. The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/iter8/v1.1.1/templates/notify/_payload-github.tpl) is configured with `event_type` set to `iter8`. This indicates that once the `repository_dispatch` has been sent, only workflows on the default branch with `on.repository_dispatch.types` set to `iter8` will be triggered.

3. Create a GitHub [personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) for the `token` parameter.
4. Ensure that you have a Kubernetes cluster and the [`kubectl` CLI](https://kubernetes.io/docs/reference/kubectl/). You can create a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
5. Install the Iter8 controller

    --8<-- "docs/getting-started/install.md"
    
6. Deploy the sample HTTP service in the Kubernetes cluster.
```shell
kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
kubectl expose deploy httpbin --port=80
```
7. Launch the performance test using the `github` task with the appropriate values.
```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 1.1 httpbin-test iter8 \
--set "tasks={http,github}" \
--set http.url=http://httpbin.default/get \
--set github.owner=<GitHub owner> \
--set github.repo=<GitHub repository> \
--set github.token=<GitHub token>
```
8. Verify that the workflow has been triggered after the performance test has completed.

??? note "Some variations and extensions of the `github` task"
    The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/iter8/v1.1.1/templates/notify/_payload-github.tpl) sends a summary of the performance test. 
    
    In your workflow, you can read from the report and use that data for control flow or use snippets of that data in different actions. For example, you can check to see if there have been any task failures and take alternative actions.

    You do not need to use the default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/iter8/v1.1.1/templates/notify/_payload-github.tpl). You can provide your own payload by overriding the default of the `payloadTemplateURL`.