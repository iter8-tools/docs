=== "Brew"
    Install the latest stable release of the Iter8 CLI using `brew` as follows.

    ```shell
    brew tap iter8-tools/iter8
    brew install iter8
    ```
    
    ??? note "Install a specific version"
        You can install the Iter8 CLI with specific major and minor version numbers. For example, the following command installs the latest stable release of the Iter8 CLI with major `0` and minor `9`.
        ```shell
        brew tap iter8-tools/iter8
        brew install iter8@0.9
        ```

=== "Binaries"
    You can replace `v0.9.0` with [any desired Iter8 release tag](https://github.com/iter8-tools/iter8/releases).

    === "darwin-amd64 (MacOS)"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/v0.9.0/iter8-darwin-amd64.tar.gz | tar xvz -
        ```
        Move `darwin-amd64/iter8` to any directory in your `PATH`.

    === "linux-amd64"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/v0.9.0/iter8-linux-amd64.tar.gz | tar xvz -
        ```
        Move `linux-amd64/iter8` to any directory in your `PATH`.

    === "linux-386"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/v0.9.0/iter8-linux-386.tar.gz | tar xvz -
        ```
        Move `linux-386/iter8` to any directory in your `PATH`.

    === "windows-amd64"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/v0.9.0/iter8-windows-amd64.tar.gz | tar xvz -
        ```
        Move `windows-amd64/iter8.exe` to any directory in your `PATH`.


=== "Source"
    Go `1.17+` is a pre-requisite.  Replace `master` with [any desired Iter8 release tag](https://github.com/iter8-tools/iter8/releases).

    ```shell
    export TAG=master
    https://github.com/iter8-tools/iter8.git?ref=${TAG}
    cd iter8
    make install
    ```

=== "Go 1.17+"
    You can replace `latest` with [any desired Iter8 release tag](https://github.com/iter8-tools/iter8/releases).
    ```shell
    go install github.com/iter8-tools/iter8@latest
    ```
    You can now run `iter8` (from your gopath bin/ directory)
