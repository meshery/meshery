---
layout: page
title: EKS
permalink: installation/eks
---

# Quick Start with Amazon Elastic Kubernetes Service (EKS)

You may perform the steps outlined under [Managed Kubernetes](#managedk8s) following by hand or run the [generate_kubeconfig_eke.sh](./generate_kubeconfig_gke.sh) shell script using the desired ServiceAccount name and Namespace arguments, like so:

`./generate_kubeconfig_eke.sh eks-admin kube-system`

## Managed Kubernetes

In order to run Meshery in a managed Kubernetes environment, you will need to assign an existing `ServiceAccount` or create a new `ServiceAccount`:

1. Create a `ServiceAccount` with `cluster-admin` role.

``` serviceaccount.yaml

apiVersion: v1
kind: ServiceAccount
metadata:
  name: eks-admin
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: eks-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: eks-admin
  namespace: kube-system
```

1. kubectl create -f serviceaccount.yaml and run ./generate_kubeconfig_gke.sh  eks-admin kube-system  
1. KUBECONFIG is going to store in KUBECONFIG=/tmp/kube/k8s-eks-admin-kube-system-conf
1. copy the new KUBECONFIG as input to Meshery.
