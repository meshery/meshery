---
layout: default
title: MeshSync
permalink: concepts/architecture/meshsync
type: components
redirect_from: architecture/meshsync
abstract: "Meshery offers support for Kubernetes cluster and cloud state synchronization with the help of MeshSync."
language: en
display-title: "false"
list: include
---

<div style="display:flex;">
    <div style="flex: 1; margin-right: 2rem;">
        <img src="{{site.baseurl}}/assets/img/meshsync/meshsync.svg"
        style="margin-right:2rem;margin-bottom:2rem;max-width:30%;" />
    </div>
    <div style="flex: 1;">
        <h1>MeshSync</h1>
        <p>
        Managed by the <a href="{{site.baseurl}}/concepts/architecture/operator">Meshery Operator</a>, MeshSync is a custom Kubernetes controller that provides tiered discovery and continual synchronization with Meshery Server as to the state of the Kubernetes clusters and their workloads.
        </p>
    </div>
</div>


### Key Features

- **Greenfield and Brownfield Support**: MeshSync discovers and identifies infrastructure whether you're starting from scratch (greenfield) with Meshery performing the initial deployment of your infrastructure or bringing in Meshery to manage your existing infrastructure (brownfield).
- **Real-time / Event-driven Status**: MeshSync subscribes for updates from your platform and listens for state changes, understanding that changes to managed infrastructure may be done out of band of Meshery.
- **Scalable and Performant**: Meshery's event-driven approach ensures speed and scalability. You have control over the depth of object discovery to manage large clusters efficiently. MeshSync's working snapshot of the state of each cluster under management is stored in the Server's local database and continuously refreshed.

## Discovery
MeshSync supports both greenfield and brownfield discovery of infrastructure. Greenfield discovery manages infrastructure created and managed entirely by Meshery, while brownfield discovery identifies separately created infrastructure.

### Brownfield: Discovering existing resources

The resources that are present inside the cluster are discovered efficiently with the help of pipelines. The data is constructed in a particular format specific to Meshery and published across to different parts of the architecture.

### Greenfield: Tracking newly created resources

 Meshery earmarks infrastucture for which it is the orginal lifecycle manager. In other words, Meshery tags the resources it creates. In Kubernetes deployments, earmarking is performed using annotations, notably the key/value pair:

`designs.meshery.io: <design-id>`

The propagation of the labels and annotations to the native k8s resources would be the responsibility of the workload/trait implementor.
The following annotations are added to resources that are created by Meshery Server.


```yaml
Labels:
 - resource.pattern.meshery.io/id=<uuid> # unique identifier for the design
```

## Identifying Infrastructure under Management

Supplied by Meshery Server, MeshSync uses composite fingerprints to uniquely and positively identify managed infrastructure, capturing essential characteristics like versions and configurations.

### Composite Fingerprints

Fingerprinting, the process of positively identifying and classifying resources, is performed using a set of pre-defined attributes that have been designated as unique to that type of resource. For example:

 - Prometheus typically offers metrics on 9090/tcp, but not always.
 - Prometheus is typically deployed from a prebuilt container offered by the open source project, but not always.

See Connnection State Management for additional information.

Fingerprinting a service mesh is the act of uniquely identifying managed infrastructure, their versions and other specific characteristics.

As a guiding principal, each set of composite fingerprints uses the same identifiers that each element management tool uses to identify itself (e.g., istioctl version).

Creating a composite set of keys involves using the builder pattern. For example:

- Images
- CRDs
- Deployment


## Configuration

### Subscribing to events/changes on every component

The informer in MeshSync actively listens to changes in resources and updates them in real time based on the informer configuration in the CRD.

## Subscription Status and Health


### Flushing MeshSync

### Sythentic Test of MeshSync

*TODO: Include example of how to invoke this built-in check.*
# Scalability and Performance

One Meshery Operator and one MeshSync are deployed to each Kuberentes cluster under management.

## Tiered Discovery

Kubernetes clusters may grow very large with thousands of objects on them. The process of positively identifying and classifying resources by type, aligning them with Meshery's object model can be intense. Discovery tiers (for speed and scalability of MeshSync) successively refine the process of  infrasturcture identification (see [Composite Prints](#composite-fingerprints)). 

For efficient management of large Kubernetes clusters, MeshSync uses tiered discovery. This approach progressively refines the identification of relevant infrastructure, optimizing the speed and scalability of MeshSync. You have control over the depth of object discovery, enabling you to strike the right balance between granularity and performance for efficient cluster management.
## Event-Driven Implementation

Meshery's event-driven approach ensures high-speed operations, making it suitable for managing both small and large clusters. [Meshery Broker](./broker) uses NATS as the messaging bus to ensure continuous communication between MeshSync and Meshery Server. In case of connectivity interruptions, MeshSync data is persisted in NATS topics.

# MeshSync FAQs

## How to configure MeshSync's resource discovery behavior: Can specific, "uninteresting" resources be blacklisted?  

MeshSync is managed by [Meshery Operator]({{site.baseurl}}/concepts/architecture/operator), which watches for changes on the `meshsync` CRD for changes and updates the deployed MeshSync instance accordingly. You can blacklist specific Kubernetes resources from being discovered and watched by MeshSync. In order to identify the list of one or more resources for MeshSync to ignore, update the `meshsync` CRD using kubectl:

- Download the CRD with `kubectl get crd meshsyncs.meshery.layer5.io -o yaml > meshsync.yaml`
- Open the downloaded file and edit the field `informer_config` to blacklist all the types of resources that you don't want updates from.
- Apply the new definition with `kubectl apply -f meshsync.yaml`

# Roadmap

## Non-Kubernetes Deployments

Even if you're not using Kubernetes, Meshery empowers you to manage your infrastructure efficiently, providing a unified solution for different deployment environments.

# Recap

MeshSync maintains an up-to-date snapshot of your cluster, ensuring you always have an accurate view of your infrastructure. This snapshot is refreshed in real-time through event-based updates. Whether you're starting fresh or adopting Meshery into existing setups, MeshSync supports both greenfield and brownfield discovery of your environment.

# Suggested Reading

{% include suggested-reading.html diffName="true" isDiffTag="true" diffTag=tag %}
{% include related-discussions.html tag="mesheryctl" %}