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
| Operator version | `operator.version` | Not propagated - see [Known gaps](#known-gaps) | Nothing |

`deploymentMode` is one of `operator` (Meshery Operator installs MeshSync and
Broker into the cluster) or `embedded` (MeshSync runs in-process inside Meshery
Server; nothing is installed into the cluster). The built-in default is
`embedded`.

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
| `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` | Server env / `viper` | Process config | Last-resort fallback when no layer sets a deployment mode. Logged at boot. |
| MeshSync deployment mode (wizard) | Connection Wizard -> **MeshSync Mode** step | `connection.metadata.meshsync_deployment_mode` | Registration-time and reconfigure-time picker. Writes the legacy key directly via `PUT /api/integrations/connections/{id}/meshsync-mode`. |
| MeshSync deployment mode (per-context, at import) | Connection Wizard -> kubeconfig context review | Same key, per context | Set on every Kubernetes connection at registration. |
| Operator status | GraphQL `changeOperatorStatus` | Cluster state | Deploys/undeploys the operator directly, bypassing the layered document. |

`meshsync_deployment_mode` is the **legacy materialization** of
`operator.deploymentMode`. Every consumer that predates the layered document -
the connection state machine, the header status chips, the kubeconfig flows -
reads that key, so `UpdateConnectionControllersConfig` rewrites it on every
update using the full precedence chain.

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

1. **`operator.version` is accepted but never propagated.** The form offers it
   and the server stores and validates it, but `ApplyControllersConfigToCluster`
   writes only the MeshSync and Broker custom resources and the MeshSync
   Deployment. Nothing reads `operator.version`. A user who sets it sees no
   effect and no error.

2. **A server-wide `operator.deploymentMode` change cannot reach an existing
   connection.** `applyControllersConfigToConnection` resolves the mode as
   `metadata.meshsync_deployment_mode` first, falling back to the resolved
   document only when that key is undefined. Every Kubernetes connection has
   that key written at registration, so the fallback is unreachable and the
   fanned-out re-apply silently keeps the old mode.

3. **The wizard's mode picker and the layered override diverge.** The wizard
   writes `metadata.meshsync_deployment_mode` alone, leaving any
   `controllers_config.operator.deploymentMode` untouched. The controllers
   editor then reports an "Override" chip for a mode the connection is not
   running.

## Test coverage

End-to-end coverage lives in `ui/tests/e2e/controllers-config.spec.ts` and is
keyed to the **Operator, MeshSync & Broker Settings** Test Plan Test Group via
the Allure `testGroup` label, the same mechanism the Connection Lifecycle report
uses (see [Contributing to Meshery's CLI tests]({{< relref "cli/tests.md" >}})).

Unit coverage:

- `server/models/connections/controllers_config_test.go` - precedence, atomic
  collection merge, metadata round-trip, validation.
- `server/models/controllers_config_apply_test.go` - cluster propagation,
  withdrawal, restart-only-on-watch-list-change.
- `ui/components/configuration/__tests__/ControllersConfigForm.test.tsx` -
  tri-state inherit/override semantics and conditional field visibility.

Before adding a setting, add its row above, its propagation assertion, and its
end-to-end case. A setting whose effect a user cannot observe is a defect; so is
one that fails silently.
