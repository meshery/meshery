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