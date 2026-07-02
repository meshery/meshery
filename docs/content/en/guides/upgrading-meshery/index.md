---
title: Upgrading Meshery
description: >-
  Step-by-step procedure for upgrading Meshery and all of its components —
  CLI, Server, Operator, MeshSync, and Broker — on Docker and Kubernetes
  deployments.
categories: [installation]
weight: 10
aliases:
- /guides/upgrade
---

This guide walks through upgrading a running Meshery deployment. Meshery is a
composition of components that upgrade in a specific order: the CLI first, the
Server second, and the components on managed clusters — Meshery Operator,
MeshSync, and Broker — **automatically, driven by the Server**. You do not
upgrade the Operator by hand.

For background on which components version together, see the
[Upgrade Guide]({{< ref "installation/upgrades/index.md" >}}). For production
practices (pinned versions, upgrade-friendly probes, rollback rehearsal), see
the [Operational Readiness Checklist]({{< ref "installation/production/operational-readiness-checklist.md" >}}).

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
re-applies the `meshery-operator` Helm chart **at the chart version matching
the Server release**. That single `helm upgrade`, performed by the Server:

1. runs the chart's CRD update Job, which server-side-applies the current
   `Broker` and `MeshSync` CRD schemas (Helm alone never updates CRDs on
   upgrade — the Job is what delivers schema changes to live clusters);
2. rolls the Operator Deployment to the operator image version pinned in that
   chart;
3. the Operator then reconciles MeshSync and Broker to their expected
   versions and configuration.

No action is required on managed clusters. See
[How Meshery Server manages Meshery Operator]({{< ref "installation/upgrades/index.md#how-meshery-server-manages-meshery-operator" >}})
for the mechanics and caveats.

{{% alert title="Do not hand-upgrade the Operator on Server-managed clusters" color="warning" %}}
A manual `helm upgrade` of `meshery-operator` (or a hand-edited image tag) on
a cluster that Meshery Server manages is a stopgap at best: the Server's
reconciliation re-applies its own pinned chart version and will revert your
change. The durable way to get a newer Operator is to upgrade Meshery Server.
{{% /alert %}}

## Step 4: Verify

```bash
# Server and components
mesheryctl system status

# Operator image now matches the release bundled with your Server version
kubectl -n meshery get deploy meshery-operator \
  -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'

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
  older operator chart on managed clusters.
- CRD schemas are not rolled back by `helm rollback` directly. When the
  rolled-back Server reconnects to managed clusters, it re-applies the older
  operator chart as a Helm upgrade, and that chart's CRD update Job re-applies
  the older schemas. Stored objects remain readable throughout, because served
  versions stay identical across current schema revisions.

{{< related-discussions tag="meshery" >}}
