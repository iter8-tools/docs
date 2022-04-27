---
template: main.html
---

# Experiment namespace and group

Iter8 experiments like [`load-test-http`](../../tutorials/load-test-http/kubernetesusage.md) and [`load-test-grpc`](../../tutorials/load-test-grpc/kubernetesusage.md) can be run within Kubernetes. Such experiments are launched within a Kubernetes namespace, and associated with a unique *group* within that namespace.

For example, consider the following invocation:

```shell
iter8 k launch -c load-test-http -g hbin \
--set url=http://httpbin.default \
--set SLOs.http/latency-mean=50
```

In the above invocation, the `iter8 k launch` implicitly specifies the namespace as `default`, and explicitly specifies the group as `hbin`. If the group name is not specified explicitly, then it is set to `default`.

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

1.  Run multiple experiments concurrently within a Kubernetes namespace. These experiments may be associated with the same app or with different apps.
2.  Replace a currently running experiment in Kubernetes with a new one. When you invoke `iter8 k launch`, any previous experiment runs within the group are wiped out and replaced with a new run.

## How groups work

Under the covers, Iter8 implements each experiment group as a [Helm release](https://helm.sh/docs/glossary/#release) and each new experiment run within the group as an [update](https://helm.sh/docs/glossary/#release-number-release-version) of that release.