---
layout: page
title: GKE
permalink: installation/gke
---

# Quick Start with Google Kubernetes Engine (GKE)

- Make your way to meshery.ui, login with your user details and head over to the local port, found at <code>localhost:3000</code>.
- Download the token by clicking the "Get Token" option in the dropdown menu under your User Account avatar.
- Utilize the token to run the following command:

| command           | flag                | function                                                     | Usage                     |
|:------------------|:-------------------:|:-------------------------------------------------------------|:--------------------------|
|                   | --system config     | configures Meshery with the kubeconfig, generated with the help of user details, to provide cluster access for public clouds(GKE). | `mesheryctl system config gke --token "PATH TO TOKEN"` |

Once configured, head over to the Quick Start Guide and continue with the steps outlined for <a href = "https://github.com/layer5io/meshery/blob/d75b6eec23021a1d8d3a38ec890fb10d44af3a35/docs/pages/installation/gke.md#manual-configuration">GKE</a>.

## **Manual Configuration**

Follow the below mentioned steps to set up manually:

You may perform the steps outlined under [Managed Kubernetes](#managedk8s), following by hand or run the [generate_kubeconfig_gke.sh](./generate_kubeconfig_gke.sh) shell script using the desired ServiceAccount name and Namespace arguments, like so:

`./generate_kubeconfig_gke.sh cluster-admin-sa-gke default`

Having run this script, fire up the meshery server and head over to the local port, usually found at <code>localhost:3000/settings</code> or <code>localhost:9081/settings</code>. This is where you configure your settings on the adaptor(Istio etc).
<br>Supply the generated file <code><b> config-cluster-admin-sa-gke-default.yaml </b></code> under kube-config.

### **Managed Kubernetes**
In order to run Meshery in a managed Kubernetes environment, you will need to assign an existing `ServiceAccount` or create a new `ServiceAccount`:

1. Create a `ServiceAccount` with `cluster-admin` role.
2. Get secret name from `ServiceAccount`.
3. Extract CA certificate and user token from the secret.
4. Generate new kubeconfig yaml file to use as input to Meshery.
