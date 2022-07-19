---
template: main.html
---

# Namespaces and groups for Kubernetes experiments

[Kubernetes experiments](../../getting-started/concepts.md#kubernetes-experiments) are launched within a [namespace](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/), and are associated with a unique *group* within that namespace.

For example, consider the following invocation:

```shell
iter8 k launch -g hbin \
--set "tasks={http,assess}" \
--set http.url=http://httpbin.default/get \
--set assess.SLOs.upper.http/latency-mean=50
```

In the above invocation, the `iter8 k launch` command implicitly specifies the namespace as `default`, and explicitly specifies the group as `hbin`. If the group name is not specified explicitly, then it is set to `default`. The namespace can be specified explicitly using the `-n` or `--namespace` flags (see [here](../commands/iter8_k_launch.md#options-inherited-from-parent-commands)).

The following example illustrates the relationship between namespaces, groups, and experiments.

```shell
.
├── namespace1
│   ├── group-a
│   │   └── experiment
│   ├── group-b
│   │   └── experiment
│   └── group-c
│       └── experiment
├── namespace2
│   ├── group-a
│   │   └── experiment
│   ├── group-b
│   │   └── experiment
│   └── group-c
│       └── experiment
└── namespace3
    └── group-x
        └── experiment
```

## Use-cases

1.  Run multiple experiments concurrently within a Kubernetes namespace by associating them with distinct groups. These experiments may be associated with the same app or with different apps.
2.  Replace a currently running experiment in Kubernetes with a new one. When you invoke `iter8 k launch`, any previous experiment executions within the group is wiped out and replaced with a fresh experiment that starts to execute.

## How groups work

Under the covers, Iter8 implements each experiment group as a [Helm release](https://helm.sh/docs/glossary/#release) and each new experiment run within the group as an [update](https://helm.sh/docs/glossary/#release-number-release-version) of that release.