---
layout: page
title: Minikube
permalink: installation/platforms/minikube
---

# Quick Start with Minikube

Below are instructions to generate config file for Minikube cluster which will be used in Meshery configuration.

## Compatibility

The following minimum component versions are required:

| Name                                                                                  | Version |
| :------------------------------------------------------------------------------------ | :-----: |
| [Minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)                  |  1.0.0  |
| [Kubernetes](https://istio.io/docs/setup/kubernetes/prepare/platform-setup/minikube/) | 1.14.1  |
| [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)                    | 1.14.1  |

## Steps

Perform the following steps in order.

### 1. Start minikube

```
minikube start --cpus 4 --memory 8192 --kubernetes-version=v1.14.1
```

<i>Note: minimum memory required is --memory=4096 (for Istio deployments only)</i>

<i>Note: If you are using docker driver, after completing meshery installation steps execute below command to establish connectivity between Meshery and Kubernetes server.
</i>

```
docker network connect bridge meshery_meshery_1
```

### 2. Configure Meshery to use minikube

1. Login to Meshery. Under your user profile, click `Get Token`.
2. Use `mesheryctl` to configure Meshery to use minikube. Execute:

```sh
mesheryctl system config minikube -t ~/Downloads/auth.json
```

### Alternative 2. Manually generate and load kubeconfig file for Meshery to use

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

Note: Make sure "current-context" is set to "minikube".

### 3. Finish up

Follow the rest of Meshery [installation](/docs/pages/installation/index.md) steps.
