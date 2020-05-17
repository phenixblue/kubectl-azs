# azs Kubernetes CLI

The `azs` CLI is a utlity that can be used as a standalone binary or as a `kubectl` plugin to list Kubernetes objects by failure domain (commonly referred to as an Availability Zone). This utility requires a valid kubeconfig context and a given label on the cluster nodes to delineate the failure domain(s), and RBAC permissions to list Node resources. By default the `failure-domain.kubernetes.io/zone` and `failure-domain.beta.kubernetes.io/zone` labels are used. 

This was an opportunity for me to start working with golang and is meant to be a simple example for building a kubectl plugin.

The utility is written in [go](https://golang.org) and uses the [cobra](https://github.com/spf13/cobra), and [goreleaser](https://goreleaser.com) projects.

## Usage (Standalone)

```
The "azs" utility is a tool to list Kubernetes objects by Availability Zone. The utility can
be used standalone or as a "kubectl" plugin

Usage:
  azs [flags]
  azs [command]

Available Commands:
  help        Help about any command
  nodes       List Kubernetes nodes by AZ
  pods        List Kubernetes pods by AZ
  version     Prints the version of the "azs" utility

Flags:
      --context string      The name of the kubeconfig context to use
  -h, --help                help for azs
      --kubeconfig string   Kubernetes configuration file
  -l, --label string        The target label that defines the Availability Zone on nodes. The default value also includes
                            the beta version of the same label (default "failure-domain.kubernetes.io/zone")
  -n, --namespace string    The Namespace to use when listing Pods

Use "azs [command] --help" for more information about a command.
```

### Use as a kubectl plugin

Copy binary into your path

```bash
$ cp ./azs /usr/local/bin/kubectl-azs
```

### List Availability Zones using default label

```bash
$ kubectl azs                                                          
AZ
0
1
2
```

### List Nodes and show Availability Zones using default label

```bash
$ kubectl azs nodes                                                        
NODE NAME                           AZ
aks-agentpool-94661123-vmss000000   0
aks-agentpool-94661123-vmss000001   1
aks-agentpool-94661123-vmss000002   2
```

### List Nodes and show Availability Zones using custom label

```bash
$ kubectl azs nodes -l example.io/custom-az-label
NODE NAME                           AZ
aks-agentpool-94661123-vmss000000   az1
aks-agentpool-94661123-vmss000001   az2
aks-agentpool-94661123-vmss000002   az3
```

### List Pods and show Availability Zones using default label

```bash
$ kubectl azs pods                                                        
NAME                                        NAMESPACE     NODE                                AZ
nginx                                       default       aks-agentpool-94661123-vmss000000   0
coredns-544d979687-dzsnh                    kube-system   aks-agentpool-94661123-vmss000000   0
coredns-544d979687-w4z55                    kube-system   aks-agentpool-94661123-vmss000002   2
coredns-autoscaler-6b69c49744-r4qd8         kube-system   aks-agentpool-94661123-vmss000002   2
dashboard-metrics-scraper-867cf6588-2pcmh   kube-system   aks-agentpool-94661123-vmss000002   2
kube-proxy-5fkl6                            kube-system   aks-agentpool-94661123-vmss000000   0
kube-proxy-jbdp6                            kube-system   aks-agentpool-94661123-vmss000002   2
kube-proxy-tm7cx                            kube-system   aks-agentpool-94661123-vmss000001   1
kubernetes-dashboard-7f7676f7b5-sz6dt       kube-system   aks-agentpool-94661123-vmss000002   2
metrics-server-6df44d5569-5czfg             kube-system   aks-agentpool-94661123-vmss000002   2
omsagent-4fjxj                              kube-system   aks-agentpool-94661123-vmss000000   0
omsagent-jm87x                              kube-system   aks-agentpool-94661123-vmss000001   1
omsagent-qbhz9                              kube-system   aks-agentpool-94661123-vmss000002   2
omsagent-rs-7bf46448f4-5z2v5                kube-system   aks-agentpool-94661123-vmss000001   1
tunnelfront-6446d9d9cc-fdrnv                kube-system   aks-agentpool-94661123-vmss000002   2
```

### List Pods in a specific namespace and show Availability Zones using default label

```bash
$ kubectl azs pods -n default                                   
NAME                                        NAMESPACE     NODE                                AZ
nginx                                       default       aks-agentpool-94661123-vmss000000   0
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

