---
template: main.html
---

# Running experiments in the local environment

Iter8 can be used to run experiments in the local environment as in the following example. In contrast to Kubernetes experiment launch, the local launch **does not** use the `k` subcommand.

```shell
iter8 launch \
--set "tasks={http,assess}" \
--set http.url=https://httpbin.org/get \
--set assess.SLOs.upper.http/latency-mean=50 \
--set assess.SLOs.upper.http/error-count=0
```

Local experiments are useful for development and debugging. Further, assuming that the HTTP and gRPC services can be reached from the local environment, they can be tested using local experiments even without access to the cluster where they may be hosted.

Use the [`iter8 assert`](../../user-guide/commands/iter8_assert.md) and [`iter8 report`](../../user-guide/commands/iter8_report.md) commands (without the `k` subcommand) to assert and report results from local experiments.