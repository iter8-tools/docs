---
template: main.html
---

# Running experiments in the local environment

Iter8 can be used to run experiments in your local environment as in the following example.

```shell
iter8 launch \
--set "tasks={http,assess}" \
--set http.url=https://httpbin.org/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0
```

Local experiments are especially useful for development and debugging. Local experiments can also be used for load testing and SLO validation of HTTP and gRPC services -- as long as the services are reachable, you can test them with Iter8 even if you do not have access to the cluster where those services might be hosted.

You can use the [`iter8 assert`](../../user-guide/commands/iter8_assert.md) and [`iter8 report`](../../user-guide/commands/iter8_report.md) commands (without the `k` subcommand) to assert and report results from local experiments.