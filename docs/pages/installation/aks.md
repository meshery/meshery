---
layout: page
title: AKS
permalink: installation/platforms/aks
---

# Quick Start with Azure Kubernetes Service (AKS)

## Managed Kubernetes
The following set of instructions expects you to have created a AKS cluster in your resource group.

- Install [Azure CLI(`az`)](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli), and login
to your azure account using `az login`.
- After successfull login, you have to select the subscription with which your AKS is associated with
```shell script
az account set --subscription {{SUBSCRIPTION_ID}}
```
- Get the kubeconfig from your AKS cluster
```shell script
az aks get-credentials --resource-group {{RESOURCE_GROUP}} --name {{AKS_SERVICE_NAME}}
```
- Set your cluster context and check your cluster-info
```shell script
kubectl set-context {{AKS_SERVICE_NAME}}
kubectl cluster-info
```
- Install Meshery into your AKS cluster 
```shell script
kubectl create ns meshery
helm repo add meshery https://meshery.io/charts/
helm install meshery --namespace meshery meshery/meshery
```
- Port-forward meshery pod 
```shell script
export POD_NAME=$(kubectl get pods --namespace meshery -l "app.kubernetes.io/name=meshery,app.kubernetes.io/instance=meshery" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace meshery port-forward $POD_NAME 8080:8080
```
- Congrats, you have got Meshery up & running on your AKS cluster. Navigate to `localhost:8080/provider` to login to meshery
and start managing your service mesh deployments.