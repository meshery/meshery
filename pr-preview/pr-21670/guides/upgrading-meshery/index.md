# Upgrading Meshery

> Step-by-step procedure for upgrading Meshery and all of its components — CLI, Server, Operator, MeshSync, and Broker — on Docker and Kubernetes deployments.

Source: /pr-preview/pr-21670/guides/upgrading-meshery/

This guide walks through upgrading a running Meshery deployment. Meshery is a
composition of components that upgrade in a specific order: the CLI first, the
Server second, and the components on managed clusters — Meshery Operator,
MeshSync, and Broker — **automatically, driven by the Server**. You do not
upgrade the Operator by hand.

For background on which components version together, see the
[Upgrade Guide](/pr-preview/pr-21670/installation/upgrades/). For production
practices (pinned versions, upgrade-friendly probes, rollback rehearsal), see
the [Operational Readiness Checklist](/pr-preview/pr-21670/installation/production/operational-readiness-checklist/).

## Before you begin

- Note your release channel (`stable` for production; `edge` for testing).
- Record current versions so you can verify the upgrade:

```bash
mesheryctl version
kubectl -n meshery get deploy meshery-operator \
  -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'
```

## Step 1: Upgrade `mesheryctl`

Use the package manager you installed with:

```bash
brew upgrade mesheryctl        # Homebrew
scoop update mesheryctl        # Scoop
curl -L https://meshery.io/install | DEPLOY_MESHERY=false bash -   # Bash
```

## Step 2: Upgrade Meshery Server

### Docker deployments

```bash
mesheryctl system update    # pull latest images per your release channel
mesheryctl system restart   # apply them to the running deployment
```

### Kubernetes deployments (Helm)

Pin an explicit chart version rather than tracking latest, and use the
upgrade-friendly probe values so Server pods are not killed while reloading
capabilities:

```bash
helm repo update meshery
helm upgrade meshery meshery/meshery --namespace meshery \
  --version <target-version> \
  -f values-upgrade.yaml \
  --wait --timeout 10m
```

Keep your copy of
[`values-upgrade.yaml`](https://github.com/meshery/meshery/blob/master/install/kubernetes/helm/meshery/values-upgrade.yaml)
version-controlled alongside your own values.

## Step 3: Managed-cluster components upgrade themselves

When the upgraded Meshery Server (re)connects to a managed cluster, it
re-applies the `meshery-operator` Helm chart. The version it asks for tracks
the Server release, and before installing anything the Server checks that
version against the chart repository's published index - so what is actually
deployed is **always a chart the repository carries**. If the matching chart
has not been published yet (chart publishing trails Server releases), the
newest published release is used instead, and the substitution is reported in
the events feed rather than made silently. See
[How Meshery Server manages Meshery Operator](/pr-preview/pr-21670/installation/upgrades/#how-meshery-server-manages-meshery-operator)
for the full resolution rules, including the minimum chart version and how to
pin one yourself.

That single `helm upgrade`, performed by the Server:

1. runs the chart's CRD update Job, which server-side-applies the current
   `Broker` and `MeshSync` CRD schemas (Helm alone never updates CRDs on
   upgrade — the Job is what delivers schema changes to live clusters);
2. rolls the Operator Deployment to the operator image version pinned in that
   chart;
3. the Operator then reconciles MeshSync and Broker to their expected
   versions and configuration - each a pinned release, never a moving
   `stable-latest` tag, so an Operator pod restart never changes what runs.

No action is required on managed clusters.

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Do not hand-upgrade the Operator on Server-managed clusters</div>


A manual `helm upgrade` of `meshery-operator` (or a hand-edited image tag) on
a cluster that Meshery Server manages is a stopgap at best: the Server's
reconciliation re-applies the chart version it resolves and will revert your
change. The durable way to get a newer Operator is to upgrade Meshery Server;
to hold one connection at a specific chart version, set `operator.version` on
that connection - see
[Choosing the chart version yourself](/pr-preview/pr-21670/guides/troubleshooting/meshery-operator-meshsync/#choosing-the-chart-version-yourself).
</div>


If the Operator does not come back after the upgrade, its status card carries
the reason - a chart version that could not be resolved surfaces there and in
the connection's Diagnostics rather than disappearing. A failure caused by a
transient chart-repository outage clears on its own: redeploy the Operator from
the connection's actions and the Server re-resolves the version.

## Step 4: Verify

```bash
# Server and components
mesheryctl system status

# The Meshery Operator image actually running in the cluster
kubectl -n meshery get deploy meshery-operator \
  -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'

# The chart release that installed it (chart version and app version)
helm -n meshery list --filter meshery-operator

# CRDs are current (v1alpha2 storage) and healthy
kubectl get crds brokers.meshery.io meshsyncs.meshery.io

# Broker and MeshSync are reconciled and ready
kubectl -n meshery get brokers,meshsyncs
kubectl -n meshery get statefulset/meshery-nats deployment/meshery-meshsync
```

In Meshery UI, confirm the cluster connection shows the Operator, MeshSync,
and Broker as connected under **Settings → Environment**.

## Rolling back

```bash
helm rollback meshery --namespace meshery    # Kubernetes deployments
```

Rolling back the Server is low-risk for data: durable state lives with your
Remote Provider, and the local database is a rebuildable cache. Two notes:

- The Operator follows the Server: after a rollback, the Server re-applies the
  operator chart it resolves for the rolled-back release. Rolling back far
  enough that the Server would name a chart older than the minimum deployable
  version normally still lands a working Operator - the Server raises that
  derived request to the oldest published chart at or above the minimum and
  says so in the events feed. An explicit `operator.version` on a connection is
  never raised, so a pin below the minimum survives the rollback unchanged.
- CRD schemas are not rolled back by `helm rollback` directly. When the
  rolled-back Server reconnects to managed clusters, it re-applies the older
  operator chart as a Helm upgrade, and that chart's CRD update Job re-applies
  the older schemas. Stored objects remain readable throughout, because served
  versions stay identical across current schema revisions.

<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/design-meshery-mcp-server-architecture-and-registration-interface/7954" target="_blank" rel="noopener noreferrer">
            Design: Meshery MCP Server architecture and registration interface
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/the-meshery-mcp-server-foundation-is-up-lets-agree-on-what-to-build-next/7952" target="_blank" rel="noopener noreferrer">
            The Meshery MCP Server foundation is up, let&#39;s agree on what to build next
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/new-intro-topic/7975" target="_blank" rel="noopener noreferrer">
            New intro topic
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-mcp-server-poc-an-ai-agent-managing-kubernetes-through-mcp/7974" target="_blank" rel="noopener noreferrer">
            Meshery MCP Server POC: an AI agent managing Kubernetes through MCP
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/approach-for-context-window-aware-retrieval-in-the-ai-adapter-issue-20994/7963" target="_blank" rel="noopener noreferrer">
            Approach for context-window-aware retrieval in the AI Adapter (Issue #20994)
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/there-is-some-error-in-running-the-localhost-of-layer5-in-my-server/7818" target="_blank" rel="noopener noreferrer">
            There is some error in running the localhost of layer5 in my server
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/rfc-aligning-the-foundation-for-the-meshery-mcp-server/7913" target="_blank" rel="noopener noreferrer">
            RFC: Aligning the Foundation for the Meshery MCP Server
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-development-meeting-august-12th-2026/7948" target="_blank" rel="noopener noreferrer">
            Meshery Development Meeting | August 12th, 2026
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
