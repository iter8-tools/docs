=== "Text"
    ```shell
    iter8 report
    ```

    ??? note "The text report looks like this"
        ```shell

        Experiment summary:
        *******************

          Experiment completed: true
          No task failures: true
          Total number of tasks: 2
          Number of completed tasks: 2

        Whether or not service level objectives (SLOs) are satisfied:
        *************************************************************

          SLO Conditions        |Satisfied
          --------------        |---------
          istio/error-rate <= 0 |true
          

        Latest observed values for metrics:
        ***********************************

          Metric                           |value
          -------                          |-----
          istio/error-count                |0.00
          istio/error-rate                 |0.00
          istio/le500ms-latency-percentile |2.00
          istio/mean-latency               |120.05
          istio/request-count              |2.00
        ```

=== "HTML"
    ```shell
    iter8 report -o html > report.html # view in a browser
    ```

    ??? note "The HTML report looks like this"
        ![HTML report](images/report.html.png)