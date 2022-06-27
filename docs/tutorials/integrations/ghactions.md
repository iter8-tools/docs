---
template: main.html
---

# Use Iter8 in a GitHub Actions workflow

Install the latest version of the Iter8 CLI using `iter8-tools/iter8@v0.11`. Once installed, the Iter8 CLI can be used as documented in various tutorials. For example:

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
