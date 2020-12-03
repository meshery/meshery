---
layout: page
title: AKS
permalink: installation/platforms/aks
type: installation
language: en
list: include
image: /docs/assets/img/platforms/aks.svg
---

{% include installation_prerequisites.html %}

**To set up and run Meshery on AKS**:

- Connect Meshery to your AKS cluster
    * [Meshery CLI (mesheryctl)](#connect-meshery-to-azure-kubernetes-cluster)
    * [Azure CLI (az)](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)

- [Install Meshery on your AKS cluster](#install-meshery-into-your-aks-cluster)
- [Access Meshery's UI](#port-forward-to-the-meshery-ui)

### **Connect Meshery to Azure Kubernetes Cluster**

The following set of instructions expects you to have created a AKS cluster in your resource group Configure Meshery to connect to your AKS cluster by executing:

```
$ mesheryctl system config aks
```

#### **Manual Steps**

Alternatively, you may execute the following steps to manually configure Meshery to connect to your AKS cluster.

- Install [Azure CLI(az)](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli), and login
to your azure account using *az login*.
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

### **Install Meshery into your AKS cluster** 

```sh
$ kubectl create ns meshery
$ helm repo add meshery https://meshery.io/charts/
$ helm install meshery --namespace meshery meshery/meshery
```

### **Port forward to the Meshery UI**

```sh
export POD_NAME=$(kubectl get pods --namespace meshery -l "app.kubernetes.io/name=meshery,app.kubernetes.io/instance=meshery" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace meshery port-forward $POD_NAME 8080:8080
```

Congratulations. Meshery should now be running in your AKS cluster and the Meshery UI should be locally accessible. Navigate to [http://localhost:8080](http://localhost:8080] to log into Meshery and start managing your service mesh deployments.