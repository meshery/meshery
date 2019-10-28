---
layout: page
title: GKE
permalink: installation/gke
---

# Quick Start with Google Kubernetes Engine (GKE)
You may perform the steps outlined under [Managed Kubernetes](#managedk8s) following by hand or run the [generate_kubeconfig_gke.sh](./generate_kubeconfig_gke.sh) shell script using the desired ServiceAccount name and Namespace arguments, like so:

`./generate_kubeconfig_gke.sh cluster-admin-sa-gke default`

Having run this script, supply the generated file `config-cluster-admin-sa-gke-default.yaml` in your Meshery settings page.

## Configuration
This is where you configure your settings on the adaptor(Istio etc) and other things 

## Managed Kubernetes
In order to run Meshery in a managed Kubernetes environment, you will need to assign an existing `ServiceAccount` or create a new `ServiceAccount`:

1. Create a `ServiceAccount` with `cluster-admin` role.
1. Get secret name from `ServiceAccount`.
1. Extract CA certificate and user token from the secret.
1. Generate new kubeconfig yaml file to use as input to Meshery.