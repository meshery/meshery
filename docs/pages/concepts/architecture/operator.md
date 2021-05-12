---
layout: default
title: Operator
permalink: concepts/architecture/operator
type: concepts
redirect_from: architecture/operator
abstract: "Meshery Operator controls and manages the lifeycle of components deployed inside a kubernetes cluster"
language: en
list: include
---

Meshery operator is a kubernetes controller manager which is conceptually known as a kubernetes operator. This manages the lifecycle of every meshery component that is deployed or running inside the cluster.

## List of Controllers that is managed by the operator

### Broker Controller

Meshery broker is one of the core components of the meshery architecture. This controller manages the lifecycle of broker that meshery uses for data streaming across the cluster and the outside world.

See the "[Meshery Broker]({{site.baseurl}}/architecture/broker)" for more information.

### MeshSync Controller

MeshSync Controller manages the lifecycle of MeshSync that is deployed for resource synchronization for the cluster.

See the "[MeshSync]({{site.baseurl}}/architecture/meshsync)" for more information.

### Operator FAQs

#### When is Meshery Operator deployed and when is it deleted?  
As a Kubernetes custom controller, Meshery Operator is provisioned and deprovisioned when Meshery Server is connected to or disconnected from Kubernetes cluster. Meshery Server connections to Kubernetes clusters are controlled using Meshery Server clients: `mesheryctl` or Meshery UI.  This behavior described below is consistent whether your Meshery deployment is using Docker or Kubernetes as the platform to host the Meshery deployment.

**Meshery CLI**
`mesheryctl` initiates connection to Kubernetes cluster when `mesheryctl system start` is executed and disconnects when `mesheryctl system stop` is executed. This behavior is consistent whether your Meshery deployment is using Docker or Kubernetes as the platform to host the Meshery deployment.

**Meshery UI**
Meshery UI offers more granular control over the deployment of Meshery Operator in that you can remove Meshery Operator from a Kubernetes cluster without disconnecting Meshery Server from the Kubernetes cluster. You can control the deployment of Meshery Operator using the on/off switch found in the Meshery Operator section of  Settings.
