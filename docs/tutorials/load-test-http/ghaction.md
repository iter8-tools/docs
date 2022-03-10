---
template: main.html
---

# Benchmark and Validate using Iter8 Github Action

Iter8 experiments can be run as part of a Github workflow using the [Iter8 action](https://github.com/iter8-tools/iter8-action).

### Usage overview

The Iter8 action runs an Iter8 experiment by reference to the experiment chart and, optionally, a chart repository. The experiment is configured by the definition of a `values.yaml` file.

To use the action to benchmark and validate an HTTP service, use the [load-test-http chart](https://github.com/iter8-tools/hub/tree/main/charts/load-test-http). For example:

```yaml
# Configure experiment using (Helm) values.yaml file
- run: |
    cat << EOF > myvalues.yaml
        url: http://localhost:8080
    EOF
# Run Iter8 experiment
- uses: iter8-tools/iter8-action@v1
  with:
    chart: load-test-http
    valuesFile: myvalues.yaml
```

If the experiment has an error or any service-level objectives (SLOs) cannot be validated, the action will fail. The exception is when the option `validateSLOs` is `false`. In this case, only an execution error will result in failure.

#### Log output

For each execution of the Iter8 action, the output of the action includes the following for reference:

- Version of iter8 (output of `iter8 version`)
- The experiment run in yaml
- An experiment report (the output of `iter8 report`)
- Assessmet of success (output of `iter8 assert -c completed -c nofailures -c slos`)

### Running sample workflow

1. Fork the Iter8 docs repo: <https://github.com/iter8-tools/docs> to your organization, *myorg*.

2. If necessary, enable workflows on the your forked repository (by navigating to <https://github.com/myorg/docs/actions>).

3. Run the workflow:

    - Navigate to the **Actions** tab of your forked repository: <https://github.com/myorg/docs/actions>.
    - Select the workflow **end-to-end tests**.
    - Click the **Run workflow** button.

4. When the workflow has completed, there will be a new entry for the execution. Select the latest, then the entry for the **local httpbin tests** job. The log for each execution of the Iter8 action, can be inspected by inspecting the steps labeled *Run iter8-tools/iter8-action@v1*.

### Reference

- A full list of options for the Iter8 action are [here](https://github.com/iter8-tools/iter8-action/tree/v1).
- A full list of the options for the `load-test-http` chart are [here](https://github.com/iter8-tools/hub/tree/main/charts/load-test-http).
