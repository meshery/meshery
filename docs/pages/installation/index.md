---
layout: page
title: Installation Guide
permalink: installation
---
# Quick Start 
Getting Meshery up and running on a locally on Docker-enabled system is easy. Use the Meshery command line interface, `mesheryctl`, to start Meshery on any of its [supported platforms](#compatibility-matrix).

### Using Docker on Mac or Linux
**Homebrew**

* Install `mesheryctl` and run Meshery on Mac with homebrew:

**Installation Commands**
To install `mesheryctl`, execute the following commands.
```
brew tap layer5io/tap
brew install mesheryctl
mesheryctl start
```
**Upgrading**
To upgrade `mesheryctl`, execute the following command.
```
brew upgrade mesheryctl
```

Example output of a successful upgrade:
```
➜  ~ brew upgrade mesheryctl
==> Upgrading 1 outdated package:
layer5io/tap/mesheryctl 0.3.2 -> 0.3.4
==> Upgrading layer5io/tap/mesheryctl
==> Downloading https://github.com/layer5io/meshery/releases/download/v0.3.4/mesheryctl_0.3.4_Darwin_x86_64.zip
==> Downloading from https://github-production-release-asset-2e65be.s3.amazonaws.com/157554479/17522b00-2af0-11ea-8aef-cbfe8
######################################################################## 100.0%
🍺  /usr/local/Cellar/mesheryctl/0.3.4: 5 files, 10.2MB, built in 4 seconds
Removing: /usr/local/Cellar/mesheryctl/0.3.2... (5 files, 10.2MB)
Removing: /Users/lee/Library/Caches/Homebrew/mesheryctl--0.3.2.zip... (3.9MB)
==> Checking for dependents of upgraded formulae...
==> No dependents found!

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
