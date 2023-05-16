=== "GitHub Actions"
    Install the latest stable release of the Iter8 CLI in your GitHub Actions workflow as follows.

    ```yaml
    - name: Install Iter8
      run: GOBIN=/usr/local/bin go install github.com/iter8-tools/iter8@v0.14
    ```
