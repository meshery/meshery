---
title: Registering a Connection
aliases:
- /guides/infrastructure-management/connection-wizard
- /guides/infrastructure-management/infrastructure-discovery
- /connection-wizard
categories: [infrastructure]
weight: -5
description: Use the Connection Wizard to create and update Connections - Kubernetes clusters, Grafana, Prometheus, and more - in your Meshery deployment.
---

A [Connection]({{< ref "concepts/logical/connections/index.md" >}}) is how Meshery tracks and manages a resource - a Kubernetes cluster, a Grafana instance, a Prometheus server, and [many more]({{< ref "extensions/models/_index.md" >}}). The **Connection Wizard** is the guided, in-UI way to register a new Connection or reconfigure an existing one, without hand-editing YAML or memorizing API payloads.

This guide covers creating and updating Connections with the wizard. For what a Connection _is_, the states it moves through, and how it is managed over time, see the canonical references:

- [Connections]({{< ref "concepts/logical/connections/index.md" >}}) - what Connections are and their full state lifecycle.
- [Credentials]({{< ref "concepts/logical/credentials.md" >}}) - how Meshery authenticates to a Connection.
- [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}) - operating Connections after they are registered.

{{% alert color="info" title="Discovered vs. manually registered Connections" %}}
Meshery learns about Connections two ways. **Managed** Connections (for example, the resources inside a Kubernetes cluster) are auto-discovered by [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) and arrive already in the [Discovered]({{< ref "concepts/logical/connections/index.md#state-discovered" >}}) state. **Unmanaged** Connections (for example, a standalone Grafana or Prometheus) are added by you. The Connection Wizard is how you do the latter - and how you bring a Kubernetes cluster under management in the first place.
{{% /alert %}}

## Before you begin

- A running Meshery deployment. See the [Quick Start]({{< ref "installation/quick-start/index.md" >}}) if you do not have one yet.
- Permission to add Connections. The wizard is permission-gated: adding a Kubernetes cluster requires the **add cluster** permission, and other Connection kinds require the **connect metrics** permission. If you lack both, the **Create Connection** button is disabled. See [Roles and Permissions]({{< ref "reference/extensibility/authorization/index.md" >}}).
- For an authenticated Connection (most Grafana/Prometheus instances), the access token or credential you intend to use. You can paste it during the wizard or reuse an existing [Credential]({{< ref "concepts/logical/credentials.md" >}}).

## Launching the Connection Wizard

1. Open the **Connections** page in Meshery (**Lifecycle → Connections**).
2. Click **Create Connection**.

The wizard opens as a modal. The set of Connection kinds you can create is driven by the [connection definitions]({{< ref "project/contributing/models/connections" >}}) registered in your Meshery Server's [Registry]({{< ref "concepts/logical/registry.md" >}}). Out of the box this includes **Kubernetes**, **Grafana**, and **Prometheus**; your deployment may offer more. If a kind you need is missing, a contributor can add it - see [Contributing a Connection]({{< ref "project/contributing/models/connections" >}}).

## Creating a Connection

Most Connections follow the same generic flow. Each step is rendered from the connection definition itself, so the exact fields you see depend on the kind you choose.

1. **Choose Connection.** Pick the kind of Connection to create (for example, Grafana). Kinds you do not have permission to add are shown but cannot be selected.
2. **Configure Connection.** Fill in the Connection's details - typically the endpoint URL and an optional friendly name. Required fields are validated before you can continue. For a Grafana Connection, for instance, you supply the Grafana endpoint (e.g. `http://grafana.example:3000`).
3. **Associate Credential.** Provide the secret Meshery will use to authenticate. You can either:
   - **Reuse an existing credential** - the list is filtered to credentials that match the Connection's kind, or
   - **Create a new credential** - enter the token, API key, or `username:password` and give it a name (it defaults to the Connection's name).

   You may also choose to **skip credential verification**, which registers the Connection without first probing reachability - useful when the target is not reachable yet but you still want it on record. This step is omitted entirely for kinds that do not define a credential (and for Kubernetes, whose kubeconfig _is_ its credential - see below).
4. **Review & Create.** Confirm the summary and click **Create Connection**. Meshery registers the Connection and immediately attempts to connect to it.
5. **Done.** On success, the Connection becomes a first-class resource, listed in the [Connections]({{< ref "concepts/logical/connections/index.md" >}}) table and ready to use.

{{% alert color="info" title="What 'Create' actually does" %}}
Creating a Connection performs two transitions in sequence: it **registers** the Connection (recording it and its credential) and then **connects** to it (verifying reachability and beginning management). A reachable Connection lands in the [Connected]({{< ref "concepts/logical/connections/index.md#state-connected" >}}) state; if Meshery cannot reach it - or you skipped verification - it remains [Registered]({{< ref "concepts/logical/connections/index.md#state-registered" >}}). You can drive further transitions later from the Connections table.
{{% /alert %}}

### Credentials

Credentials entered in the wizard are persisted as first-class, named [Credentials]({{< ref "concepts/logical/credentials.md" >}}), encrypted at rest, and reusable across other Connections. Meshery never exposes them in logs or API responses. To learn how Meshery interprets a credential's secret (Basic auth vs. bearer token vs. anonymous), see [Credentials]({{< ref "concepts/logical/credentials.md" >}}) and the [Telemetry authentication note](https://docs.meshery.io/guides/telemetry/).

## Importing a Kubernetes cluster

Kubernetes uses a dedicated flow because a single kubeconfig can describe many clusters and its kubeconfig also serves as its credential.

1. **Choose Connection** → **Kubernetes**.
2. **Import Kubeconfig.** Upload a kubeconfig file. Meshery parses it and lists the contexts it contains, indicating which are reachable - nothing is persisted yet.
3. **Select contexts.** Choose which contexts to import. For each, you can override the Connection name and choose a [MeshSync deployment mode](#meshsync-deployment-mode).
4. **Review Import.** Confirm your selection and import. Meshery creates one Connection per selected context and reports the outcome, grouped into connected, registered, ignored, and errored buckets.

Each imported cluster is created as a Kubernetes Connection that [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) keeps in sync. From there, Meshery can deploy and operate workloads, visualize the cluster, and more. See [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}). If the cluster already runs workloads you did not deploy through Meshery, see [Bringing Existing Infrastructure Under Meshery Management]({{< ref "guides/infrastructure-management/managing-existing-infrastructure.md" >}}) for what to evaluate before importing and how discovery treats pre-existing resources.

{{% alert color="dark" title="Who can access an imported cluster?" %}}
A Kubernetes Connection is owned by the user who imported it and is private until you explicitly share it - by assigning it to an [Environment]({{< ref "concepts/logical/environments.md" >}}) and that environment to a [Workspace]({{< ref "concepts/logical/workspaces.md" >}}). See the sharing FAQ under [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}).
{{% /alert %}}

### MeshSync deployment mode

When you import or reconfigure a Kubernetes cluster, you choose how [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) - the component that keeps Meshery's view of the cluster's resources up to date - runs:

- **Operator** - installs the [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}) into the cluster. MeshSync runs in-cluster and streams resource changes to Meshery in real time.
- **Embedded** - runs MeshSync from within Meshery Server. Nothing is installed into the cluster; discovery happens out-of-cluster. This is the default.

Switching the mode later makes Meshery redeploy MeshSync accordingly (see [Updating a Connection](#updating-a-connection)). For the behavioral trade-offs between the two modes - cluster footprint, permissions, network requirements, and what each mode gives up - and for every other setting of these components, see [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}}).

## Updating a Connection

The wizard also reconfigures an already-registered Connection. From the [Connections]({{< ref "concepts/logical/connections/index.md" >}}) table, open a Connection's action menu and choose **Configure**. The wizard opens in configure mode and presents only the post-registration steps relevant to that kind.

For a Kubernetes Connection, this is where you change the [MeshSync deployment mode](#meshsync-deployment-mode). Selecting a different mode and clicking **Apply** makes Meshery undeploy MeshSync and redeploy it in the newly selected mode (Operator or Embedded) for that cluster.

{{% alert color="info" title="Changing a Connection's state" %}}
Configuring a Connection is distinct from transitioning its **state** (Connected, Ignored, Disconnected, and so on). State transitions - and the rules governing which are allowed - are driven by the connection definition and performed from the status control on the Connections table. See [States and the Lifecycle of Connections]({{< ref "concepts/logical/connections/index.md#states-and-the-lifecycle-of-connections" >}}).
{{% /alert %}}

<!-- The Telemetry pages (guides/telemetry/*) ship in meshery/meshery#20161. Until that
     merges, these are absolute docs.meshery.io links so this page does not break the Hugo
     build, since an unresolved ref shortcode fails the build. Convert them to ref-shortcode
     links once the Telemetry pages exist on master. -->
## Using Connections for Telemetry

Grafana and Prometheus Connections you register with the wizard power Meshery's [Telemetry](https://docs.meshery.io/guides/telemetry/) views. Once such a Connection reaches the **Connected** or **Registered** state, it becomes selectable in the Telemetry connection picker, where you can:

- Browse and render your existing dashboards - see [Grafana Dashboards](https://docs.meshery.io/guides/telemetry/grafana-dashboards).
- Explore metrics and save PromQL panels - see [Prometheus Metrics](https://docs.meshery.io/guides/telemetry/prometheus-metrics).

## Registering Connections from the CLI

Prefer the terminal? `mesheryctl` can create, list, view, and delete Connections too. See [`mesheryctl connection`]({{< ref "reference/references/mesheryctl/connection/_index.md" >}}).

## Related

- [Connections]({{< ref "concepts/logical/connections/index.md" >}}) - concepts and state lifecycle.
- [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}}) - tuning the components that keep an imported cluster in sync.
- [Credentials]({{< ref "concepts/logical/credentials.md" >}}) - authentication for Connections.
- [Environments]({{< ref "concepts/logical/environments.md" >}}) and [Workspaces]({{< ref "concepts/logical/workspaces.md" >}}) - grouping and sharing Connections.
- [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}) - lifecycle operations.
- [Contributing a Connection]({{< ref "project/contributing/models/connections" >}}) - add a new Connection kind to the wizard.

{{< discuss >}}
