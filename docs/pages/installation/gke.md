---
layout: page
title: GKE
permalink: installation/gke
---

# Quick Start with Google Kubernetes Engine (GKE)

- Navigate to the Meshery UI, login with your user details and head over to the local port (e.g. `http://localhost:9081`).
- Download the authentication token by clicking the "Get Token" option in the dropdown menu under your User Account avatar.
- Utilize the token to run the following command:

| command           | flag                | function                                                     | Usage                     |
|:------------------|:-------------------:|:-------------------------------------------------------------|:--------------------------|
|                   | --system config     | configures Meshery with the kubeconfig, generated with the help of user details, to provide cluster access for public clouds(GKE). | `mesheryctl system config gke --token "PATH TO TOKEN"` |

Once configured, you can proceed to the Quick Start Guide and continue with the steps outlined for [GKE](/docs/installation/gke#manual-configuration).

## **Manual Configuration**

Follow the below mentioned steps to set up manually:

You may perform the steps outlined under [Managed Kubernetes](#managedk8s), following by hand or run the [generate_kubeconfig_gke.sh](./generate_kubeconfig_gke.sh) shell script using the desired ServiceAccount name and Namespace arguments, like so:

`./generate_kubeconfig_gke.sh cluster-admin-sa-gke default`

- Once the script is complete, you may proceed to start the Meshery server and navigate to the local port (e.g. `http://localhost:9081/settings`).
- Start with the GKE-compatible configuration by executing `mesheryctl system start` and supply the generated file `**config-cluster-admin-sa-gke-default.yaml**` under kube-config.

### **Managed Kubernetes**
In order to run Meshery in a managed Kubernetes environment, you will need to assign an existing `ServiceAccount` or create a new `ServiceAccount`:

1. Create a `ServiceAccount` with `cluster-admin` role.
2. Get secret name from `ServiceAccount`.
3. Extract CA certificate and user token from the secret.
4. Generate new kubeconfig yaml file to use as input to Meshery.
