---
title: Networking & Connectivity
linkTitle: Networking
description: >-
  Network ports and directional flows, ingress and Emissary configuration,
  secure WebSocket support, Broker exposure, egress to Remote Providers, the
  OAuth callback URL, and network policies for production Meshery.
categories: [installation]
weight: 40
aliases:
- /installation/production/networking
---

Getting networking right is what makes the difference between a Meshery
deployment that "works on my cluster" and one that reliably manages a fleet from
behind an ingress. This page covers the ports and directions traffic flows, how
to expose Meshery securely, how the Broker is reached (especially
out-of-cluster), and the egress Meshery needs.

## Network ports

Meshery uses the following ports across its components. The Broker ports beyond
the client port are for clustering/routing features and are not all in use in
current releases.

{{< network-ports >}}

In addition, optional [Meshery Adapters](/concepts/architecture/adapters) listen
on their own ports (for example `10000/tcp` and up), and the Meshery Operator's
`kube-rbac-proxy` listens on `8443/tcp` in-cluster. Inside the pod, Meshery
Server listens on `8080/tcp`, which the Kubernetes Service exposes as
`9081/tcp`.

## Directional traffic flows

Plan firewall rules and network policies around the direction each connection is
initiated, not just the ports.

| Source | Destination | Port(s) | Purpose | Notes |
| :--- | :--- | :--- | :--- | :--- |
| Clients (UI/browser, `mesheryctl`) | Meshery Server | 9081/tcp (HTTP/API/GraphQL); 80/tcp (WebSocket) | UI, REST, GraphQL, live updates | Front with ingress/TLS in production. |
| Meshery Server | Managed cluster Kubernetes API | 443/tcp (typically) | Discovery, deploy/undeploy, Operator lifecycle | Uses kubeconfig context (out-of-cluster) or ServiceAccount (in-cluster). |
| Meshery Server | Meshery Broker | 4222/tcp | Event/data streaming from each cluster | Reaches ClusterIP in-cluster; reaches an exposed endpoint out-of-cluster. |
| Meshery Server | Remote Provider | 443/tcp (egress) | Authentication, capabilities, durable state | Must be allowed through egress firewalls/proxies. |
| Meshery Server | Meshery Adapters | adapter ports (e.g. 10000+/tcp) | Capability operations | Only when adapters are deployed. |
| Operations/monitoring | Meshery Broker | 8222/tcp | Broker HTTP monitoring endpoint | Optional; for observability. |

{{% alert title="Direction matters for out-of-cluster" color="info" %}}
When Meshery Server is out-of-cluster, the Server initiates connections **to**
each managed cluster's API server and Broker. Those endpoints must be reachable
from wherever Meshery Server runs—through cloud load balancers, VPN/peering, or
NodePort exposure. See [Broker exposure](#exposing-the-broker) below.
{{% /alert %}}

## Exposing Meshery securely (ingress, TLS, WebSocket)

In production, do not expose the raw `LoadBalancer` Service. Front Meshery with
an ingress controller that terminates TLS and proxies both HTTP and WebSocket
traffic.

### Ingress with the Helm chart

The chart includes an `ingress` block (disabled by default). Enable it and set
the class, hosts, annotations, and TLS:

```yaml
ingress:
  enabled: true
  ingressClassName: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: meshery.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: meshery-tls
      hosts:
        - meshery.example.com
```

### Secure WebSocket support

Meshery relies on WebSocket connections for live updates in the UI. Your ingress
**must** be configured to upgrade and proxy WebSocket traffic, or the UI will
load but fail to receive real-time updates.

- With **Emissary-Ingress (Ambassador)**, ensure the `Mapping` allows WebSocket
  upgrades (Emissary supports `allow_upgrade` for `websocket`) and routes both
  the HTTP and WebSocket paths to the Meshery Service. Terminate TLS at the edge
  and use `wss://` for secure WebSockets.
- With **NGINX Ingress**, the controller handles `Connection: Upgrade` /
  `Upgrade: websocket` headers automatically for HTTP/1.1 backends; ensure proxy
  read/send timeouts are long enough for long-lived connections.
- Whatever the controller, terminate TLS at the edge so browser traffic is
  `https://`/`wss://` end-to-end.

{{% alert title="WebSocket is not optional" color="warning" %}}
A common production misconfiguration is an ingress that proxies HTTP but drops
WebSocket upgrades. The symptom is a UI that loads but never updates in real
time. Validate WebSocket upgrade support explicitly during rollout.
{{% /alert %}}

## The OAuth callback URL

When Meshery sits behind an ingress, reverse proxy, or load balancer, the
external URL users reach is not the address Meshery sees internally. For the
Remote Provider OAuth flow to complete, set the callback URL to Meshery's
**external** address:

- Environment variable: `MESHERY_SERVER_CALLBACK_URL`
- Helm value: `env.MESHERY_SERVER_CALLBACK_URL`

```bash
helm install meshery meshery/meshery --namespace meshery --create-namespace \
  --set env.MESHERY_SERVER_CALLBACK_URL=https://meshery.example.com
```

The callback path Meshery uses is `/api/user/token`, so the effective callback
becomes `https://meshery.example.com/api/user/token`. If login redirects fail or
loop when behind a proxy, a missing or incorrect callback URL is the usual
cause. See
[Authentication, Authorization &amp; Identity](/installation/production/authentication-and-identity)
and the
[environment variables reference](/installation/advanced/environment-variables).

## Exposing the Broker

The [Meshery Broker](/concepts/architecture/broker) is reached by Meshery Server
on `4222/tcp`. How it is exposed depends on topology:

- **In-cluster Server.** The Server reaches the Broker over its ClusterIP; no
  external exposure is needed.
- **Out-of-cluster Server.** The Broker must be reachable from outside the
  cluster. The Operator publishes the Broker's reachable endpoint into the
  Broker custom resource `status`, selecting an external endpoint from, in
  order of preference: the LoadBalancer hostname, the LoadBalancer IP, the
  kubeconfig host with the NodePort, the ClusterIP with the cluster port, or a
  worker node IP with the NodePort. Ensure your Broker Service type
  (`LoadBalancer` or `NodePort`) and any firewalls make that endpoint reachable
  from the Server.

```yaml
# Example Broker status populated by the Operator (out-of-cluster)
status:
  endpoint:
    external: a1b2c3...-elb.us-east-1.amazonaws.com:4222
    internal: 10.96.49.130:4222
```

Confirm the cluster supports the Service type you rely on (`LoadBalancer` or
`NodePort`) and that the resulting `external` endpoint is reachable from Meshery
Server. This is the most common multi-cloud connectivity pitfall—see
[Multi-Cluster &amp; Multi-Cloud](/installation/production/multi-cluster-and-multi-cloud).

## Egress requirements

Meshery Server typically needs outbound connectivity to:

- **Each managed cluster's Kubernetes API server** (out-of-cluster topologies).
- **The Remote Provider** (e.g. `https://cloud.meshery.io`) for authentication,
  capabilities, and durable state—over `443/tcp`.
- **Container registries and content sources** for images and bundled content
  during install/upgrade. You can reduce or disable some content downloads with
  `SKIP_DOWNLOAD_CONTENT` and `SKIP_DOWNLOAD_EXTENSIONS` in tightly egress-
  restricted environments.

In restricted networks, allow-list the Remote Provider and registry endpoints,
or run behind an egress proxy. Verify that provider egress is not silently
blocked—failed provider reachability shows up as authentication problems.

## Network policies

Default-deny network policies are a strong production baseline. When you apply
them, explicitly allow:

1. **Ingress to Meshery Server** on `8080` (in-cluster, behind the Service/
   ingress) from your ingress controller, and to the Broker on `4222` from
   Meshery Server.
2. **Egress from Meshery Server** to each cluster's API server, to the Broker,
   to the Remote Provider (`443`), and to any adapters.
3. **Intra-namespace** traffic among Meshery Server, Operator, MeshSync, and
   Broker as required for discovery and eventing.

Scope policies per namespace (deploy Meshery into a dedicated `meshery`
namespace) and pair them with the RBAC and security-context hardening in
[Security Hardening](/installation/production/security-hardening).

## Connectivity validation

After deployment, validate connectivity end-to-end:

- `mesheryctl system check` runs pre- and post-deployment health checks,
  including connectivity. See the
  [reference](/reference/mesheryctl/system/check).
- In the UI, the connection chip for each cluster reflects live connectivity;
  the Broker and Operator/MeshSync follow the connection's lifecycle.
- For Broker/Operator/MeshSync connectivity problems, use the
  [Operator &amp; MeshSync troubleshooting guide](/guides/troubleshooting/meshery-operator-meshsync).

{{< related-discussions tag="meshery" >}}
