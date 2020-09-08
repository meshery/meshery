---
layout: page
title: WSL2
permalink: installation/platforms/wsl2
---

# Quick Start with WSL2
Below are instructions to configure WSL2 and install a Kubernetes cluster on top.

## Compatibility
The following minimum Windows build versions are required:

| Name   | Version |
|:------ |:-------:|
| [WSL2](https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-18917) | build 18917 |
| [Custom Kernel](https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-18945) | build 18945 |
| [Kernel with K8s required modules](https://docs.microsoft.com/en-us/windows/wsl/release-notes#build-19013) | build 19013 |

## Steps
Perform the following steps in order.
### 1. [Install WSL](https://docs.microsoft.com/en-us/windows/wsl/wsl2-install)
```Powershell
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Windows-Subsystem-Linux
Enable-WindowsOptionalFeature -Online -FeatureName VirtualMachinePlatform
Restart-Computer
wsl --set-default-version 2
```

<i>Note: setting the default to version to 2 will avoid doing it for every distro created</i>

### 2. [Install a new distro](https://docs.microsoft.com/en-us/windows/wsl/install-win10#install-your-linux-distribution-of-choice)
In this how-to, Ubuntu 18.04 will be the distro used, but feel free to use another distro.

### 3. Install Docker Engine, CLI and Compose
The first components required are Docker Engine, Docker CLI (command line interface) and Docker Compose.

[Docker Engine and CLI](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
```bash
sudo apt-get update
sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo service docker start
sudo usermod -aG $USER
```

[Docker Compose](https://docs.docker.com/compose/install/)
```bash
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 4. Install a Kubernetes cluster
Once Docker is installed, the next step will be to install a Kubernetes cluster.
In this how-to, [K3d](https://github.com/rancher/k3d) will be used as it relys only on Docker.

```bash
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
k3d cluster create
export KUBECONFIG="$(k3d kubeconfig get 'k3s-default')"
```

### 5. Finish up

Follow the rest of Meshery [installation](/docs/installation) steps.
