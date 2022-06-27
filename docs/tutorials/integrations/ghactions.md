---
template: main.html
---

# Use Iter8 in a GitHub Actions workflow

Install the latest version of Iter8 CLI using `iter8-tools/iter8@v0.11`. Once installed, Iter8 can be used as documented in the various tutorials. For example:

```yaml
# install Iter8 CLI
- uses: iter8-tools/iter8@v0.11
# launch a local experiment
- run: |
    iter8 launch --set "tasks={http}" --set http.url=http://httpbin.org/get
# launch an experiment inside Kubernetes;
# this assumes that your Kubernetes cluster is accessible 
# from the GitHub Actions pipeline
- run: |
    iter8 k launch --set "tasks={http}" \
    --set http.url=http://httpbin.org/get \
    --set runner=job
```
