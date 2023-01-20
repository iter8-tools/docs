---
template: main.html
---

# AutoX: Auto Experimentation

AutoX, short for "auto experimentation", allows Iter8 to detect new versions of your services and automatically trigger new experiments, allowing you to tests your services as soon as you push out a new version.

![AutoX](images/autox.png)

To be more exact, AutoX detects changes to a particular set of labels from a particular Kubernetes resource and will execute experiments based on those labels. For example, when a deployment is updated and its version label is bumped, AutoX can spin up a new SLO test to see if the new version satisifies requirements.

<!-- 
Is it clear: trigger vs trigger labels? Should triggers be called trigger resource instead?
-->

In order for AutoX to function, you must specify a trigger and a set of experiment charts. The trigger specifies the Kubernetes resource object that AutoX should watch and the experiment charts specify the Iter8 experiments AutoX should launch. 

The trigger is a combination of name, namespace, and GVR (group, version, resource). When a particular set of labels of a matching resource is changed **and** the resource continues to have an `iter8.tools/autox-group` label (referred to as the AutoX label for the remainder of the article) then AutoX will launch the experiments. The labels that need to be changed are the following: `app.kubernetes.io/name"`, `app.kubernetes.io/version"`, and `iter8.tools/track"` (referred to as the trigger labels for the remainder of the article).

<!-- 
For trigger labels, do we really expect users to change any of the 3?
Name label is immutable? Do we expect users to change the track label?

If so, then it's really just the version label, right?
-->

The experiment charts come from [Iter8 Hub](https://github.com/iter8-tools/hub). Iter8 Hub primarily contains Iter8 experiments but it also contains other experiments such as Litmus-based chaos injection experiments.

<p align='center'>
  <img alt-text="AutoX flowchart" src="../images/flowchart.png" width="50%" />
</p>

***

???+ warning "Before you begin"
    1. Try [your first experiment](../../getting-started/your-first-experiment.md). Understand the main [concepts](../../getting-started/concepts.md) behind Iter8 experiments.

## Setup Kubernetes cluster with ArgoCD

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

See [here](https://argo-cd.readthedocs.io/en/stable/getting_started/#1-install-argo-cd) for more information.

## Setup Kubernetes cluster with Iter8 AutoX

```bash
helm repo add iter8 https://iter8-tools.github.io/hub/
helm install autox-httpbin iter8/autox \
--set 'releaseGroupSpecs.httpbin.trigger.name=httpbin' \
--set 'releaseGroupSpecs.httpbin.trigger.namespace=default' \
--set 'releaseGroupSpecs.httpbin.trigger.group=apps' \
--set 'releaseGroupSpecs.httpbin.trigger.version=v1' \
--set 'releaseGroupSpecs.httpbin.trigger.resource=deployments' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.name=iter8' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.values.tasks={ready,http,assess}' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.values.ready.deploy=httpbin' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.values.ready.service=httpbin' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.values.ready.timeout=60s' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.values.http.url=http://httpbin.default/get' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.values.assess.SLOs.upper.http/latency-mean=50' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.values.assess.SLOs.upper.http/error-count=0' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.version=0.12.2' \
--set 'releaseGroupSpecs.httpbin.releaseSpecs.iter8.values.runner=job'
```

As mentioned previously, the input to AutoX is a trigger and a set of experiment charts. Here, the trigger is a Kubernetes resource object with the name `httpbin`, namespace `default`, and GVR `apps`, `deployments`, and `v1`, meaning that AutoX will watch any resource that meets the description. When a changed is made to the trigger labels of a matching resource and the resource continues to have the AutoX label, then AutoX will install the Helm charts.

In this case, there is only one experiment chart to install. The specified experiment chart is pointing to an Iter8 experiment, specifically an HTTP SLO validation test on the `httpbin` service. 

This experiment is composed of three tasks, `ready`, `http`, and `assess`. The `ready` task will ensure that the `httpbin` deployment and service are running. The `http` task will make requests to the specified URL and will collect latency and error-related metrics. Lastly, the `assess` task will ensure that the mean latency is less than 50 milliseconds and the error count is 0. In addition, the runner is set to job as this will be a [single-loop experiment](https://iter8.tools/0.11/getting-started/concepts/#iter8-experiment).

## Create application

Now, we will create the `httpbin` deployment and service.

```bash
kubectl create deployment httpbin --image=kennethreitz/httpbin --port=80
kubectl expose deployment httpbin --port=80
```

## Apply labels

In the previous step, we created an `apps/v1` deployment with the name `httpbin` in the `default` namespace, which matches the trigger that we configured for AutoX. However, to enable AutoX for the Kubernetes resource object, we need to assign it the AutoX label.

```bash
kubectl label deployment httpbin iter8.tools/autox=true
```

## Observe automatic experiment

After you have assigned the AutoX label, an experiment should start.

You can now use `iter8` commands in order to check the status and the results of the experiment. Note that you need to specify an experiment group via the `-g` option. The experiment group for AutoX experiments is in the form `autox-<release group spec name>-<release spec name>`. In this case, it would be `autox-httpbin-iter8`. 

The following command allows you to check the status of the experiment. If the experiment does not immediately start, try waiting a minute.

```bash
iter8 k assert -c nofailure -c slos -g autox-httpbin-iter8
```

??? note "Sample output from assert"
    ```
    INFO[2023-01-11 14:43:45] inited Helm config                           
    INFO[2023-01-11 14:43:45] experiment has no failure                    
    INFO[2023-01-11 14:43:45] SLOs are satisfied                           
    INFO[2023-01-11 14:43:45] all conditions were satisfied  
    ```

    We can see in the sample output that the experiment has completed and all SLOs and conditions were satisfied.

And the following command allows you to check the results of the experiment.

```bash
iter8 k report -g autox-httpbin-iter8
```

??? note "Sample output from assert"
    ```
    Experiment summary:
    *******************

    Experiment completed: true
    No task failures: true
    Total number of tasks: 4
    Number of completed tasks: 4

    Whether or not service level objectives (SLOs) are satisfied:
    *************************************************************

    SLO Conditions                 | Satisfied
    --------------                 | ---------
    http/error-count <= 0          | true
    http/latency-mean (msec) <= 50 | true
    

    Latest observed values for metrics:
    ***********************************

    Metric                     | value
    -------                    | -----
    http/error-count           | 0.00
    http/error-rate            | 0.00
    http/latency-max (msec)    | 25.11
    http/latency-mean (msec)   | 5.59
    http/latency-min (msec)    | 1.29
    http/latency-p50 (msec)    | 4.39
    http/latency-p75 (msec)    | 6.71
    http/latency-p90 (msec)    | 10.40
    http/latency-p95 (msec)    | 13.00
    http/latency-p99 (msec)    | 25.00
    http/latency-p99.9 (msec)  | 25.10
    http/latency-stddev (msec) | 4.37
    http/request-count         | 100.00
    ```

    In the sample output, we can see an experiment summary, a list of SLOs and whether they were satisfied or not, as well as any additional metrics that were collected as part of the experiment.

## Push new version of the application

Now that AutoX is watching the `httpbin` deployment, any change that we make to its trigger labels will cause AutoX to trigger a new experiment.

We will trigger a new experiment by adding a new trigger label to the `httpbin` deployment.

```bash
kubectl label deployment httpbin app.kubernetes.io/version=1.0.0
```

## Observe new automatic experiment

Check to see if a new experiment should have started. Refer to [Observe automatic experiment](#observe-automatic-experiment) for the necessary commands.

## Next steps

You can continue to modify the trigger labels of the `httpbin` deployment to trigger new experiments. For example, you can continue to bump the `app.kubernetes.io/version` label to trigger the HTTP SLO validation test. Now it is easy to know if your app meets basic HTTP performance requirements.

It is possible to use AutoX to conduct more complex experiments. Iter8 experiments are composed from discrete tasks so you can design an experiment that best fits your use case. For example, instead of using the `httpbin` task, you can use `abn` task in order to run an A/B/n experiment. You can also run experiments that are not from Iter8. For example, there is also a experiment chart for a Litmus Chaos chaos experiment. Lastly, you can supply multiple experiments charts so AutoX will deploy a suite of experiments to run whenever you update your trigger resource.

Note that if you delete the AutoX label or if you delete the `httpbin` development, then AutoX will also delete the respective experiments.

## Clean up

<!-- What about deleting Argo CD? -->

```bash
helm delete autox-httpbin
kubectl delete deployment/httpbin
kubectl delete service/httpbin
```