---
template: main.html
---

# GitHub Actions

There are two ways that you can use Iter8 with GitHub Actions. You can [run Iter8 CLI within a GitHub Actions workflow](#use-iter8-in-a-github-actions-workflow) and you can also [use Iter8 to trigger a GitHub Actions workflow](#use-iter8-to-trigger-a-github-actions-workflow) from an experiment.

# Use Iter8 in a GitHub Actions workflow

Install the latest version of the Iter8 CLI using `iter8-tools/iter8@v0.12`. Once installed, the Iter8 CLI can be used as documented in various tutorials. For example:

```yaml linenums="1"
# install Iter8 CLI
- uses: iter8-tools/iter8@v0.12
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

# Trigger a GitHub Actions workflow from an Iter8 experiment

Iter8 provides a [`github` task](../..user-guide/tasks/github) that sends a [repository_dispatch](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#repository_dispatch) which can trigger the workflows in the default branch of a GitHub repository.

## Example

In this example, you will run the [Your First Experiment](../../getting-started/your-first-experiment.md) but at the end of the experiment, Iter8 will trigger a workflow on GitHub.

In this simple example, the workflow will simply print out the experiment report that it will receive with the `repository_dispatch`. In a more sophisticated scenario, the workflow could, for example, read from the experiment report and based on whether or not task failured or metrics did not meet SLOs, determine what to do next. In a GitOps scenario, it could make changes to the Git repo, promote winners, or restart pipelines among other things.

To summarize what will happen, you will create a new GitHub repo, add a workflow that will respond to the `github` task, set up and run an experiment, and check if the workflow was triggered.

The `github` task requires the name of a repo, the name of the owner, as well as an authentication token in order to send the `repository_dispatch`. To see a full list of the `github` task parameters, see [here](../../user-guide/tasks/github.md#parameters).

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

    Note that this workflow has one job that will print out the `client_payload`. The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/iter8-0.12.2/templates/notify/_payload-github.tpl) is configured with `client_payload` set to experiment report. This means that this job will simply print out the entire experiment report.

    Also note that the `on.repository_dispatch.types` is set to `iter8`. The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/iter8-0.12.2/templates/notify/_payload-github.tpl) is configured with `event_type` set to `iter8`. This indicates that once the `repository_dispatch` has been sent, only workflows on the default branch with `on.repository_dispatch.types` set to `iter8` will be triggered.

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
    1. The default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/iter8-0.12.2/templates/notify/_payload-github.tpl) sends the entirety of the experiment report. In your workflow, you can read from the report and use that data for control flow or use snippets of that data in different actions. For example, you can check to see if there have been any task failures during the experiment and perform different actions.
    2. You do not need to use the default `github` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/iter8-0.12.2/templates/notify/_payload-github.tpl). You can provide your own payload by overriding the default of the `payloadTemplateURL`. For example, instead of sending the entirety of the experiment report, you can create a payload template that only sends a subset.
    3. Try a [multi-loop experiment](../../getting-started/concepts.md#runner) with an [`if` parameter](../../user-guide/tasks/github.md#if-parameter) to control when the `github` task is run. 
    
        A multi-loop experiment will allow you to run the tasks on a recurring basis, allowing you to monitor your app over a course of time. For example:

        ```shell
        iter8 k launch \
        --set "tasks={http,assess,github}" \
        --set http.url=http://127.0.0.1/get \
        --set assess.SLOs.upper.http/latency-mean=50 \
        --set assess.SLOs.upper.http/error-count=0 \
        --set github.owner=<GitHub owner> \
        --set github.repo=<GitHub repository> \
        --set github.token=<GitHub token> \
        --set runner=cronjob \
        --set cronjobSchedule="*/1 * * * *"
        ```

        This will run `http`, `assess`, and `github` tasks every minute. If you would like to run the `github` task only during the 10th loop, use the `if` parameter.

        ```diff
          iter8 k launch \
          --set "tasks={http,assess,github}" \
          --set http.url=http://127.0.0.1/get \
          --set assess.SLOs.upper.http/latency-mean=50 \
          --set assess.SLOs.upper.http/error-count=0 \
          --set github.owner=<GitHub owner> \
          --set github.repo=<GitHub repository> \
          --set github.token=<GitHub token> \
        + --set github.if="Result.NumLoops == 10"
          --set runner=cronjob \
          --set cronjobSchedule="*/1 * * * *"
        ```
