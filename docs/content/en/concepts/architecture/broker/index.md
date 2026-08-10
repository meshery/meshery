---
title: Broker
description: "Meshery broker component facilitates data streaming between kubernetes cluster components and outside world."
aliases:
- /architecture/broker/
---

Meshery Broker provides data streaming across independent components of
Meshery, whether those components are running inside or outside of the
Kubernetes cluster. It is a [NATS](https://nats.io) message broker, deployed
and managed by [Meshery Operator]({{< ref "concepts/architecture/operator/index.md" >}})
from a `Broker` custom resource: the Operator reconciles the resource into a
NATS `StatefulSet`, a client `Service` (plus a headless `Service` for peer
discovery), and its configuration, following the official NATS Helm chart's
topology. Client access is authenticated with a token that the Operator
generates into a Secret — no static credentials ship with the deployment.

[![Meshery Log Viewer](./images/meshery-log-viewer.svg
)](./images/meshery-log-viewer.svg)

## Declarative service networking

How the Broker is exposed on the network is part of the `Broker` resource's
spec, and every field is **reconfigurable on a live Broker**: the Operator
updates the Service in place (no pod deletion), re-derives
`status.endpoint`, and re-reconciles MeshSync so it reconnects to the new
address.

```yaml
apiVersion: meshery.io/v1alpha2
kind: Broker
metadata:
  name: meshery-broker
  namespace: meshery
spec:
  size: 1
  service:
    type: NodePort                  # ClusterIP (default) | NodePort | LoadBalancer
    annotations:                    # cloud LB hints, MetalLB pools, internal-LB switches
      service.beta.kubernetes.io/aws-load-balancer-internal: "true"
    loadBalancerClass: metallb      # LoadBalancer type only
    loadBalancerSourceRanges:       # LoadBalancer type only
      - 10.0.0.0/8
    externalEndpointOverride: nats.example.com:4222   # pin the advertised endpoint
```

- **`type`** — `ClusterIP` keeps the Broker cluster-internal (the default;
  suitable for in-cluster Meshery Server, kind, minikube, and bare metal).
  `NodePort` or `LoadBalancer` expose it to an out-of-cluster Meshery Server.
- **`externalEndpointOverride`** — advertises a fixed `host:port` instead of
  the derived one, for Brokers behind an ingress/gateway, NAT, or in
  air-gapped topologies.
- Invalid combinations (for example `loadBalancerClass` without
  `type: LoadBalancer`) are rejected at admission by validation rules on the
  CRD.

The Operator publishes the resulting addresses on the resource's status:

```yaml
status:
  endpoint:
    external: 1e2cd156-0123456789.us-south-3.elb.cloud-provider.com:4222
    internal: 10.96.49.130:4222
```

See [How does the operator expose information about broker endpoints?]({{< ref "concepts/architecture/operator/index.md#how-does-the-operator-expose-information-about-broker-endpoints" >}})
for the endpoint selection order. For the Broker's place in the full
configuration surface of Meshery's in-cluster components - including when to
choose each service type - see
[Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md#configuring-meshery-broker" >}}).

## Common tasks

**Change how the Broker is exposed — on a live Broker:**

```bash
kubectl -n meshery patch broker meshery-broker --type merge \
  -p '{"spec":{"service":{"type":"NodePort"}}}'
kubectl -n meshery get broker meshery-broker -o jsonpath='{.status.endpoint}{"\n"}'
```

**Pin the advertised endpoint** (ingress/NAT in front of the Broker):

```bash
kubectl -n meshery patch broker meshery-broker --type merge \
  -p '{"spec":{"service":{"externalEndpointOverride":"nats.example.com:4222"}}}'
```

**Check Broker health:**

```bash
kubectl -n meshery get broker meshery-broker -o jsonpath='{.status.conditions}{"\n"}'
kubectl -n meshery get statefulset meshery-nats
```

**Scale the Broker:**

```bash
kubectl -n meshery patch broker meshery-broker --type merge -p '{"spec":{"size":3}}'
```

### Broker FAQs

#### How many Brokers can run?
It is recommended to run one Broker instance (one `Broker` resource) per
Kubernetes cluster. The instance itself can be scaled through `spec.size`
(1–10 replicas) based on the data volume in the cluster.

#### What does an HA configuration look like?
We leverage Kubernetes functionality for high availability: the Broker runs as
a StatefulSet whose pods are restarted automatically on failure, with the
headless Service providing stable peer identities, and Meshery Operator
continuously reconciles the Broker back to its declared state. The Operator
reports readiness through the `Broker` resource's status conditions.

#### What stateful characteristics does the Broker have?
All messages published to the Broker are persisted in-memory within the
Broker instance until consumed. Persistent volumes/disk are not used by
default (NATS JetStream is available in the underlying deployment but
disabled by default).

#### How is access to the Broker secured?
Client connections authenticate with a token that Meshery Operator generates
and stores in a Kubernetes Secret in the Broker's namespace. MeshSync and
Meshery Server obtain the token from the same source; no credentials are
committed to images or manifests. Use `spec.service.loadBalancerSourceRanges`
to restrict network access when exposing the Broker through a load balancer.

#### How do I know if the Broker is working? How do I troubleshoot the Broker?
To check if your Broker instance is running smoothly (it's deployed as a Kubernetes StatefulSet), follow these quick checks:

- Confirm that the Broker pods are running and `status.conditions` on the
  `Broker` resource report ready.
- Confirm `status.endpoint` is populated, and that the external endpoint is
  reachable from wherever Meshery Server runs (only `ClusterIP` requires no
  external exposure — use `NodePort`/`LoadBalancer`, or an
  `externalEndpointOverride`, for out-of-cluster Servers).
- Make sure MeshSync's `BROKER_URL` points at the Broker's current endpoint
  (the Operator maintains this automatically).

Still seeing issues? The **[Meshery Troubleshooting Guide]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}})** covers common problems with the Broker, MeshSync, and Operator — and offers clear steps to resolve them.
