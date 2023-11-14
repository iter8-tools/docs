```shell
helm upgrade --install --repo https://iter8-tools.github.io/iter8 --version 0.18 iter8 controller \
--set clusterScoped=true
```

For additional install options, see [Iter8 Installation](https://iter8.tools/0.18/user-guide/topics/install/).