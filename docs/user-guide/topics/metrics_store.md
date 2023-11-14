---
template: main.html
---

# Metrics store

One of Iter8's key advantages is that it incorporates its own metrics store simplifying the set up and execution of A/B/n and performance tests. Iter8 currently supports the following databases:

- BadgerDB
- Redis

Iter8 uses BadgerDB by default. Note, however, that BadgerDB is not suitable for prodcution use and is only suitable for a single instance of Iter8. Support for other databases are in the works. See [below](#contribute-a-new-metrics-store-implementation) for details on how to contribute additional implementations.

## Using Redis as the metrics store 

We assume that Redis is deployed. For example, for a basic deployment:

```shell
kubectl create deploy redis --image=redis/redis-stack:latest --port=6379
kubectl expose deploy redis --port=6379
```

Run Iter8 with the metrics store implementation set to `redis` and specify its endpoint:

```shell
helm upgrade --install --repo https://iter8-tools.github.io/iter8 --version 0.18 iter8 controller  \
--set clusterScoped=true \
--set metrics.implementation=redis \
--set metrics.redis.addresss=redis:6379
```

## Contribute a new metrics store implementation

To contribute a new metrics store implementation:

1. Create an [issue](https://github.com/iter8-tools/iter8/issues) for discussion.

2. Submit a pull request on the [Iter8 project](https://github.com/iter8-tools/iter8) with the following updates:

    - Create sub-folder in [storage](https://github.com/iter8-tools/iter8/tree/master/storage) and provide an implementation of this [interface](https://github.com/iter8-tools/iter8/blob/master/storage/interface.go) including test cases.

    - Add a new case to [metrics.GetClient()](https://github.com/iter8-tools/iter8/blob/master/metrics/client.go).

    - Update [go.mod](https://github.com/iter8-tools/iter8/blob/master/go.mod) and [go.sum](https://github.com/iter8-tools/iter8/blob/master/go.sum) if needed.

    - Update the default Helm chart configuration [values.yaml](https://github.com/iter8-tools/iter8/blob/master/charts/controller/values.yaml) and bump the chart version in [Chart.yaml](https://github.com/iter8-tools/iter8/blob/master/charts/controller/Chart.yaml).

    - Please also consider including your information in our list of [adopters](https://github.com/iter8-tools/iter8/blob/master/ADOPTERS.md).

3. Submit a second pull request on the Iter8 [docs project](https://github.com/iter8-tools/docs) updating the list of available implementations.

4. Alert the project reviewers on [Slack](https://join.slack.com/t/iter8-tools/shared_invite/zt-awl2se8i-L0pZCpuHntpPejxzLicbmw) `#development` channel.
