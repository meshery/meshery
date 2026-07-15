---
title: Configuring Meshery Operator, MeshSync, and Broker
categories: [infrastructure]
weight: -4
description: Tune how Meshery deploys and operates Meshery Operator, MeshSync, and Meshery Broker on your clusters - deployment modes, discovery scope, secret handling, service networking, and the behavioral impact of every setting.
---

Every Kubernetes cluster that Meshery manages is served by up to three components: [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}), [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}), and [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}). Their defaults work out of the box, but production clusters routinely need tuning: narrowing discovery scope, keeping Secret contents out of Meshery's database, exposing the Broker across network boundaries, or skipping in-cluster components entirely.

This guide walks through the full configuration surface of all three components: what each setting does, how to change it, its default, and the behavioral impact of changing it. It is written for platform engineers operating real clusters with `kubectl`, `mesheryctl`, and the Meshery UI.

## How the three components fit together

Meshery Server does not watch your clusters directly. It delegates that work:

1. **Meshery Server** establishes a [Connection]({{< ref "concepts/logical/connections/index.md" >}}) to each cluster and, in operator mode, installs Meshery Operator into it.
2. **Meshery Operator** runs in the cluster and reconciles two custom resources: a `Broker` and a `MeshSync`. It turns the `Broker` resource into a NATS StatefulSet and Services, turns the `MeshSync` resource into the MeshSync Deployment, and injects the Broker's address into MeshSync as the `BROKER_URL` environment variable.
3. **MeshSync** discovers the cluster's resources with Kubernetes informers and publishes each discovered object to the Broker on the NATS subject `meshery.meshsync.core`.
4. **Meshery Server** subscribes to the Broker and persists the stream into its database, giving you a continuously refreshed view of the cluster.

One Operator, one MeshSync, and one Broker serve each managed cluster. For architectural depth, see the concept pages for the [Operator]({{< ref "concepts/architecture/operator/index.md" >}}), [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}), and [Broker]({{< ref "concepts/architecture/broker/index.md" >}}); this guide focuses on configuration and its consequences.

## Choosing a deployment mode

The first and most consequential choice is *where MeshSync runs*. Meshery supports two deployment modes per Kubernetes connection:

- **Operator mode** - Meshery Server installs Meshery Operator into the cluster. The Operator deploys the Broker and MeshSync. MeshSync streams events to the Broker, and Meshery Server subscribes to the Broker over NATS.
- **Embedded mode** (the default) - MeshSync runs as a library inside the Meshery Server process. Nothing is installed into the managed cluster, and no Broker is involved: discovery events flow to the server over in-process channels.

### Behavioral impacts of each mode

| Consideration | Operator mode | Embedded mode |
| --- | --- | --- |
| Cluster footprint | Operator, MeshSync, and Broker (NATS) workloads in the `meshery` namespace | None |
| Cluster permissions needed | Enough to install the Operator and its RBAC | Read access via the connection's kubeconfig |
| Transport | NATS, via Meshery Broker | In-process channels inside Meshery Server |
| Network requirement | Meshery Server must be able to reach the Broker endpoint (see [Broker service networking](#broker-service-networking)) | Meshery Server must be able to reach the Kubernetes API server |
| Discovery scope configuration | `MeshSync` custom resource `watch-list` | Built-in default watch list, unless the MeshSync CRD and custom resource exist in the cluster (see [note below](#watch-list-in-embedded-mode)) |
| Resource cost | Paid by the managed cluster | Paid by the Meshery Server host; each embedded connection adds informer caches to the server process |
| Survives Meshery Server restarts | In-cluster discovery keeps running; the server re-subscribes on reconnect | Discovery restarts with the server |

### Where the mode is set, and precedence

The deployment mode is resolved per connection, in this order:

1. **Per-connection setting** - the `meshsync_deployment_mode` entry in the connection's metadata. You set this in the UI when [importing a cluster]({{< ref "guides/infrastructure-management/registering-a-connection.md#meshsync-deployment-mode" >}}) and can change it later from the Connections page (open the connection's action menu and choose **Configure**). Switching modes redeploys MeshSync accordingly: operator to embedded undeploys the Operator and starts the in-server routine; embedded to operator stops the in-server routine and installs the Operator.
2. **Server-wide default** - the `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` environment variable on Meshery Server (`operator` or `embedded`). Connections without an explicit per-connection mode use this value.
3. **Built-in default** - `embedded`.

Additionally, setting `DISABLE_OPERATOR=true` on Meshery Server prevents the server from deploying Meshery Operator to any cluster, regardless of connection settings. Use it for environments where installing cluster components is prohibited.

{{% alert color="info" title="Precedence in one line" %}}
Connection metadata overrides the server-wide default, which overrides the built-in default: <code>meshsync_deployment_mode</code> (per connection) &gt; <code>MESHSYNC_DEFAULT_DEPLOYMENT_MODE</code> (server) &gt; <code>embedded</code> (built in).
{{% /alert %}}

Two lifecycle behaviors are worth knowing when you change modes or connections:

- **Disconnecting a connection undeploys only what Meshery installed** (Operator, Broker, MeshSync). Your workloads are never touched, and the discovered inventory is retained in Meshery's database.
- **Deleting a connection flushes its discovered inventory**, and only when it is the last connection referencing that cluster - multiple kubeconfig contexts pointing at the same cluster are refcounted.

## Configuring Meshery Operator

Meshery Operator itself has little to configure by design: it is a controller whose job is to reconcile the `Broker` and `MeshSync` custom resources described below.

| Setting | Mechanism | Default | Behavioral impact |
| --- | --- | --- | --- |
| Deploy or do not deploy | Per-connection deployment mode; Settings page Operator switch; `DISABLE_OPERATOR=true` on the server | Deployed for operator-mode connections | Without the Operator, no Broker or MeshSync runs in-cluster; embedded mode (or no discovery) applies |
| Operator version | Tracks the Meshery Server release via the `meshery-operator` Helm chart | Matches your server version | Upgrading Meshery Server upgrades the Operator; manual operator upgrades on server-managed clusters are reverted by the server's reconciliation |
| Namespace | Fixed | `meshery` | Operator, Broker, and MeshSync objects live in the `meshery` namespace |

Meshery Server installs the Operator from the `meshery-operator` Helm chart at [meshery.io/charts](https://meshery.io/charts), pinned to the chart version matching the server release. To upgrade the Operator, upgrade Meshery Server; see [How Meshery Server manages Meshery Operator]({{< ref "installation/upgrades/index.md#how-meshery-server-manages-meshery-operator" >}}).

You can also toggle the Operator per cluster without disconnecting the cluster: the Meshery Operator section of **Settings** in Meshery UI provides an on/off switch, and `kubectl` shows you what the Operator has deployed:

```bash
kubectl -n meshery get deploy meshery-operator
kubectl -n meshery get brokers,meshsyncs
```

## Configuring MeshSync

MeshSync's configuration spans three mechanisms: fields on the `MeshSync` custom resource, environment variables on the MeshSync Deployment, and CLI flags. The table below is the complete surface as of current releases; the sections that follow explain each in detail.

| Setting | Mechanism | Default | Takes effect |
| --- | --- | --- | --- |
| Watched resource types and event types | `MeshSync` CR `spec.watch-list` (`whitelist` or `blacklist`) | Built-in watch list of common resource types | After a MeshSync pod restart |
| Namespaces included in output | `--outputNamespaces` CLI flag | All namespaces | Rollout (args change restarts pods) |
| Resource types included in output | `--outputResources` CLI flag | All watched types | Rollout (args change restarts pods) |
| Secret value redaction | `MESHSYNC_REDACT_SECRETS` env var | Off | Rollout (env change restarts pods) |
| Broker content deduplication | `MESHSYNC_BROKER_CONTENT_DEDUP` env var | Off | Rollout (env change restarts pods) |
| Log verbosity | `DEBUG` env var | Info level | Rollout (env change restarts pods) |
| Image version | `MeshSync` CR `spec.version` | `stable-latest` | Operator rolls the Deployment |
| Replica count | `MeshSync` CR `spec.size` (1-10) | `1` | Operator scales the Deployment |
| Broker to publish to | `MeshSync` CR `spec.broker.native` or `spec.broker.custom.url` | The cluster's `meshery-broker` | Operator re-reconciles `BROKER_URL` |

{{% alert color="dark" title="Discovery is read-only; exec and log streaming are not" %}}
MeshSync's discovery never creates, modifies, or deletes anything in your cluster - it only lists and watches. Two <em>on-request</em> operations served over the Broker go further: pod log streaming (reads <code>pods/log</code>) and interactive exec sessions (opens a shell in a target container via <code>pods/exec</code>). These power Meshery's log viewer and interactive terminal, run only when a client requests them, and are never part of discovery (file/snapshot mode does not start them at all). Factor them into your RBAC review if your security posture distinguishes read from exec.
{{% /alert %}}

### Scoping discovery with the watch-list

The watch-list is MeshSync's primary discovery-scope control. It lives on the `MeshSync` custom resource (`meshery-meshsync` in the `meshery` namespace) under `spec.watch-list`, and it is read by the MeshSync agent itself at startup. Note that the Operator deploys MeshSync but does not act on the watch-list; the agent reads its own custom resource, which is why changes require a pod restart rather than waiting for operator reconciliation.

The watch-list holds exactly one of two keys:

- **`whitelist`** - a JSON array of objects, each naming a resource type and the event types to publish for it. Only listed types are watched.
- **`blacklist`** - a JSON array of resource-type strings to exclude from MeshSync's default watch list. Everything else in the default list is watched.

Resource types are identified as `<plural>.<version>.<group>`; core-group resources end with a trailing dot (for example, `pods.v1.` or `deployments.v1.apps`). Event types are `ADDED`, `MODIFIED`, and `DELETED`.

To exclude noisy or uninteresting types, blacklist them:

```bash
kubectl -n meshery edit meshsync meshery-meshsync
```

```yaml
spec:
  watch-list:
    data:
      blacklist: '["events.v1.","replicasets.v1.apps"]'
```

To watch only a specific set of types (and only specific events per type), use a whitelist instead:

```yaml
spec:
  watch-list:
    data:
      whitelist: '[{"Resource":"namespaces.v1.","Events":["ADDED","MODIFIED","DELETED"]},{"Resource":"deployments.v1.apps","Events":["ADDED","MODIFIED","DELETED"]},{"Resource":"pods.v1.","Events":["MODIFIED"]}]'
```

Then restart MeshSync for the new scope to take effect:

```bash
kubectl -n meshery rollout restart deploy/meshery-meshsync
```

Behavioral details worth knowing:

- **Whitelist and blacklist are mutually exclusive.** Setting both, or supplying a list that fails to parse, is a configuration error: MeshSync logs a warning and runs with its full built-in pipeline set - it watches *more*, not less, when the watch-list is invalid. Check MeshSync's logs after changing the watch-list.
- **An explicit whitelist is authoritative.** If a whitelist matches nothing, MeshSync watches nothing; it does not silently fall back to defaults.
- **Custom resources are handled dynamically.** Independent of the watch-list, MeshSync watches the cluster's CustomResourceDefinitions and rebuilds its discovery pipeline when a CRD is added or removed, so newly installed CRD-backed resource types begin to be tracked without a restart.
- **The default watch list covers common types** such as namespaces, pods, deployments, services, nodes, secrets, configmaps, persistent volumes, and more. When no `meshery-meshsync` custom resource exists at all (standalone and typical embedded-mode runs), MeshSync uses this built-in default list.
- **Secrets are watched by default** in every configuration path. See [Redacting Secret contents](#redacting-secret-contents).

<a id="watch-list-in-embedded-mode"></a>
{{% alert color="info" title="Watch-list in embedded mode" %}}
MeshSync looks for its watch-list in the <code>meshery-meshsync</code> custom resource of the target cluster, in <em>any</em> deployment mode. In embedded mode against a cluster that has never run Meshery Operator, that custom resource (and its CRD) does not exist, so MeshSync uses its built-in default watch list. To customize discovery scope for an embedded-mode connection, install the MeshSync CRD and create the <code>meshery-meshsync</code> resource in the target cluster.
{{% /alert %}}

### Filtering output by namespace or resource type

Two CLI flags narrow what MeshSync emits, and they behave differently:

- `--outputResources` - a comma-separated list of resource types to include (for example, `pod,deployment,service`). Matching is case-insensitive, and singular or plural forms both match. This filter is applied *before informers are registered*: excluded resource types are never watched at all, so it reduces API server watch load and MeshSync memory as well as publish volume. Functionally it is a flag-based counterpart to a watch-list whitelist.
- `--outputNamespaces` - a comma-separated list of namespaces to include (for example, `default,production`). Namespace matching is case-sensitive. This is a *publish-time* filter: informers still watch all namespaces, and events from other namespaces are dropped just before writing. Note that when this filter is set, cluster-scoped resources (nodes, namespaces, persistent volumes, and so on) are also filtered out of the output, because they belong to no namespace in the list.

These flags are most commonly used with MeshSync's file output mode and the [`kubectl meshsync snapshot`]({{< ref "extensions/extensions/kubectl-meshsync-snapshot/index.md" >}}) plugin. In operator mode, the MeshSync Deployment runs the binary without arguments by default; to apply output filters there, add container args to the Deployment:

```bash
kubectl -n meshery patch deploy meshery-meshsync --type json -p '[
  {"op":"add","path":"/spec/template/spec/containers/0/args",
   "value":["--outputNamespaces","default,production"]}
]'
```

{{% alert color="dark" title="The Operator owns the MeshSync Deployment" %}}
Meshery Operator reconciles the MeshSync Deployment with Server-Side Apply. Fields the Operator manages (image, command, resources, probes) are converged back to the declared state of the <code>MeshSync</code> custom resource, so do not hand-edit those. Fields the Operator does not manage, such as extra environment variables and container args, persist across reconciliation. Prefer expressing intent through the custom resource wherever a field exists for it.
{{% /alert %}}

### Redacting Secret contents

MeshSync watches `secrets.v1.` by default, and Secret objects are published to the Broker and persisted in Meshery's database with their payloads intact. In security-sensitive environments you have two options:

1. **Do not discover Secrets at all**: blacklist `secrets.v1.` in the [watch-list](#scoping-discovery-with-the-watch-list).
2. **Discover Secrets but redact their values**: set `MESHSYNC_REDACT_SECRETS=true` on the MeshSync Deployment.

```bash
kubectl -n meshery set env deploy/meshery-meshsync MESHSYNC_REDACT_SECRETS=true
```

With redaction enabled, MeshSync replaces every value in a Secret's `data` and `stringData` fields with the placeholder `[REDACTED]` before publishing. Keys are preserved, so Meshery still shows that a Secret exists and which keys it defines - only the sensitive values are withheld. ConfigMaps are not affected. Redaction is off by default; accepted truthy values are `true` and `1`. Changing an environment variable on the Deployment triggers a rollout automatically. (Requires MeshSync v1.0.2 or newer; older images ignore the variable.)

### Deduplicating Broker traffic

Kubernetes controllers sometimes rewrite objects without meaningful change, and MeshSync republishes on each such event. MeshSync already suppresses updates whose `resourceVersion` has not changed; for clusters where objects churn with identical content under new resource versions, an additional opt-in deduplicator is available:

```bash
kubectl -n meshery set env deploy/meshery-meshsync MESHSYNC_BROKER_CONTENT_DEDUP=true
```

When enabled, MeshSync keeps an in-memory content hash (SHA-256) per resource UID on the broker output path and skips publishing an `ADDED`/`MODIFIED` event whose payload is byte-identical to the last published one. `DELETED` events are always published (and evict the UID from the cache, bounding memory). The full object is always sent when a publish does happen; the wire format never changes to a delta.

It is off by default for a reason: the dedup cache persists across informer resyncs, and a resync is also a recovery path that intentionally republishes everything. With dedup enabled, a recovery resync will not republish objects that have not changed since they were last published. Enable it when reduced Broker and database churn is worth that trade-off. Accepted truthy values follow Go's `strconv.ParseBool` (`true`, `1`, `t`, and so on). (Requires MeshSync v1.0.2 or newer; older images ignore the variable.)

### Debug logging

Set `DEBUG=true` (or `1`) on the MeshSync Deployment to switch from info-level to debug-level logging, including caller information:

```bash
kubectl -n meshery set env deploy/meshery-meshsync DEBUG=true
```

Debug logs show discovery pipeline construction, informer registration per resource type, and broker connectivity attempts - the first place to look when resources are missing from Meshery's view. Remember to unset it afterward (`kubectl -n meshery set env deploy/meshery-meshsync DEBUG-`); debug volume on a busy cluster is substantial.

### Version and replica count

The `MeshSync` custom resource declares the desired MeshSync version and scale:

```yaml
apiVersion: meshery.io/v1alpha1
kind: MeshSync
metadata:
  name: meshery-meshsync
  namespace: meshery
spec:
  version: stable-latest   # image tag for meshery/meshsync
  size: 1                  # replicas, 1-10
  broker:
    native:
      name: meshery-broker
      namespace: meshery
```

- **`spec.version`** maps to the container image `meshery/meshsync:<version>`. The default is `stable-latest`. Tags ending in `-latest` are pulled on every pod start (`imagePullPolicy: Always`); pinned tags are pulled only if not present. Pin a specific version in production if you need to control exactly when MeshSync behavior changes.
- **`spec.size`** sets Deployment replicas, validated to the range 1-10, defaulting to 1. MeshSync replicas do not coordinate: there is no leader election or work sharding, so every replica watches the entire cluster and publishes its own copy of every event, multiplying API server watch load, Broker traffic, and Meshery Server's ingest and database write load without adding discovery capacity. Keep `size: 1` for normal operation; the [work-queue design](#what-is-coming) is the roadmap direction for scaling discovery throughput.
- **`spec.broker`** selects where MeshSync publishes: `native` points at a `Broker` resource by name and namespace (the Operator resolves its endpoint and injects `BROKER_URL`), while `custom.url` points MeshSync at an externally managed NATS verbatim.

The `MeshSync` resource's `status` reports where MeshSync is publishing (`publishing-to`) and reconciliation `conditions`. MeshSync deliberately does not write its running version back into `spec.version`: the spec is your (and the Operator's) declaration of desired state, and the running version is advertised over the Broker instead.

### Health endpoints and probes

MeshSync serves two HTTP endpoints on port `11000` from the moment the process starts:

- `GET /healthz` - liveness. Returns `200` while the process is alive, even before (or without) a Broker connection.
- `GET /readyz` - readiness. Returns `503` until MeshSync has established its Broker connection, then `200`.

```bash
kubectl -n meshery port-forward deploy/meshery-meshsync 11000:11000 &
curl -sS http://127.0.0.1:11000/healthz
curl -sS -o /dev/null -w '%{http_code}\n' http://127.0.0.1:11000/readyz
```

`/readyz` means "connected to the Broker", not "informer caches fully primed". Right after a restart, MeshSync can be ready while its cluster snapshot is still filling in.

The operator-managed Deployment currently ships an exec-based liveness probe (it periodically runs `./meshery-meshsync -h`) and intentionally no readiness probe, since no Kubernetes Service routes traffic to MeshSync. If you run MeshSync under your own management (outside the Operator), wire HTTP probes to port `11000`: `/healthz` for liveness and, if you want Broker-connection state surfaced in pod conditions, `/readyz` for readiness.

The port is fixed at `11000` and is not configurable.

### Broker connectivity at startup

Before opening its NATS connection, MeshSync verifies the Broker is reachable by polling the Broker's HTTP monitoring endpoint, `http://<broker-host>:8222/connz`, once per second for up to 5 minutes (each probe times out after 5 seconds). While this retry loop runs, `/readyz` reports `503`. If the Broker never becomes reachable, MeshSync exits with an error; in operator mode, Kubernetes restarts the pod with backoff, so a Broker that comes up late is tolerated.

After a successful connection, reconnection is handled by the NATS client automatically (2-second reconnect intervals, up to 60 attempts). None of these timings are configurable.

If MeshSync stays not-ready, check the injected Broker address and its reachability:

```bash
kubectl -n meshery get deploy meshery-meshsync \
  -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="BROKER_URL")].value}{"\n"}'
kubectl -n meshery get broker meshery-broker -o jsonpath='{.status.endpoint}{"\n"}'
```

## Configuring Meshery Broker

Meshery Broker is NATS, deployed by the Operator from the `Broker` custom resource (`meshery-broker` in the `meshery` namespace). Clients authenticate with a token the Operator generates into the `meshery-nats-auth` Secret; MeshSync's `BROKER_URL` receives it automatically.

| Setting | Mechanism | Default | Behavioral impact |
| --- | --- | --- | --- |
| NATS version | `Broker` CR `spec.version` | Bundled NATS release (image `nats:<version>`) | Operator rolls the StatefulSet to the new image |
| Replica count | `Broker` CR `spec.size` (1-10) | `1` | Scales the NATS StatefulSet |
| Service type | `Broker` CR `spec.service.type` | `ClusterIP` | Controls whether the Broker is reachable from outside the cluster |
| Service annotations | `Broker` CR `spec.service.annotations` | None | Cloud load balancer hints, MetalLB pools, internal-LB switches |
| Load balancer class / source ranges | `spec.service.loadBalancerClass`, `spec.service.loadBalancerSourceRanges` | None | LoadBalancer-type only; CRD validation rejects them for other types |
| Advertised external endpoint | `spec.service.externalEndpointOverride` | Derived from the Service | Pins the advertised `host:port` for ingress, gateway, NAT, or air-gapped topologies |

### Broker service networking

The Broker's exposure is declarative and **reconfigurable on a live Broker**: change `spec.service`, and the Operator updates the Service in place (no pod restart), re-derives `status.endpoint`, and re-reconciles MeshSync so it reconnects to the new address.

```bash
kubectl -n meshery patch broker meshery-broker --type merge \
  -p '{"spec":{"service":{"type":"NodePort"}}}'
kubectl -n meshery get broker meshery-broker -o jsonpath='{.status.endpoint}{"\n"}'
```

Which type you need follows from where Meshery Server runs relative to the cluster:

- **Meshery Server in the same cluster**: the default `ClusterIP` is correct. The Broker acquires no external address.
- **Meshery Server outside the cluster** (Docker host, another cluster, Meshery Cloud): the server must reach the Broker's *external* endpoint. Use `NodePort` or `LoadBalancer`, or front the Broker with your own ingress/gateway and pin the advertised address with `spec.service.externalEndpointOverride: <host:port>`.

Service annotations pass load-balancer hints through to your platform - for example, requesting an internal load balancer on AWS, or a MetalLB address pool on bare metal:

```bash
kubectl -n meshery patch broker meshery-broker --type merge \
  -p '{"spec":{"service":{"type":"LoadBalancer","annotations":{"service.beta.kubernetes.io/aws-load-balancer-internal":"true"}}}}'
```

The Operator publishes the resulting addresses on `status.endpoint.internal` (always the in-cluster address and client port `4222`) and `status.endpoint.external` (empty for `ClusterIP`). The full endpoint selection order is documented in the [Operator FAQ]({{< ref "concepts/architecture/operator/index.md#how-does-the-operator-expose-information-about-broker-endpoints" >}}).

{{% alert color="warning" title="Default changed: the Broker is no longer public by default" %}}
Older Meshery Operator releases exposed the Broker as a <code>LoadBalancer</code> Service implicitly. Current releases keep the Broker cluster-internal (<code>ClusterIP</code>) unless <code>spec.service.type</code> requests otherwise. If Meshery Server runs outside the cluster and shows the Broker as unreachable after an upgrade, set <code>spec.service.type</code> to <code>NodePort</code> or <code>LoadBalancer</code>, or set <code>spec.service.externalEndpointOverride</code>.
{{% /alert %}}

Two ports matter: `4222` is the NATS client port that MeshSync and Meshery Server connect to; `8222` is the HTTP monitoring port, which MeshSync probes (`/connz`) before connecting and which the Broker pods use for their own health probes.

### Broker sizing and delivery semantics

Scale the Broker by data volume:

```bash
kubectl -n meshery patch broker meshery-broker --type merge -p '{"spec":{"size":3}}'
```

Messages are held in memory and delivered at-most-once (core NATS): if Meshery Server is disconnected when an event is published, that event is not replayed on reconnect; the state is instead rebuilt by a later resync. Durable, replayable delivery via NATS JetStream is on the [roadmap](#what-is-coming).

## Configuring from Meshery UI

Meshery UI exposes this configuration surface in two places today, with a broader configuration experience in active development.

**Available today:**

- **Deployment mode at import** - the connection wizard lets you choose Operator or Embedded per kubeconfig context when [importing a cluster]({{< ref "guides/infrastructure-management/registering-a-connection.md#meshsync-deployment-mode" >}}).
- **Deployment mode per connection** - the Connections table shows each Kubernetes connection's current mode. Open a connection's action menu, choose **Configure**, select the other mode, and **Apply**; Meshery tears down and re-establishes discovery for that connection in the new mode.
- **Operator switch and diagnostics in Settings** - the Meshery Operator section of **Settings** provides a per-cluster on/off switch for the Operator, ad hoc connectivity tests for the Operator, Broker (NATS), and MeshSync, and database reset/flush actions.

**Landing in an upcoming release** (tracked in [meshery/meshery#20487](https://github.com/meshery/meshery/issues/20487); the names below follow the implementation plan recorded on that issue and may evolve until it ships):

- **Server-wide defaults in Settings** - a new **Operator, MeshSync & Broker** tab on the Settings page, alongside Overview, Adapters, Registry, and Reset, where an administrator declares the default configuration applied to every managed cluster that does not override it: deployment mode, MeshSync version, replicas, watch-list, output filters, secret redaction, content deduplication, and debug logging, plus Broker version, replicas, and service networking.
- **Per-connection overrides on the Connections page** - a **Configure Controllers** action on each Kubernetes connection opens an **Operator, MeshSync & Broker Configuration** modal showing the connection's effective configuration, with a per-field **Inherited (server default)** vs **Override** indicator and a per-field reset back to inherit.
- **Visible precedence and safe apply** - effective configuration resolves as connection override, else server-wide default, else built-in default. For the deployment mode specifically, a default saved in Settings supersedes the `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` environment variable (which remains the deploy-time bootstrap), and `DISABLE_OPERATOR` continues to be honored. Fields whose application restarts MeshSync pods (watch-list, the environment-variable knobs, version) are labeled as such, and Meshery performs the required restart automatically.

Until that ships, the `kubectl` and environment variable mechanisms in this guide are the way to change these settings, and the per-connection deployment mode remains the one setting already manageable end to end from the UI.

## Performance tuning for large clusters

On large or busy clusters, tune in this order:

1. **Narrow the watch-list.** This is the biggest lever: it reduces informer count, API server watch load, MeshSync memory, Broker traffic, and database size all at once. Blacklist high-churn types you do not need (`events.v1.` is a common first cut), or invert to a whitelist of exactly what you care about.
2. **Mind CRD churn.** MeshSync rebuilds its discovery pipeline when a CRD is added or removed - not on every CRD modification - so ordinary controllers that update CRD metadata (for example, cert-manager's CA injector refreshing `caBundle`) do not trigger re-discovery. Installing or uninstalling many CRDs at once, however, causes repeated rebuilds; schedule bulk CRD changes accordingly.
3. **Use the output flags where the watch-list is not available.** `--outputResources` skips informer registration for excluded types entirely, making it a watch-side lever like a whitelist; `--outputNamespaces` trims publish volume (Broker traffic, database size) for namespaced objects while informers keep watching cluster-wide. Both shine in snapshot and standalone runs.
4. **Enable content deduplication** (`MESHSYNC_BROKER_CONTENT_DEDUP=true`) on clusters where controllers rewrite objects without content changes, once you have weighed the [resync trade-off](#deduplicating-broker-traffic).
5. **Let the payload slimming work for you.** MeshSync no longer publishes `metadata.managedFields`, the highest-churn, highest-bulk part of most objects, and suppresses republishes when an object's `resourceVersion` is unchanged. Both are automatic.
6. **Size the components.** The MeshSync container defaults to requests of 500m CPU / 512Mi memory with limits of 4 CPU / 4Gi memory; the initial full-cluster sync is the resource-intensive phase. Scale the Broker (`spec.size`) with event volume. Keep MeshSync at `size: 1` (see [Version and replica count](#version-and-replica-count)).
7. **Place the Broker deliberately.** Keep it `ClusterIP` unless Meshery Server genuinely needs external reachability, and restrict exposure with `loadBalancerSourceRanges` when using a load balancer.

## What applies live, and what needs a restart

| Change | Applies |
| --- | --- |
| `Broker` `spec.service.*` (type, annotations, override) | Live; Service updated in place, endpoint re-derived, MeshSync reconnected |
| `Broker` `spec.size`, `spec.version` | Operator-driven StatefulSet rollout |
| `MeshSync` `spec.version`, `spec.size`, `spec.broker` | Operator-driven Deployment update/rollout |
| `MeshSync` `spec.watch-list` | Requires a MeshSync pod restart (`kubectl -n meshery rollout restart deploy/meshery-meshsync`) |
| Env vars (`DEBUG`, `MESHSYNC_REDACT_SECRETS`, `MESHSYNC_BROKER_CONTENT_DEDUP`) | Setting them on the Deployment triggers an automatic rollout |
| CRDs added/removed in the cluster | Live; MeshSync rebuilds its pipeline automatically |
| Deployment mode (operator/embedded) | Meshery Server undeploys and redeploys MeshSync for that connection |

## Verifying behavior

After any configuration change, confirm the components are healthy and data is flowing:

```bash
# All three components at once
mesheryctl system check --operator

# Custom resource health
kubectl -n meshery get broker meshery-broker -o jsonpath='{.status.conditions}{"\n"}'
kubectl -n meshery get meshsync meshery-meshsync -o jsonpath='{.status.conditions}{"\n"}'

# MeshSync readiness (connected to Broker)
kubectl -n meshery port-forward deploy/meshery-meshsync 11000:11000 &
curl -sS -o /dev/null -w '%{http_code}\n' http://127.0.0.1:11000/readyz
```

In Meshery UI, the connection's MeshSync state should read **CONNECTED**; the meanings of each state, and what to do when they diverge, are covered in the [Meshery Operator, MeshSync, Broker Troubleshooting Guide]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}}).

<a id="what-is-coming"></a>
## What is coming

Several discovery and delivery controls are designed but not yet shipped. They are listed here so you can plan; none of the knobs below exist in current releases.

- **Tiered, user-configurable discovery** ([meshery/meshsync#577](https://github.com/meshery/meshsync/issues/577)) - replaces flat discovery with a funnel: a cheap broad first pass, deeper inspection of what was found, and infrastructure-specific tiers activated only when the relevant platform is detected. Discovery depth becomes a per-resource setting, and watch-list changes are picked up live rather than on restart.
- **Rate-limited work queue** ([meshery/meshsync#578](https://github.com/meshery/meshsync/issues/578)) - a bounded queue between informers and the output path, decoupling a slow Broker from event ingestion, coalescing rapid successive updates, and introducing worker-count and queue-size knobs.
- **Durable delivery via NATS JetStream** ([meshery/meshsync#580](https://github.com/meshery/meshsync/issues/580)) - at-least-once, replayable event delivery with configurable stream retention on the `Broker` resource, closing the gap described in [delivery semantics](#broker-sizing-and-delivery-semantics).
- **Periodic reconciliation** ([meshery/meshsync#581](https://github.com/meshery/meshsync/issues/581)) - an optional timer-driven re-list that detects and publishes only genuine drift (including missed deletes), off by default with an enforced minimum interval.

## Related

- [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}), [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}), and [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) - architecture concepts.
- [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}) - importing clusters and choosing the deployment mode.
<!-- The brownfield guide ships in a parallel PR at this agreed permalink. Convert to a
     ref-shortcode link once it exists on master, since an unresolved ref fails the Hugo build. -->
- [Bringing Existing Infrastructure Under Meshery Management](https://docs.meshery.io/guides/infrastructure-management/managing-existing-infrastructure) - discovering and managing brownfield infrastructure.
- [Meshery Operator CRDs reference]({{< ref "reference/references/meshery-operator-crds.md" >}}) - full `Broker` and `MeshSync` schema reference.
- [Troubleshooting Meshery Operator, MeshSync, and Broker]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}}) - failure scenarios and recovery.
- [`kubectl meshsync snapshot`]({{< ref "extensions/extensions/kubectl-meshsync-snapshot/index.md" >}}) - file-mode cluster snapshots without a Broker.

{{< discuss >}}
