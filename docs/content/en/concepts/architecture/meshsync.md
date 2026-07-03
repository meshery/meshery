---
title: MeshSync
description: "MeshSync ensures Meshery Server is continuously in-sync with the state of infrastructure under management."
display_title: false
aliases:
- /architecture/meshsync/
---

<div style="display:flex;align-items:center;">
    <div style="flex: 4;">
        <h1>MeshSync</h1>
        <p>
        MeshSync is a custom Kubernetes controller that provides tiered discovery and continual synchronization with Meshery Server as to the state of managed multi-cloud and cloud native infrastructure. It operates in one of two modes: operator or embedded. When it runs in operator mode, it is managed by the <a href="{{< ref "concepts/architecture/operator/index.md" >}}">Meshery Operator</a>.
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

Meshery earmarks infrastructure for which it is the original lifecycle manager. In other words, Meshery tags the resources it creates. In Kubernetes deployments, earmarking is performed using annotations, notably the key/value pair:

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

See Connection State Management for additional information.

Fingerprinting is the act of uniquely identifying managed infrastructure, their versions and other specific characteristics.

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

### Synthetic Test of MeshSync

_TODO: Include example of how to invoke this built-in check._

# Scalability and Performance

One Meshery Operator and one MeshSync are deployed to each Kubernetes cluster under management.

## Tiered Discovery

Kubernetes clusters may grow very large with numerous objects within them. The process of positively identifying and classifying resources by type and aligning them with Meshery's object model can be intense, and discovery tiers (for speed and scalability of MeshSync) successively refine the process of infrastructure identification (see [Composite Fingerprints](#composite-fingerprints)).

For efficient management of large Kubernetes clusters, MeshSync uses tiered discovery. This approach progressively refines the identification of relevant infrastructure, optimizing the speed and scalability of MeshSync. You have control over the depth of object discovery, enabling you to strike the right balance between granularity and performance for efficient cluster management.

## Event-Driven Implementation

Meshery's event-driven approach ensures high-speed operations, making it suitable for managing both small and large clusters. [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) uses NATS as the messaging bus to ensure continuous communication between MeshSync and Meshery Server. In case of connectivity interruptions, MeshSync data is persisted in NATS topics.

## Broker connection

In operator mode, [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}})
wires MeshSync to the Broker: it derives the Broker's address from the
`Broker` resource's `status.endpoint` and injects it into the MeshSync
Deployment as the `BROKER_URL` environment variable — always a
`nats://host:port` URL. Because the Operator watches the Broker, a change to
the Broker's endpoint (for example after
[reconfiguring its service networking]({{< ref "concepts/architecture/broker/index.md#declarative-service-networking" >}}))
re-reconciles MeshSync so it reconnects to the new address automatically. To
point MeshSync at an externally managed NATS instead, set
`spec.broker.custom.url` on the `MeshSync` resource.

# MeshSync deployment mode

MeshSync operates in one of two modes: operator or embedded.

## Operator mode (default)

When it runs in operator mode, it is managed by the <a href="{{< ref "concepts/architecture/operator/index.md" >}}">Meshery Operator</a>.

## Embedded mode

When it runs in embedded mode, it is integrated into the Meshery server as a library and no additional resources are deployed to the managed cluster.

## Mode selection and switch

The user selects the deployment mode when creating a new Kubernetes connection (submitting a kube config). The selection is applied to all contexts from the submitted config.

The user can switch the deployment mode per connection on the connections list page.

When the deployment mode is switched from operator to embedded: the operator is undeployed from the managed cluster, and the MeshSync library routine is started inside the Meshery server for the managed cluster.

When the deployment mode is switched from embedded to operator: the MeshSync library routine is stopped for the managed cluster, and the operator is deployed to the managed cluster.

# Common tasks

**Check MeshSync health:**

```bash
kubectl -n meshery get meshsync meshery-meshsync -o jsonpath='{.status.conditions}{"\n"}'
kubectl -n meshery rollout status deploy/meshery-meshsync
```

**Verify which Broker MeshSync is connected to:**

```bash
kubectl -n meshery get deploy meshery-meshsync \
  -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="BROKER_URL")].value}{"\n"}'
```

**Trigger a fresh discovery** (MeshSync re-discovers on start):

```bash
kubectl -n meshery rollout restart deploy/meshery-meshsync
```

# MeshSync FAQs

## How to configure MeshSync's resource discovery behavior: Can specific, "uninteresting" resources be blacklisted?

Yes. MeshSync reads its discovery filter from the `watch-list` section of the
`MeshSync` **custom resource** (not the CRD schema) — `meshery-meshsync` in
the `meshery` namespace. The `whitelist` and `blacklist` keys each hold a
JSON-encoded list; resources are identified as `<plural>.<version>.<group>`
(core-group resources end with a trailing dot, e.g. `pods.v1.`).

To ignore specific resource types, edit the custom resource:

```bash
kubectl -n meshery edit meshsync meshery-meshsync
```

```yaml
spec:
  watch-list:
    data:
      blacklist: '["events.v1.","replicasets.v1.apps"]'
```

Alternatively, a whitelist inverts the behavior — only the listed resources
(and event types) are watched:

```yaml
spec:
  watch-list:
    data:
      whitelist: '[{"Resource":"namespaces.v1.","Events":["ADDED","MODIFIED","DELETED"]},{"Resource":"pods.v1.","Events":["MODIFIED"]}]'
```

Restart MeshSync (`kubectl -n meshery rollout restart deploy/meshery-meshsync`)
for the new filter to take effect.


{{% alert color="info" title="Still seeing issues?" %}}
Check the [**Meshery Troubleshooting Guide**]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}}) for help with common issues.
{{% /alert %}}


# Roadmap

## Non-Kubernetes Deployments

Even if you're not using Kubernetes, Meshery empowers you to manage your infrastructure efficiently, providing a unified solution for different deployment environments.

# Recap

MeshSync maintains an up-to-date snapshot of your cluster, ensuring you always have an accurate view of your infrastructure. This snapshot is refreshed in real-time through event-based updates. Whether you're starting fresh or adopting Meshery into existing setups, MeshSync supports both greenfield and brownfield discovery of your environment.

