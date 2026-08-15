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

### When a controller reports UNKOWN

Any of the three can instead report **UNKOWN** (spelled that way on the wire). It is not a health state: it means Meshery made no observation of that controller on this cluster, so it is asserting nothing about it. All three report it at once when the connection's kubeconfig could not be read or its Kubernetes client could not be created, since nothing about the cluster was observable. The cards stay visible on purpose - the reason is in the connection's [Diagnostics](#diagnostics-in-the-connection-detail-view) and in the events feed, not in the status itself.

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
1. **Situation:** The `meshery-operator` pod never becomes ready after connecting a cluster. Its `kube-rbac-proxy` container is in `ImagePullBackOff`, and/or its `manager` container crash-loops with `open /tmp/k8s-webhook-server/serving-certs/tls.crt: no such file or directory`.
   1. **Probable cause:** An old Meshery Operator Helm chart was installed. See [Meshery Operator will not start: ImagePullBackOff and a missing webhook certificate](#meshery-operator-will-not-start-imagepullbackoff-and-a-missing-webhook-certificate).
1. **Situation:** Terminating stuck resources during Meshery Operator upgrade.
   1. **Probable cause:** During upgrade, the Meshery Operator may fail to upgrade smoothly if the previous version's MeshSync or Broker resources are stuck in a terminating state (typically due to finalizers).
   2. **Remediation:** Follow these steps to clean up stuck resources and safely apply the upgrade:
      1. Upgrade the `mesheryctl` CLI binary first if it is outdated, as `mesheryctl system update` only updates Meshery container images and manifests (see the [Upgrading mesheryctl]({{< ref "installation/upgrades/index.md#upgrading-meshery-cli" >}}) guide).
      2. Check if the `MeshSync` and `Broker` custom resources are in a terminating state by inspecting their deletion timestamps:
         ```bash
         kubectl get meshsync meshery-meshsync -n meshery -o jsonpath='{.metadata.deletionTimestamp}'
         kubectl get broker meshery-broker -n meshery -o jsonpath='{.metadata.deletionTimestamp}'
         ```
      3. For any resource confirmed to be terminating (having a `deletionTimestamp`), patch only that specific resource to clear its finalizers. **Do not** patch healthy resources, as this disrupts the Operator's normal cleanup path:
         ```bash
         # Run only if meshery-meshsync is stuck in Terminating
         kubectl patch meshsync meshery-meshsync -n meshery -p '{"metadata":{"finalizers":null}}' --type=merge

         # Run only if meshery-broker is stuck in Terminating
         kubectl patch broker meshery-broker -n meshery -p '{"metadata":{"finalizers":null}}' --type=merge
         ```
      4. Pull the latest Meshery container images and apply the Helm manifests upgrade:
         ```bash
         mesheryctl system update
         ```
      5. Restart the Meshery components:
         ```bash
         mesheryctl system restart
         ```
      6. Verify the operator components initialize successfully:
         * Run the system check:
           ```bash
           mesheryctl system check --operator
           ```
         * Confirm that the Broker workload pod (`meshery-nats-*` or `meshery-broker-*` depending on the Operator version) and other operator pods are running and ready:
           ```bash
           kubectl get pods -n meshery
           ```

## Meshery Operator will not start: ImagePullBackOff and a missing webhook certificate

**Signature.** After adding a Kubernetes connection, the operator never reaches **DEPLOYED**. Both of these appear together:

```bash
kubectl -n meshery get pods
# meshery-operator-...   1/2   ImagePullBackOff   (kube-rbac-proxy)
#                              CrashLoopBackOff   (manager)

kubectl -n meshery describe pod -l app=meshery-operator | grep -A2 kube-rbac-proxy
# Failed to pull image "registry.k8s.io/kubebuilder/kube-rbac-proxy:v0.16.0"
# ...older charts name the same sidecar under the other registry:
# Failed to pull image "gcr.io/kubebuilder/kube-rbac-proxy:v0.16.0"

kubectl -n meshery logs deploy/meshery-operator -c manager
# open /tmp/k8s-webhook-server/serving-certs/tls.crt: no such file or directory
```

**Cause.** Both symptoms come from one thing: an **old `meshery-operator` Helm chart**, from before `v1.0.51`.

- Those charts ship a `kube-rbac-proxy` sidecar next to the manager. Charts around `v0.8.214` and above - including the chart a helm-installed Meshery Server asks for, which is the usual case - name it `registry.k8s.io/kubebuilder/kube-rbac-proxy:v0.16.0`; older charts such as `v0.8.180` name the same image as `gcr.io/kubebuilder/kube-rbac-proxy:v0.16.0`. Match on either. When that image does not pull, the container sits in `ImagePullBackOff` and the Pod never becomes ready.
  - Why the pull fails is not the same story for both registries, and only one of them is settled: the `gcr.io` copy of `v0.16.0` is gone (the repository's tag list is empty and the tag returns `404`), while the `registry.k8s.io` copy still resolves from an unrestricted network. So if your cluster reports the pull failing on `registry.k8s.io`, treat it as something about that cluster's access to the registry rather than as the image having been withdrawn. The remedy below does not depend on which it is.
- They also set no `ENABLE_WEBHOOKS` on the manager and mount no serving certificate. Current operator images treat an unset `ENABLE_WEBHOOKS` as *enabled*, so the manager looks for a certificate that the old chart never created and crash-loops.

`v1.0.51` is the oldest published chart confirmed to render without the sidecar and with `ENABLE_WEBHOOKS=false` - the conversion webhook opt-in, off by default - and every chart above it does the same. Both conditions behind the symptoms - the sidecar, and the absent `ENABLE_WEBHOOKS` - were confirmed by rendering the published archives of `v1.0.40` and of the contiguous run `v1.0.41` through `v1.0.50` (`v1.0.47` was never published); charts older than `v1.0.40` were not rendered, so treat "older than `v1.0.51`" as the boundary rather than a claim about any specific ancient release.

**Remedy.** Get a chart at or above `v1.0.51` into the cluster.

If Meshery Server deployed the operator for you (the usual case: you added a kubeconfig and Meshery installed the operator), simply **upgrade Meshery Server**. The version a Server derives from its own release is floored: the Server resolves it against what the repository actually publishes and raises anything older than `v1.0.51` to the oldest published chart at or above that boundary, telling you it did so in the events feed. Then reconnect the cluster.

The floor covers that derived version and nothing else, so two cases still leave an old chart in place:

- **You pinned one.** An `operator.version` you set is honored as written and is never raised, so a connection pinned below `v1.0.51` keeps reinstalling that chart until you change or clear the pin. See [Choosing the chart version yourself](#choosing-the-chart-version-yourself).
- **Nothing published reaches the floor.** If the repository carries no chart at or above `v1.0.51`, the newest published chart is deployed anyway, and the events feed says so and warns that the Operator may not become ready.

If you installed the operator chart yourself with Helm:

```bash
helm repo add meshery https://meshery.github.io/meshery.io/charts
helm repo update
helm upgrade --install meshery-operator meshery/meshery-operator \
  --namespace meshery --create-namespace
```

`helm repo update` matters on its own: a stale local repository cache will happily reinstall the same broken chart.

If the operator Pod is still wedged after the upgrade, delete it so the new spec takes effect immediately:

```bash
kubectl -n meshery rollout restart deploy/meshery-operator
```

### Confirming which Meshery Operator version is deployed

The chart version and the operator image tag are different numbers; the image tag is what tells you which operator is actually running:

```bash
# The operator image actually running in the cluster
kubectl -n meshery get deploy meshery-operator \
  -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'

# The Helm chart release that installed it (chart version and app version)
helm -n meshery list --filter meshery-operator
```

A `kube-rbac-proxy` container in the output of the following command means the chart predates `v1.0.51`, whatever version it claims:

```bash
kubectl -n meshery get deploy meshery-operator \
  -o jsonpath='{range .spec.template.spec.containers[*]}{.name}{"\n"}{end}'
```

## Choosing the chart version yourself

Meshery Server normally deploys the operator chart that matches its own release, falling back to the newest published chart when that one is not published yet (chart publishing trails Meshery Server releases). To pin a specific chart version for one connection, set **`operator.version`** in the connection's controllers configuration.

The value must be a chart version that the repository publishes, for example `v1.0.64` (the leading `v` is optional - `1.0.64` names the same chart). A moving tag such as `stable-latest`, or a version that is not published, is rejected with a visible error rather than being silently replaced. Clear the field to go back to tracking the Meshery Server release.

A pin is deliberate, so it also escapes the minimum-version floor that guards the derived version: `operator.version` below `v1.0.51` is installed as written rather than raised. Every published chart checked below that boundary leaves the Operator Pod unable to become ready (see [Meshery Operator will not start](#meshery-operator-will-not-start-imagepullbackoff-and-a-missing-webhook-certificate)), so choosing a chart that deploys is yours to get right.

Release candidates are the one thing Meshery will never pick for you: when it falls back to the newest published chart, or raises an old one to the oldest chart known to deploy, it skips any version carrying a prerelease suffix such as `v1.0.66-rc.1`. Naming a prerelease in `operator.version` deploys it exactly as asked.

If the chart repository cannot be reached, Meshery still reports the operator's status and image version for an operator that is already installed - only installing or upgrading it is withheld, and the reason appears in the connection's diagnostics and in the events feed.

## When the chart version cannot be resolved

A version that cannot be resolved - an `operator.version` naming something the repository does not publish, or a chart repository Meshery could not read at all - stops the install *before* Helm is called. Nothing partial is applied to the cluster.

**The failure is visible, not silent.** The Operator's status card carries the error, and the connection's [Diagnostics](#diagnostics-in-the-connection-detail-view) report `operator_deploy_failed` with the underlying cause; the same cause appears in the events feed. The card does not vanish and the operator does not sit at a blank status while Meshery quietly gives up.

**Retrying is the remedy for a transient outage.** Meshery re-resolves the version on a user-initiated deploy, so a chart repository that was briefly unreachable needs no reconnect and no restart: redeploy the Operator from the connection's actions and it self-heals. Check first that Meshery Server has outbound access to `https://meshery.github.io/meshery.io/charts` ([egress requirements]({{< ref "installation/production/networking-and-connectivity.md#egress-requirements" >}})).

If the error names your own pin, the pin is the problem - correct `operator.version` to a published version or clear it to go back to tracking the Meshery Server release. An explicit pin is never quietly swapped for a working one.

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

Click a Kubernetes connection's row in the Connections table to open its detail view. Below the Operator / MeshSync / Broker status chips, a **Diagnostics** section lists actionable problems and remediation, derived from the live controller status and Meshery's actual Broker connection. Every code it reports, and the remedy for each, is listed in [Diagnostics]({{< ref "guides/infrastructure-management/kubernetes-connection-lifecycle.md#diagnostics" >}}).

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

