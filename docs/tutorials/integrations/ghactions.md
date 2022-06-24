---
template: main.html
---

# Use Iter8 in a GitHub Actions workflow

Install the latest version of Iter8 CLI using the GitHub Action `iter8-tools/iter8@v0.11`. Once installed, Iter8 CLI can be used as documented in the [tutorials](../../getting-started/your-first-experiment.md) and [user guide](../../user-guide/commands/iter8.md). For example:

```yaml
- uses: iter8-tools/iter8@v0.11 # install Iter8
- run: |
    iter8 version
    iter8 launch -c load-test-http \
    --set url=http://httpbin.org/get
```
