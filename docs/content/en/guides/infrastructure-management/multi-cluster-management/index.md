---
title: "Tutorial: Multi-cluster Management"
description: "A step-by-step tutorial on how to manage and operate multiple Kubernetes clusters from a single Meshery instance."
categories: [infrastructure]
---

# Tutorial: Multi-cluster Management with Meshery

This tutorial walks you through the process of connecting and managing multiple Kubernetes clusters. Meshery acts as a multi-cluster management plane, allowing you to operate disparate infrastructure from a single deployment.

## Prerequisites

To complete this tutorial, you will need:
- Meshery installed and running (see [Installation](/installation)).
- A `kubeconfig` file containing one or more cluster contexts.

## Connecting Multiple Clusters

Meshery can discover clusters automatically via your environment or you can add them manually. A single `kubeconfig` file can contain multiple contexts, all of which can be imported into Meshery at once.

### Approach 1: Using `mesheryctl`

For users who prefer the command line, `mesheryctl` provides a streamlined way to point Meshery to your cluster configurations.

**Step 1: Check your current contexts**
Verify that your local machine has the desired contexts:
```bash
kubectl config get-contexts
```

**Step 2: Point Meshery to your kubeconfig**
Use the following command to configure Meshery with your clusters:
```bash
mesheryctl system config kubernetes --kubeconfig /path/to/your/kubeconfig
```

### Approach 2: Manual Addition via UI

You can also manage your cluster connections directly from the Meshery UI:
1. Navigate to **Settings** in the left sidebar.
2. Click on the **Environment** tab.
3. Under the **Kubernetes** section, click the **(+)** button to upload a new `kubeconfig` or select a specific context from an existing configuration.

## Managing Contexts

Once your clusters are connected, they appear as **Contexts** in Meshery. You can switch the active cluster by using the context switcher located at the top-right of the Meshery UI.

### Viewing Cluster Health
Navigate to the **Dashboard** to see a high-level overview of the status of all your connected contexts.

## Deploying Designs Across Clusters

When deploying designs, you can target specific clusters:
1. Go to **Designs** and select a design.
2. In the **Deploy** dialog, select the target context(s) where the design should be applied.
3. Click **Deploy**.

## Troubleshooting

If you encounter issues connecting to remote clusters, please refer to the central [Troubleshooting Guide](/guides/troubleshooting) for detailed network and permission checks.
