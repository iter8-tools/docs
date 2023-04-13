---
template: main.html
---

# Canary Model Rollouts

This tutorial shows how Iter8 can be used to implement a canary rollout of ML models. In a cnaary rollout, inference requests whose metadata match certain requirements are sent to the candidate version. The remaining requests go to the primary, or initial, version. Iter8 enables a canary rollout by automatically configuring the network to distribute inference requests.

After a one time initialization step, the end user merely deploys candidate models, evaluates them and either promotes the candiate or deletes it. Optionally, the end user can modify the percentage of inference requests being sent to the candidate. Iter8 automatically handles the underlying network configuration.

In this tutorial, we use the Istio service mesh to distribute inference requests between different versions of a model.

