---
layout: page
title: EKS
permalink: installation/eks
---

# Quick Start with Amazon Elastic Kubernetes Service (EKS)

## Managed Kubernetes
In order to run Meshery in a managed Kubernetes environment, you will need to assign an existing `ServiceAccount` or create a new `ServiceAccount`:

1. Create a `ServiceAccount` with `cluster-admin` role.
1. Get secret name from `ServiceAccount`.
1. Extract CA certificate and user token from the secret.
1. Generate new kubeconfig yaml file to use as input to Meshery.
