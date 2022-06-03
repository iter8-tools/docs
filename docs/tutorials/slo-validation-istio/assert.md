Assert that the experiment completed without failures, and all SLOs are satisfied.
```shell
iter8 assert -c completed -c nofailure -c slos
```

The `iter8 assert` command asserts if the experiment result satisfies conditions that are specified. If assert conditions are satisfied, it exits with code `0`; else, it exits with code `1`. Assertions are especially useful inside CI/CD/GitOps pipelines.

??? note "Sample output from assert"
    ```shell
    INFO[2021-11-10 09:33:12] experiment completed
    INFO[2021-11-10 09:33:12] experiment has no failure                    
    INFO[2021-11-10 09:33:12] SLOs are satisfied                           
    INFO[2021-11-10 09:33:12] all conditions were satisfied
    ```