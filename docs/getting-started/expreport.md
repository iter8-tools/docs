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
          Total number of tasks: 1
          Number of completed tasks: 1

        Latest observed values for metrics:
        ***********************************

          Metric                              |value
          -------                             |-----
          built-in/http-error-count           |0.00
          built-in/http-error-rate            |0.00
          built-in/http-latency-max (msec)    |203.78
          built-in/http-latency-mean (msec)   |17.00
          built-in/http-latency-min (msec)    |4.20
          built-in/http-latency-p50 (msec)    |10.67
          built-in/http-latency-p75 (msec)    |12.33
          built-in/http-latency-p90 (msec)    |14.00
          built-in/http-latency-p95 (msec)    |15.67
          built-in/http-latency-p99 (msec)    |202.84
          built-in/http-latency-p99.9 (msec)  |203.69
          built-in/http-latency-stddev (msec) |37.94
          built-in/http-request-count         |100.00
        ```

=== "HTML"
    ```shell
    iter8 k report -o html > report.html # view in a browser
    ```

    ??? note "The HTML report looks like this"
        ![HTML report](images/report.html.png)