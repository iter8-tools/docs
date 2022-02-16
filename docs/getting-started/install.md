---
template: main.html
title: Install Iter8
hide:
- toc
---

# Install Iter8 CLI

=== "Brew"
    Install the latest stable release of the Iter8 CLI using `brew` as follows.

    ```shell
    brew tap iter8-tools/iter8
    brew install iter8
    ```
    
    ???+ note "Install a specific version"
        You can install the Iter8 CLI with specific major and minor version numbers. For example, the following command installs the latest stable release of the Iter8 CLI with major `0` and minor `9`.
        ```shell
        brew tap iter8-tools/iter8
        brew install iter8@0.9
        ```

=== "Binaries"
    Replace `${TAG}` below with the [latest or any desired Iter8 release tag](https://github.com/iter8-tools/iter8/releases). For example,
    ```shell
    export TAG=v0.9.0
    ```

    === "darwin-amd64 (MacOS)"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-darwin-amd64.tar.gz | tar xvz -
        ```
        Move `darwin-amd64/iter8` to any directory in your `PATH`.

    === "linux-amd64"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-linux-amd64.tar.gz | tar xvz -
        ```
        Move `linux-amd64/iter8` to any directory in your `PATH`.

    === "linux-386"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-linux-386.tar.gz | tar xvz -
        ```
        Move `linux-386/iter8` to any directory in your `PATH`.

    === "windows-amd64"
        ```shell
        wget -qO- https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-windows-amd64.tar.gz | tar xvz -
        ```
        Move `windows-amd64/iter8.exe` to any directory in your `PATH`.


=== "Source"
    Build the Iter8 CLI from source as follows. Go `1.17+` is a pre-requisite.
    ```shell
    # you can replace master with a specific tag such as v0.9.0
    export TAG=master
    https://github.com/iter8-tools/iter8.git?ref=${TAG}
    cd iter8
    make install
    ```

=== "Go 1.17+"
    Install the latest stable release of the Iter8 CLI using `go 1.17+` as follows.

    ```shell
    go install github.com/iter8-tools/iter8@latest
    ```
    You can now run `iter8` (from your gopath bin/ directory)

    ???+ note "Install a specific version"
        You can also install the Iter8 CLI with a specific tag. For example, the following command installs version `0.9.0` of the Iter8 CLI.
        ```shell
        go install github.com/iter8-tools/iter8@v0.9.0
        ```


