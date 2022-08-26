---
template: main.html
---

# notify

Send a notification (HTTP request) to another web service.

The `notify` task is used to create task templates that notify particular services. 

For example, the `slack` task sends the experiment report in a message to a particular channel using a [incoming webhook](https://api.slack.com/messaging/webhooks). The `github` task triggers GitHub workflows via a [repository_dispatch](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#repository_dispatch) and sends the experiment report.

## Usage Example
In this example, the `slack` task sends a message to a Slack channel using the Slack hook. See [here](../../tutorials/integrations/slack.md#use-iter8-to-send-a-message-to-a-slack-channel) for more information.

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

In this example, the `github` task triggers GitHub workflows in a particular GitHub repository. See [here](../../tutorials/integrations/ghactions.md#use-iter8-to-trigger-a-github-actions-workflow) for more information.

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

## Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| url | string | URL of the notification hook |
| method | string | HTTP method that needs to be used |
| params | map[string]string | Set of HTTP parameters that need to be sent |
| headers | map[string]string | Set of HTTP headers that need to be sent |
| payloadTemplateURL | string | URL of the request payload template that should be used |
| softFailure | bool | Indicates the task and experiment should not fail if the task cannot successfully send the request |

## Payload

If the request requires a payload, then the `payloadTemplateURL` should be used. The `payloadTemplateURL` is a URL to a template of a payload.

For example, this is a [URL](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-slack.tpl) of a payload template that can be used to send Slack messages.

Additionally, this is a [URL](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-github.tpl) of a payload template that can be used to trigger GitHub workflows. 

The template has access to the experiment report via `.Report`. Iter8 also enables [Sprig](https://pkg.go.dev/github.com/Masterminds/sprig) template functions that can be used to transform the report in various ways.

For example,

```
Raw JSON:
{{ .Report | toJson }}

Formatted JSON:
{{ .Report | toPrettyJson }}

Escaped JSON:
{{ regexReplaceAll "\"" (regexReplaceAll "\n" (.Report | toPrettyJson) "\\n") "\\\""}}
```