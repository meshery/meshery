---
layout: default
title: Managing Multiple Kubernetes Clusters
abstract: Learn how to connect, configure, and manage multiple Kubernetes clusters using Meshery.
permalink: guides/infrastructure-management/multi-cluster-management
type: guides
category: infrastructure-management
language: en
---

# Tutorial: Managing Multiple Kubernetes Clusters with Meshery

Meshery is designed natively for multi-cluster operations. This tutorial walks you through connecting multiple Kubernetes clusters to Meshery, switching between contexts, and deploying designs across cluster boundaries.

## Learning Objectives

By the end of this tutorial, you will be able to:

- Connect multiple Kubernetes clusters to Meshery
- Switch between cluster contexts
- Deploy designs to specific or multiple clusters

## Prerequisites

Before you begin, ensure you have:

- Meshery installed and running
- `mesheryctl` installed on your local machine
- Two or more Kubernetes clusters available
- `kubectl` configured with access to your clusters
- A `kubeconfig` file containing contexts for each cluster (multiple contexts can be defined in a single kubeconfig)

## Step 1: Connect Your Clusters to Meshery

Meshery can discover clusters automatically if it has access to your `kubeconfig`, or you can add them manually through the UI.

### Option A: Using mesheryctl

Use `mesheryctl` to configure your local environment and register your cluster contexts with Meshery Server. Run the following command, specifying the path to your `kubeconfig`:

```bash
mesheryctl system config kubernetes --file ~/.kube/config
```

This reads the contexts defined in your `kubeconfig`, configures your local environment, and registers them with Meshery. If your clusters are already defined in your default `kubeconfig` (`~/.kube/config`), Meshery may auto-detect them on startup.

> **Note:** You can define multiple cluster contexts in a single `kubeconfig` file. `mesheryctl` will register all available contexts with Meshery.

### Option B: Manual Addition via UI

1. Open Meshery UI and navigate to **Settings > Environments**.
2. Click **Add Cluster**.
3. Upload your `kubeconfig` file. All contexts within the file will be made available.
4. Select the clusters you wish to connect.

## Step 2: Verify Cluster Connections

After connecting your clusters, verify that Meshery has successfully established connections:

- In Meshery UI, go to **Settings > Environments**.
- Confirm each cluster appears and shows a `Connected` status.

Alternatively, use the CLI:

```bash
mesheryctl system check
```

This runs a health check and reports the connectivity status of each registered cluster.

## Step 3: Switch Between Cluster Contexts

When working with multiple clusters, you can target a specific cluster for operations:

- In Meshery UI, use the cluster selector in the top navigation to switch your active context.
- Set your current context using `kubectl`:

```bash
kubectl config use-context <context-name>
```

Then use mesheryctl system context switch <context-name> to set the active context for Meshery operations.

## Step 4: Deploy a Design to Multiple Clusters

Meshery enables you to deploy the same design to one or more clusters simultaneously:

1. Open [Kanvas](/extensions/kanvas) (Meshery's visual designer).
2. Create or open an existing design.
3. Click **Deploy** and select one or more target clusters from the cluster list.
4. Confirm the deployment. Meshery will orchestrate the rollout across the selected clusters.

## Troubleshooting

If you encounter connectivity issues between Meshery and your clusters, refer to the [Meshery Operator, MeshSync, and Broker Troubleshooting Guide](/guides/troubleshooting/meshery-operator-meshsync). That guide covers the most common scenarios including connection failures, MeshSync sync issues, and Broker communication problems.

For installation-related problems, see [Troubleshooting Meshery Installations](/guides/troubleshooting/installation).

## Summary

You have connected multiple Kubernetes clusters to Meshery, verified their status, and learned how to target specific clusters when deploying designs. Meshery's multi-cluster support allows you to manage infrastructure across environments from a single control plane.

## Next Steps

- Explore [Environments](/concepts/logical/environments) to group and manage clusters as a unit.
- Learn about [Kanvas](/extensions/kanvas) for collaborative visual design across clusters.
- Review the [Meshery CLI reference](/reference/mesheryctl) for advanced `mesheryctl` commands.
