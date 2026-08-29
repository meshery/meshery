# Registering a Connection

> Use the Connection Wizard to create and update Connections - Kubernetes clusters, Grafana, Prometheus, and more - in your Meshery deployment.

Source: /pr-preview/pr-21670/guides/infrastructure-management/registering-a-connection/

A [Connection](/pr-preview/pr-21670/concepts/logical/connections/) is how Meshery tracks and manages a resource - a Kubernetes cluster, a Grafana instance, a Prometheus server, and [many more](/pr-preview/pr-21670/extensions/models/). The **Connection Wizard** is the guided, in-UI way to register a new Connection or reconfigure an existing one, without hand-editing YAML or memorizing API payloads.

This guide covers creating and updating Connections with the wizard. For what a Connection _is_, the states it moves through, and how it is managed over time, see the canonical references:

- [Connections](/pr-preview/pr-21670/concepts/logical/connections/) - what Connections are and their full state lifecycle.
- [Credentials](/pr-preview/pr-21670/concepts/logical/credentials/) - how Meshery authenticates to a Connection.
- [Managing Connections](/pr-preview/pr-21670/guides/infrastructure-management/lifecycle-management/) - operating Connections after they are registered.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Discovered vs. manually registered Connections</div>


Meshery learns about Connections two ways. **Managed** Connections (for example, the resources inside a Kubernetes cluster) are auto-discovered by [MeshSync](/pr-preview/pr-21670/concepts/architecture/meshsync/) and arrive already in the [Discovered](/pr-preview/pr-21670/concepts/logical/connections/#state-discovered) state. **Unmanaged** Connections (for example, a standalone Grafana or Prometheus) are added by you. The Connection Wizard is how you do the latter - and how you bring a Kubernetes cluster under management in the first place.
</div>


## Before you begin

- A running Meshery deployment. See the [Quick Start](/pr-preview/pr-21670/installation/quick-start/) if you do not have one yet.
- Permission to add Connections. The wizard is permission-gated: adding a Kubernetes cluster requires the **add cluster** permission, and other Connection kinds require the **connect metrics** permission. If you lack both, the **Create Connection** button is disabled. See [Roles and Permissions](/pr-preview/pr-21670/reference/extensibility/authorization/).
- For an authenticated Connection (most Grafana/Prometheus instances), the access token or credential you intend to use. You can paste it during the wizard or reuse an existing [Credential](/pr-preview/pr-21670/concepts/logical/credentials/).

## Launching the Connection Wizard

1. Open the **Connections** page in Meshery (**Lifecycle → Connections**).
2. Click **Create Connection**.

The wizard opens as a modal. The set of Connection kinds you can create is driven by the [connection definitions](/pr-preview/pr-21670/project/contributing/models/connections/) registered in your Meshery Server's [Registry](/pr-preview/pr-21670/concepts/logical/registry/). Out of the box this includes **Kubernetes**, **Grafana**, **Prometheus**, **Artifact Hub**, and **GitHub**; your deployment may offer more. Because the list is registry-driven, the wizard surfaces each kind - its name, description, icon, and its Configure and Credential fields - straight from the definition, so no kind is hard-coded into the UI. If a kind you need is missing, a contributor can add it - see [Contributing a Connection](/pr-preview/pr-21670/project/contributing/models/connections/).

## Creating a Connection

Most Connections follow the same generic flow. Each step is rendered from the connection definition itself, so the exact fields you see depend on the kind you choose.

1. **Choose Connection.** Pick the kind of Connection to create (for example, Grafana). Kinds you do not have permission to add are shown but cannot be selected.
2. **Configure Connection.** Fill in the Connection's details - typically the endpoint URL and an optional friendly name. Required fields are validated before you can continue. For a Grafana Connection, for instance, you supply the Grafana endpoint (e.g. `http://grafana.example:3000`).
3. **Associate Credential.** Provide the secret Meshery will use to authenticate. You can either:
   - **Reuse an existing credential** - the list is filtered to credentials that match the Connection's kind, or
   - **Create a new credential** - enter the token, API key, or `username:password` and give it a name (it defaults to the Connection's name).

   You may also choose to **skip credential verification**, which registers the Connection without first probing reachability - useful when the target is not reachable yet but you still want it on record. This step is omitted entirely for kinds that do not define a credential (and for Kubernetes, whose kubeconfig _is_ its credential - see below).
4. **Review & Create.** Confirm the summary and click **Create Connection**. Meshery registers the Connection and immediately attempts to connect to it.
5. **Done.** On success, the Connection becomes a first-class resource, listed in the [Connections](/pr-preview/pr-21670/concepts/logical/connections/) table and ready to use.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">What 'Create' actually does</div>


Creating a Connection performs two transitions in sequence: it **registers** the Connection (recording it and its credential) and then **connects** to it (verifying reachability and beginning management). A reachable Connection lands in the [Connected](/pr-preview/pr-21670/concepts/logical/connections/#state-connected) state; if Meshery cannot reach it - or you skipped verification - it remains [Registered](/pr-preview/pr-21670/concepts/logical/connections/#state-registered). You can drive further transitions later from the Connections table.
</div>


### Credentials

Credentials entered in the wizard are persisted as first-class, named [Credentials](/pr-preview/pr-21670/concepts/logical/credentials/), encrypted at rest, and reusable across other Connections. Meshery never exposes them in logs or API responses. To learn how Meshery interprets a credential's secret (Basic auth vs. bearer token vs. anonymous), see [Credentials](/pr-preview/pr-21670/concepts/logical/credentials/) and the [Telemetry authentication note](https://docs.meshery.io/guides/telemetry/).

## Importing a Kubernetes cluster

Kubernetes uses a dedicated flow because a single kubeconfig can describe many clusters and its kubeconfig also serves as its credential.

1. **Choose Connection** → **Kubernetes**.
2. **Import Kubeconfig.** Upload a kubeconfig file. Meshery parses it and lists the contexts it contains, indicating which are reachable - nothing is persisted yet.
3. **Select contexts.** Choose which contexts to import. For each, you can override the Connection name and choose a [MeshSync deployment mode](#meshsync-deployment-mode).
4. **Review Import.** Confirm your selection and import. Meshery creates one Connection per selected context and reports the outcome, grouped into connected, registered, ignored, and errored buckets.

Each imported cluster is created as a Kubernetes Connection that [MeshSync](/pr-preview/pr-21670/concepts/architecture/meshsync/) keeps in sync. From there, Meshery can deploy and operate workloads, visualize the cluster, and more. See [Managing Connections](/pr-preview/pr-21670/guides/infrastructure-management/lifecycle-management/). If the cluster already runs workloads you did not deploy through Meshery, see [Bringing Existing Infrastructure Under Meshery Management](/pr-preview/pr-21670/guides/infrastructure-management/managing-existing-infrastructure/) for what to evaluate before importing and how discovery treats pre-existing resources.

<div class="alert alert-dark" role="alert"><div class="h4 alert-heading" role="heading">Who can access an imported cluster?</div>


A Kubernetes Connection is owned by the user who imported it and is private until you explicitly share it - by assigning it to an [Environment](/pr-preview/pr-21670/concepts/logical/environments/) and that environment to a [Workspace](/pr-preview/pr-21670/concepts/logical/workspaces/). See the sharing FAQ under [Managing Connections](/pr-preview/pr-21670/guides/infrastructure-management/lifecycle-management/).
</div>


### MeshSync deployment mode

When you import or reconfigure a Kubernetes cluster, you choose how [MeshSync](/pr-preview/pr-21670/concepts/architecture/meshsync/) - the component that keeps Meshery's view of the cluster's resources up to date - runs:

- **Operator** - installs the [Meshery Operator](/pr-preview/pr-21670/concepts/architecture/operator/) into the cluster. MeshSync runs in-cluster and streams resource changes to Meshery in real time.
- **Embedded** - runs MeshSync from within Meshery Server. Nothing is installed into the cluster; discovery happens out-of-cluster. This is the default.

Switching the mode later makes Meshery redeploy MeshSync accordingly (see [Updating a Connection](#updating-a-connection)). For the behavioral trade-offs between the two modes - cluster footprint, permissions, network requirements, and what each mode gives up - and for every other setting of these components, see [Configuring Meshery Operator, MeshSync, and Broker](/pr-preview/pr-21670/guides/infrastructure-management/configuring-operator-meshsync-broker/).

## Connecting a source (Artifact Hub, GitHub)

Some Connections are **sources** - external registries and repositories that Meshery reads from to generate [Models](/pr-preview/pr-21670/concepts/logical/models/) and components, and (for repositories) to import designs kept under version control. They follow the same generic flow described above; only their fields differ.

### Artifact Hub

Register an [Artifact Hub](https://artifacthub.io) instance to source Helm charts and other CNCF artifacts, then generate Meshery models and components from the CRDs they carry.

- **Configure Connection.** Supply the **Artifact Hub Endpoint** - the base URL of the instance (for example `https://artifacthub.io`). Point this at your own deployment to source from a self-hosted Artifact Hub. An optional friendly **Connection Name** is also accepted.
- **Associate Credential.** Artifact Hub's public API needs no authentication - leave both credential fields empty to use it. To raise your rate limit or reach private content, provide an **API Key ID** _and_ its **API Key Secret**. The two are paired: supply both or neither.

### GitHub

Register a [GitHub](https://github.com) repository to generate models and components from the Kubernetes manifests and CRDs it holds, and to import designs stored alongside the code they describe.

- **Configure Connection.** Supply the **Repository URL** Meshery reads from (for example `https://github.com/meshery/meshery`); GitHub Enterprise Server URLs are accepted. Optionally set a **Branch** or tag (leave empty for the repository's default branch) and a friendly **Connection Name**.
- **Associate Credential.** Provide an **Access Token** (a personal access token or GitHub App installation token). It is required for private repositories; for public repositories it is optional and raises the API rate limit.

For both kinds you can **skip credential verification** to register the source without first probing it - useful when it is not reachable yet but you still want it on record. Once connected, a source's packages or manifests become available for [model and component generation](/pr-preview/pr-21670/guides/configuration-management/importing-models/).

## Updating a Connection

The wizard also reconfigures an already-registered Connection. From the [Connections](/pr-preview/pr-21670/concepts/logical/connections/) table, open a Connection's action menu and choose **Configure**. The wizard opens in configure mode and presents only the post-registration steps relevant to that kind.

For a Kubernetes Connection, this is where you change the [MeshSync deployment mode](#meshsync-deployment-mode). Selecting a different mode and clicking **Apply** makes Meshery undeploy MeshSync and redeploy it in the newly selected mode (Operator or Embedded) for that cluster.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Changing a Connection's state</div>


Configuring a Connection is distinct from transitioning its **state** (Connected, Ignored, Disconnected, and so on). State transitions - and the rules governing which are allowed - are driven by the connection definition and performed from the status control on the Connections table. See [States and the Lifecycle of Connections](/pr-preview/pr-21670/concepts/logical/connections/#states-and-the-lifecycle-of-connections).
</div>


<!-- The Telemetry pages (guides/telemetry/*) ship in meshery/meshery#20161. Until that
     merges, these are absolute docs.meshery.io links so this page does not break the Hugo
     build, since an unresolved ref shortcode fails the build. Convert them to ref-shortcode
     links once the Telemetry pages exist on master. -->
## Using Connections for Telemetry

Grafana and Prometheus Connections you register with the wizard power Meshery's [Telemetry](https://docs.meshery.io/guides/telemetry/) views. Once such a Connection reaches the **Connected** or **Registered** state, it becomes selectable in the Telemetry connection picker, where you can:

- Browse and render your existing dashboards - see [Grafana Dashboards](https://docs.meshery.io/guides/telemetry/grafana-dashboards).
- Explore metrics and save PromQL panels - see [Prometheus Metrics](https://docs.meshery.io/guides/telemetry/prometheus-metrics).

## Registering Connections from the CLI

Prefer the terminal? `mesheryctl` can create, list, view, and delete Connections too. See [`mesheryctl connection`](/pr-preview/pr-21670/reference/references/mesheryctl/connection/).

## Related

- [Connections](/pr-preview/pr-21670/concepts/logical/connections/) - concepts and state lifecycle.
- [Configuring Meshery Operator, MeshSync, and Broker](/pr-preview/pr-21670/guides/infrastructure-management/configuring-operator-meshsync-broker/) - tuning the components that keep an imported cluster in sync.
- [Credentials](/pr-preview/pr-21670/concepts/logical/credentials/) - authentication for Connections.
- [Environments](/pr-preview/pr-21670/concepts/logical/environments/) and [Workspaces](/pr-preview/pr-21670/concepts/logical/workspaces/) - grouping and sharing Connections.
- [Managing Connections](/pr-preview/pr-21670/guides/infrastructure-management/lifecycle-management/) - lifecycle operations.
- [Contributing a Connection](/pr-preview/pr-21670/project/contributing/models/connections/) - add a new Connection kind to the wizard.

<div class="alert alert-dark" role="alert">
  <h4 class="alert-heading">Discussion Forum</h4>
  <p>Don't find an answer to your question here? Ask on the <a href="https://discuss.meshery.io/">Discussion Forum</a>.</p>
</div>
