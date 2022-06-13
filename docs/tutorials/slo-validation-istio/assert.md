Assert that the experiment encountered no failures, and all SLOs are satisfied. 

```shell
iter8 k assert -c nofailure -c slos
```

??? note "Sample output from assert"
    ```shell
    INFO[2021-11-10 09:33:12] experiment has no failure                    
    INFO[2021-11-10 09:33:12] SLOs are satisfied                           
    INFO[2021-11-10 09:33:12] all conditions were satisfied
    ```