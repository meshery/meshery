---
layout: page
title: EKS
permalink: installation/eks
---

# Quick Start with Amazon Elastic Kubernetes Service (EKS)

## Managed Kubernetes
In order to run Meshery in a managed Kubernetes environment, you will need to assign an existing `ServiceAccount` or create a new `ServiceAccount`:

Note: Make sure you are able to access EKS with kubectl by follwing <a href="https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html" target="_blank"> EKS guide</a>







### 1. Create a `ServiceAccount` with `cluster-admin` role.



 &nbsp;&nbsp;&nbsp; a. Create a service account


```
kubectl create serviceaccount sa-1
```
&nbsp;&nbsp;&nbsp; b. Adding/Binding `cluster-admin` role to new service account `sa-1`



```
kubectl create clusterrolebinding sa-1-binding --clusterrole=cluster-admin \
 --serviceaccount=default:sa-1
 ```

### 2. Get secret name from `ServiceAccount`.

&nbsp;&nbsp;&nbsp; a. Get secret name

```
$ kubectl get secrets

NAME                           TYPE                                  DATA   AGE
default-token-fnfjp            kubernetes.io/service-account-token   3      95d
sa-1-token-5z9xj               kubernetes.io/service-account-token   3      66m

```
Note: Here my secret Name is **sa-1-token-5z9xj** 

&nbsp;&nbsp;&nbsp; b. Get secret/token

```

$ kubectl describe secret  sa-1-token-5z9xj 
Name:         sa-1-token-5z9xj
Namespace:    default
Labels:       <none>
Annotations:  kubernetes.io/service-account.name: sa-1
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
kubectl config set-credentials sa-1 --token=XXXXX

o/p:User "sa-1" set.
```

&nbsp;&nbsp;&nbsp; b. Set current context to our new service account `sa-1`

```
kubectl config set-context --current --user=sa-1

o/p:
Context "aws" modified.
 ```

 &nbsp;&nbsp;&nbsp; c. Generate kubeconfig yaml file to use as input to Meshery.

 ```
 kubectl config view --minify --flatten >  config_aws_eks.yaml 
 ```


