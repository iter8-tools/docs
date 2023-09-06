??? note "Sample performance test logs"
    ```shell
      time=2023-09-01 19:31:40 level=info msg=task 2: ready: started
      time=2023-09-01 19:31:40 level=info msg=task 2: ready: completed
      time=2023-09-01 19:31:40 level=info msg=task 3: http: started
      {"ts":1693596700.013507,"level":"info","file":"httprunner.go","line":100,"msg":"Starting http test","run":"0","url":"http://httpbin.default/get","threads":"4","qps":"8.0","warmup":"parallel","conn-reuse":""}
      {"ts":1693596712.606946,"level":"info","file":"periodic.go","line":832,"msg":"T001 ended after 12.534760214s : 25 calls. qps=1.9944537887591696"}
      {"ts":1693596712.616122,"level":"info","file":"periodic.go","line":832,"msg":"T002 ended after 12.544591006s : 25 calls. qps=1.9928907995519867"}
      {"ts":1693596712.623089,"level":"info","file":"periodic.go","line":832,"msg":"T003 ended after 12.551572714s : 25 calls. qps=1.9917822706086104"}
      {"ts":1693596712.628555,"level":"info","file":"periodic.go","line":832,"msg":"T000 ended after 12.557040548s : 25 calls. qps=1.9909149695293316"}
      {"ts":1693596712.629567,"level":"info","file":"periodic.go","line":564,"msg":"Run ended","run":"0","elapsed":"12.557657464s","calls":"100","qps":"7.963268649959411"}
      time=2023-09-01 19:31:52 level=info msg=task 3: http: completed
    ```
