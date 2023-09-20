---
template: main.html
---

# Use Iter8 to send a message to a Slack channel

Iter8 provides a [`slack` task](../../user-guide/tasks/slack.md)  that sends a message to a Slack channel using a [webhook](https://api.slack.com/messaging/webhooks).

## Example

In this example, you will run the [Your first performance test](../../getting-started/first-performance.md) but at the end of the performance test, Iter8 will send a message on Slack. 

The message will simply contain a summary of the performance test in text form. However, you can easily construct a more sophisticated message by providing your own payload template.

This task could provide important updates on a performance test over Slack, for example a summary at the end of the test.

To summarize what will happen, you will create a new channel on Slack and configure a webhook, set up and run a performance test, and check if a message was sent to the channel.

The `slack` task requires the URL of the Slack webhook. To see a full list of the `github` task parameters, see [here](../../user-guide/tasks/slack.md#parameters).

1. Create a new channel in your Slack organization.
2. Create a Slack app, enable incoming webhooks, and create a new incoming webhook. See [here](https://api.slack.com/messaging/webhooks).
3. Ensure that you have a Kubernetes cluster and the [`kubectl` CLI](https://kubernetes.io/docs/reference/kubectl/). You can create a local Kubernetes cluster using tools like [Kind](https://kind.sigs.k8s.io/) or [Minikube](https://minikube.sigs.k8s.io/docs/).
4. Install the Iter8 controller

    --8<-- "docs/tutorials/installiter8controller.md"
    
5. Deploy the sample HTTP service in the Kubernetes cluster.
```shell
kubectl create deploy httpbin --image=kennethreitz/httpbin --port=80
kubectl expose deploy httpbin --port=80
```
6. Launch the performance test with the `slack` task with the appropriate values.
```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 0.18 httpbin-test iter8 \
--set "tasks={http,slack}" \
--set http.url=http://httpbin.default/get \
--set slack.url=<Slack webhook> \
--set slack.method=POST
```
7. Verify that the message has been sent after the performance test has completed.

??? note "Some variations and extensions of the `slack` task"
    The default `slack` task [payload](https://raw.githubusercontent.com/iter8-tools/iter8/v0.18.3/templates/notify/_payload-slack.tpl) sends a summary of the performance test.

    However, you do not need to use the default payload. You can provide your own payload by overriding the default of the `payloadTemplateURL`.

    For example, you can also use Slack's [Block Kit](https://api.slack.com/block-kit/building) to create more sophisticated messages. You can use markdown, create different sections, or add interactivity, such as buttons.