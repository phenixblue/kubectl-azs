# azs Kubernetes CLI


## About

The `azs` CLI is a utlity that can be used as a standalone binary or as a `kubectl` plugin to list Kubernetes objects by availability zone. This utility requires
`kubectl` to be installed, a valid kubeconfig referenced via the `KUBECONFIG` environment variable, and the `failure-domain.beta.kubernetes.io/zone`
label to be set for the cluster nodes.

This was an opportunity for me to start working with golang and is meant to be a simple example for building a kubectl plugin.

The utility is written in [go](https://golang.org) and uses the [cobra](https://github.com/spf13/cobra), [dep](https://github.com/golang/dep), and [goreleaser](https://goreleaser.com) projects.



## Usage (Standalone)

- List nodes/pods by Avalability Zone

    ```bash
    $ azs nodes
    $ azs pods
    ```

## Usage (kubectl Plugin)

- Copy binary into your path

    ```bash
    $ cp ./azs /usr/local/bin/kubectl-azs
    ```


- List nodes/pods by Availability Zone

    ```bash
    $ kubectl azs nodes
    $ kubectl azs pods
    ```

## Developers

### Build a new release

- Setup GitHub Personal Access Token

    ```bash
    $ export GITHUB_TOKEN="XXXXXXXXXXXXXXXXXXXXXXXXXX"
    ```

- Edit the version in `main.go`

- Tag new version and push

    ```bash
    $ git tag -a v0.0.1 -m "First release"
    $ git push origin v0.0.1
    ```

- Build new binaries

    ```bash
    $ goreleaser --rm-dist
    ```

