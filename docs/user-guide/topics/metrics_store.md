---
template: main.html
---

# Metrics store

One of Iter8's key advantages is that it incorporates its own metrics store simplifying the set up and execution of A/B/n and performance tests. This metrics store, based on BadgerDB is not suitable for production use. It is suitable only for a single instance of Iter8.

To support production use cases, the metrics store can be replaced with another distribued database such as Redis. Currently Iter8 supports:

- BadgerDB
- Redis

See [below](#contribute-a-new-metrics-store-implementation) for details on how to contribute additional implementations.

## Using Redis as the metrics store 

We assume that Redis is deployed. For example, for a basic deployment:

```shell
kubectl create deploy redis --image=redis/redis-stack:latest --port=6379
kubectl expose deploy redis --port=6379
```

Run Iter8 identifying the metrics store implementaion as `redis` and identify the endpoint:

```shell
helm upgrade --install --repo https://iter8-tools.github.io/iter8 --version 0.18 iter8 controller  \
--set clusterScoped=true \
--set metrics.implementation=redis \
--set metrics.redis.addresss=redis:6379
```

## Contribute a new metrics store implementation

To contribute a new metrics store implmentation:

1. Create an [issue](https://github.com/iter8-tools/iter8/issues) for discussion.

2. Submit a pull request on the [Iter8 project](https://github.com/iter8-tools/iter8) with the following updates:

    - Create sub-folder in [storage](https://github.com/iter8-tools/iter8/tree/master/storage) and provide an implementation of this [interface](https://github.com/iter8-tools/iter8/blob/master/storage/interface.go) including test cases.

    - Add a new case to [metrics.GetClient()](https://github.com/iter8-tools/iter8/blob/master/metrics/client.go).

    - Update [go.mod](https://github.com/iter8-tools/iter8/blob/master/go.mod) and [go.sum](https://github.com/iter8-tools/iter8/blob/master/go.sum) if needed.

    - Update the default Helm chart configuration [values.yaml](https://github.com/iter8-tools/iter8/blob/master/charts/controller/values.yaml) and bump the chart version in [Chart.yaml](https://github.com/iter8-tools/iter8/blob/master/charts/controller/Chart.yaml).

    - Please also consider including your information in our list of [adopters](https://github.com/iter8-tools/iter8/blob/master/ADOPTERS.md).

3. Submit a second pull request on the Iter8 [docs project](https://github.com/iter8-tools/docs) updating the list of available implementations.

4. Alert the project reviewers on [slack](https://join.slack.com/t/iter8-tools/shared_invite/zt-awl2se8i-L0pZCpuHntpPejxzLicbmw) `#development` channel.
