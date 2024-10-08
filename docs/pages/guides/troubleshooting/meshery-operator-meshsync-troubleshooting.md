---
layout: default
title: Meshery Operator, MeshSync, Broker Troubleshooting Guide
permalink: guides/troubleshooting/meshery-operator-meshsync
language: en
abstract: Comprehensive guidance for troubleshooting Meshery Operator, MeshSync and Broker deployments under various scenarios.
type: guides
category: troubleshooting
---

{% include alert.html type="info" title="What is Meshery Operator?" content="<a href='/concepts/architecture/operator'>Meshery Operator</a> controls and monitors the lifecycle of components deployed inside Meshery-managed Kubernetes clusters. Learn more about <a href='/concepts'>Meshery's architecture</a>." %}

This guide offers comprehensive for troubleshooting instructions for [Meshery Operator]({{site.baseurl}}/concepts/architecture/operator) and its custom controllers, [MeshSync]({{site.baseurl}}/concepts/architecture/meshsync) and [Broker]({{site.baseurl}}/concepts/architecture/broker). Follow the steps outlined in this document to ensure a smooth Meshery deployment.

First, understand the [Meshery Operator Deployment Scenarios](#meshery-operator-deployment-scenarios) and the [Status of Meshery Operator, MeshSync, and Meshery Broker](#understanding-the-status-of-meshery-operator-meshsync-and-meshery-broker) to identify the deployment model fitting of your environment. Then, follow the guidance under the respective scenario to troubleshoot accordingly.

{% include alert.html type="dark" title="Meshery Error Code Reference" content="Have specific error with an error code? See the <a href='/reference/error-codes'>Meshery Error Code Reference</a> for probable cause and suggested remediations." %}

## Understanding the Status of Meshery Operator, MeshSync, and Meshery Broker

Each Meshery Operator controller offers a health status that you can use to understand their current health in your deployment. Their health statuses and meanings are described below.of MeshSync and Meshery Broker.

### MeshSync Health Status

- **ENABLED:** Custom Resource present. MeshSync Controller is not connected to Broker.
- **DEPLOYED:** Custom Resource present. MeshSync Controller is present but the state is not RUNNING or ERRDISABLE, though
- **RUNNING:** MeshSync pod present and in a running state.
- **CONNECTED:** Deployed and connected to Broker.
- **UNDEPLOYED:** Custom Resource not present.

### Meshery Broker Health Status

- **DEPLOYED:** External IP not exposed OR External IP exposed but Meshery Server is not connected as a client to Broker hence data is not being published.
- **UNDEPLOYED:** Custom Resource not deployed.
- **CONNECTED:** Deployed, sending data to Meshery Server.

## Meshery Operator Deployment Scenarios

Because Meshery is versatile in its deployment models, there are different of scenarios in which you may need to troubleshoot the health of Meshery Operator. Identify the deployment model fitting of your environment and follow the guidance under the respective scenario to troublshoot accordingly.

### In-Cluster Deployment

<!-- Meshery Operator, MeshSync, and Broker are deployed in the same cluster as Meshery Server. This is the default deployment scenario when using `mesheryctl system start` or `make run-local`. -->

Whether using [`mesheryctl system start`]({{site.baseurl}}/installation), [`helm install`]({{site.baseurl}}/installation/kubernetes/helm) or `make run-local`, Meshery Server will automatically connect to any available Kubernetes clusters found in your kubeconfig (under `$HOME/.kube/config`). Once connected, operator, broker(NATS) and meshsync will automatically get deployed in the same clusters.

If everything is fine, by viewing the connection in Meshery UI, MeshSync should be in **CONNECTED:** state. Otherwise, check the Operator's pod logs:

`kubectl logs <meshery-operator-pod> -n meshery`

### Out-of-Cluster Deployment

1. Meshery Server is deployed on any Docker host (- Meshery Server is deployed on a Docker host, and Meshery Operator is deployed on a Kubernetes cluster).
   _or_
2. Meshery is managing multiple clusters, some of which are not the cluster unto which Meshery Server is deployed.

## Common Failure Scenarios

Some common failure situations that Meshery users might face are described below.

1. **Situation:** No deployment of Meshery Operator, MeshSync, and Broker.
   1. **Probable cause:** Meshery Server cannot connect to Kubernetes cluster; cluster unreachable or kubeconfig without proper permissions needed to deploy Meshery Operator; Kubernetes config initialization issues.
1. **Situation:** Meshery Operator with MeshSync and Broker deployed, but Meshery Server is not receiving data from MeshSync or data the [Meshery Database]({{site.baseurl}}/concepts/architecture/database) is stale.
   1. **Probable cause:** 
   2. Meshery Server lost subscription to Meshery Broker; Broker server not exposed to external IP; MeshSync not connected to Broker; MeshSync not running; Meshery Database is stale.
   3. The SQL database in Meshery serves as a cache for cluster state. A single button allows users to dump/reset the Meshery Database.
   4. Orphaned MeshSync and Broker controllers - Meshery Operator is not present, but MeshSync and Broker controllers are running.

## Operating Meshery without Meshery Operator

Meshery Operator, MeshSync, and Broker are crucial components in a Meshery deployment. Meshery can function without them, but some functions of Meshery will be disable / unusable. Whether Meshery Operator is initially deployed via `mesheryctl` command or via Meshery Server, you can monitor the health of the Meshery Operator deployment using either the CLI or UI clients.

## Verifying the Status of Meshery Operator, MeshSync, and Meshery Broker

## Troubleshooting using Meshery CLI

The following commands are available to troubleshoot Meshery Operator, MeshSync, and Broker.

**Meshery Server and Adapters**

- `mesheryctl system status` - Displays the status of Meshery Server and Meshery Adapters.

**Meshery Operator, MeshSync, and Broker**

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

### Synthetic Test for Ensuring Change in Cluster State

Initiate a synthetic check to verify a fully functional Operator deployment, testing MeshSync/Broker connectivity.

- Empty database shows the main-cluster node.
- Corrupt database triggers an error snackbar with a link to the Settings screen.
- Disconnected Kubernetes displays MeshSync logo pulsating when data is received.

<div class="section">
Future Enhancements for Troubleshooting:

- NATS/MeshSync not running prompts a review of available operations in the Settings panel.

</div>

## See Also

- [Troubleshooting Meshery Installations](/guides/troubleshooting/installation)
- [Troubleshooting Errors while running Meshery](guides/troubleshooting/meshery-server)

{% include related-discussions.html tag="meshery" %}

