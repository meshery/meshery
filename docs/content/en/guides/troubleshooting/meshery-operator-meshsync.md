---
title: Meshery Operator, MeshSync, Broker Troubleshooting Guide
description: Comprehensive guidance for troubleshooting Meshery Operator, MeshSync and Broker deployments under various scenarios.
categories: [troubleshooting]
---

{{% alert color="info" title="What is Meshery Operator?" %}}
<a href='{{< ref "concepts/architecture/operator/index.md" >}}'>Meshery Operator</a> controls and monitors the lifecycle of components deployed inside Meshery-managed Kubernetes clusters. Learn more about <a href='{{< ref "concepts/_index.md" >}}'>Meshery's architecture</a>.
{{% /alert %}}

This guide offers comprehensive troubleshooting instructions for [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}}) and its custom controllers, [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) and [Broker]({{< ref "concepts/architecture/broker/index.md" >}}). Follow the steps outlined in this document to ensure a smooth Meshery deployment.

First, understand the [Meshery Operator Deployment Scenarios](#meshery-operator-deployment-scenarios) and the [Status of Meshery Operator, MeshSync, and Meshery Broker](#understanding-the-status-of-meshery-operator-meshsync-and-meshery-broker) to identify the deployment model fitting of your environment. Then, follow the guidance under the respective scenario to troubleshoot accordingly.

{{% alert color="dark" title="Meshery Error Code Reference" %}}
Have specific error with an error code? See the <a href='{{< ref "reference/references/error-codes.md" >}}'>Meshery Error Code Reference</a> for probable cause and suggested remediations.
{{% /alert %}}

## Understanding the Status of Meshery Operator, MeshSync, and Meshery Broker

Each Meshery Operator controller offers a health status that you can use to understand its current health in your deployment. These statuses are computed by Meshery Server from what it observes of the Operator, MeshSync, and Broker; their meanings are described below.

### Meshery Operator Health Status

- **DEPLOYED:** Operator deployment rollout is done, pod is in a ready state, old pod (if any) has been terminated.
- **DEPLOYING:** Operator deployment is present, but its rollout is in progress. Pod is not yet in ready state, or old pod has not been terminated.
- **NOTDEPLOYED:** Operator deployment is not present in the cluster.

### MeshSync Health Status

- **ENABLED:** Custom Resource present. MeshSync Controller is not connected to Broker.
- **DEPLOYED:** Custom Resource present. MeshSync Controller is present but the state is not RUNNING or ERRDISABLE, though
- **RUNNING:** MeshSync pod present and in a running state.
- **CONNECTED:** Deployed and connected to Broker.
- **UNDEPLOYED:** Custom Resource not present.

### Meshery Broker Health Status

- **DEPLOYED:** External IP not exposed OR External IP exposed but Meshery Server is not connected as a client to Broker hence data is not being published.
- **UNDEPLOYED:** Custom Resource not deployed.
- **CONNECTED:** Deployed, sending data to Meshery Server.

## Meshery Operator Deployment Scenarios

Because Meshery is versatile in its deployment models, there are different scenarios in which you may need to troubleshoot the health of Meshery Operator. Identify the deployment model fitting your environment and follow the guidance under the respective scenario to troubleshoot accordingly.

### In-Cluster Deployment

<!-- Meshery Operator, MeshSync, and Broker are deployed in the same cluster as Meshery Server. This is the default deployment scenario when using `mesheryctl system start` or `make run-local`. -->

Whether using [`mesheryctl system start`]({{< ref "installation/_index.md" >}}), [`helm install`]({{< ref "installation/kubernetes/helm.md" >}}) or `make run-local`, Meshery Server will automatically connect to any available Kubernetes clusters found in your kubeconfig (under `$HOME/.kube/config`). Once connected, operator, broker(NATS) and meshsync will automatically get deployed in the same clusters.

If everything is fine, by viewing the connection in Meshery UI, MeshSync should be in **CONNECTED:** state. Otherwise, check the Operator's pod logs:

`kubectl logs <meshery-operator-pod> -n meshery`

### Out-of-Cluster Deployment

1. Meshery Server is deployed on any Docker host (- Meshery Server is deployed on a Docker host, and Meshery Operator is deployed on a Kubernetes cluster).
   _or_
2. Meshery is managing multiple clusters, some of which are not the cluster unto which Meshery Server is deployed.

In this model Meshery Server must both **reach** and **authenticate to** the in-cluster Meshery Broker (NATS):

- **Reachability.** The Broker is usually exposed as `ClusterIP` only, which is not reachable from outside the cluster. When Meshery Server runs out-of-cluster it **automatically establishes a self-healing port-forward** to the Broker's NATS pod through the Kubernetes API server (like `kubectl port-forward`, using the credentials Meshery already holds) and connects over it - no manual step required. This is on by default out-of-cluster, skipped automatically in-cluster, and can be disabled with `MESHERY_MANAGED_BROKER_PORTFORWARD=false` (after which you must provide your own path, e.g. `kubectl port-forward -n meshery svc/meshery-nats 4222:4222`, or expose the Broker via NodePort/LoadBalancer).
- **Authentication.** The Operator provisions NATS with token authentication (secret `meshery-nats-auth`); Meshery Server reads that token and presents it automatically. Without it the Broker rejects the connection with an authorization violation.

For the full walkthrough of the Kubernetes connection lifecycle, its components, and these connectivity behaviors, see [Kubernetes Connection Lifecycle]({{< ref "guides/infrastructure-management/kubernetes-connection-lifecycle.md" >}}).

## Common Failure Scenarios

Some common failure situations that Meshery users might face are described below.

1. **Situation:** No deployment of Meshery Operator, MeshSync, and Broker.
   1. **Probable cause:** Meshery Server cannot connect to Kubernetes cluster; cluster unreachable or kubeconfig without proper permissions needed to deploy Meshery Operator; Kubernetes config initialization issues.
1. **Situation:** Meshery Operator with MeshSync and Broker deployed, but Meshery Server is not receiving data from MeshSync or data the [Meshery Database]({{< ref "concepts/architecture/database/index.md" >}}) is stale.
   1. **Probable cause:** 
   2. Meshery Server lost subscription to Meshery Broker; Broker server not exposed to external IP; MeshSync not connected to Broker; MeshSync not running; Meshery Database is stale.
   3. The SQL database in Meshery serves as a cache for cluster state. A single button allows users to dump/reset the Meshery Database.
   4. Orphaned MeshSync and Broker controllers - Meshery Operator is not present, but MeshSync and Broker controllers are running.
   5. **Broker unreachable / not authenticated (out-of-cluster Meshery):** the Broker is `ClusterIP`-only and unreachable, or Meshery is not presenting the NATS token. See [Out-of-Cluster Deployment](#out-of-cluster-deployment); the connection's [Diagnostics](#diagnostics-in-the-connection-detail-view) will report `broker_unreachable` with remediation.
1. **Situation:** The `meshery-nats` (Broker) pod is in `CrashLoopBackOff` and never becomes ready.
   1. **Probable cause:** Some Meshery Operator versions inject the NATS token into `nats.conf` **unquoted** (`token: $NATS_TOKEN`). When the generated token happens to look like a number, the NATS config parser rejects it and the pod crash-loops. Confirm with `kubectl logs -n meshery meshery-nats-0 -c nats` (look for a `variable reference for 'NATS_TOKEN' ... could not be parsed` error). The fix belongs in the Operator (quote it: `token: "$NATS_TOKEN"`); redeploying often generates a token that parses.

## Operating Meshery without Meshery Operator

Meshery Operator, MeshSync, and Broker are crucial components in a Meshery deployment. Meshery can function without them, but some functions of Meshery will be disable / unusable. Whether Meshery Operator is initially deployed via `mesheryctl` command or via Meshery Server, you can monitor the health of the Meshery Operator deployment using either the CLI or UI clients.

## Verifying the Status of Meshery Operator, MeshSync, and Meshery Broker

## Troubleshooting using Meshery CLI

The following commands are available to troubleshoot Meshery Operator, MeshSync, and Broker.

**Meshery Server and Adapters**

- `mesheryctl system status` - Displays the status of Meshery Server and Meshery Adapters.

**Meshery Operator, MeshSync, and Broker**

- `mesheryctl system check` - Displays the status of Meshery Operator, MeshSync, and Broker.

## Troubleshooting using Meshery UI

Based on discussed scenarios, the UI exposes tools to perform the following actions:

- (Re)deploy Operator, MeshSync, Broker.
- Uninstall and Install MeshSync, Broker, Operator.
- Reset Database.
- Ad hoc Connectivity Test for Operator, Meshery Broker, MeshSync.
- Reconnect Meshery Server to Meshery Broker.
- Ad hoc Connectivity Test for Kubernetes context.
- Rediscover kubeconfig, delete, (re)upload kubeconfig.

### Diagnostics in the connection detail view

Click a Kubernetes connection's row in the Connections table to open its detail view. Below the Operator / MeshSync / Broker status chips, a **Diagnostics** section lists actionable problems and remediation, derived from the live controller status and Meshery's actual Broker connection:

| Code | Meaning | Remediation |
|---|---|---|
| `connection_inactive` | No active session for the connection yet. | Connect the cluster. |
| `operator_not_deployed` | The Operator is not deployed (operator mode). | Reconnect the cluster, or switch MeshSync mode; ensure Meshery can create resources in the `meshery` namespace. |
| `broker_unreachable` | The Broker is up but Meshery cannot reach/authenticate to it. | The managed port-forward normally handles reachability out-of-cluster; otherwise port-forward the Broker, expose it via NodePort/LoadBalancer, or run Meshery in-cluster. |

The same data is available at `GET /api/system/controllers/diagnostics?connectionId=<id>`.

### Synthetic Test for Ensuring Change in Cluster State

Initiate a synthetic check to verify a fully functional Operator deployment, testing MeshSync/Broker connectivity.

- Empty database shows the main-cluster node.
- Corrupt database triggers an error snackbar with a link to the Settings screen.
- Disconnected Kubernetes displays MeshSync logo pulsating when data is received.

<div class="section">
Future Enhancements for Troubleshooting:

- NATS/MeshSync not running prompts a review of available operations in the Settings panel.

</div>

## Inspecting MeshSync Directly

When the CLI and UI clients don't explain *why* data is missing or stale, inspect the MeshSync pod directly.

**Read MeshSync logs** (enable debug logging for detail):

```bash
kubectl -n meshery logs deploy/meshery-meshsync
# For verbose output, set DEBUG=true on the Deployment and let it restart:
kubectl -n meshery set env deploy/meshery-meshsync DEBUG=true
```

**Check liveness and readiness** (MeshSync serves these on port `11000`):

```bash
kubectl -n meshery port-forward deploy/meshery-meshsync 11000:11000 &
curl -sS http://127.0.0.1:11000/healthz    # liveness
curl -sS -o /dev/null -w '%{http_code}\n' http://127.0.0.1:11000/readyz   # 200 == connected to Broker
```

{{% alert color="info" title="What readiness does and does not mean" %}}
<code>/readyz</code> returns <code>200</code> once MeshSync has connected to the Broker, <strong>not</strong> once its informer caches have finished priming. Immediately after a (re)start MeshSync may report ready while its cluster snapshot is still filling in. If Meshery shows a partial cluster right after a restart, give discovery a moment or trigger a fresh discovery with <code>kubectl -n meshery rollout restart deploy/meshery-meshsync</code>.
{{% /alert %}}

**Verify the Broker is reachable from MeshSync.** On startup MeshSync runs a connectivity test against the Broker's monitoring endpoint (`http://<broker-host>:8222/connz`) before opening its NATS client; a failure here appears in the MeshSync logs and blocks readiness. Confirm the `BROKER_URL` value and that the Broker Service is reachable:

```bash
kubectl -n meshery get deploy meshery-meshsync \
  -o jsonpath='{.spec.template.spec.containers[0].env[?(@.name=="BROKER_URL")].value}{"\n"}'
```

### Behaviors that commonly explain missing or churning data

- **A new or changed CRD triggers a full re-discovery.** MeshSync watches the cluster's CustomResourceDefinitions and rebuilds its informers when the CRD set changes. On clusters where controllers rewrite CRDs frequently (for example, cert-manager's CA injector updating CRD `caBundle` fields), this can cause repeated re-discovery and transient load or gaps. If you observe this, scope discovery with a whitelist (see the [MeshSync configuration FAQ]({{< ref "concepts/architecture/meshsync.md#meshsync-faqs" >}})).
- **Secrets are discovered by default.** MeshSync watches `secrets.v1.`, and the Secret objects it forwards to Meshery Server include their `data` and `stringData` payload. Those Secret contents are therefore transmitted over the Broker and persisted in the Meshery Database. In security-sensitive environments, either blacklist `secrets.v1.` (or use a whitelist that omits it) to keep Secrets out of discovery entirely, or set `MESHSYNC_REDACT_SECRETS=true` on the MeshSync Deployment to keep discovering Secrets while replacing their values with `[REDACTED]` (keys are preserved). See [Redacting Secret contents]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md#redacting-secret-contents" >}}).
- **Discovery is watch-driven with no periodic re-list.** MeshSync relies on the Kubernetes watch stream rather than polling. If you suspect the in-memory snapshot has drifted, force a re-list with `kubectl -n meshery rollout restart deploy/meshery-meshsync` or reset the Meshery Database from the UI.

## See Also

- [Kubernetes Connection Lifecycle]({{< ref "guides/infrastructure-management/kubernetes-connection-lifecycle.md" >}})
- [Troubleshooting Meshery Installations]({{< ref "guides/troubleshooting/installation.md" >}})
- [Troubleshooting Errors while running Meshery]({{< ref "guides/troubleshooting/meshery-server.md" >}})

{{< related-discussions tag="meshery" >}}

