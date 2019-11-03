---
layout: page
title: Installation Guide
permalink: installation
---

# Platform Compatibility <a name="compatibility-matrix"></a>

Meshery's compatibility has been confirmed with the following platforms:

| Platform      | Version       |
| -------------:|:-------------|   
| [Docker Engine](/docs/installation/docker) | 19.x and above |
| [Docker Desktop](/docs/installation/docker) | 2.0.x and above |
| [EKS](/docs/installation/eks) | 1.12.x and above |
| [GKE](/docs/installation/gke) | 1.14.x and above |
| [Kubernetes](/docs/installation/kubernetes) | 1.12.x and above |
| [Minikube](/docs/installation/minikube) | 1.2.x and above |

## What is `mesheryctl`?
`mesheryctl` is a command line interface to manage a Meshery deployment. `mesheryctl` allows you to control Meshery's lifecycle with commands like `start`, `stop`, `status`, `cleanup`. Running `cleanup` will remove all active container instanaces, prune pulled images and remove any local volumes crated by starting Meshery.

## Quick Start 
Download `mesheryctl`. Install Meshery on your local machine by running the following:

```
sudo curl -L https://git.io/meshery -o /usr/local/bin/meshery
sudo chmod a+x /usr/local/bin/meshery
meshery start
```
Upon starting Meshery successfully, instructions to access Meshery will be printed on the sceen.
