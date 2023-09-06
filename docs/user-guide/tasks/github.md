---
template: main.html
---

# github

Trigger GitHub workflows via a [repository_dispatch](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#repository_dispatch).

A `repository_dispatch` will trigger workflows in the default branch of the GitHub repository. By default, an summary of the performance test will also be sent.

## Usage Example

```shell
helm upgrade --install \
--repo https://iter8-tools.github.io/iter8 --version 0.16 httpbin-test iter8 \
--set "tasks={http,github}" \
--set http.url=http://httpbin.default/get \
--set github.owner=<GitHub owner> \
--set github.repo=<GitHub repository> \
--set github.token=<GitHub token>
```

See [here](../../tutorials/integrations/ghactions.md#use-iter8-to-trigger-a-github-actions-workflow) for a more in-depth tutorial.

## Parameters

| Name | Type | Required | Default value | Description |
| ---- | ---- | -------- | ------------- | ----------- |
| owner | string | Yes | N/A | Owner of the GitHub repository |
| repo | string | Yes | N/A | GitHub repository |
| token | string | Yes | N/A | Authorization token |
| payloadTemplateURL | string | No | [https://raw.githubusercontent.com/iter8-tools/iter8/v0.16.5/templates/notify/_payload-github.tpl](https://raw.githubusercontent.com/iter8-tools/iter8/v0.16.5/templates/notify/_payload-github.tpl) | URL to a payload template |
| softFailure | bool | No | true | Indicates the performance test should not fail if the task cannot successfully send the request |

## Default payload

A `repository_dispatch` requires a payload that contains the type of the event. 

The [default payload template](https://raw.githubusercontent.com/iter8-tools/iter8/v0.16.5/templates/notify/_payload-github.tpl) will set the `event_type` to `iter8`. In addition, it will also provide a performance test summary in the `client_payload`, which means that this data will be accessible in the GitHub workflow via `${{ toJson(github.event.client_payload) }}`.

However, if you would like to use a different payload template, simply set a `payloadTemplateURL` and Iter8 will not use the default.