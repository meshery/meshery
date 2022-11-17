---
layout: default
title: Supported Platforms
permalink: installation/platforms
type: installation
display-title: "false"
redirect_from: installation/platforms/
language: en
list: exclude
---

# Supported Platforms<a name="compatibility-matrix"></a>

Meshery deploys as a set of Docker containers, which can be deployed to either a Docker host or Kubernetes cluster. See the complete list of supported platforms in the table below. With service meshes having sprung to life in the context of Kubernetes, so too, can Meshery’s deployment models be characterized in the context of Kubernetes. A given deployment of Meshery can be described as either an _in-cluster_ or an _out-of-cluster_ deployment. Meshery deploys as a stand-alone, management plane on a Docker host (_out-of-cluster_) or as a management plane in a Kubernetes cluster (_in-cluster_).

## Platform Compatibility Matrix

{% include alert.html type="info" title="<a href='/project/compatibility-matrix'>Complete Compatibility Matrix and Testing Dashboard</a>" content="For a complete compatibility matrix and project test status dashboard, see Meshery <a href='/project/compatibility-matrix'>Compatibility Matrix</a>" %}

Find installation instructions for the Supported Platforms in the compatibility table.

| Platform                                                                                                                                                               |         Version         |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :---------------------: |
| <img src="/assets/img/platforms/docker.svg" width="20" height="20" /> [Docker]({{ site.baseurl }}/installation/platforms/docker)                                       |                         |
| &nbsp;&nbsp;&nbsp; <img src="/assets/img/platforms/docker.svg" width="20" height="20" /> [Docker Engine]({{ site.baseurl }}/installation/platforms/docker)             |     19.x and above      |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/docker.svg" width="20" height="20" /> [Docker Desktop]({{ site.baseurl }}/installation/platforms/docker)             |     2.0.x and above     |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/docker.svg" width="20" height="20" /> [Docker Extension]({{ site.baseurl }}/installation/platforms/docker-extension) |     2.0.x and above     |
| <img src="/assets/img/platforms/kubernetes.svg" width="20" height="20" /> [Kubernetes]({{ site.baseurl }}/installation/platforms/kubernetes)                           |    1.12.x and above     |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/aks.svg" width="20" height="20" /> [AKS]({{ site.baseurl }}/installation/platforms/aks)                              |                         |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/eks.png" width="20" height="20" /> [EKS]({{ site.baseurl }}/installation/platforms/eks)                              |    1.12.x and above     |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/gke.png" width="20" height="20" /> [GKE]({{ site.baseurl }}/installation/platforms/gke)                              |    1.14.x and above     |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/helm.svg" width="20" height="20" /> [Helm]({{ site.baseurl }}/installation/platforms/kubernetes#using-helm)          |                         |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/kind.png" width="20" height="20" /> [KinD]({{ site.baseurl }}/installation/platforms/kind)                           |         v0.7.0          |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/kubesphere.png" width="20" height="20" /> [Kubesphere]({{ site.baseurl }}/installation/platforms/kubesphere)               |     v3.3     | 
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/minikube.png" width="20" height="20" /> [Minikube]({{ site.baseurl }}/installation/platforms/minikube)               |     1.2.x and above     |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/openshift.svg" width="20" height="20" /> OpenShift                                                                   |       In Progress       |
| <img src="/assets/img/platforms/apple.svg" width="20" height="20" vertical-align="middle" /> [Mac]({{ site.baseurl }}/installation#mac-or-linux)                       |                         |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/homebrew.png" width="20" height="20" /> [Mac - Homebrew]({{ site.baseurl }}/installation#mac-or-linux)               | macOS 10.12 - 10.15, 11 |
| &nbsp;&nbsp;&nbsp; [Scoop]({{ site.baseurl }}/installation/platforms/scoop)                                                                                            |                         |
| &nbsp;&nbsp;&nbsp;<img src="/assets/img/platforms/wsl2.png" width="20" height="20" /> [Windows]({{ site.baseurl }}/installation/platforms/windows)                     |  Build 18917 and above  |
| <img src="/assets/img/platforms/raspberry-pi.png" width="20" height="20" /> Raspberry Pi                                                                               |       In Progress       |
