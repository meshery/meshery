---
layout: default
title: Meshery Operator, MeshSync, Broker Troubleshooting Guide
permalink: guides/troubleshooting/meshery-operator-meshsync
language: en
abstract: This documentation provides comprehensive guidance on troubleshooting in Meshery Operator, MeshSync and Broker, ensuring you can address common issues efficiently.
type: guides
category: troubleshooting
---

{% include alert.html type="dark" title="Meshery Error Code Reference" content="Have specific error with an error code? See the <a href='/reference/error-codes'>Meshery Error Code Reference</a> for probable cause and suggested remediations." %}

There are common issues Meshery users may face while operating the [Meshery Operator]({{site.baseurl}}/concepts/architecture/operator/) and its custom controllers, [MeshSync]({{site.baseurl}}/concepts/architecture/meshsync) and [Broker]({{site.baseurl}}/concepts/architecture/broker), that can be resolved by performing specific actions. This documentation aims to empower users by providing a set of troubleshooting tools and actions.

## Understanding the Status of Meshery Operator, MeshSync, and Meshery Broker

The following table describes the various states of MeshSync and Meshery Broker and their implications.

**MeshSync:**

- **ENABLED:** Custom Resource present. MeshSync Controller is not connected to Broker.
- **DEPLOYED:** Custom Resource present. MeshSync Controller is present but the state is not RUNNING or ERRDISABLE, though
- **RUNNING:** MeshSync pod present and in a running state.
- **CONNECTED:** Deployed and connected to Broker.
- **UNDEPLOYED:** Custom Resource not present.

**Meshery Broker:**

- **DEPLOYED:** External IP not exposed OR External IP exposed but Meshery Server is not connected as a client to Broker hence data is not being published.

- **UNDEPLOYED:** Custom Resource not deployed.
- **CONNECTED:** Deployed, sending data to Meshery Server.

## Meshery Operator Deployment Scenarios

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
    1. Probable cause: Meshery Server lost subscription to Meshery Broker; Broker server not exposed to external IP; MeshSync not connected to Broker; MeshSync not running; Meshery Database is stale.
    2. The SQL database in Meshery serves as a cache for cluster state. A single button allows users to dump/reset the Meshery Database.
1. Orphaned MeshSync and Broker controllers - Meshery Operator is not present, but MeshSync and Broker controllers are running.

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

This documentation provides comprehensive guidance on troubleshooting in Meshery, ensuring users can address common issues efficiently.

{% if page.suggested-reading != false and page.title and page.type and page.category and page.url %}
{% include_cached suggested-reading.html  title=page.title type=page.type category=page.category url=page.url language="en" %}
{% endif %}

{% include related-discussions.html tag="meshery" %}