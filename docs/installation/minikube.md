---
layout: page
title: Minikube
parent: Installation
permalink: /installation/minikube
nav_order: 1
---
## Table of contents
{: .no_toc }

1. TOC
{:toc}

---
# Quick Start with Minikube
Below are instructions to generate config file for Minikube cluster which will be used in Meshery configuration.

## Prerequisites
Below versions were successfully tested:
<table style="color:#FFF;text-align:center;">
<th>Version</th><th>Name</th>
<tr><td>1.0.0</td><td><a href="https://kubernetes.io/docs/tasks/tools/install-minikube/">Minikube</a></td></tr>
<tr><td>1.14.1</td><td><a href="https://istio.io/docs/setup/kubernetes/prepare/platform-setup/minikube/">Kubernetes cluster</a></td></tr>
<tr><td>1.14.1</td><td><a href="https://kubernetes.io/docs/tasks/tools/install-kubectl/">Kubectl</a></td></tr>
</table>

## Steps
<div style="font-size:1.25em;">1. Start minikube:</div>
```
minikube start --cpus 4 --memory 8192 --kubernetes-version=v1.14.1
```

<i>Note: minimum memory required is --memory=4096</i>
<div style="font-size:1.25em;">2. Generate config file</div>
This configuration file will be used by Meshery.

```
kubectl config view --minify --flatten > config_minikube.yaml
```
```
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: < cert shortcutted >
    server: https://192.168.99.100:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate-data: <cert shortcutted >
    client-key-data: < key shortcutted >
```
Note: Make sure `current-context` is set to "minikube".

<div style="font-size:1.25em;">3. Finish up</div>

Follow the rest of Meshery [installation](../installation.md) steps.
