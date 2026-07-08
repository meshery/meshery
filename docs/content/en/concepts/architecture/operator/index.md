---
title: Operator
description: "Meshery Operator controls and manages the lifecycle of components deployed inside a kubernetes cluster"
display_title: false
aliases:
- /architecture/operator/
---

# Meshery Operator <img src="./images/B203EFA85E89491B.png" width="30" height="35" style="display:inline"/>

Meshery Operator is a Kubernetes operator (built on the
[Kubebuilder](https://book.kubebuilder.io) `go/v4` framework) that deploys and
manages the lifecycle of two Meshery components critical to Meshery's
operation of Kubernetes clusters: [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}})
and [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}). Deploy one
Meshery Operator per Kubernetes cluster under management — whether Meshery
Server is deployed inside or outside of the clusters under management.

## Deployments

It is recommended to deploy one Meshery Operator per cluster.

[![Meshery Operator and MeshSync](./images/meshery-operator-and-meshsync.svg
)](./images/meshery-operator-and-meshsync.svg)

### Initialization Sequence

[![Meshery Operator and MeshSync](./images/meshery-operator-deployment-sequence.svg
)](./images/meshery-operator-deployment-sequence.svg)

### How Meshery Server installs and upgrades the Operator

Meshery Server installs the `meshery-operator` Helm chart from
[meshery.io/charts](https://meshery.io/charts) into the `meshery` namespace of
each managed cluster, requesting the chart version that **matches the Meshery
Server release**. Upgrading Meshery Server is therefore what upgrades the
Operator: the Server re-applies the newer chart, the chart's CRD update Job
refreshes the CRD schemas, and the Operator Deployment rolls to the operator
version pinned in that chart. Manual operator upgrades on Server-managed
clusters are reverted by the Server's reconciliation — see
[How Meshery Server manages Meshery Operator]({{< ref "installation/upgrades/index.md#how-meshery-server-manages-meshery-operator" >}}).

## Custom resources

The Operator owns two CRDs, each reconciled into concrete workloads:

| CRD | Custom resource | Reconciles into |
|---|---|---|
| `brokers.meshery.io` | `meshery-broker` | NATS `StatefulSet`, client + headless `Service`s, configuration |
| `meshsyncs.meshery.io` | `meshery-meshsync` | MeshSync `Deployment` (with the Broker endpoint injected) |

Both CRDs serve two API versions — `v1alpha1` (compatibility) and `v1alpha2`
(storage) — with identical schemas, so either version can be read and written
interchangeably. Reconciliation uses Kubernetes Server-Side Apply with a
stable field manager: the Operator converges owned objects to their desired
state without fighting other controllers, watches the objects it owns
(including the Broker's Services), and reports health through standard status
`conditions` on each custom resource.

To customize what these custom resources declare - MeshSync's discovery scope,
version, and scale, and the Broker's version, scale, and service networking -
see [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}}).

## Controllers managed by Meshery Operator

### Broker Controller

Meshery Broker is one of the core components of the Meshery architecture. This controller manages the lifecycle of the NATS-backed broker that Meshery uses for data streaming across the cluster and the outside world, including declarative, in-place reconfiguration of the Broker's service networking.

See [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}) for more information.

### MeshSync Controller

MeshSync Controller manages the lifecycle of MeshSync that is deployed for resource synchronization for the cluster. It derives the Broker's address from the `Broker` resource's status and injects it into MeshSync as `BROKER_URL` (a `nats://host:port` URL); when the Broker's endpoint changes, the controller re-reconciles MeshSync so it reconnects.

See [MeshSync]({{< ref "concepts/architecture/meshsync.md" >}}) for more information.

## Common tasks

**Check Operator health:**

```bash
kubectl -n meshery rollout status deploy/meshery-operator
kubectl -n meshery get brokers,meshsyncs
kubectl -n meshery get broker meshery-broker -o jsonpath='{.status.conditions}'
```

**Find the running Operator version:**

```bash
kubectl -n meshery get deploy meshery-operator \
  -o jsonpath='{.spec.template.spec.containers[0].image}{"\n"}'
```

**Scrape Operator metrics:** the Operator serves Prometheus metrics over TLS
on port `8443` with authentication and authorization. Bind your scraper's
ServiceAccount to the `meshery-metrics-reader` ClusterRole and scrape the
`meshery-operator` Service.

**Upgrade the Operator:** upgrade Meshery Server — see the
[Upgrading Meshery guide]({{< ref "guides/upgrading-meshery/index.md" >}}).

## Operator FAQs

### When is Meshery Operator deployed and when is it deleted?
As a Kubernetes custom controller, Meshery Operator is provisioned and deprovisioned when Meshery Server is connected to or disconnected from Kubernetes cluster. Meshery Server connections to Kubernetes clusters are controlled using Meshery Server clients: `mesheryctl` or Meshery UI.  This behavior described below is consistent whether your Meshery deployment is using Docker or Kubernetes as the platform to host the Meshery deployment.

**Meshery CLI**

`mesheryctl` initiates connection to Kubernetes cluster when `mesheryctl system start` is executed and disconnects when `mesheryctl system stop` is executed. This behavior is consistent whether your Meshery deployment is using Docker or Kubernetes as the platform to host the Meshery deployment.

**Meshery UI**

Meshery UI offers more granular control over the deployment of Meshery Operator in that you can remove Meshery Operator from a Kubernetes cluster without disconnecting Meshery Server from the Kubernetes cluster. You can control the deployment of Meshery Operator using the on/off switch found in the Meshery Operator section of  Settings.

### Does the Meshery Operator use an SDK or framework?
Yes. Meshery Operator is built on [Kubebuilder](https://book.kubebuilder.io)
(`go/v4` layout) with controller-runtime; the Operator SDK is used to produce
its OLM bundle for Operator Lifecycle Manager installations.

### What happens to the CRDs when the Operator is removed?
The `brokers.meshery.io` and `meshsyncs.meshery.io` CRDs — and any `Broker`
and `MeshSync` objects — deliberately survive operator removal, so
reinstalling the Operator picks up where it left off with no data loss.
Deleting the CRDs (which deletes every custom resource of those types) is an
explicit manual step: `kubectl delete crd brokers.meshery.io meshsyncs.meshery.io`.

### How does the operator expose information about broker endpoints?

During Broker reconciliation, the operator derives endpoints **purely from
the Broker's client Service** (its spec and status — no network probing) and
populates the `status` field of the `brokers/meshery-broker` Custom Resource
(CR):

For example, for an in-cluster Meshery deployment:

```yaml
status:
  endpoint:
    external: localhost:31670
    internal: 10.96.49.130:4222
```

For example, for an out-of-cluster Meshery deployment:

```yaml
status:
  endpoint:
    external: 1e2cd15619524f569e695f648ae7c74e-0123456789.us-south-3.elb.cloud-provider.com:4222
    internal: 10.96.49.130:4222
```

The internal endpoint is always the Service's `ClusterIP` + client port.

The external endpoint is selected in order of precedence:

- `spec.service.externalEndpointOverride` on the `Broker` resource, verbatim
  (for brokers fronted by an ingress/gateway, NAT, or air-gapped topologies)
- LoadBalancer ingress hostname (if present) + client port
- LoadBalancer ingress IP + client port
- Kubernetes API server host + the Service's `nodePort` (NodePort services)
- `ClusterIP` + client port (cluster-internal only)

Because derivation is driven by watches on the Service, changing the Broker's
service networking (for example `ClusterIP` → `NodePort`) recomputes
`status.endpoint` automatically, and MeshSync is re-reconciled to connect to
the new address.

### Troubleshooting Meshery Operator and Related Components

To verify that your Meshery Operator and related components are functioning properly, perform the following checks:

- Ensure the Operator pods are running.
- Confirm that your cluster has appropriate RBAC permissions set.
- Validate that Meshery Server is able to communicate with Meshery Operator.

If you're seeing issues with **Meshery Operator**, **MeshSync**, or the **Broker**, refer to the [Meshery Troubleshooting Guide]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}}).

Whether you're facing installation issues, resource syncing failures, or Broker communication problems, the guide walks you through how to identify and fix them effectively.
