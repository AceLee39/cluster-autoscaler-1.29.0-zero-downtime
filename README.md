# Kubernetes Autoscaler Zero downtime support  **"single replica"** 
1. Base on https://github.com/kubernetes/autoscaler/archive/refs/tags/cluster-autoscaler-1.29.0.tar.gz

How to compile : cd to cluster-autoscaler folder
make build

flag customized
Please note :
config : command in the autoscaler deployment yaml
```
command:
- ./cluster-autoscaler 
- --v=4 
- --stderrthreshold=info 
- --cloud-provider=civo 
- --nodes=1:10:workers 
- --skip-nodes-with-local-storage=false 
- --skip-nodes-with-system-pods=false 
- --rolling-restart-label-app=test
```
the pod should with label app=test which config in the above command rolling-restart-label-app=test

2. https://aws.amazon.com/cn/blogs/china/use-the-readiness-gate-to-solve-the-service-interruption-caused-by-the-rolling-upgrade-of-eks-pod

   2.1 deployment yaml config with lifecycle preStop

```
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
    lifecycle:
      preStop:
        exec:
           command: ["/bin/bash", "-c", "sleep 120"]
```
   2.2 使用 Readiness Gate



-----------------------------------------------------------------------------------------------------------------------------------------------------
[![Release Charts](https://github.com/kubernetes/autoscaler/actions/workflows/release.yaml/badge.svg)](https://github.com/kubernetes/autoscaler/actions/workflows/release.yaml) [![Tests](https://github.com/kubernetes/autoscaler/actions/workflows/ci.yaml/badge.svg)](https://github.com/kubernetes/autoscaler/actions/workflows/ci.yaml) [![GoDoc Widget]][GoDoc]

This repository contains autoscaling-related components for Kubernetes.

## What's inside

[Cluster Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler) - a component that automatically adjusts the size of a Kubernetes
Cluster so that all pods have a place to run and there are no unneeded nodes. Supports several public cloud providers. Version 1.0 (GA) was released with kubernetes 1.8.

[Vertical Pod Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler) - a set of components that automatically adjust the
amount of CPU and memory requested by pods running in the Kubernetes Cluster. Current state - beta.

[Addon Resizer](https://github.com/kubernetes/autoscaler/tree/master/addon-resizer) - a simplified version of vertical pod autoscaler that modifies
resource requests of a deployment based on the number of nodes in the Kubernetes Cluster. Current state - beta.

[Charts](https://github.com/kubernetes/autoscaler/tree/master/charts) - Supported Helm charts for components above.

## Contact Info

Interested in autoscaling? Want to talk? Have questions, concerns or great ideas?

Please join us on #sig-autoscaling at https://kubernetes.slack.com/, or join one
of our weekly meetings. See [the Kubernetes Community Repo](https://github.com/kubernetes/community/blob/master/sig-autoscaling/README.md) for more information.

## Getting the Code

Fork the repository in the cloud:

1. Visit https://github.com/kubernetes/autoscaler
1. Click Fork button (top right) to establish a cloud-based fork.

The code must be checked out as a subdirectory of `k8s.io`, and not `github.com`.

```shell
mkdir -p $GOPATH/src/k8s.io
cd $GOPATH/src/k8s.io
# Replace "$YOUR_GITHUB_USERNAME" below with your github username
git clone https://github.com/$YOUR_GITHUB_USERNAME/autoscaler.git
cd autoscaler
```

Please refer to Kubernetes [Github workflow guide] for more details.

[GoDoc]: https://godoc.org/k8s.io/autoscaler
[GoDoc Widget]: https://godoc.org/k8s.io/autoscaler?status.svg
[Github workflow guide]: https://github.com/kubernetes/community/blob/master/contributors/guide/github-workflow.md
