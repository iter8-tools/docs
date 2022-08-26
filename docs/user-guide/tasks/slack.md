---
template: main.html
---

# slack

Send the experiment report in a message to a Slack channel using a [incoming webhook](https://api.slack.com/messaging/webhooks). 

## Usage Example

```shell
iter8 launch \
--noDownload \
--set "tasks={http,assess,slack}" \
--set http.url=http://127.0.0.1/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0 \
--set slack.url=<Slack webhook> \
--set slack.method=POST \
--set runner=job
```

See [here](../../tutorials/integrations/slack.md#use-iter8-to-send-a-message-to-a-slack-channel) for a more in-depth tutorial.

## Parameters

| Name | Type | Required | Default value | Description |
| ---- | ---- | -------- | ------------- | ----------- |
| url | string | Yes | N/A | URL to the Slack webhook |
| payloadTemplateURL | string | No | [https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-slack.tpl](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-slack.tpl) | URL to a payload template |
| softFailure | bool | No | true | Indicates the task and experiment should not fail if the task cannot successfully send the request |

## Default payload

The payload will determine what will be contained in the Slack message. The [default payload templae](https://raw.githubusercontent.com/iter8-tools/hub/main/templates/_payload-slack.tpl) of the `slack` task is to send the experiment report in text form.

However, if you would like to use a different payload template, simply set a `payloadTemplateURL` and Iter8 will not use the default.