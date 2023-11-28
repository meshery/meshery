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

## Operating Meshery without Meshery Operator

Meshery Operator, MeshSync, and Broker are crucial components in a Meshery deployment. Meshery can function without them, but some functions of Meshery will be disable / unusable. Whether Meshery Operator is initially deployed via `mesheryctl` command or via Meshery Server, you can monitor the health of the Meshery Operator deployment using either the CLI or UI clients. 

### Scenarios in which Meshery Operator is not deployed

## Concepts

### Operator, MeshSync, and Broker

Meshery's architecture involves the following scenarios during startup:

**Single Cluster Deployment:**

- Meshery Operator with MeshSync and Broker deployed.
- Meshery Operator not deployed, but MeshSync and Broker are.
- Orphaned MeshSync and Broker controllers.
- No deployment of Meshery Operator, MeshSync, and Broker.

**Multiple Cluster Scenarios:**
- Meshery Operator, MeshSync, and Broker are crucial components.
- Meshery can function without them but has limited use cases.
- Meshery Server ensures health even if mesheryctl initiates provisioning.

## SQL Database (Cache)

The SQL database in Meshery serves as a cache for cluster state. A single button allows users to dump/reset the Meshery Database.

## Meshery Extension

### Designer Mode

Upon Meshery extension's first load, a GET request initializes the MeshMap plugin. Errors are classified into two types: plugin not found or built on a different version. MeshMap loads, and Designer is functional if no errors occur.

### Visualizer Mode

GraphQL queries fetch header data and view data for the Visualizer canvas. Checks ensure data types and properties are correct, enabling canvas display. If no clusters are connected, a modal prompts the user to select one.

## Status of MeshSync and Meshery Broker

### Present Scenario:

**MeshSync:**
- **ENABLED:** Custom Resource deployed.
- **NOT ACTIVE:** Custom Resource not deployed.

**MesheryBroker:**
- **CONNECTED:** Custom Resource deployed, NATS connection established.
- **NOT ACTIVE:** Custom Resource not deployed or no NATS connection.

### Proposed:

**MeshSync:**
- **ENABLED:** Custom Resource present, MeshSync may/may not be connected.
- **DEPLOYED:** MeshSync pod present and healthy.
- **RUNNING:** MeshSync pod present and in a running state.
- **CONNECTED:** Deployed and connected to Broker.
- **UNDEPLOYED:** Custom Resource not present.

**Meshery Broker:**
- **DEPLOYED:** Custom Resource deployed, External IP exposed.
- **UNDEPLOYED:** Custom Resource deployed or External IP not exposed.
- **CONNECTED:** Deployed, sending data to Meshery Server.

## Fault Scenarios

Common issues Meshery users face frequently, requiring tooling exposure in the UI:

1. Not receiving data from MeshSync.
2. Stale data in the Meshery Database.
3. Operator, Broker, and/or MeshSync not deployed.
4. Kubernetes Config initialization issues.

## Troubleshooting Toolkit in Meshery UI

Based on discussed scenarios, the UI exposes tools to perform the following actions:

- (Re)deploy Operator, MeshSync, Broker.
- Uninstall and Install MeshSync, Broker, Operator.
- Reset Database.
- Ad hoc Connectivity Test for Operator, Meshery Broker, MeshSync.
- Reconnect Meshery Server to Meshery Broker.
- Ad hoc Connectivity Test for Kubernetes context.
- Rediscover kubeconfig, delete, (re)upload kubeconfig.

## Synthetic Test for Ensuring Change in Cluster State

Initiate a synthetic check to verify a fully functional Operator deployment, testing MeshSync/Broker connectivity.

## Fault Scenarios Explained for MeshMap Extension

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
- **REMEDIATION:** Pull/build the latest MeshMap, update Mesheryctl and Meshery.

## Desired Behavior

- Empty database shows the main-cluster node.
- Corrupt database triggers an error snackbar with a link to the Settings screen.
- Disconnected Kubernetes displays MeshSync logo pulsating when data is received.
- NATS/MeshSync not running prompts a review of available operations in the Settings panel.

---

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