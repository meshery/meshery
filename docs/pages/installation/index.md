---
layout: page
title: Installation Guide
permalink: installation
---
# Quick Start 
Getting Meshery up and running on a locally on Docker-enabled system is easy. Use the Meshery command line interface, `mesheryctl`, to start Meshery on any of its [supported platforms](#compatibility-matrix).

### Using Docker on Mac or Linux
**Homebrew**
1. Download `mesheryctl` and install Meshery on Docker on Mac or Linux with homebrew:

**Installation Commands**
```
brew tap layer5io/tap
brew install mesheryctl
mesheryctl start

```
**Upgradation Command**
```
brew upgrade mesheryctl

```

**Bash**
1. Download `meshery` and install Meshery on Docker on Mac or Linux with homebrew:

```
curl -L https://git.io/meshery | bash -
```

### Using Docker on Windows
1. Visit [Meshery releases](https://github.com/layer5io/meshery/releases/latest), download and unzip the `mesheryctl` utility. Add `mesheryctl` to your PATH for ease of use.

```
./mesheryctl start
```

2. Upon starting Meshery successfully, instructions to access Meshery will be printed on the sceen.

## What is `mesheryctl`?
`mesheryctl` is a command line interface to manage a Meshery deployment. `mesheryctl` allows you to control Meshery's lifecycle with commands like `start`, `stop`, `status`, `cleanup`. Running `cleanup` will remove all active container instanaces, prune pulled images and remove any local volumes crated by starting Meshery.

# Getting Started with other Platforms<a name="compatibility-matrix"></a>
The table below provides installation instructions for each platform that Meshery is compatibile with:

| Platform      | Version       |
| -------------:|:-------------|   
| [Docker Engine](/docs/installation/docker) | 19.x and above |
| [Docker Desktop](/docs/installation/docker) | 2.0.x and above |
| [EKS](/docs/installation/eks) | 1.12.x and above |
| [GKE](/docs/installation/gke) | 1.14.x and above |
| [Kubernetes](/docs/installation/kubernetes) | 1.12.x and above |
| [Minikube](/docs/installation/minikube) | 1.2.x and above |
| [WSL2](/docs/installation/wsl2) | Build 18917 and above |
