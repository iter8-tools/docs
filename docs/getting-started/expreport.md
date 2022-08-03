=== "Text"
    ```shell
    iter8 k report
    ```

    ??? note "The text report looks like this"
        ```shell
        Experiment summary:
        *******************

          Experiment completed: true
          No task failures: true
          Total number of tasks: 4
          Number of completed tasks: 4

        Whether or not service level objectives (SLOs) are satisfied:
        *************************************************************

          SLO Conditions                 |Satisfied
          --------------                 |---------
          http/error-count <= 0          |true
          http/latency-mean (msec) <= 50 |true
          

        Latest observed values for metrics:
        ***********************************

          Metric                     |value
          -------                    |-----
          http/error-count           |0.00
          http/error-rate            |0.00
          http/latency-max (msec)    |21.18
          http/latency-mean (msec)   |6.17
          http/latency-min (msec)    |1.91
          http/latency-p50 (msec)    |5.78
          http/latency-p75 (msec)    |7.89
          http/latency-p90 (msec)    |9.50
          http/latency-p95 (msec)    |11.33
          http/latency-p99 (msec)    |20.59
          http/latency-p99.9 (msec)  |21.12
          http/latency-stddev (msec) |3.27
          http/request-count         |100.00
        ```

=== "HTML"
    ```shell
    iter8 k report -o html > report.html # view in a browser
    ```

    ??? note "The HTML report looks like this"
        ![HTML report](images/report.html.png)