# Kubernetes Connection Lifecycle

> What happens end-to-end when Meshery manages a Kubernetes cluster - the connection states, the in-cluster components (Operator, Broker/NATS, MeshSync), MeshSync deployment modes, how Meshery reaches and authenticates to the Broker, the built-in Diagnostics, and how to operate and troubleshoot it all.

---

LLMS index: [llms.txt](/pr-preview/pr-22672/llms.txt)

---

This guide explains what actually happens when Meshery manages a **Kubernetes**
[Connection](/pr-preview/pr-22672/concepts/logical/connections/): the components
Meshery deploys, how it stays in sync with your cluster, the connectivity it
depends on, and the knobs and diagnostics you use to operate it. For the generic
connection state machine that applies to every kind of connection, see
[Connections](/pr-preview/pr-22672/concepts/logical/connections/); for how to
register one in the first place, see
[Registering a Connection](/pr-preview/pr-22672/guides/infrastructure-management/registering-a-connection/).

## What happens when you connect a cluster

Meshery automatically discovers the clusters in your kubeconfig
(`$HOME/.kube/config`) and, for each, tracks a Connection through its
[lifecycle states](/pr-preview/pr-22672/concepts/logical/connections/#states-and-the-lifecycle-of-connections)
- `Discovered` → `Registered` → **`Connected`**, and later `Disconnected` /
`Deleted` as appropriate.

When a Kubernetes connection transitions to **Connected**, Meshery (depending on
the [MeshSync deployment mode](#meshsync-deployment-modes)) brings up the pieces
needed to observe the cluster and begins streaming live cluster state into the
[Meshery Database](/pr-preview/pr-22672/concepts/architecture/database/). From
that point the Connections table and the cluster views reflect what is really
running in the cluster.

## The in-cluster components

In **operator** mode, three components run inside the managed cluster
(namespace `meshery`):

| Component | What it is | What it does |
|---|---|---|
| **Meshery Operator** | A Kubernetes operator | Deploys and reconciles the Broker and MeshSync; owns their lifecycle. |
| **Meshery Broker** | A [NATS](/pr-preview/pr-22672/concepts/architecture/broker/) server (`meshery-nats`) | The message bus. MeshSync publishes cluster state to it; Meshery Server subscribes to it. |
| **MeshSync** | A custom controller | Watches the cluster (multi-tier discovery) and publishes changes to the Broker. |

Each reports a status that surfaces on the connection (for example the Broker and
MeshSync move to **CONNECTED** once Meshery is actually receiving data). See
[Understanding the status](/pr-preview/pr-22672/guides/troubleshooting/meshery-operator-meshsync/#understanding-the-status-of-meshery-operator-meshsync-and-meshery-broker).

### Which Operator version gets installed

Meshery Server installs the Operator as a Helm chart from
[meshery.io/charts](https://meshery.io/charts). The version it asks for tracks
the Meshery Server release, but the Server checks that version against the
repository's published index **before** installing, so it never asks Helm for a
chart that does not exist:

- **Not published yet** (chart publishing trails Server releases): the newest
  published chart is installed instead.
- **Too old to run**: a version below the oldest chart verified to deploy is
  raised to the oldest published chart at or above that boundary.
- **Release candidates** are never selected automatically.

Neither substitution is silent - each is reported in the events feed, naming
the version deployed and why. To choose the version yourself for one
connection, set `operator.version`; an explicit pin is honored exactly, and is
refused with a visible error if it names a version the repository does not
publish or a moving tag such as `stable-latest`. Details, including how to
confirm what is deployed, are in
[Choosing the chart version yourself](/pr-preview/pr-22672/guides/troubleshooting/meshery-operator-meshsync/#choosing-the-chart-version-yourself).

What the Operator then deploys is pinned by the Operator itself: MeshSync and
the Broker run specific released versions with pull policies to match, so
nothing drifts underneath a restart.

## MeshSync deployment modes

Meshery can run MeshSync in one of two modes. You can set it per connection, and
the default can be set with the `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` server
environment variable.

### Operator mode

The Operator, Broker (NATS) and MeshSync run **inside the cluster**. MeshSync
publishes to the in-cluster Broker and Meshery Server subscribes to it. This is
the scalable model - it isolates discovery from Meshery Server and works well
across many clusters - but it requires Meshery Server to have a network path to
the in-cluster Broker (see [Broker connectivity](#broker-connectivity)).

### Embedded mode (default)

MeshSync runs **in-process inside Meshery Server** (using an internal, in-memory
message channel instead of NATS), watching the cluster directly through the
Kubernetes API that Meshery already has access to. No Operator, NATS, or
port-forward is involved - so it "just works" out of the box for a single,
directly-reachable cluster. The trade-off is that discovery consumes Meshery
Server's own resources and it separates less cleanly across many clusters.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Switching modes</div>


Switch a connection's MeshSync mode from the Kubernetes connection's actions in
the UI. Under the hood this calls
`POST /api/integrations/connections/{connectionId}/actions` with
`{"action":"setMeshsyncMode","mode":"operator"|"embedded"}`; Meshery tears down
the old setup and redeploys MeshSync in the new mode.
</div>


## Broker connectivity

This section applies to **operator mode** only (embedded mode has no in-cluster
Broker). For Meshery Server to receive MeshSync data it must be able to both
**reach** and **authenticate to** the Broker.

### Authentication (token)

The Meshery Operator provisions NATS with token authentication and stores the
token in the `meshery-nats-auth` secret. Meshery Server reads that token and
presents it when connecting; without it the Broker rejects the connection with an
authorization violation. This is automatic - there is nothing to configure.

### Reachability

- **Meshery running in-cluster** (Helm / `mesheryctl system start --platform
  kubernetes`): the Broker's `ClusterIP` is directly reachable. Nothing to do.
- **Meshery running out-of-cluster** (a Docker host, your laptop, or managing a
  _different_ cluster): the Broker is typically exposed as `ClusterIP` only, which
  is **not reachable from outside the cluster**. Meshery needs a path to it.

### Managed port-forward (automatic)

When Meshery Server runs **out-of-cluster**, it automatically establishes a
**self-healing port-forward** to the Broker's NATS pod through the Kubernetes API
server (the same mechanism as `kubectl port-forward`, using the credentials
Meshery already holds) and connects to NATS over it. It re-establishes the tunnel
if the pod restarts or the connection drops, and tears it down when the connection
is disconnected or deleted. **No manual `kubectl port-forward` is required.**

| Behavior | Detail |
|---|---|
| Default | **On** when Meshery runs out-of-cluster. |
| In-cluster | **Skipped** automatically (the `ClusterIP` is directly reachable). |
| Disable | Set `MESHERY_MANAGED_BROKER_PORTFORWARD=false` on Meshery Server. |

When disabled, Meshery falls back to the resolved Broker endpoint (and
`127.0.0.1` / `host.docker.internal` on the Broker port), so you can provide your
own path, for example:

```bash
kubectl port-forward -n meshery svc/meshery-nats 4222:4222
```

## Diagnostics

Click a Kubernetes connection's row in the Connections table to open its detail
view. Alongside the Operator / MeshSync / Broker status chips, a **Diagnostics**
section reports actionable problems and remediation, computed from the live
controller status and Meshery's actual Broker connection. The codes it reports:

| Code | Meaning | What to do |
|---|---|---|
| `connection_inactive` | No active session for the connection yet. | Connect the cluster. |
| `operator_deploy_failed` | Meshery tried to deploy the Operator and could not; the underlying cause is carried in the diagnostic. | Read the cause. An unreadable kubeconfig or missing permissions in the `meshery` namespace need fixing first; an unresolvable chart version means the requested `operator.version` is not published, or the chart repository was unreachable. Redeploying the Operator re-resolves the version, so a transient repository outage clears on retry. |
| `operator_not_deployed` | The Operator isn't deployed (operator mode). | Reconnect the cluster, or switch MeshSync mode; ensure Meshery can create resources in the `meshery` namespace. |
| `broker_unreachable` | The Broker is up but Meshery can't reach/authenticate to it. | Make the Broker reachable (the managed port-forward normally handles this; otherwise port-forward it, expose it via NodePort/LoadBalancer, or run Meshery in-cluster). |
| `broker_networking` | Informational: how Meshery reaches the Broker - managed port-forward, in-cluster ClusterIP, or a direct endpoint - and the address in use. | Nothing to fix; it reports the transport behind a Broker that Meshery can see. |

The same diagnostics are available programmatically at
`GET /api/system/controllers/diagnostics?connectionId=<id>`.

## Operating a connection

From the connection in Meshery UI you can:

- **Switch MeshSync mode** (operator ↔ embedded) - see [MeshSync deployment modes](#meshsync-deployment-modes).
- **(Re)deploy or uninstall** Operator / MeshSync / Broker.
- **Reconnect** Meshery Server to the Broker, or run an ad-hoc connectivity test.
- **Rediscover / re-upload** the kubeconfig, or **disconnect / delete** the
  connection (which also tears down the managed port-forward and MeshSync data
  handler).
- **Reset the Meshery Database** if the cached cluster state looks stale.

## Known kinks

- **NATS pod in `CrashLoopBackOff` right after deploy.** Some Meshery Operator
  versions inject the NATS token into `nats.conf` **unquoted**
  (`token: $NATS_TOKEN`). When the randomly-generated token happens to look like a
  number (e.g. starts with digits and contains an `e`), the NATS config parser
  rejects it and the pod crash-loops, so the Broker never becomes ready. Check
  with `kubectl logs -n meshery meshery-nats-0 -c nats`; the fix belongs in the
  operator (quote it: `token: "$NATS_TOKEN"`). Redeploying often generates a token
  that parses.
- **"MeshSync running but not connected to Broker."** Almost always a Broker
  reachability or authentication problem for an out-of-cluster Meshery - see
  [Broker connectivity](#broker-connectivity). The Diagnostics section will say
  `broker_unreachable` with remediation.

## See Also

- [Connections - states and lifecycle](/pr-preview/pr-22672/concepts/logical/connections/)
- [Registering a Connection](/pr-preview/pr-22672/guides/infrastructure-management/registering-a-connection/)
- [Meshery Operator, MeshSync, Broker Troubleshooting](/pr-preview/pr-22672/guides/troubleshooting/meshery-operator-meshsync/)
- [MeshSync](/pr-preview/pr-22672/concepts/architecture/meshsync/) &middot; [Meshery Operator](/pr-preview/pr-22672/concepts/architecture/operator/) &middot; [Meshery Broker](/pr-preview/pr-22672/concepts/architecture/broker/)

<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/new-contributor-intro-vedit-agrawal-bughunterx2101-looking-to-pair-with-a-meshmate-for-meshery-ui-ai-adapter/8172" target="_blank" rel="noopener noreferrer">
            New Contributor Intro: Vedit Agrawal (@BugHunterX2101) | Looking to pair with a MeshMate for Meshery UI &amp; AI Adapter
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-seeking-an-experienced-meshery-contributor-as-a-guide/8166" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate: Seeking an Experienced Meshery Contributor as a Guide
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-to-reactivate-my-slack-account/8168" target="_blank" rel="noopener noreferrer">
            Request to Reactivate My Slack Account
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-frontend-backend-cloud-native-development/8161" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate — Frontend, Backend &amp; Cloud-Native Development
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-for-a-meshmate-active-meshery-contributor/8158" target="_blank" rel="noopener noreferrer">
            Request for a MeshMate — Active Meshery Contributor
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/bug-new-catalog-publish-request-email-renders-with-oversized-social-icons-and-squashed-content-on-mobile/8138" target="_blank" rel="noopener noreferrer">
            [Bug] &#34;New Catalog publish request&#34; email renders with oversized social icons and squashed content on mobile
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/how-to-present-an-issue-pr-in-a-development-meeting-feedback-wanted/8107" target="_blank" rel="noopener noreferrer">
            How to Present an Issue/PR in a Development Meeting: feedback wanted
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newbie-here-excited-to-join/8052" target="_blank" rel="noopener noreferrer">
            Newbie here -- Excited to join
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
