---
layout: page
title: EKS
permalink: installation/eks
---

# Quick Start with Amazon Elastic Kubernetes Service (EKS)

Make your way to Meshery and download the token by clicking the "Get Token" option in the dropdown menu. Utilize the token to run the following command:

| command           | flag                | function                                                     | Usage                     |
|:------------------|:-------------------:|:-------------------------------------------------------------|:--------------------------|
|                   | --system config     | configures Meshery with the kubeconfig, generated with the help of user details, to provide cluster access for public clouds(EKS). | `mesheryctl system config gke --token "PATH TO TOKEN"` |

## **Manual Configuration**


Follow the below mentioned steps to set up manually:

### **Managed Kubernetes**
In order to run Meshery in a managed Kubernetes environment, you will need to assign an existing `ServiceAccount` or create a new `ServiceAccount`:

Note: Make sure you are able to access EKS with kubectl by following the <a href="https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html" target="_blank"> EKS guide</a>


### 1. Create a `ServiceAccount` with `cluster-admin` role

 &nbsp;&nbsp;&nbsp; a. Create a service account

```
kubectl create serviceaccount meshery
```
&nbsp;&nbsp;&nbsp; b. Adding/Binding `cluster-admin` role to new service account `meshery`

```
kubectl create clusterrolebinding meshery-binding --clusterrole=cluster-admin \
 --serviceaccount=default:meshery
 ```

### 2. Get secret name from `ServiceAccount`.

&nbsp;&nbsp;&nbsp; a. Get secret name

```
$ kubectl get secrets

NAME                           TYPE                                  DATA   AGE
default-token-fnfjp            kubernetes.io/service-account-token   3      95d
meshery-token-5z9xj               kubernetes.io/service-account-token   3      66m

```
Note: Here the secret name is **meshery-token-5z9xj**

&nbsp;&nbsp;&nbsp; b. Get secret/token

```

$ kubectl describe secret  sa-1-token-5z9xj
Name:         meshery-token-5z9xj
Namespace:    default
Labels:       <none>
Annotations:  kubernetes.io/service-account.name: meshery
              kubernetes.io/service-account.uid: 397XXX-XXX-XXXX-XXXXX-XXXXX

Type:  kubernetes.io/service-account-token

Data
====
ca.crt:     1025 bytes
namespace:  7 bytes
token:      XXXhbGciOiJSUXXXX

```


### 3. Generate new kubeconfig yaml file to use as input to Meshery.

&nbsp;&nbsp;&nbsp; a.Set config Credential using above generate `token`.

```
kubectl config set-credentials meshery --token=XXXXX

o/p:User "meshery" set.
```

&nbsp;&nbsp;&nbsp; b. Set current context to our new service account `meshery`

```
kubectl config set-context --current --user=meshery

o/p:
Context "aws" modified.
 ```

 &nbsp;&nbsp;&nbsp; c. Generate kubeconfig yaml file to use as input to Meshery.

 ```
 kubectl config view --minify --flatten >  config_aws_eks.yaml
 ```
