---
template: main.html
---

# Use Iter8 in a GitHub Actions workflow

Install the latest version of Iter8 using the GitHub Action `iter8/iter8@v0.10`. Once installed, Iter8 can be used as documented (see [user guide](../../user-guide/commands/iter8.md)) in `run` actions. For example:

```yaml
- uses: iter8/iter8@v0.10 # install Iter8
- run: |
    iter8 version
    iter8 launch -c load-test-http \
    --set url=http://httpbin.org/get
```
