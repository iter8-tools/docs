---
template: main.html
---

# Automated Experiments: AutoX

AutoX, short for "automated experiments", allows Iter8 to detect changes to your Kubernetes resources objects and automatically start new experiments, allowing you to test your applications as soon as you release a new version.

To configure AutoX, you will need to specify a set of experiment groups and, for each group, the Kubernetes resource object (trigger object) that you expect AutoX to watch and one or more experiments to be performed in response to new versions of this object.

The trigger object is specified by providing the name, namespace, and the group-version-resource (GVR) metadata of the trigger object.

See the following example:

```bash
helm install autox autox --repo https://iter8-tools.github.io/hub/ --version 0.1.6 \
--set 'groups.myApp.trigger.name=myApp' \
--set 'groups.myApp.trigger.namespace=default' \
--set 'groups.myApp.trigger.group=apps' \
--set 'groups.myApp.trigger.version=v1' \
--set 'groups.myApp.trigger.resource=deployments' \
--set 'groups.myApp.specs.iter8-http.name=iter8' \
--set 'groups.myApp.specs.iter8-http.values.tasks={ready,http,assess}' \
--set 'groups.myApp.specs.iter8-http.values.ready.deploy=myApp' \
--set 'groups.myApp.specs.iter8-http.values.ready.service=myApp' \
--set 'groups.myApp.specs.iter8-http.values.ready.timeout=60s' \
--set 'groups.myApp.specs.iter8-http.values.http.url=http://myApp.default/get' --set 'groups.myApp.specs.iter8-http.values.assess.SLOs.upper.http/error-count=0' --set 'groups.myApp.specs.iter8-http.values.assess.SLOs.upper.http/latency-mean=50' \
--set 'groups.myApp.specs.iter8-http.version=0.13.0' \
--set 'groups.myApp.specs.iter8-http.values.runner=job'
```

In this example, there is only one experiment group named `myApp` (`groups.myApp...`), and within that group, there is the trigger object definition (`groups.myApp.trigger...`) and a single experiment spec named `iter8-http` (`groups.myApp.specs.iter8...`).

***

In this next example, we have augmented the previous example with an additional experiment spec.

```bash
  helm install autox autox --repo https://iter8-tools.github.io/hub/ --version 0.1.6 \
  --set 'groups.myApp.trigger.name=myApp' \
  --set 'groups.myApp.trigger.namespace=default' \
  --set 'groups.myApp.trigger.group=apps' \
  --set 'groups.myApp.trigger.version=v1' \
  --set 'groups.myApp.trigger.resource=deployments' \
  --set 'groups.myApp.specs.iter8-http.name=iter8' \
  --set 'groups.myApp.specs.iter8-http.values.tasks={ready,http,assess}' \
  --set 'groups.myApp.specs.iter8-http.values.ready.deploy=myApp' \
  --set 'groups.myApp.specs.iter8-http.values.ready.service=myApp' \
  --set 'groups.myApp.specs.iter8-http.values.ready.timeout=60s' \
  --set 'groups.myApp.specs.iter8-http.values.http.url=http://myApp.default/get' \
  --set 'groups.myApp.specs.iter8-http.values.assess.SLOs.upper.http/error-count=0' \
  --set 'groups.myApp.specs.iter8-http.values.assess.SLOs.upper.http/latency-mean=50' \
  --set 'groups.myApp.specs.iter8-http.version=0.13.0' \
  --set 'groups.myApp.specs.iter8-http.values.runner=job' \
  --set 'groups.myApp.specs.iter8-grpc.values.tasks={ready,grpc,assess}' \
  --set 'groups.myApp.specs.iter8-grpc.values.ready.deploy=myApp' \
  --set 'groups.myApp.specs.iter8-grpc.values.ready.service=myApp' \
  --set 'groups.myApp.specs.iter8-grpc.values.ready.timeout=60s' \
  --set 'groups.myApp.specs.iter8-grpc.values.grpc.host=...' \
  --set 'groups.myApp.specs.iter8-grpc.values.grpc.call=...' \
  --set 'groups.myApp.specs.iter8-grpc.values.grpc.protoURL=...' \
  --set 'groups.myApp.specs.iter8-grpc.values.assess.SLOs.upper.grpc/error-rate=0' \
  --set 'groups.myApp.specs.iter8-grpc.values.assess.SLOs.upper.grpc/latency/latency-mean=50' \
  --set 'groups.myApp.specs.iter8-grpc.values.runner=job'
```

Now, when a new version of the trigger is released, AutoX will relaunch not only an HTTP SLO validation test but also a GRPC SLO validation test.

***

A trigger object must have a `app.kubernetes.io/version` label (version label. This label is used to identify new versions of the trigger object, which will cause AutoX to relaunch experiments.

If the trigger does not have a version label, the AutoX will attempt to remove any preexisting experiments.

<p align='center'>
  <img alt-text="AutoX flowchart" src="../images/flowchart.png" width="60%" />
</p>
