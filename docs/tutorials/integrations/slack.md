---
template: main.html
---

# Use Iter8 to send a message to a Slack channel

Iter8 provides a [`slack` task](../..user-guide/tasks/slack)  that sends a message to a Slack channel using a [webhook](https://api.slack.com/messaging/webhooks).

## Example

In this example, you will run the [Your First Experiment](../../getting-started/your-first-experiment.md) but at the end of the experiment, Iter8 will send a message on Slack. 

The message will simply contain the experiment report in text form. However, you can easily construct a more sophisticated message by providing your own payload template.

This task could provide important updates on an experiment over Slack, for example a summary at the end of an experiment.

To summarize what will happen, you will create a new channel on Slack and configure a webhook, set up and run an experiment, and check if a message was sent to the channel.

The `slack` task requires the URL of the Slack webhook. To see a full list of the `github` task parameters, see [here](../../user-guide/tasks/slack.md#parameters).

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
--set "tasks={http,assess,slack}" \
--set http.url=http://127.0.0.1/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0 \
--set slack.url=<Slack webhook> \
--set slack.method=POST \
--set runner=job
```
7. Verify that the message has been sent after the experiment has completed.

??? note "Some variations and extensions of the `slack` task"
    The default `slack` task [payload](https://raw.githubusercontent.com/iter8-tools/iter8/v0.11.10/charts/iter8/templates/_payload-slack.tpl) sends the entirety of the experiment report.

    However, you do not need to use the default payload. You can provide your own payload by overriding the default of the `payloadTemplateURL`.

    For example, you can create a payload that selectively print out parts of the experiment report instead of the whole thing.

    You can also use Slack's [Block Kit](https://api.slack.com/block-kit/building) in order to make your message more sophisticated. You can use markdown, create different sections, or add interactivity, such as buttons.