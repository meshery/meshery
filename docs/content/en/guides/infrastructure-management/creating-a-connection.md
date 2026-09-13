---
title: Creating a Connection
aliases:
- /guides/infrastructure-management/registering-a-connection
- /guides/infrastructure-management/connection-wizard
- /guides/infrastructure-management/infrastructure-discovery
- /connection-wizard
categories: [infrastructure]
weight: -5
description: Create and update Connections - Kubernetes clusters, Grafana, Prometheus, and more - with the Connection Wizard in the Meshery UI or with mesheryctl.
---

A [Connection]({{< ref "concepts/logical/connections/index.md" >}}) is how Meshery tracks and manages a resource - a Kubernetes cluster, a Grafana instance, a Prometheus server, and [many more]({{< ref "extensions/models/_index.md" >}}). Meshery supports creating and managing Connections through either the **Connection Wizard** in the Meshery UI or the **`mesheryctl` CLI**.

The **Connection Wizard** is the guided, in-UI way to create a new Connection or reconfigure an existing one, without hand-editing YAML or memorizing API payloads. Users who prefer working from the terminal or automating workflows can create and manage **Kubernetes** Connections using [`mesheryctl connection`]({{< ref "reference/references/mesheryctl/connection/_index.md" >}}). Grafana, Prometheus, and other non-Kubernetes kinds are created from the Meshery UI.

This guide covers both supported connection creation workflows. For what a Connection _is_, the states it moves through, and how it is managed over time, see the canonical references:

- [Connections]({{< ref "concepts/logical/connections/index.md" >}}) - what Connections are and their full state lifecycle.
- [Credentials]({{< ref "concepts/logical/credentials.md" >}}) - how Meshery authenticates to a Connection.
- [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}) - operating Connections after they are created.

{{% alert color="info" title="Discovered vs. manually registered Connections" %}}
Meshery learns about Connections two ways. **Managed** Connections (for example, the resources inside a Kubernetes cluster) are auto-discovered by [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) and arrive already in the [Discovered]({{< ref "concepts/logical/connections/index.md#state-discovered" >}}) state. **Unmanaged** Connections (for example, a standalone Grafana or Prometheus) are added by you. The Connection Wizard is how you do the latter - and how you bring a Kubernetes cluster under management in the first place.
{{% /alert %}}

## Before you begin

- A running Meshery deployment. See the [Quick Start]({{< ref "installation/quick-start/index.md" >}}) if you do not have one yet.
- Permission to add Connections. The wizard is permission-gated: adding a Kubernetes cluster requires the **add cluster** permission, and other Connection kinds require the **connect metrics** permission. If you lack both, the **Create Connection** button is disabled. See [Roles and Permissions]({{< ref "reference/extensibility/authorization/index.md" >}}).
- For an authenticated Connection (most Grafana/Prometheus instances), the access token or credential you intend to use. You can paste it during the wizard or reuse an existing [Credential]({{< ref "concepts/logical/credentials.md" >}}).

## Creating a Connection

Meshery supports two workflows for creating Connections. Choose the workflow that best fits your environment.

{{< tabs id="connection-creation-tabs" >}}

Meshery UI | fa fa-desktop

## Launching the Connection Wizard

1. Open the **Connections** page in Meshery (**Lifecycle → Connections**).

<a href="../images/connections-page.png">
  <img
    src="../images/connections-page.png"
    alt="Connections page in Meshery"
    style="width:60%; max-width:800px;" />
</a>

2. Click **Create Connection**.

The wizard opens as a modal. The set of Connection kinds you can create is driven by the [connection definitions]({{< ref "project/contributing/models/connections" >}}) registered in your Meshery Server's [Registry]({{< ref "concepts/logical/registry.md" >}}). Out of the box this includes **Kubernetes**, **Grafana**, and **Prometheus**; your deployment may offer more. If a kind you need is missing, a contributor can add it - see [Contributing a Connection]({{< ref "project/contributing/models/connections" >}}).

<a href="../images/connection-wizard-choose-connection.png">
  <img src="../images/connection-wizard-choose-connection.png"
       alt="Choose a connection type in the Connection Wizard"
       style="width:60%; max-width:800px;">
</a>

## Generic connection flow

Most Connections follow the same generic flow. Each step is rendered from the connection definition itself, so the exact fields you see depend on the kind you choose.

1. **Choose Connection.** Pick the kind of Connection to create (for example, Grafana). Kinds you do not have permission to add are shown but cannot be selected.
2. **Configure Connection.** Fill in the Connection's details - typically the endpoint URL and an optional friendly name. Required fields are validated before you can continue. For a Grafana Connection, for instance, you supply the Grafana endpoint (e.g. `http://grafana.example:3000`).

<a href="../images/connection-wizard-configure-connection.png">
  <img src="../images/connection-wizard-configure-connection.png"
       alt="Configure a connection in the Connection Wizard"
       style="width:60%; max-width:800px;">
</a>

3. **Associate Credential.** Provide the secret Meshery will use to authenticate. You can either:

<a href="../images/connection-wizard-associate-credential.png">
  <img src="../images/connection-wizard-associate-credential.png"
       alt="Associate a credential in the Connection Wizard"
       style="width:60%; max-width:800px;">
</a>

   - **Reuse an existing credential** - the list is filtered to credentials that match the Connection's kind, or
   - **Create a new credential** - enter the token, API key, or `username:password` and give it a name (it defaults to the Connection's name).

   You may also choose to **skip credential verification**, which records the Connection without first probing reachability - useful when the target is not reachable yet but you still want it on record. This step is omitted entirely for kinds that do not define a credential (and for Kubernetes, whose kubeconfig _is_ its credential - see below).

4. **Review & Create.** Confirm the summary and click **Create Connection**. Meshery creates the Connection and attempts to connect to it (verify reachability, verify credentials, and begin management when possible).

<a href="../images/connection-wizard-review-create.png">
  <img src="../images/connection-wizard-review-create.png"
       alt="Review connection details before creating"
       style="width:60%; max-width:800px;">
</a>

5. **Done.** The Connection is added to the [Connections]({{< ref "concepts/logical/connections/index.md" >}}) table and is ready to use.

{{% alert color="info" title="What 'Create' actually does" %}}
Creating a Connection records it (and its credential) and then tries to **connect** - verifying reachability and beginning management. A reachable Connection lands in the [Connected]({{< ref "concepts/logical/connections/index.md#state-connected" >}}) state. If Meshery cannot reach the target - or you skipped verification - the Connection remains available in the table so you can connect later (for example after the endpoint is up). You can drive further state transitions from the Connections table. See [States and the Lifecycle of Connections]({{< ref "concepts/logical/connections/index.md#states-and-the-lifecycle-of-connections" >}}).
{{% /alert %}}

### Credentials

Credentials entered in the wizard are persisted as first-class, named [Credentials]({{< ref "concepts/logical/credentials.md" >}}), encrypted at rest, and reusable across other Connections. Meshery never exposes them in logs or API responses. To learn how Meshery interprets a credential's secret (Basic auth vs. bearer token vs. anonymous), see [Credentials]({{< ref "concepts/logical/credentials.md" >}}) and the [Telemetry authentication note]({{< ref "guides/telemetry/_index.md" >}}).

## Importing a Kubernetes cluster

Kubernetes uses a dedicated flow because a single kubeconfig can describe many clusters and its kubeconfig also serves as its credential. The wizard separates **discovering** what is in the file from **importing** selected contexts as Connections, and separates **importing** a Connection from **connecting** to it (interacting with the cluster).

1. **Choose Connection** → **Kubernetes**.
2. **Import Kubeconfig.** Upload a kubeconfig file, then click **Discover Contexts**. Meshery reads the contexts inside that file and checks which are reachable. A context that cannot be reached is shown as **not found** (unreachable) - discovery has not failed; Meshery simply cannot talk to that API server yet. Nothing is imported as a Connection until you continue.

<a href="../images/connection-wizard-import-kubeconfig.png">
  <img src="../images/connection-wizard-import-kubeconfig.png"
       alt="Import a kubeconfig file in the Connection Wizard"
       style="width:60%; max-width:800px;">
</a>

3. **Review Contexts.** Choose which clusters (contexts) to import, rename them if you like, and set each context's [MeshSync deployment mode](#meshsync-deployment-mode). Use the checkbox **Connect reachable clusters after import** when you want Meshery to start managing reachable clusters immediately. Importing records the Connection; connecting is what starts interaction with the cluster. Unreachable clusters can still be imported and connected later once the API server is reachable.

<a href="../images/connection-wizard-review-contexts.png">
  <img
    src="../images/connection-wizard-review-contexts.png"
    alt="Review Kubernetes contexts before importing"
    style="width:60%; max-width:800px;">
</a>

4. Click **Import.** Meshery creates one Connection per selected context and shows the **Done** step with each name and status chip. Typical outcomes:

   - **[Connected]({{< ref "concepts/logical/connections/index.md#state-connected" >}})** - reachable and connected after import (when connect-after-import is enabled).
   - **[Discovered]({{< ref "concepts/logical/connections/index.md#state-discovered" >}})** - recorded for use; not yet connected (for example connect-after-import was off, or connect did not complete).
   - **[Not Found]({{< ref "concepts/logical/connections/index.md#state-not-found" >}})** - Meshery could not reach the cluster (same idea as **not found** during Discover Contexts).

<a href="../images/connection-wizard-import-complete.png">
  <img
    src="../images/connection-wizard-import-complete.png"
    alt="Kubernetes import completed successfully"
    style="width:60%; max-width:800px;">
</a>

Each imported cluster is created as a Kubernetes Connection. Once connected, [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) keeps Meshery's view of the cluster's resources in sync. From there, Meshery can deploy and operate workloads, visualize the cluster, and more. See [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}). If the cluster already runs workloads you did not deploy through Meshery, see [Bringing Existing Infrastructure Under Meshery Management]({{< ref "guides/infrastructure-management/managing-existing-infrastructure.md" >}}) for what to evaluate before importing and how discovery treats pre-existing resources.

{{% alert color="dark" title="Who can access an imported cluster?" %}}
A Kubernetes Connection is owned by the user who imported it and is private until you explicitly share it - by assigning it to an [Environment]({{< ref "concepts/logical/environments.md" >}}) and that environment to a [Workspace]({{< ref "concepts/logical/workspaces.md" >}}). See the sharing FAQ under [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}).
{{% /alert %}}

### MeshSync deployment mode

When you import or reconfigure a Kubernetes cluster, you choose how [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) - the component that keeps Meshery's view of the cluster's resources up to date - runs:

- **Operator** - installs the [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}) into the cluster. MeshSync runs in-cluster and streams resource changes to Meshery in real time.
- **Embedded** - runs MeshSync from within Meshery Server. Nothing is installed into the cluster; discovery happens out-of-cluster. This is the default.

Switching the mode later makes Meshery redeploy MeshSync accordingly (see [Updating a Connection](#updating-a-connection)). For the behavioral trade-offs between the two modes - cluster footprint, permissions, network requirements, and what each mode gives up - and for every other setting of these components, see [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}}).

## Updating a Connection

The wizard also reconfigures an already-created Connection. From the [Connections]({{< ref "concepts/logical/connections/index.md" >}}) table, open a Connection's action menu and choose **Configure**. The wizard opens in configure mode and presents only the post-creation steps relevant to that kind.

For a Kubernetes Connection, this is where you change the [MeshSync deployment mode](#meshsync-deployment-mode). Selecting a different mode and clicking **Apply** makes Meshery undeploy MeshSync and redeploy it in the newly selected mode (Operator or Embedded) for that cluster.

<!-- tab -->

mesheryctl | fa fa-terminal

## Creating a Connection with mesheryctl

If you haven't already, install and configure `mesheryctl` by following the [Meshery CLI installation guide]({{< ref "installation/_index.md" >}}).

Once mesheryctl is installed, authenticate with `mesheryctl system login` and ensure Meshery Server is running and you have a valid Kubernetes context configured. You can then create and manage **Kubernetes** Connections from the terminal.

`mesheryctl connection create` supports Kubernetes provider types only (`aks`, `eks`, `gke`, `minikube`). To create Grafana, Prometheus, or other non-Kubernetes Connections, use the **Meshery UI** Connection Wizard.

The CLI is well suited for terminal-based workflows, scripting, and automation where using the Connection Wizard is not required.

Before creating a Kubernetes Connection with `mesheryctl`, ensure that you've completed the prerequisites for your Kubernetes provider. For example, GKE requires the `gcloud` CLI and EKS requires the AWS CLI. See the [Quick Start with Kubernetes]({{< ref "installation/kubernetes/_index.md" >}}) guide for provider-specific prerequisites and setup instructions.

### Create a Connection

Create a Kubernetes Connection by specifying the Kubernetes provider type.

```bash
mesheryctl connection create --type minikube
```

You can also create Connections for other supported Kubernetes providers:

```bash
mesheryctl connection create --type aks
mesheryctl connection create --type eks
mesheryctl connection create --type gke
```

If you're authenticating with a token file, use the `--token` flag. See the [`mesheryctl connection`]({{< ref "reference/references/mesheryctl/connection/_index.md" >}}) reference for supported authentication options.

### Verify Creation

List Connections to verify that your Connection has been created.

```bash
mesheryctl connection list
```

### Inspect a Connection

View detailed information about a Connection by name or ID.

```bash
mesheryctl connection view <connection-name>
```

### Remove a Connection

Delete a Connection when it is no longer needed.

```bash
mesheryctl connection delete <connection-id>
```

For additional commands, flags, supported providers, and examples, see [`mesheryctl connection`]({{< ref "reference/references/mesheryctl/connection/_index.md" >}}).

{{< /tabs >}}

{{% alert color="info" title="Changing a Connection's state" %}}
Configuring a Connection is distinct from transitioning its **state** (for example Connected, Discovered, Disconnected, Not Found, or Deleted). State transitions - and the rules governing which are allowed - are driven by the connection definition and performed from the status control on the Connections table. See [States and the Lifecycle of Connections]({{< ref "concepts/logical/connections/index.md#states-and-the-lifecycle-of-connections" >}}).
{{% /alert %}}

## Using Connections for Telemetry

Grafana and Prometheus Connections you create with the wizard power Meshery's [Telemetry]({{< ref "guides/telemetry/_index.md" >}}) views. Once such a Connection reaches the **Connected** state (or is otherwise available for use), it becomes selectable in the Telemetry connection picker, where you can:

- Browse and render your existing dashboards - see [Grafana Dashboards]({{< ref "guides/telemetry/grafana-dashboards.md" >}}).
- Explore metrics and save PromQL panels - see [Prometheus Metrics]({{< ref "guides/telemetry/prometheus-metrics.md" >}}).

## Related

- [Connections]({{< ref "concepts/logical/connections/index.md" >}}) - concepts and state lifecycle.
- [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}}) - tuning the components that keep an imported cluster in sync.
- [Credentials]({{< ref "concepts/logical/credentials.md" >}}) - authentication for Connections.
- [Environments]({{< ref "concepts/logical/environments.md" >}}) and [Workspaces]({{< ref "concepts/logical/workspaces.md" >}}) - grouping and sharing Connections.
- [Managing Connections]({{< ref "guides/infrastructure-management/lifecycle-management/index.md" >}}) - lifecycle operations.
- [Contributing a Connection]({{< ref "project/contributing/models/connections" >}}) - add a new Connection kind to the wizard.

{{< discuss >}}
