---
layout: default
title: Supported Platforms
permalink: installation/platforms
type: installation
display-title: "false"
language: en
list: exclude
---

# Supported Platforms<a name="compatibility-matrix"></a>
Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster. See the complete list of supported platforms in the table below. With service meshes having sprung to life in the context of Kubernetes, so too, can Meshery’s deployment models be characterized in the context of Kubernetes. A given deployment of Meshery can be described as either an _in-cluster_ or an _out-of-cluster_ deployment. Meshery deploys as a stand-alone, management plane on a Docker host (_out-of-cluster_) or as a management plane in a Kubernetes cluster (_in-cluster_). 


## Platform Compatibility Matrix
Find installation instructions for the Supported Platforms in the compatibility table.

| Platform                                                                                                                                              |        Version        |
| ----------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------------: |
| <img src="/assets/img/platforms/docker.svg" width="20" height="20" /> [Docker](/docs/installation/platforms/docker)                              |                       |
| &nbsp;&nbsp;&nbsp; <img src="/assets/img/platforms/docker.svg" width="20" height="20" /> [Docker Engine](/docs/installation/platforms/docker)    |    19.x and above     |
| <img src="/assets/img/platforms/kubernetes.svg" width="20" height="20" /> [Kubernetes](/docs/installation/platforms/kubernetes)                  |   1.12.x and above    |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/aks.svg" width="20" height="20" /> [AKS](/docs/installation/platforms/aks)                     |                       |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/docker.svg" width="20" height="20" /> [Docker Desktop](/docs/installation/platforms/docker)    |    2.0.x and above    |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/eks.png" width="20" height="20" /> [EKS](/docs/installation/platforms/eks)                     |   1.12.x and above    |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/gke.png" width="20" height="20" /> [GKE](/docs/installation/platforms/gke)                     |   1.14.x and above    |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/helm.svg" width="20" height="20" /> [Helm](/docs/installation/platforms/kubernetes#using-helm) |                       |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/kind.png" width="20" height="20" /> [KinD](/docs/installation/platforms/kind)                  |        v0.7.0         |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/minikube.png" width="20" height="20" /> [Minikube](/docs/installation/platforms/minikube)      |    1.2.x and above    |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/openshift.svg" width="20" height="20" /> OpenShift                                             |      In Progress      |
| <img src="/assets/img/platforms/apple.svg" width="20" height="20" vertical-align="middle" /> [Mac](/docs/installation#mac-or-linux)              |                       |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/homebrew.png" width="20" height="20" /> [Mac - Homebrew](/docs/installation#mac-or-linux)      | macOS 10.12 - 10.15, 11 |
| &nbsp;&nbsp;&nbsp; [Scoop](/docs/installation#windows)                                                                                                |                       |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/wsl2.png" width="20" height="20" /> [Windows](/docs/installation/platforms/windows)                  | Build 18917 and above |
| <img src="/assets/img/platforms/raspberry-pi.png" width="20" height="20" /> Raspberry Pi                                                         |      In Progress      |
