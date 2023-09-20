---
template: main.html
---

# slack

Send an performance test summary in a message to a Slack channel using a [incoming webhook](https://api.slack.com/messaging/webhooks). 

## Usage Example

```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 0.18 httpbin-test iter8 \
--set "tasks={http,slack}" \
--set http.url=http://httpbin.default/get \
--set slack.url=<Slack webhook> \
--set slack.method=POST
```

See [here](../../tutorials/integrations/slack.md#use-iter8-to-send-a-message-to-a-slack-channel) for a more in-depth tutorial.

## Parameters

| Name | Type | Required | Default value | Description |
| ---- | ---- | -------- | ------------- | ----------- |
| url | string | Yes | N/A | URL to the Slack webhook |
| payloadTemplateURL | string | No | [https://raw.githubusercontent.com/iter8-tools/iter8/v0.18.3/templates/notify/_payload-slack.tpl](https://raw.githubusercontent.com/iter8-tools/iter8/v0.18.3/templates/notify/_payload-slack.tpl) | URL to a payload template |
| softFailure | bool | No | true | Indicates the performance test should not fail if the task cannot successfully send the request |

## Default payload

The payload will determine what will be contained in the Slack message. The [default payload template](https://raw.githubusercontent.com/iter8-tools/iter8/v0.18.3/templates/notify/_payload-slack.tpl) of the `slack` task is to send a performance test summary in text form.

However, if you would like to use a different payload template, simply set a `payloadTemplateURL` and Iter8 will not use the default.