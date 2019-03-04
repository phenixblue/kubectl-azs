# azs Kubernetes CLI


## About

The `azs` CLI is a utlity that can be used as a standalone binary or as a `kubectl` plugin to list Kubernetes objects by availability zone. This utility requires
`kubectl` to be installed, a valid kubeconfig referenced via the `KUBECONFIG` environment variable, and the `failure-domain.beta.kubernetes.io/zone`
label to be set for the cluster nodes.

The utility is written in `go` and uses the `cobra` and `goreleaser` projects.


## Usage (Standalone)

- List nodes/pods by Avalability Zone
    ```
    $ azs nodes
    $ azs pods
    ```

## Usage (kubectl Plugin)
- Copy binary into your path
    ```
    $ cp ./azs /usr/local/bin/kubectl-azs
    ```
- List nodes/pods by Availability Zone
    ```
    $ kubectl azs azs
    $ kubectl azs pods
    ```

## Developers

### Build a new release

- Setup GitHub Personal Access Token

    ```
    $ export GITHUB_TOKEN="XXXXXXXXXXXXXXXXXXXXXXXXXX"
    ```

- Edit the version in `main.go`

- Tag new version and push

    ```
    $ git tag -a v0.0.1 -m "First release"
    $ git push origin v0.0.1
    ```

- Build new binaries

    ```
    $ goreleaser --rm-dist
    ```

