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
Have specific error with an error code? See the <a href='{{< ref "reference/error-codes.md" >}}'>Meshery Error Code Reference</a> for probable cause and suggested remediations.
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

### When a controller reports UNKNOWN

Any of the three can instead report **UNKNOWN** (spelled that way on the wire). It is not a health state: it means Meshery made no observation of that controller on this cluster, so it is asserting nothing about it. All three report it at once when the connection's kubeconfig could not be read or its Kubernetes client could not be created, since nothing about the cluster was observable. The cards stay visible on purpose - the reason is in the connection's [Diagnostics](#diagnostics-in-the-connection-detail-view) and in the events feed, not in the status itself.

## Meshery Operator Deployment Scenarios

Because Meshery is versatile in its deployment models, there are different scenarios in which you may need to troubleshoot the health of Meshery Operator. Identify the deployment model fitting your environment and follow the guidance under the respective scenario to troubleshoot accordingly.

### In-Cluster Deployment

In this scenario, Meshery Server runs in the same Kubernetes cluster as Meshery Operator, MeshSync, and Broker. This is what you get when you deploy Meshery to Kubernetes with [`mesheryctl system start`]({{< ref "installation/_index.md" >}}) or [`helm install`]({{< ref "installation/kubernetes/helm.md" >}}). When Meshery Server connects to a cluster in operator mode, it deploys Meshery Operator there, unless the server runs with `DISABLE_OPERATOR=true`. The Operator then deploys Broker (NATS) and MeshSync.

When everything is healthy, the cluster's connection in Meshery UI shows MeshSync as **CONNECTED**. If it does not, check the Operator's logs:

```bash
kubectl logs deploy/meshery-operator -n meshery -c manager
```

### Out-of-Cluster Deployment

This scenario applies when either of these is true:

- Meshery Server runs on a Docker host, or on your workstation with `make server`, and Meshery Operator runs in a Kubernetes cluster.
- Meshery manages several clusters, and some of them are not the cluster that Meshery Server runs in.

In this model Meshery Server must both **reach** and **authenticate to** the in-cluster Meshery Broker (NATS):

- **Reachability.** The Broker is usually exposed as `ClusterIP` only, which is not reachable from outside the cluster. When Meshery Server runs out-of-cluster it **automatically establishes a self-healing port-forward** to the Broker's NATS pod through the Kubernetes API server (like `kubectl port-forward`, using the credentials Meshery already holds) and connects over it - no manual step required. This is on by default out-of-cluster, skipped automatically in-cluster, and can be disabled with `MESHERY_MANAGED_BROKER_PORTFORWARD=false` (after which you must provide your own path, e.g. `kubectl port-forward -n meshery svc/meshery-nats 4222:4222`, or expose the Broker via NodePort/LoadBalancer).
- **Authentication.** The Operator provisions NATS with token authentication (secret `meshery-nats-auth`); Meshery Server reads that token and presents it automatically. Without it the Broker rejects the connection with an authorization violation.

For the full walkthrough of the Kubernetes connection lifecycle, its components, and these connectivity behaviors, see [Kubernetes Connection Lifecycle]({{< ref "guides/infrastructure-management/kubernetes-connection-lifecycle.md" >}}).

## Common Failure Scenarios

Some common failure situations that Meshery users might face are described below.

1. **Situation:** No deployment of Meshery Operator, MeshSync, and Broker.
   1. **Probable cause:** Meshery Server cannot connect to Kubernetes cluster; cluster unreachable or kubeconfig without proper permissions needed to deploy Meshery Operator; Kubernetes config initialization issues.
1. **Situation:** Meshery Operator with MeshSync and Broker deployed, but Meshery Server is not receiving data from MeshSync or data the [Meshery Database]({{< ref "concepts/architecture/database/index.md" >}}) is stale.
   1. **Probable causes:**
      - Meshery Server lost its subscription to Meshery Broker, MeshSync is not connected to Broker, or MeshSync is not running.
      - The Meshery Database is stale. Meshery Server's database is a cache of cluster state, and you can reset it from Meshery UI.
      - MeshSync and Broker are orphaned: Meshery Operator is gone, but MeshSync and Broker are still running.
      - **Broker unreachable / not authenticated (out-of-cluster Meshery):** the Broker is `ClusterIP`-only and unreachable, or Meshery is not presenting the NATS token. See [Out-of-Cluster Deployment](#out-of-cluster-deployment); the connection's [Diagnostics](#diagnostics-in-the-connection-detail-view) will report `broker_unreachable` with remediation.
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

If Meshery Server deployed the operator for you (the usual case: you added a kubeconfig and Meshery installed the operator), **reconnect the cluster** - after upgrading Meshery Server, if it predates this check. On every connect, Meshery Server reads the chart version of the `meshery-operator` Helm release in the `meshery` namespace, and if it is older than `v1.0.51` it upgrades the release to the chart it resolves for its own release, which is never older than `v1.0.51`. It does this whether the Operator is crash-looping or happens to be running. The events feed reports the upgrade (`meshery-server-1487`), naming the chart it found and the chart it installed. A release at or above `v1.0.51` is never touched.

The upgrade is a Helm upgrade, and Helm never upgrades a chart's CRDs on its own. Every chart from `v1.0.51` on carries a pre-upgrade Job (`meshery-operator-crd-update`) that applies the current `brokers.meshery.io` and `meshsyncs.meshery.io` CRDs - including the `v1alpha2` version that old charts such as `v1.0.0` never installed - before the Operator is changed. If that Job fails, Helm stops the upgrade before touching the Operator Deployment, so you are never left with a new Operator on old CRDs.

Earlier Meshery Server releases floored only the chart they were about to install, so an Operator they found already installed was left on its old chart however many times the Server was upgraded. That is what the check above closes.

Some cases still leave an old chart in place, and each is reported rather than silent:

- **You pinned one.** An `operator.version` you set is honored as written and is never raised. If the pin is below `v1.0.51`, an installed release below it is left alone, and the events feed and the connection's diagnostics report it (`meshery-server-1488`) on every connect until you change or clear the pin. A pin at `v1.0.51` or newer is what an old release is upgraded to. See [Choosing the chart version yourself](#choosing-the-chart-version-yourself).
- **Nothing published reaches the floor.** If the repository carries no chart at or above `v1.0.51`, a fresh install gets the newest published chart with a warning that the Operator may not become ready, and an installed release below `v1.0.51` is left alone and reported (`meshery-server-1489`) - upgrading it to another chart that cannot run would fix nothing.
- **The upgrade failed.** The events feed carries Helm's reason (`meshery-server-1491`), and the next connect tries again. If the CRD update Job is the cause, read it with `kubectl -n meshery logs job/meshery-operator-crd-update`.
- **Meshery cannot read the release** (`meshery-server-1490`), usually because the identity in the connection's kubeconfig may not list Secrets in the `meshery` namespace, where Helm keeps its release records.
- **The Operator is not the `meshery-operator` release.** An Operator installed as part of the Meshery Helm chart belongs to the `meshery` release and is upgraded with it (`helm upgrade meshery ...`); one applied without Helm has no chart version for Meshery to read. Upgrade either yourself, as below.

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

A pin is deliberate, so it also escapes the minimum-version floor that guards the derived version: `operator.version` below `v1.0.51` is installed as written rather than raised, and an Operator already installed below `v1.0.51` is not upgraded while the pin stands - Meshery reports that in the events feed instead (`meshery-server-1488`). Every published chart checked below that boundary leaves the Operator Pod unable to become ready (see [Meshery Operator will not start](#meshery-operator-will-not-start-imagepullbackoff-and-a-missing-webhook-certificate)), so choosing a chart that deploys is yours to get right.

Release candidates are the one thing Meshery will never pick for you: when it falls back to the newest published chart, or raises an old one to the oldest chart known to deploy, it skips any version carrying a prerelease suffix such as `v1.0.66-rc.1`. Naming a prerelease in `operator.version` deploys it exactly as asked.

If the chart repository cannot be reached, Meshery still reports the operator's status and image version for an operator that is already installed - only installing or upgrading it is withheld, and the reason appears in the connection's diagnostics and in the events feed.

## When the chart version cannot be resolved

A version that cannot be resolved - an `operator.version` naming something the repository does not publish, or a chart repository Meshery could not read at all - stops the install *before* Helm is called. Nothing partial is applied to the cluster.

**The failure is visible, not silent.** The Operator's status card carries the error, and the connection's [Diagnostics](#diagnostics-in-the-connection-detail-view) report `operator_deploy_failed` with the underlying cause; the same cause appears in the events feed. The card does not vanish and the operator does not sit at a blank status while Meshery quietly gives up.

**Retrying is the remedy for a transient outage.** Meshery re-resolves the version on a user-initiated deploy, so a chart repository that was briefly unreachable needs no reconnect and no restart: redeploy the Operator from the connection's actions and it self-heals. Check first that Meshery Server has outbound access to `https://meshery.github.io/meshery.io/charts` ([egress requirements]({{< ref "installation/production/networking-and-connectivity.md#egress-requirements" >}})).

If the error names your own pin, the pin is the problem - correct `operator.version` to a published version or clear it to go back to tracking the Meshery Server release. An explicit pin is never quietly swapped for a working one.

## Operating Meshery without Meshery Operator

Meshery Operator, MeshSync, and Broker are crucial components in a Meshery deployment. Meshery works without them, but it cannot discover or track the state of your clusters. Whether `mesheryctl` or Meshery Server deployed Meshery Operator, you can check its health from either the CLI or the UI.

## Troubleshooting using Meshery CLI

The following commands are available to troubleshoot Meshery Operator, MeshSync, and Broker.

**Meshery Server and Adapters**

- `mesheryctl system status` shows the status of Meshery Server and Meshery Adapters.

**Meshery Operator, MeshSync, and Broker**

- `mesheryctl system check --operator` shows whether Meshery Operator, MeshSync, and Broker are running in the cluster of your kubeconfig's current context. It works only when your `mesheryctl` context's platform is `kubernetes`.

## Troubleshooting using Meshery UI

Meshery UI lets you do the following from a Kubernetes connection:

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

{{< related-discussions tag="meshery" >}}

