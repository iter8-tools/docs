Assert that the experiment encountered no failures, and all SLOs are satisfied. 

```shell
iter8 k assert -c nofailure -c slos
```

Note that because these are looping experiments without a defined end, do not use the `completed` assert condition.

The `iter8 assert` command asserts if the experiment result satisfies conditions that are specified. If assert conditions are satisfied, it exits with code `0`; else, it exits with code `1`. Assertions are especially useful inside CI/CD/GitOps pipelines.

??? note "Sample output from assert"
    ```shell
    INFO[2021-11-10 09:33:12] experiment completed
    INFO[2021-11-10 09:33:12] experiment has no failure                    
    INFO[2021-11-10 09:33:12] SLOs are satisfied                           
    INFO[2021-11-10 09:33:12] all conditions were satisfied
    ```