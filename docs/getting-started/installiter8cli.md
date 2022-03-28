=== "Brew"
    Install the latest stable release of the Iter8 CLI using `brew` as follows.

    ```shell
    brew tap iter8-tools/iter8
    brew install iter8@0.9
    ```
    
=== "Binaries"
    Replace `${TAG}` below with [any desired Iter8 release tag](https://github.com/iter8-tools/iter8/releases).

    === "darwin-amd64 (MacOS)"
        ```shell
        wget https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-darwin-amd64.tar.gz
        tar -xvf iter8-darwin-amd64.tar.gz
        ```
        Move `darwin-amd64/iter8` to any directory in your `PATH`.

    === "linux-amd64"
        ```shell
        wget https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-linux-amd64.tar.gz
        tar -xvf iter8-linux-amd64.tar.gz
        ```
        Move `linux-amd64/iter8` to any directory in your `PATH`.

    === "linux-386"
        ```shell
        wget https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-linux-386.tar.gz
        tar -xvf iter8-linux-386.tar.gz
        ```
        Move `linux-386/iter8` to any directory in your `PATH`.

    === "windows-amd64"
        ```shell
        wget https://github.com/iter8-tools/iter8/releases/download/${TAG}/iter8-windows-amd64.tar.gz
        tar -xvf iter8-windows-amd64.tar.gz
        ```
        Move `windows-amd64/iter8.exe` to any directory in your `PATH`.
