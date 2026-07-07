---
title: Bringing Existing Infrastructure Under Meshery Management
linkTitle: Managing Existing Infrastructure
categories: [infrastructure]
weight: -4
description: What happens when you point Meshery at a cluster that is already running workloads - how brownfield discovery works, what to evaluate before you connect, the end-to-end process, and the caveats.
---

Most clusters do not arrive empty. They carry months or years of accumulated Deployments, CRDs, operators, and Secrets, deployed by pipelines that predate Meshery and that will keep running after it. This guide is for the platform engineer bringing that kind of estate - a brownfield - under Meshery's management: what connecting a live cluster actually does, what it deliberately does not do, and what to evaluate before you connect.

The short version: greenfield versus brownfield is a non-event for Meshery's discovery model. [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) performs read-only discovery of whatever the cluster already contains. Meshery does not require that resources be created through it, does not modify what it finds, and continues to discover resources created outside of it for as long as the cluster is connected. That last part is not an edge case - it is the operating model.

For the underlying concepts, see the canonical references:

- [Connections]({{< ref "concepts/logical/connections/index.md" >}}) - what a Connection is and its full state lifecycle.
- [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) - the discovery component, its architecture, and its FAQs.
- [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}) - the Connection Wizard mechanics, step by step.
- [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}) - operating Connections after they are registered.

{{% alert color="info" title="What 'management' means at each stage" %}}
Discovery populates Meshery's **inventory**: a continuously refreshed, read-only snapshot of the cluster. Management **actions** - deploying or undeploying a [Design]({{< ref "concepts/logical/designs.md" >}}), transitioning a Connection's state - are separate, explicit, and user-initiated. Connecting a cluster never mutates the workloads already running on it.
{{% /alert %}}

## How Meshery treats an existing cluster

When MeshSync starts against a cluster, its informers first **list** every resource type in scope - delivering each pre-existing object to Meshery as a discovered resource - and then **watch** for subsequent changes. The initial listing is the brownfield snapshot; the watch keeps it current. There is no requirement that a resource carry any particular label or annotation to be discovered.

Meshery does not brand what it deploys, either: resources created by deploying a Meshery [Design]({{< ref "concepts/logical/designs.md" >}}) reach the cluster carrying exactly the metadata the design declares - Meshery Server does not inject identifying labels or annotations at deploy time - and discovery treats them identically to resources that predate Meshery. Your inventory is one uniform, continuously refreshed view of the cluster, whatever created its contents.

## What happens when you connect a live cluster

Connecting an existing cluster is the same flow as connecting any cluster - the [Connection Wizard]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}) or `mesheryctl`. Mechanically, this is what takes place:

1. **You upload a kubeconfig and select contexts.** Meshery parses the file and lists its contexts, indicating which are reachable; nothing is persisted yet. Each context you select becomes its own [Connection]({{< ref "concepts/logical/connections/index.md" >}}).
2. **Each Connection enters the state machine.** New Kubernetes Connections are created in the [Discovered]({{< ref "concepts/logical/connections/index.md#state-discovered" >}}) state. Reachable ones transition automatically through [Registered]({{< ref "concepts/logical/connections/index.md#state-registered" >}}) (Meshery verifies reachability against the cluster's `/livez` endpoint and registers the cluster's component models) to [Connected]({{< ref "concepts/logical/connections/index.md#state-connected" >}}) in the same import. Unreachable contexts remain Discovered until you act on them.
3. **MeshSync starts, in the deployment mode you chose.** Per context, you choose how MeshSync runs (recorded as the `meshsync_deployment_mode` metadata on the Connection):
   - **Embedded** (the default): MeshSync runs as a library inside Meshery Server, authenticating with the credentials in the kubeconfig you uploaded. Nothing is installed into the cluster - no Operator, no Broker, no CRDs.
   - **Operator**: Meshery Server installs the [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}) into the cluster (a Helm chart install into the `meshery` namespace). The operator in turn deploys [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) (NATS) and MeshSync, and injects the Broker address into MeshSync as `BROKER_URL`. Administrators can also disable operator deployment server-wide with the `DISABLE_OPERATOR` environment variable on Meshery Server.
4. **Initial discovery runs, then continuous watch.** Every resource type in MeshSync's watch list is listed, and every existing object in scope lands in Meshery's inventory. From then on, changes stream in event by event - including changes made by your existing pipelines, entirely outside Meshery.
5. **The inventory becomes visible.** Discovered resources appear in Meshery's dashboard and resource views, and are queryable at `GET /api/system/meshsync/resources`.

## The Connection lifecycle for a brownfield estate

Meshery models a Connection's lifecycle with explicit states - Discovered, Registered, Connected, Ignored, Maintenance, Disconnected, Deleted, and Not Found - documented in [States and the Lifecycle of Connections]({{< ref "concepts/logical/connections/index.md#states-and-the-lifecycle-of-connections" >}}). For an estate you already operate, the practical readings are:

- **Discovered** - Meshery knows the cluster exists. Nothing has been verified, nothing deployed, nothing collected.
- **Registered** - reachability and usability are verified; Meshery holds the Connection for your administrative decision. Still purely observational.
- **Connected** - the only state in which Meshery does anything to the cluster, and what it does is bounded: in operator mode it installs its own components (Operator, Broker, MeshSync); in embedded mode it installs nothing. Discovery begins. Your workloads are read, not touched.
- **Ignored** - you keep the cluster in Meshery's field of view but explicitly opt it out of management; Meshery stops re-discovering it.
- **Disconnected** - Meshery lost the ability to communicate with the cluster (network, expired credentials, or the cluster went away). Previously collected inventory is retained but goes stale.
- **Deleted** - the Connection and its collected inventory are removed from Meshery. The cluster itself is unaffected.

Two honest footnotes, verified against Meshery Server's Kubernetes state machine: the **Maintenance** state, while defined for Connections generally, is not currently part of the implemented Kubernetes transition set; and **Not Found** is where a Connection lands when registration cannot verify the cluster (you can delete it or re-register once the cluster is reachable).

## Before you connect: an evaluation checklist

### RBAC: what Meshery's components can actually do

MeshSync's code performs no writes against your cluster - it lists and watches. Its effective *reach*, and the enforcement of that read-only posture, depend on the deployment mode:

- **Embedded mode** runs discovery under the identity in the kubeconfig you upload. Its RBAC is that identity's RBAC - no more, no less. For a locked-down cluster, this is the least-privilege path: create a dedicated read-only user or ServiceAccount context scoped to exactly the namespaces and resource types you want discovered, and upload that. If the identity cannot read a resource type, that type is simply not discovered.
- **Operator mode** installs the meshery-operator Helm chart, whose ClusterRole is deliberately narrow - no wildcards. It grants create/delete/get/list/patch/update/watch on core `configmaps`, `secrets`, and `services`, and on `deployments` and `statefulsets` (apps) - the resources the operator manages to run Broker and MeshSync - plus full control of its own `brokers.meshery.io` and `meshsyncs.meshery.io` custom resources and `create` on `tokenreviews` and `subjectaccessreviews`.

  Two consequences follow. First, MeshSync's Deployment runs under the same `meshery-operator` ServiceAccount - there is no separate, read-only MeshSync role - so MeshSync's read-only behavior is a property of its code, not of a dedicated RBAC boundary. If you require RBAC-enforced read-only discovery, use embedded mode with a read-only identity. Second, discovery in operator mode is bounded by that same role: resource types it does not authorize (Pods, Nodes, Namespaces, CRDs, and much of MeshSync's default watch list) are not readable, and therefore not discovered, until a cluster administrator explicitly grants additional `get`/`list`/`watch` permissions to the `meshery-operator` ServiceAccount. Nothing self-escalates; widening discovery is a deliberate act.
- **In-cluster Meshery Server** (the `meshery` Helm chart) is the opposite end of the spectrum: its ServiceAccount is bound to a ClusterRole granting `*` on `*`. If Meshery Server itself runs in a cluster, discovery of that local cluster is effectively unrestricted, and the server holds write access it uses for deploying designs. Factor this into which cluster hosts your Meshery Server.

### Network reachability

- **Meshery Server must reach the cluster's API server.** Registration verifies reachability (a `/livez` probe); if the API endpoint in your kubeconfig is private, the Connection will not progress past Discovered.
- **In operator mode, Meshery Server must also reach the Broker.** The Broker's client Service defaults to `ClusterIP` - reachable only from inside the cluster. That works when Meshery Server runs in the same cluster; an out-of-cluster server (your laptop, another cluster, a SaaS deployment) needs the Broker exposed via the `Broker` custom resource's `spec.service.type` (`NodePort` or `LoadBalancer`) or an explicit `spec.service.externalEndpointOverride` for ingress and NAT topologies. The operator publishes the resolved address in the Broker's `status.endpoint`.
- **Embedded mode sidesteps the Broker entirely** - discovery traffic flows through Meshery Server's existing kubeconfig connection. On restrictive networks, this is the simplest topology.

For the specific service-networking and endpoint knobs, see [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md#broker-service-networking" >}}).

### Scale: the initial snapshot on a large estate

The first sync lists every object of every watched resource type. On a large brownfield cluster, this causes a burst of load on the API server, a burst of traffic to the Broker (operator mode), and a burst of writes into Meshery's database. There are no rate or depth tuning knobs in MeshSync today; **scoping is the lever**:

- Narrow *what is watched* with the whitelist or blacklist in the `MeshSync` custom resource's `watch-list` (see the [MeshSync FAQs]({{< ref "concepts/architecture/meshsync.md#meshsync-faqs" >}}) for exact syntax).
- Narrow *what is published* with MeshSync's `--outputNamespaces` and `--outputResources` filters.
- Expect re-discovery on CRD churn: installing or removing a CRD triggers MeshSync to rebuild its pipeline and re-list. Estates with heavy CRD turnover resync more often.

The full set of scoping and deployment knobs, and where each is configured, is covered in [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md#configuring-meshsync" >}}).

### Sensitive data: Secrets are discovered by default

MeshSync's default watch list includes `secrets.v1.`, and Secret **values are published to Meshery unredacted by default**. On a brownfield cluster this means credentials that predate Meshery flow into its inventory the moment discovery starts. Decide your posture before connecting, not after:

- **Redact**: set `MESHSYNC_REDACT_SECRETS=true` on the MeshSync process to replace every value in a discovered Secret's `data` with `[REDACTED]` while preserving the keys - the inventory still shows that Secrets exist and how they are shaped, without their contents.
- **Exclude**: blacklist `secrets.v1.` in the `watch-list` so Secrets are not discovered at all.
- **Deny by RBAC**: in embedded mode, upload an identity that has no `get`/`list`/`watch` on Secrets; in operator mode, do not grant Secret read access beyond what the chart requires. What cannot be read cannot be published.

### Multi-cluster and multi-user realities

- Each kubeconfig context becomes its own Connection, and in operator mode each connected cluster gets its own Operator, Broker, and MeshSync. Plan RBAC and network exposure per cluster, not once. See [Multi-Cluster and Multi-Cloud]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}}).
- A Kubernetes Connection is owned by the user who imported it and is private until explicitly shared - by assigning it to an [Environment]({{< ref "concepts/logical/environments.md" >}}) and that Environment to a [Workspace]({{< ref "concepts/logical/workspaces.md" >}}). On a team, decide up front who imports production clusters and into which Environments they land.
- Deployment mode is per Connection and switchable later: changing it tears down the old mode and brings up the new one (undeploying the operator when moving to embedded, and vice versa).

## Walkthrough: bringing a live cluster under management

1. **(Optional) Rehearse discovery offline.** MeshSync's file output mode captures a point-in-time snapshot of a cluster to local YAML with nothing installed and no Broker involved - a zero-commitment preview of exactly what discovery would collect. The [kubectl meshsync snapshot]({{< ref "extensions/extensions/kubectl-meshsync-snapshot/index.md" >}}) plugin packages this.
2. **Connect the cluster.** In the Meshery UI (**Lifecycle → Connections → Create Connection → Kubernetes**), upload your kubeconfig, select the contexts to import, and choose each context's MeshSync deployment mode - embedded if you are not ready to install anything into a production cluster, operator for in-cluster, event-streamed discovery. Full steps: [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}). For managed clouds, `mesheryctl system config aks|eks|gke|minikube` uploads the context for you.
3. **Verify component health.** Confirm the Connection shows **Connected** in the Connections table. In operator mode, verify the in-cluster components:

   ```bash
   mesheryctl system check --operator
   ```

   This checks the Meshery Operator, Broker, and MeshSync deployments and the Broker's advertised endpoint. For deeper, kubectl-level checks (MeshSync's `/healthz` and `/readyz` probes, `BROKER_URL` wiring), see [MeshSync's common tasks]({{< ref "concepts/architecture/meshsync.md#common-tasks" >}}) and the [troubleshooting guide]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}}).
4. **Watch the inventory populate.** Your estate appears in Meshery's dashboard and resource views as the initial listing completes - existing resources first, then live updates as your clusters change. On a large estate, give the first sync time to finish before judging completeness. From the CLI:

   ```bash
   mesheryctl connection list --kind kubernetes
   ```
5. **Scope discovery to what you care about.** If the inventory is noisier than useful - or the initial sync is heavier than you want - narrow the watch list and output filters (see the [checklist above](#scale-the-initial-snapshot-on-a-large-estate)) and restart MeshSync to apply.
6. **Organize and share.** Assign the Connection to an [Environment]({{< ref "concepts/logical/environments.md" >}}), and the Environment to a [Workspace]({{< ref "concepts/logical/workspaces.md" >}}), to make the cluster available to your team with scoped access.
7. **Operate deliberately.** Discovery has populated the inventory; management actions remain yours to initiate:
   - Meshery can render the discovered state of a cluster as a [Design]({{< ref "concepts/logical/designs.md" >}}) on the fly (the resource API supports `?asDesign=true`) for visualization and evaluation. There is no import source that persists a design *from* live cluster state today - designs are imported from Kubernetes manifests, Helm charts, Docker Compose files, or existing design files. If your estate's manifests live in Git, import those directly: `mesheryctl design import -f ./manifests.yaml -s "Kubernetes Manifest"` (see [design import]({{< ref "reference/references/mesheryctl/design/import.md" >}})).
   - Deploying a design creates or updates the resources that the design declares; undeploying removes them. These actions apply only to what the design declares - they never extend to the rest of the discovered inventory.

## Caveats and limitations

- **Discovery is read-only, with one active-operation exception.** MeshSync itself never creates, updates, or deletes cluster resources. Separately from discovery, Meshery's pod-level troubleshooting features (interactive `exec` sessions and log streaming, available in operator mode via the Broker) are active operations inside your workloads - user-initiated, on demand, and only possible if MeshSync's ServiceAccount is granted `create` on `pods/exec` and `get` on `pods/log`, which the current chart RBAC does not include by default.
- **Inventory freshness is best-effort between syncs.** MeshSync publishes over core NATS (at-most-once delivery). If Meshery Server is disconnected when an event fires, that event is not replayed; the snapshot is rebuilt on the next resync or re-list. Treat the inventory as an operational cache, not an audit log.
- **Disconnecting is non-destructive; deleting flushes.** Disconnecting a cluster (or undeploying Meshery's components) leaves every workload untouched and retains the collected inventory - which then goes stale. Deleting the Connection removes the collected inventory from Meshery, though only when it is the last Connection referencing that cluster (several kubeconfig contexts can point at the same cluster). Either way, the cluster itself is unaffected.
- **Drift is surfaced, not prevented.** Resources changed or created outside Meshery keep flowing into the inventory - that is the design. A design you deployed reflects what you declared; the inventory reflects what is observed. Meshery does not block out-of-band changes, and a design does not "adopt" the live resources it happens to describe.
- **Operator-mode discovery breadth is RBAC-bounded.** As noted [above](#rbac-what-mesherys-components-can-actually-do), the current meshery-operator chart authorizes reads on fewer resource types than MeshSync's default watch list requests; resource types without granted read access are not discovered until you extend the ServiceAccount's permissions.
- **Watch-list changes require a MeshSync restart**, and the Maintenance connection state is not currently implemented for Kubernetes Connections (both noted in the sections above).

## Related

- [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}}) - every deployment and discovery-scoping knob in one place.
- [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}) - the Connection Wizard, step by step.
- [Connections]({{< ref "concepts/logical/connections/index.md" >}}) - concepts and the state lifecycle.
- [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}) - lifecycle operations on registered Connections.
- [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}), [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}), and [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) - component architecture.
- [Multi-Cluster and Multi-Cloud]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}}) - running Meshery across many clusters.

{{< discuss >}}
