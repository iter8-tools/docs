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

          SLO Conditions                   |Satisfied
          --------------                   |---------
          grpc/error-rate <= 0             |true
          grpc/latency/mean (msec) <= 50   |true
          grpc/latency/p90 (msec) <= 100   |true
          grpc/latency/p97.5 (msec) <= 200 |true
          

        Latest observed values for metrics:
        ***********************************

          Metric                    |value
          -------                   |-----
          grpc/error-count          |0.00
          grpc/error-rate           |0.00
          grpc/latency/mean (msec)  |21.48
          grpc/latency/p90 (msec)   |34.00
          grpc/latency/p97.5 (msec) |37.00
          grpc/request-count        |200.00
        ```

=== "HTML"
    ```shell
    iter8 report -o html > report.html # view in a browser
    ```

    ??? note "The HTML report looks like this"
        ![HTML report](images/report.html.png)