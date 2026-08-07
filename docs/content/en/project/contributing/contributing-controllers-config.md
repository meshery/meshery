---
title: "Contributing to Operator, MeshSync & Broker Settings"
description: >
  The complete surface of settings that govern Meshery Operator, MeshSync and
  Meshery Broker behavior: where each one lives, what it writes, what reads it,
  and what happens when it changes.
categories: [Contributing]
weight: 27
---

Meshery Operator, MeshSync and Meshery Broker are configured through a layered
document, `MesheryControllersConfig`, defined in
[meshery/schemas](https://github.com/meshery/schemas) as the `controllers_config`
construct (`v1alpha1`). This page is the authoritative map of that surface for
contributors: every setting, its storage, its reader, and its observable effect.

Keep it current. A setting added without a row here is a setting nobody can
find, and a row without an end-to-end test is a setting nobody can trust.

## The three layers

Configuration resolves through a fixed precedence chain. Every field is
independently layered - leaving a field unset at one layer defers to the next.

| Precedence | Layer | Stored in | Set from |
| --- | --- | --- | --- |
| 1 (highest) | Per-connection override | `connection.metadata.controllers_config` | Connections table -> row actions -> **Configure Operator, MeshSync & Broker** |
| 2 | Server-wide default | `system_settings` row `controllers_config_defaults` | Settings -> **Operator, MeshSync & Broker** tab |
| 3 (lowest) | Built-in default | Compiled in: `connections.BuiltInControllersConfig()` | Not settable |

`connections.ResolveControllersConfig` returns two documents:

- **merged** - only fields explicitly set at layer 1 or 2. This is what
  propagates to the cluster. Fields absent here are *withdrawn* on the next
  apply, reverting to the operator's or chart's own value.
- **effective** - merged overlaid onto the built-in defaults. This is the
  complete resolved view the API returns to clients.

The UI mirrors the built-in defaults in `BUILT_IN_CONTROLLERS_CONFIG`
(`ui/components/configuration/ControllersConfigForm.tsx`). The two must agree;
they are asserted against each other by unit test.

### Collections merge atomically

Scalars merge per leaf. Collections - `meshsync.watchList`,
`meshsync.outputNamespaces`, `meshsync.outputResources`,
`broker.service.annotations`, `broker.service.loadBalancerSourceRanges` - merge
whole: a layer that sets one replaces the lower layer's value entirely.
Element-wise merging of a whitelist or a namespace filter would produce a scope
neither layer asked for.

## The settings

Every row is a distinct behavior. "Propagates to" names the cluster object the
server writes; "Observable as" names what the user can actually see change.

### Meshery Operator

| Setting | Wire path | Propagates to | Observable as |
| --- | --- | --- | --- |
| Deployment mode | `operator.deploymentMode` | Deploy/undeploy of the operator; MeshSync data pipeline teardown and reattach | Connection detail "MeshSync Deployment Mode"; operator/MeshSync/broker status chips in the header |
| Operator version | `operator.version` | `meshery-operator` Helm release chart version | Operator version in the header status chip; `meshery-operator` pod image |

`deploymentMode` is one of `operator` (Meshery Operator installs MeshSync and
Broker into the cluster) or `embedded` (MeshSync runs in-process inside Meshery
Server; nothing is installed into the cluster). The built-in default is
`embedded`.

### Where `operator.version` is read

`connections.OperatorChartVersionFromControllersConfig` is the single reader.
It feeds `MesheryControllersHelper.operatorDeploymentConfig`, which overlays the
resolved value onto the `controllers.OperatorDeploymentConfig` assembled once at
boot by `models.NewOperatorDeploymentConfig`. That struct is what the meshkit
operator controller handler is constructed with, and its `MesheryReleaseVersion`
is the chart version passed to the `meshery-operator` Helm release.

Two consequences shape the code around it:

- The handler captures the chart version **at construction**, so
  `SetControllersConfig` must run before `AddCtxControllerHandlers`. Every call
  site orders them that way.
- A version changed on a live connection therefore cannot take effect by
  re-applying the document. `MesheryControllersHelper.ReconcileOperatorChartVersion`
  re-attaches the handlers and re-runs the Helm release (installed with
  `UpgradeIfInstalled`, so this is an in-place upgrade). It is a no-op in
  embedded mode, when the operator is disabled server-wide or explicitly
  undeployed for the context, and when the resolved version already matches the
  attached handler's - so a chart-version change never resurrects an operator
  the user turned off.

Leaving the field unset at every layer resolves to `""`, and the boot-time chart
version applies unchanged: a connection on Inherit behaves exactly as it did
before the field was wired up.

### The chart version that reaches Helm is pinned, not assumed

`operatorDeploymentConfig` returns the version that was *asked for*. Nothing yet
guarantees such a chart exists, and two everyday cases guarantee it does not:

- Chart publishing is decoupled from Meshery Server releases and trails them, so
  a current server routinely names a chart version the repository has not
  published. Helm then fails with a chart-not-found error and the operator never
  deploys.
- Charts below `models.MinimumOperatorChartVersion` are published but cannot run
  (retired `kube-rbac-proxy` image, no webhook certificate). Because the derived
  version tracks the server's own release, an old server would install one of
  those on every cluster it is ever pointed at, forever.

`MesheryControllersHelper.pinnedOperatorDeploymentConfig` is therefore the only
deployment config allowed to reach Helm. It lists what the repository publishes
(`helpers/utils.PublishedChartVersions`, a TTL-cached `index.yaml` read) and
hands that to `models.ResolveOperatorChartVersion`, which distinguishes two
sources:

- **Derived** (`OperatorChartVersionDerived`) - the boot-time server release.
  Nobody chose it, so it may be corrected: unpinned or unpublished falls back to
  the newest published chart, and anything below the floor is raised to the
  oldest published chart at or above the floor. Every correction returns a
  one-sentence reason that is logged and emitted as a warning event, so no
  substitution is silent.
- **Requested** (`OperatorChartVersionRequested`) - an explicit
  `operator.version`. Never substituted. A moving tag or an unpublished version
  fails through `setOperatorError`, which the connection-diagnostics API and the
  events feed surface to the user. The floor does not apply, so deliberately
  pinning an old chart still works.

Membership in the published set is decided by semver comparison, not by string
equality, because Helm treats `1.0.64` and `v1.0.64` as the same version.
Resolution always returns the repository's own spelling of the matched release -
that is the string handed to Helm and the string
`attachedOperatorChartVersion` is compared against - and a spelling difference
is not a substitution, so it produces no reason.

An unreadable repository index fails rather than guessing: without the index
there is no way to know which versions exist. Operator lifecycle is then
withheld for that connection (no `MesheryOperator` handler is attached) while
the Broker and MeshSync handlers, which need no chart version, still attach.
`collectControllersStatus` synthesizes an `OPERATOR` row carrying `UNKOWN`
whenever no operator handler is attached, so the Operator card stays visible
rather than disappearing from a snapshot the client applies wholesale; the
reason is carried by the `operator_deploy_failed` diagnostic, which reads
`GetOperatorError`.

`NewOperatorDeploymentConfig` consequently leaves the chart version empty for an
unstamped build instead of asking the GitHub releases API for the newest server
tag. The newest GitHub release is not the newest published chart, so that call
spent a rate-limited request to produce a version that frequently does not exist
in the chart repository; an empty version resolves to the newest published chart,
which is what was wanted.

`ReconcileOperatorChartVersion` compares **pinned against pinned**.
`attachedOperatorChartVersion` always records a published version, so comparing
the raw request against it would differ forever whenever the request is being
substituted - re-running the Helm upgrade on every reconcile.

### MeshSync

| Setting | Wire path | Propagates to | Observable as |
| --- | --- | --- | --- |
| Version | `meshsync.version` | MeshSync CR `spec.version` | MeshSync version in the header status chip; MeshSync pod image |
| Replicas | `meshsync.replicas` | MeshSync CR `spec.size` | MeshSync pod count |
| Watch list | `meshsync.watchList` | MeshSync CR `spec.watch-list.data.{whitelist,blacklist}` + rolling restart | Which resources appear under Dashboard -> Resources |
| Output namespaces | `meshsync.outputNamespaces` | MeshSync Deployment arg `--outputNamespaces` | Which namespaces' resources are published |
| Output resources | `meshsync.outputResources` | MeshSync Deployment arg `--outputResources` | Which kinds are published |
| Secret redaction | `meshsync.redactSecrets` | MeshSync Deployment env `MESHSYNC_REDACT_SECRETS` | Secret values redacted in Dashboard -> Resources |
| Broker content dedup | `meshsync.brokerContentDedup` | MeshSync Deployment env `MESHSYNC_BROKER_CONTENT_DEDUP` | Fewer duplicate MeshSync events in the notification center |
| Debug logging | `meshsync.debugLogging` | MeshSync Deployment env `DEBUG` | MeshSync pod logs |

`watchList` accepts **at most one** of `whitelist` or `blacklist`; both set is
rejected. MeshSync reads its watch-list at startup only, so a watch-list change
also stamps `meshery.io/restarted-at` on the MeshSync Deployment pod template,
which rolls the pods.

Be precise about what that guarantee covers. Every setting in this section that
propagates to the Deployment - the three env vars and the two output-filter args
- lives in the **pod template**, so changing any of them triggers an ordinary
rolling update by definition; that is Kubernetes, not something Meshery chooses.
What the annotation handling guarantees is narrower: `meshery.io/restarted-at`
is refreshed *only* when the watch-list changes, and the previously-applied
value is otherwise carried forward unchanged, so an apply that changes nothing
in the pod template does not roll pods gratuitously. A watch-list change is the
only thing that forces a restart when the template is otherwise identical.

Version, replicas and watch-list are operator-mode only: they live on the
MeshSync custom resource, which embedded-mode clusters never install.

### Meshery Broker

| Setting | Wire path | Propagates to | Observable as |
| --- | --- | --- | --- |
| Version | `broker.version` | Broker CR `spec.version` | Broker version in the header status chip |
| Replicas | `broker.replicas` | Broker CR `spec.size` | NATS statefulset replica count |
| Service type | `broker.service.type` | Broker CR `spec.service.type` | `meshery-broker` Service type; broker endpoint |
| Service annotations | `broker.service.annotations` | Broker CR `spec.service.annotations` | Annotations on the broker client Service |
| Load balancer class | `broker.service.loadBalancerClass` | Broker CR `spec.service.loadBalancerClass` | Which LB controller provisions the Service |
| Load balancer source ranges | `broker.service.loadBalancerSourceRanges` | Broker CR `spec.service.loadBalancerSourceRanges` | Which CIDRs may reach the broker |
| External endpoint override | `broker.service.externalEndpointOverride` | Broker CR `spec.service.externalEndpointOverride` | The advertised broker endpoint (connection diagnostics) |

`loadBalancerClass` and `loadBalancerSourceRanges` are valid only when
`service.type` is `LoadBalancer`; the server rejects them otherwise, and the
form clears them when the effective type changes away from `LoadBalancer`.
Service changes reconcile in place and must not restart broker pods.

## What each deployment mode can apply

Meshery Operator manages MeshSync and Meshery Broker, so the deployment mode is
not one setting among many - it decides which of the others can reach anything.
In `embedded` mode Meshery Server runs MeshSync in-process and installs nothing
into the cluster: no operator release, no MeshSync Deployment, no Broker, no
`meshery.io` custom resources.

| Setting | `operator` | `embedded` |
| --- | --- | --- |
| `operator.deploymentMode` | applies | applies |
| `operator.version` | applies | inert - no operator release is installed |
| `meshsync.outputNamespaces`, `meshsync.outputResources` | applies | applies - passed to the in-process `libmeshsync` run |
| `meshsync.version`, `meshsync.replicas`, `meshsync.watchList` | applies | inert - MeshSync CR absent |
| `meshsync.redactSecrets`, `meshsync.brokerContentDedup`, `meshsync.debugLogging` | applies | inert - these are env on the MeshSync Deployment; embedded MeshSync takes them from the Meshery Server process environment |
| every `broker.*` | applies | inert - no Broker on the cluster |

The server states this at apply time through
`ControllersConfigApplyResult.Skipped`. The UI states it *before* the save:
`ui/components/configuration/deploymentMode.ts` is the single client-side
statement of the same structure (`takesEffectIn`), and the editors are governed
by it - the per-connection editor renders a setting the effective mode cannot
apply as inert, marks it **Not applied**, and says why in the form body rather
than in a tooltip. Values already stored for such a setting are kept and shown
as dormant with a control to clear them; they become live again if the
connection moves to `operator` mode.

The server-wide defaults editor annotates instead of disabling: its mode is only
what *inheriting* connections get, and a connection that overrides the mode to
`operator` uses every value stored there.

That decision needs the mode the connection actually runs, which is why
`GET /api/integrations/connections/{connectionId}/controllers/config` resolves
`effective.operator.deploymentMode` through
`connections.ResolveConnectionControllersConfig` rather than by merging the two
editable layers: a connection whose mode comes from the materialized cache or
from `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` would otherwise be described as running
the built-in `embedded` mode, and every applicability decision made from that
value would be wrong. `merged` is deliberately left untouched - it is what
propagates to the cluster and must keep carrying only explicitly-set fields.

## Server-side apply and withdrawal

Every cluster write uses server-side apply under the `meshery-server` field
manager (`server/models/controllers_config_apply.go`). The applied document
always describes the *complete* set of fields Meshery Server owns, so:

- Setting a field takes ownership of it.
- Clearing it at every layer withdraws it on the next apply, reverting to the
  operator's or chart's own value.

The operator applies the same objects under its own field manager and never sets
these env names, args, or the restart annotation, so ownership stays disjoint.

Absent custom resources or Deployment (operator not deployed yet, embedded-mode
cluster) are **skipped and reported** in `ControllersConfigApplyResult.Skipped`,
not treated as errors: the configuration re-applies when the connection
reconnects.

## Non-layered settings that govern the same components

These are not part of `controllers_config` but govern the same three components.
Contributors changing controller behavior must account for them.

| Setting | Where | Storage | Notes |
| --- | --- | --- | --- |
| `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` | Server env / `viper` | Process config | Fallback when no layer and no materialization sets a deployment mode. Logged at boot. |
| MeshSync deployment mode (wizard) | Connection Wizard -> **MeshSync Mode** step | `controllers_config.operator.deploymentMode` | Reconfigure-time picker. Writes the override through `connections.SetDeploymentModeOverride` via `POST /api/integrations/connections/{id}/actions` (`setMeshsyncMode`). |
| MeshSync deployment mode (per-context, at import) | Connection Wizard -> kubeconfig context review | Same field, per context | Written only when the user actually picks one; otherwise the connection inherits. |
| Operator status | GraphQL `changeOperatorStatus` | Cluster state | Deploys/undeploys the operator directly, bypassing the layered document. |

### One store for the deployment mode

`connection.metadata.meshsync_deployment_mode` used to carry two different
facts under one key - the user's explicit per-connection choice *and* a
materialized cache of the resolved mode - and that ambiguity was the root cause
of two separate defects. The two are now separated
(`server/models/connections/deployment_mode_resolution.go`):

- **`controllers_config.operator.deploymentMode`** is the only store of the
  explicit choice. Every entry point - the wizard's MeshSync Mode step, the
  kubeconfig import picker, the controllers editor - writes it through
  `connections.SetDeploymentModeOverride`, so no two controls can disagree
  about what a connection is set to.
- **`meshsync_deployment_mode`** is only the materialization of the *resolved*
  mode, written through `connections.MaterializeMeshsyncDeploymentMode` and
  read by the consumers that predate the layered document (the connection state
  machine, the header status chips, the kubeconfig flows).

`connections.ResolveDeploymentMode` is the single decision point, and it reports
the layer it resolved from:

| Precedence | Layer | Reported as |
| --- | --- | --- |
| 1 (highest) | Layered document (per-connection override over server-wide default) | `layeredConfig` |
| 2 | Materialized `meshsync_deployment_mode` | `legacyConnectionMetadata` |
| 3 | `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` | `serverEnvDefault` |
| 4 (lowest) | Compiled-in default (`embedded`) | `builtIn` |

Ranking the materialization *below* the layered document is what lets a
server-wide default reach an existing connection: every Kubernetes connection
has that key written at registration, so while it sat on top no fan-out could
ever change a mode. It survives as a compatibility floor for connections
registered before the layered document existed - their recorded mode is kept
when, and only when, no layer sets one.

A server-wide default change is not just a document to re-apply.
`reconcileInheritedDeploymentMode` persists the refreshed materialization and
drives the same undeploy/redeploy path the wizard uses, so an inheriting
connection actually switches modes on the cluster.

## API

| Method | Path | Handler |
| --- | --- | --- |
| `GET` | `/api/system/controllers/config` | `GetControllersDefaultConfig` |
| `PUT` | `/api/system/controllers/config` | `UpdateControllersDefaultConfig` |
| `GET` | `/api/integrations/connections/{connectionId}/controllers/config` | `GetConnectionControllersConfig` |
| `PUT` | `/api/integrations/connections/{connectionId}/controllers/config` | `UpdateConnectionControllersConfig` |

Both `PUT`s validate before persisting (`connections.ValidateControllersConfig`)
and return the stored document. The connection `PUT` additionally applies the
resolved configuration to that connection's cluster; the system `PUT` fans the
re-apply out to every tracked Kubernetes connection.

Writes are asynchronous with respect to the cluster: the override persists even
if the apply fails, and the failure surfaces as an event rather than an HTTP
error, because the configuration re-applies on the next connect.

## Known gaps

Tracked here so they are not rediscovered. Each is a defect, not a design.

1. **The apply result is invisible.** `ApplyControllersConfigToCluster` returns
   which custom resources it patched, whether MeshSync was restarted, and a
   `Skipped[]` list with human-readable reasons. That result only reaches event
   metadata; both editors report "saved" either way, so a save that changed
   nothing on the cluster looks identical to one that changed everything.

## Test coverage

Unit coverage that exists today:

- `server/models/connections/controllers_config_test.go` - precedence, atomic
  collection merge, metadata round-trip, validation, and the single reader of
  `operator.version`.
- `server/models/connections/deployment_mode_resolution_test.go` - deployment
  mode precedence (a server-wide default reaching an inheriting connection),
  the wizard/editor convergence on one store, and the effective document
  reporting the mode a connection actually runs.
- `server/models/operator_chart_version_test.go` - `operator.version` selecting
  the operator Helm chart version, the layering it resolves through, and the
  cases where a chart-version reconcile must not touch the cluster.
- `server/models/operator_chart_pinning_test.go` - pinning the resolved version
  to one the repository publishes: the floor, the unpublished-release fallback,
  moving tags never reaching Helm, explicit requests failing loudly, and the
  pinned-against-pinned reconcile comparison.
- `server/helpers/utils/helm_chart_repo_test.go` - the `index.yaml` read: chart
  version extraction, structured failures, TTL reuse, that failures are not
  cached, that the lock is not held across the fetch, that concurrent misses
  single-flight, and that the cached catalogue is never handed out by
  reference.
- `server/handlers/controllers_status_handler_test.go` - that the `OPERATOR`
  row is present in the status snapshot even when no operator handler is
  attached.
- `ui/components/configuration/__tests__/deploymentMode.test.ts` - which
  settings each deployment mode can apply, and how the mode governing each
  editor is resolved and attributed to a layer.
- `ui/components/configuration/__tests__/ControllersConfigForm.test.tsx` - the
  rendered editor: tri-state inherit/override, the conditional LoadBalancer-only
  fields, and the deployment-mode gating (inert-and-explained on a connection,
  annotated-but-live on the server-wide defaults).
- `server/models/controllers_config_apply_test.go` - cluster propagation with
  fake clients: per-setting custom-resource and Deployment patch contents,
  withdrawal when a field is cleared at every layer,
  restart-only-on-watch-list-change, and the `Skipped[]` reporting.
- `server/handlers/controllers_config_handlers_test.go` - the four endpoints:
  validation rejection, the layered response shape, clearing an override, and
  that an unreadable override is skipped rather than reconciled.

End-to-end coverage lives in `ui/tests/e2e/controllers-config.spec.ts`, keyed to
the **Operator, MeshSync & Broker Settings** Test Plan Test Group via the Allure
`testGroup` label - the same mechanism the Connection Lifecycle report uses (see
[Contributing to Meshery's CLI tests]({{< relref "cli/tests.md" >}})). It covers
every wire path in the tables above plus the precedence chain, the inherit
round-trip, validation rejection, and the mode gating.

What remains uncovered:

- Cluster propagation end to end. The spec's watch-scope and broker-service case
  self-skips when no Kubernetes cluster is reachable, so on an infra-less run it
  reports "skipped" rather than passing vacuously. The propagation logic itself
  is covered by the unit test above; what is untested is the real round trip
  against a live operator.

Before adding a setting, add its row above, its propagation assertion, and its
end-to-end case. A setting whose effect a user cannot observe is a defect; so is
one that fails silently.
