#Quick Start with Minikube
Below are instructions to generate config file for Minikube cluster which will be used in Meshery configuration.

##Prerequisites
Below versions were successfully tested:

| Version | Name | Details |
| --- | ------ | ------ |
|1.0.0|Minikube|[link](https://kubernetes.io/docs/tasks/tools/install-minikube/)|
|1.14.1|Kubernetes cluster|[link](https://istio.io/docs/setup/kubernetes/prepare/platform-setup/minikube/)|
|1.14.1|Kubectl|[link](https://kubernetes.io/docs/tasks/tools/install-kubectl/)|

##Steps
1. Start minikube:
```
minikube start --cpus 4 --memory 8192 --kubernetes-version=v1.14.1
```

2. Generate config file which will be used in Meshery configuration:
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
Note: Make sure current-context is set to minikube.

3. Follow Meshery [installation](../installation.md).
