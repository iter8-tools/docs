---
template: main.html
---

# Use Iter8 to send a message to a Slack channel

Iter8 provides a `slack` task that sends a message to a Slack channel using a [webhook](https://api.slack.com/messaging/webhooks).

The `slack` task has the following parameters:

| Name | Type | Required | Default value | Description |
| ---- | ---- | -------- | ------------- | ----------- |
| url | string | Yes | N/A | URL to the Slack webhook |
| payloadTemplateURL | string | No | [https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-slack.tpl](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-slack.tpl) | URL to a payload template |
| softFailure | bool | No | true | Indicates the task and experiment should not fail if the task cannot successfully send the request |

## Example

1. Create a new channel in your Slack organization.
2. Create a Slack app, enable incoming webhooks, and create a new incoming webhook. See [here](https://api.slack.com/messaging/webhooks).
3. Ensure that you have a Kubernetes cluster and the [`kubectl` CLI](https://kubernetes.io/docs/reference/kubectl/). You can create a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
4. Deploy the sample HTTP service in the Kubernetes cluster.
```shell
kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
kubectl expose deploy httpbin --port=80
```
5. Launch the experiment with the `slack` task with the appropriate values.
```shell
iter8 launch \
--noDownload \
--set "tasks={http,assess,slack}" \
--set http.url=http://127.0.0.1/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0 \
--set slack.url="<Slack webhook>" \
--set slack.method=POST \
--set runner=job
```
7. Verify that the message has been sent after the experiment has completed.

??? note "Some variations and extensions of the `slack` task"
    The default `slack` task [payload](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-slack.tpl) sends the entirety of the experiment report.

    However, you do not need to use the default payload. You can provide your own payload by overriding the default of the `payloadTemplateURL`.

    For example, you can create a payload that selectively print out parts of the experiment report instead of the whole thing.

    You can also use Slack's [Block Kit](https://api.slack.com/block-kit/building) in order to make your message more sophisticated. You can use markdown, create different sections, or add interactivity, such as buttons.