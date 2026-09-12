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
        MeshSync is a custom Kubernetes controller that performs event-driven discovery of, and continual synchronization with, Meshery Server as to the state of managed Kubernetes infrastructure. It operates in one of two modes: operator or embedded. When it runs in operator mode, it is managed by the <a href="{{< ref "concepts/architecture/operator/index.md" >}}">Meshery Operator</a>.
        </p>
    </div>
</div>

### Key Features

- **Greenfield and Brownfield Support**: MeshSync discovers infrastructure whether you're starting from scratch (greenfield) with Meshery performing the initial deployment of your infrastructure or bringing in Meshery to manage your existing infrastructure (brownfield).
- **Real-time / Event-driven Status**: MeshSync uses Kubernetes [shared informers](https://pkg.go.dev/k8s.io/client-go/informers) to subscribe for updates from the cluster and listen for state changes, understanding that changes to managed infrastructure may be done out of band of Meshery.
- **Configurable resource scope**: MeshSync watches a default set of resource types across all namespaces. You can narrow this set with a whitelist or blacklist (see [FAQs](#meshsync-faqs)) so that uninteresting resource types are not discovered.
- **Working cluster snapshot**: MeshSync's working snapshot of the state of each cluster under management is held in-memory (Kubernetes informer caches) and streamed to Meshery Server, whose local database persists it as a continuously refreshed cache.

{{% alert color="warning" title="Roadmap features" %}}
Several capabilities described in MeshSync's original design - <strong>tiered (progressively refined) discovery</strong>, <strong>composite fingerprints that span multiple entities</strong>, <strong>broker-side persistence of events</strong>, and <strong>non-Kubernetes targets</strong> - are design goals that are <strong>not yet implemented</strong>. Sections below that describe roadmap behavior are marked accordingly so that operators can distinguish current behavior from intended direction.
{{% /alert %}}

## Discovery

MeshSync supports both greenfield and brownfield discovery of infrastructure. Greenfield discovery manages infrastructure created and managed entirely by Meshery, while brownfield discovery identifies separately created infrastructure.

### Brownfield: Discovering existing resources

The resources that are present inside the cluster are discovered efficiently with the help of pipelines. The data is constructed in a particular format specific to Meshery and published across to different parts of the architecture. For a task-oriented walkthrough of connecting a cluster that already runs workloads - including the RBAC, network, scale, and sensitive-data considerations - see [Bringing Existing Infrastructure Under Meshery Management]({{< ref "guides/infrastructure-management/managing-existing-infrastructure.md" >}}).

### Greenfield: Tracking newly created resources

Resources created through Meshery - by deploying a [Design]({{< ref "concepts/logical/designs.md" >}}) - are applied to the cluster carrying exactly the metadata their design declares, and MeshSync discovers them through the same watch pipeline as every other resource.

{{% alert color="warning" title="Earmarking: not yet implemented" %}}
Automatically earmarking Meshery-created resources - tagging them with an identifying label or annotation (for example, a design identifier) so that they can be positively attributed to the design that created them - is a design goal, not current behavior. Meshery Server does not currently inject identifying labels or annotations into the resources it deploys.
{{% /alert %}}

## Identifying Infrastructure under Management

MeshSync publishes each discovered Kubernetes resource largely as-is (its metadata, labels, annotations, spec, and status) to Meshery Server. The work of classifying a resource against Meshery's object model, and of positively identifying higher-level infrastructure, is performed by **Meshery Server** using its [model registry]({{< ref "concepts/logical/registry.md" >}}), not by MeshSync itself. The one enrichment MeshSync performs inline is on `Service` resources, where it computes reachable endpoint URLs and flags whether the Service is a candidate to be promoted to a [Meshery Connection]({{< ref "concepts/logical/connections/index.md" >}}).

### Composite Fingerprints (roadmap)

{{% alert color="warning" title="Not yet implemented" %}}
Composite fingerprinting - identifying an application or platform from a set of its constituent objects - is a design goal. MeshSync today discovers and reports objects <strong>individually</strong>; it does not yet correlate an object's constituent parts into a single composite fingerprint.
{{% /alert %}}

Fingerprinting is the intended process of positively identifying and classifying resources using a set of attributes designated as unique to that type of resource. For example:

- Prometheus typically offers metrics on 9090/tcp, but not always.
- Prometheus is typically deployed from a prebuilt container offered by the open source project, but not always.

As a guiding principle, each set of composite fingerprints would use the same identifiers that each element management tool uses to identify itself (e.g., `istioctl version`), assembled via a builder pattern over signals such as container images, CRDs, and Deployments.

## Configuration

{{% alert color="info" title="Configuring MeshSync, the Operator, and the Broker" %}}
This page describes MeshSync's architecture and the mechanics of individual settings. For the complete, task-oriented configuration surface of all three in-cluster components - every setting, its default, and its behavioral impact - see <a href="{{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}}">Configuring Meshery Operator, MeshSync, and Broker</a>.
{{% /alert %}}

### Subscribing to events and changes

MeshSync registers a Kubernetes informer per watched resource type and listens for `ADDED`, `MODIFIED`, and `DELETED` events. On a `MODIFIED` event, MeshSync suppresses no-op updates by comparing the object's `resourceVersion` before republishing. Which resource types are watched, and which event types are published for each, is controlled by the `watch-list` on the `MeshSync` custom resource (see [FAQs](#meshsync-faqs)).

MeshSync also watches the cluster's `CustomResourceDefinition` set: when a CRD is added or removed, MeshSync updates its pipeline and re-runs discovery so that newly installed custom resources begin to be tracked.

### Publish and request subjects

In operator (broker) mode, MeshSync publishes resource events to the NATS subject `meshery.meshsync.core` and answers requests on `meshery.meshsync.request`. It supports request/reply for a full store snapshot (`informer-store`), an on-demand resync, its running version (`meshsync-meta`), pod log streaming (`meshery.meshsync.logs`), and interactive `exec` sessions into pods (`meshery.meshsync.exec`). These primitives back the ad hoc connectivity tests and pod-level troubleshooting available from Meshery clients.

## Health and status

When it starts, MeshSync serves two HTTP endpoints on port `11000`:

- `GET /healthz` - liveness; always returns `200` while the process is alive.
- `GET /readyz` - readiness; returns `200` once MeshSync has established its connection to the Broker (broker mode), and `503` before that.

{{% alert color="info" title="Readiness caveat" %}}
`/readyz` reflects that MeshSync has <strong>connected to the Broker</strong>, not that its informer caches have finished syncing. A <code>200</code> from <code>/readyz</code> therefore means "connected", not necessarily "cluster snapshot fully primed".
{{% /alert %}}

# Scalability and Performance

One Meshery Operator and one MeshSync are deployed to each Kubernetes cluster under management. MeshSync uses a single [`DynamicSharedInformerFactory`](https://pkg.go.dev/k8s.io/client-go/dynamic/dynamicinformer) scoped to all namespaces, registering one informer per watched resource type. Discovery is watch-driven with no periodic resync; MeshSync re-lists a resource type only on start, on an explicit resync request, or when a CRD change triggers a rebuild.

## Tiered Discovery (roadmap)

{{% alert color="warning" title="Not yet implemented" %}}
Tiered discovery is a design goal, not current behavior. MeshSync today performs <strong>flat</strong> discovery: it watches its configured set of resource types across all namespaces at full fidelity. The whitelist/blacklist described in the <a href="#meshsync-faqs">FAQs</a> is the current mechanism for controlling discovery scope; it is a coarse include/exclude of resource types, not a tunable depth.
{{% /alert %}}

Kubernetes clusters may grow very large with numerous objects within them. The intent of tiered discovery is to successively refine infrastructure identification - a cheap first pass to detect what is present, followed by progressively deeper passes only against what was found - so that discovery of large clusters stays fast and cheap. Until then, operators on large or busy clusters should scope discovery with a whitelist and be aware that CRD churn triggers a full re-discovery (see [Troubleshooting]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}})).

## Event-Driven Implementation

Meshery's event-driven approach makes it suitable for managing both small and large clusters. [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) uses NATS as the messaging bus between MeshSync and Meshery Server.

{{% alert color="warning" title="Delivery semantics" %}}
MeshSync publishes to <strong>core NATS</strong> subjects (fire-and-forget, at-most-once). Events are <strong>not</strong> persisted broker-side today, so if Meshery Server is disconnected when an event is published, that event is not replayed on reconnect - the current snapshot is instead rebuilt by re-listing the cluster on the next resync. Durable, replayable delivery (for example via NATS JetStream) is on the roadmap.
{{% /alert %}}

## Broker connection

In operator mode, [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}})
wires MeshSync to the Broker: it derives the Broker's address from the
`Broker` resource's `status.endpoint` and injects it into the MeshSync
Deployment as the `BROKER_URL` environment variable - always a
`nats://host:port` URL. Because the Operator watches the Broker, a change to
the Broker's endpoint (for example after
[reconfiguring its service networking]({{< ref "concepts/architecture/broker/index.md#declarative-service-networking" >}}))
re-reconciles MeshSync so it reconnects to the new address automatically. To
point MeshSync at an externally managed NATS instead, set
`spec.broker.custom.url` on the `MeshSync` resource.

# MeshSync deployment mode

MeshSync operates in one of two modes: operator or embedded.

## Operator mode

When it runs in operator mode, it is managed by the <a href="{{< ref "concepts/architecture/operator/index.md" >}}">Meshery Operator</a>.

## Embedded mode (default)

When it runs in embedded mode, it is integrated into the Meshery server as a library and no additional resources are deployed to the managed cluster. This is the default mode.

## Mode selection and switch

The user selects the deployment mode per context when creating a new Kubernetes connection (submitting a kubeconfig). Connections without an explicit mode use Meshery Server's `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` setting, which itself defaults to `embedded`. The trade-offs between the modes, and the full precedence rules, are covered in [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md#choosing-a-deployment-mode" >}}).

The user can switch the deployment mode per connection on the connections list page.

When the deployment mode is switched from operator to embedded: the operator is undeployed from the managed cluster, and the MeshSync library routine is started inside the Meshery server for the managed cluster.

When the deployment mode is switched from embedded to operator: the MeshSync library routine is stopped for the managed cluster, and the operator is deployed to the managed cluster.

# Output modes: broker and file

Independent of operator/embedded deployment, MeshSync can emit its output in one of two ways, selected with the `--output` flag (default `broker`):

- **`broker`** - the default. MeshSync streams resource events to NATS for Meshery Server to consume. This is how MeshSync runs in a cluster alongside Meshery Broker.
- **`file`** - MeshSync writes a point-in-time cluster snapshot to disk as Kubernetes-style YAML, with no dependency on NATS or the `MeshSync` CRD. This mode backs the [`kubectl meshsync snapshot`]({{< ref "extensions/extensions/kubectl-meshsync-snapshot/index.md" >}}) plugin and is useful for local, air-gapped, or debugging captures. It produces two files: a deduplicated snapshot (`meshery-cluster-snapshot-YYYYMMDD-00.yaml`, one entry per resource by `metadata.uid`) and, optionally, an `-extended` file containing every observed event.

File mode also honors flags to bound the capture: `--outputNamespaces` and `--outputResources` (comma-separated filters, applied to both modes), `--outputFile`, and `--stopAfter` (a duration after which MeshSync exits). These filters are applied to MeshSync's *output*; the informers themselves still watch all namespaces.

# Common tasks

**Check MeshSync health:**

```bash
kubectl -n meshery get meshsync meshery-meshsync -o jsonpath='{.status.conditions}{"\n"}'
kubectl -n meshery rollout status deploy/meshery-meshsync
```

**Probe MeshSync's liveness and readiness endpoints** (served on port `11000`):

```bash
kubectl -n meshery port-forward deploy/meshery-meshsync 11000:11000 &
curl -sS http://127.0.0.1:11000/healthz   # liveness: "ok"
curl -sS -o /dev/null -w '%{http_code}\n' http://127.0.0.1:11000/readyz   # 200 once connected to Broker
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
`MeshSync` **custom resource** (not the CRD schema) - `meshery-meshsync` in
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

Alternatively, a whitelist inverts the behavior - only the listed resources
(and event types) are watched. Only one of whitelist or blacklist may be set:

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

