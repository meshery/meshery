---
layout: default
title: Meshery Operator, MeshSync, Broker Troubleshooting Guide
permalink: guides/troubleshooting/meshery-operator-meshsync
language: en
abstract: This documentation provides comprehensive guidance on troubleshooting in Meshery Operator, MeshSync and Broker, ensuring you can address common issues efficiently.
type: guides
category: troubleshooting
---

There are common issues Meshery users may face whiile operating the [Meshery Operator]({{site.baseurl}}/concepts/architecture/operator/) and its custom controllers, [MeshSync]({{site.baseurl}}/concepts/architecture/meshsync) and [Broker]({{site.baseurl}}/concepts/architecture/broker), that can be resolved by performing specific actions. This documentation aims to empower users by providing a set of troubleshooting tools and actions.

## Status of MeshSync and Meshery Broker

The following table describes the various states of MeshSync and Meshery Broker and their implications.

**MeshSync:**

- **NOT ACTIVE:** Custom Resource not present.
- **ENABLED:** Custom Resource present. MeshSync Controller is not connected to Broker.
- **DEPLOYED:** Custom Resource present. MeshSync Controller present and healthy (what does healthy mean)?
- **RUNNING:** MeshSync pod present and in a running state.
- **CONNECTED:** Deployed and connected to Broker.
- **NOT ACTIVE:** Custom Resource not deployed or no NATS connection.
- **UNDEPLOYED:** Custom Resource not present.

**Meshery Broker:**

- **DEPLOYED:** Custom Resource deployed, External IP exposed.
- **UNDEPLOYED:** Custom Resource deployed or External IP not exposed.
- **CONNECTED:** Deployed, sending data to Meshery Server.


## Deployment Scenarios

Because Meshery is versatile in its deployment models, there are a number of scenarios in which you may need to troubleshoot the health of the operator. The following sections describe the various scenarios and the steps you can take to troubleshoot them.

### In-Cluster Deployment

Meshery Operator, MeshSync, and Broker are deployed in the same cluster as Meshery Server. This is the default deployment scenario when using `mesheryctl system start` or `make run-local`.

### Out-of-Cluster Deployment

1. Meshery Server is deployed on any Docker host (- Meshery Server is deployed on a Docker host, and Meshery Operator is deployed on a Kubernetes cluster).
_or_
2. Meshery is managing multiple clusters, some of which are not the cluster unto which Meshery Server is deployed.


## Fault Scenarios

Common failure situations that Meshery users might face are described below.

1. No deployment of Meshery Operator, MeshSync, and Broker.
    1. Probable cause: Meshery Server cannot connect to Kubernetes cluster; cluster unreachable or kubeconfig without proper permissions needed to deploy Meshery Operator; Kubernetes config initialization issues.
1. Meshery Operator with MeshSync and Broker deployed, but Meshery Server is not receiving data from MeshSync or data the [Meshery Database]({{site.baseurl}}/concepts/architecture/database) is stale.
    1. Probable cause: Meshery Server lost subscription to Meshery Broker; Broker server not expoerting external IP; MeshSync not connected to Broker; MeshSync not running; Meshery Database is stale.
    2. The SQL database in Meshery serves as a cache for cluster state. A single button allows users to dump/reset the Meshery Database.
1. Orphaned MeshSync and Broker controllers - Meshery Operator is not present, but MeshSync and Broker controllers are running.

## Operating Meshery without Meshery Operator

Meshery Operator, MeshSync, and Broker are crucial components in a Meshery deployment. Meshery can function without them, but some functions of Meshery will be disable / unusable. Whether Meshery Operator is initially deployed via `mesheryctl` command or via Meshery Server, you can monitor the health of the Meshery Operator deployment using either the CLI or UI clients.

## Meshery Extension

### Designer Mode

Upon Meshery extension's first load, a GET request initializes the MeshMap plugin. Errors are classified into two types: plugin not found or built on a different version. MeshMap loads, and Designer is functional if no errors occur.

### Visualizer Mode

GraphQL queries fetch header data and view data for the Visualizer canvas. Checks ensure data types and properties are correct, enabling canvas display. If no clusters are connected, a modal prompts the user to select one.


## Synthetic Test for Ensuring Change in Cluster State

Initiate a synthetic check to verify a fully functional Operator deployment, testing MeshSync/Broker connectivity.


## Troubleshooting using Meshery CLI

The following commands are available to troubleshoot Meshery Operator, MeshSync, and Broker.

### Meshery Server and Adapters

- `mesheryctl system status` - Displays the status of Meshery Server and Meshery Adapters.

### Meshery Operator, MeshSync, and Broker

- `mesheryctl system check` - Displays the status of Meshery Operator, MeshSync, and Broker.

## Troubleshooting using Meshery UI

Based on discussed scenarios, the UI exposes tools to perform the following actions:

- (Re)deploy Operator, MeshSync, Broker.
- Uninstall and Install MeshSync, Broker, Operator.
- Reset Database.
- Ad hoc Connectivity Test for Operator, Meshery Broker, MeshSync.
- Reconnect Meshery Server to Meshery Broker.
- Ad hoc Connectivity Test for Kubernetes context.
- Rediscover kubeconfig, delete, (re)upload kubeconfig.


## Desired Behavior

- Empty database shows the main-cluster node.
- Corrupt database triggers an error snackbar with a link to the Settings screen.
- Disconnected Kubernetes displays MeshSync logo pulsating when data is received.
- NATS/MeshSync not running prompts a review of available operations in the Settings panel.


<!-- ## Fault Scenarios Explained for MeshMap Extension

### 1. MISSING PLUGIN

- **MODE STATE:** Visualizer and Designer: Inactive, Not Interactable.
- **CAUSED WHEN:** plugin.so is not in the desired folder.
- **REMEDIATION:** Confirm plugin versions, run make dev or make prod on the local meshery-extension folder.

### 2. INCOMPATIBLE PLUGIN

- **MODE STATE:** Visualizer and Designer: Inactive, Not Interactable.
- **CAUSED WHEN:** go.mod mismatches between meshery/meshery and layer5labs/meshery-extensions.
- **REMEDIATION:** Ensure identical go.mod files and plugin version matching.

### 3. NO ACTIVE CLUSTER CONNECTIONS

- **MODE STATE:** Visualizer: Active, Not Interactable; Designer: Active, Interactable.
- **CAUSED WHEN:** No active Kubernetes cluster connections.
- **REMEDIATION:** Connect a Kubernetes cluster from settings.

### 4. MISSING MESHSYNC DATA / NATS NOT RUNNING

- **MODE STATE:** Visualizer: Active, Not Interactable; Designer: Active, Interactable.
- **CAUSED WHEN:** Meshery Broker lacks an external IP address or networking issues.
- **REMEDIATION:** Delete meshery-meshsync pod, use Docker Desktop/Kind/Minikube/external cloud provider.

### 5. CORRUPT DATABASE

- **MODE STATE:** Visualizer: Inactive, Not Interactable; Designer: Active, Interactable.
- **CAUSED WHEN:** Unable to save/query database.
- **REMEDIATION:** Use System Reset button or remove the config folder and rebuild Meshery.

### 6. INCOMPATIBLE MESHMAP AND MESHERY VERSION

- **MODE STATE:** Visualizer and Designer: Active, Interactable.
- **CAUSED WHEN:** Meshery and MeshMap version mismatches.
- **REMEDIATION:** Pull/build the latest MeshMap, update Mesheryctl and Meshery. -->


This documentation provides comprehensive guidance on troubleshooting in Meshery, ensuring users can address common issues efficiently.

#### See Also

<div class="section">
<ul>
<li><a href="{{ site.baseurl }}/reference/error-codes">Meshery Error Code Reference</a></li>
</ul>
.
</div> 

{% include discuss.html %}
<!-- {:toc} -->