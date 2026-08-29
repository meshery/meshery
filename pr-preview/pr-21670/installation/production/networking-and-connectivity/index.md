# Networking & Connectivity

> Network ports and directional flows, ingress and Emissary configuration, secure WebSocket support, CDN and edge caching of the UI, Broker exposure, egress to Remote Providers, the OAuth callback URL, and network policies for production Meshery.

Source: /pr-preview/pr-21670/installation/production/networking-and-connectivity/

Getting networking right is what makes the difference between a Meshery
deployment that "works on my cluster" and one that reliably manages a fleet from
behind an ingress. This page covers the ports and directions traffic flows, how
to expose Meshery securely, how the Broker is reached (especially
out-of-cluster), and the egress Meshery needs.

## Network ports

Meshery uses the following ports across its components. The Broker ports beyond
the client port are for clustering/routing features and are not all in use in
current releases.

<table class="table table-striped">
  <thead>
    <tr>
      <th style="text-align: left; padding: 10px; border-bottom: 2px solid #ccc;">Component</th>
      <th style="text-align: center; padding: 10px; border-bottom: 2px solid #ccc;">Port</th>
      <th style="text-align: left; padding: 10px; border-bottom: 2px solid #ccc;">Purpose</th>
    </tr>
  </thead>
  <tbody>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">Meshery Server</td>
      <td style="text-align: center; padding: 10px;">9081/tcp</td>
      <td style="padding: 10px;">UI, REST, and GraphQL APIs</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">Meshery Server</td>
      <td style="text-align: center; padding: 10px;">80/tcp</td>
      <td style="padding: 10px;">WebSocket</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;"><a href="/pr-preview/pr-21670/concepts/architecture/broker/">Meshery Broker</a></td>
      <td style="text-align: center; padding: 10px;">4222/tcp</td>
      <td style="padding: 10px;">Client communication with Meshery Server</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;"><a href="/pr-preview/pr-21670/concepts/architecture/broker/">Meshery Broker</a></td>
      <td style="text-align: center; padding: 10px;">8222/tcp</td>
      <td style="padding: 10px;">HTTP management port for monitoring Meshery Broker. Available as of Meshery v0.5.0</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;"><a href="/pr-preview/pr-21670/concepts/architecture/broker/">Meshery Broker</a></td>
      <td style="text-align: center; padding: 10px;">6222/tcp</td>
      <td style="padding: 10px;">Routing port for Broker clustering. Unused as of Meshery v0.6.0-rc-2</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;"><a href="/pr-preview/pr-21670/concepts/architecture/broker/">Meshery Broker</a></td>
      <td style="text-align: center; padding: 10px;">7422/tcp</td>
      <td style="padding: 10px;">Incoming/outgoing leaf node connections. Unused as of Meshery v0.6.0-rc-2</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;"><a href="/pr-preview/pr-21670/concepts/architecture/broker/">Meshery Broker</a></td>
      <td style="text-align: center; padding: 10px;">7522/tcp</td>
      <td style="padding: 10px;">Gateway to gateway communication. Unused as of Meshery v0.6.0-rc-2</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;"><a href="/pr-preview/pr-21670/concepts/architecture/broker/">Meshery Broker</a></td>
      <td style="text-align: center; padding: 10px;">7777/tcp</td>
      <td style="padding: 10px;">Used for Prometheus NATS Exporter. Unused as of Meshery v0.6.0-rc-2</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;"><a href="/pr-preview/pr-21670/reference/extensibility/providers/">Meshery Remote Providers</a></td>
      <td style="text-align: center; padding: 10px;">443/tcp</td>
      <td style="padding: 10px;">e.g. Meshery Cloud</td>
    </tr>
    <tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;"><a href="/pr-preview/pr-21670/guides/performance-management/managing-performance/">Meshery Perf</a></td>
      <td style="text-align: center; padding: 10px;">10013/gRPC</td>
      <td style="padding: 10px;">Performance Management</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/images/traefik-mesh.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/images/traefik-mesh.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/images/traefik-mesh.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/traefik-mesh/" style="text-decoration: none;">Meshery Adapter for Traefik Mesh</a>
      </td>
      <td style="text-align: center; padding: 10px;">10006/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/nginx-sm/images/nginx-sm.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/nginx-sm/images/nginx-sm-white.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/nginx-sm/images/nginx-sm.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/nginx-sm/" style="text-decoration: none;">Meshery Adapter for NGINX Service Mesh</a>
      </td>
      <td style="text-align: center; padding: 10px;">10010/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/nsm/images/nsm.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/nsm/images/nsm.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/nsm/images/nsm.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/nsm/" style="text-decoration: none;">Meshery Adapter for Network Service Mesh</a>
      </td>
      <td style="text-align: center; padding: 10px;">10004/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/linkerd/images/linkerd.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/linkerd/images/linkerd-white.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/linkerd/images/linkerd.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/linkerd/" style="text-decoration: none;">Meshery Adapter for Linkerd</a>
      </td>
      <td style="text-align: center; padding: 10px;">10001/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/kuma/images/kuma.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/kuma/images/kuma-white.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/kuma/images/kuma.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/kuma/" style="text-decoration: none;">Meshery Adapter for Kuma</a>
      </td>
      <td style="text-align: center; padding: 10px;">10007/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/istio/images/istio.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/istio/images/istio-white.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/istio/images/istio.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/istio/" style="text-decoration: none;">Meshery Adapter for Istio</a>
      </td>
      <td style="text-align: center; padding: 10px;">10000/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/consul/images/consul.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/consul/images/consul-white.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/consul/images/consul.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/consul/" style="text-decoration: none;">Meshery Adapter for Consul</a>
      </td>
      <td style="text-align: center; padding: 10px;">10002/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/cilium/images/cilium.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/cilium/images/cilium-white.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/cilium/images/cilium.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/cilium/" style="text-decoration: none;">Meshery Adapter for Cilium Service Mesh</a>
      </td>
      <td style="text-align: center; padding: 10px;">10012/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/app-mesh/images/app-mesh.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/app-mesh/images/app-mesh-white.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/app-mesh/images/app-mesh.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/app-mesh/" style="text-decoration: none;">Meshery Adapter for App Mesh</a>
      </td>
      <td style="text-align: center; padding: 10px;">10005/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr><tr style="border-bottom: 1px solid #eee;">
      <td style="padding: 10px;">
        <img src="/pr-preview/pr-21670/extensions/adapters/tanzu-sm/images/tanzu.svg" style="width:24px; height:24px; vertical-align: middle; margin-right: 8px;" 
             data-logo-for-dark="/pr-preview/pr-21670/extensions/adapters/tanzu-sm/images/tanzu.svg" 
             data-logo-for-light="/pr-preview/pr-21670/extensions/adapters/tanzu-sm/images/tanzu.svg" 
             loading="lazy"/>
        <a href="/pr-preview/pr-21670/extensions/adapters/tanzu-sm/" style="text-decoration: none;">Meshery Adapter for Tanzu Service Mesh</a>
      </td>
      <td style="text-align: center; padding: 10px;">10011/gRPC/gRPC</td>
      <td style="padding: 10px;">Communication with Meshery Server</td>
    </tr></tbody>
</table>


In addition, optional [Meshery Adapters](/pr-preview/pr-21670/concepts/architecture/adapters/) listen
on their own ports (for example `10000/tcp` and up), and the Meshery Operator's
authn/authz-filtered metrics endpoint listens on `8443/tcp` in-cluster. Inside the pod, Meshery
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
| Meshery Server | Meshery Helm chart repository (`meshery.github.io`) | 443/tcp (egress) | Resolving and downloading the `meshery-operator` chart | Operator mode only; see [Egress requirements](#egress-requirements). |
| Meshery Server | Meshery Adapters | adapter ports (e.g. 10000+/tcp) | Capability operations | Only when adapters are deployed. |
| Operations/monitoring | Meshery Broker | 8222/tcp | Broker HTTP monitoring endpoint | Optional; for observability. |

The `9081/tcp` (API) and `80/tcp` (WebSocket) values are Meshery Server's
component ports as listed in the [port table above](#network-ports). In a
production deployment you do not expose these directly: front them with a single
external HTTPS endpoint at the ingress, which serves both the APIs (`https://`)
and the WebSocket (`wss://`) over `443`—see [below](#exposing-meshery-securely-ingress-tls-websocket).

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Direction matters for out-of-cluster</div>


When Meshery Server is out-of-cluster, the Server initiates connections **to**
each managed cluster's API server and Broker. Those endpoints must be reachable
from wherever Meshery Server runs—through cloud load balancers, VPN/peering, or
NodePort exposure. See [Broker exposure](#exposing-the-broker) below.
</div>


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

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">WebSocket is not optional</div>


A common production misconfiguration is an ingress that proxies HTTP but drops
WebSocket upgrades. The symptom is a UI that loads but never updates in real
time. Validate WebSocket upgrade support explicitly during rollout.
</div>


## Caching the UI at the edge (CDN or caching reverse proxy)

Meshery Server serves its web UI—a static [Next.js](https://nextjs.org/)
export—directly. To make that UI safely cacheable by a CDN or caching reverse
proxy placed in front of Meshery, the Server emits **release-scoped cache
headers**: the edge can serve UI content for the lifetime of a release and pick
up new content automatically after an upgrade, **without any cache purge**. The
build/release version is the only control point.

Two tiers of response are handled differently:

| Response | `Cache-Control` | Validator | Why |
| :--- | :--- | :--- | :--- |
| **Immutable versioned assets** under `/_next/static/{chunks,css,media}/...` | `public, max-age=31536000, immutable` | — | The build content hash is in the URL, so the bytes at a given URL never change. A new release ships new hashes (new URLs); old entries simply fall out of use—no purge. |
| **HTML documents** (`/` → `index.html`, `<route>.html`) | `public, no-cache` | `ETag: "<build version>"` | HTML lives at stable URLs, so it must always be revalidated. While the release is unchanged the origin answers `If-None-Match` with `304`; a new release changes the build version → changes the `ETag` → the next revalidation returns fresh HTML. |

Because browser caches cannot be purged, HTML is never given a positive
`max-age`; it is always revalidated against the build-version `ETag`. And only
existing files receive these headers—a transient `404` (for example, a
content-hashed asset requested mid-deploy against a pod that does not have it
yet) is returned **without** cache headers, so it is never pinned as `immutable`
and cached for a year.

### Configuring the CDN or proxy

Whatever you place in front of Meshery—a cloud CDN (CloudFront, Cloudflare,
Fastly), an NGINX or Varnish caching proxy, or your ingress controller's
cache—configure it to **respect the origin's caching semantics** rather than
impose its own:

- **Honor the origin `Cache-Control` and `ETag`.** Do not strip or rewrite them,
  and do not apply a blanket TTL that overrides `no-cache` on HTML—doing so risks
  serving a stale HTML shell across an upgrade.
- **Forward conditional requests.** Pass `If-None-Match` through to the origin so
  HTML can revalidate and the Server can answer `304 Not Modified`.
- **Never cache the API or WebSocket.** Only the static UI layer carries these
  headers. API traffic (`/api/*`) and the WebSocket/event stream must pass
  through to the origin uncached; caching them breaks live updates and
  authentication.
- **Do not cache negative responses.** The Server deliberately leaves transient
  `404`s uncacheable; ensure your edge does not cache `404`s for hashed-asset
  URLs requested mid-deploy.

The headers are computed against the base-path-stripped request URL, so they
apply correctly when Meshery is served under a sub-path as well.

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">No cache purge on upgrade</div>


Because immutable asset URLs change on every build and the HTML `ETag` is the
build/release version, a Meshery upgrade is reflected at the edge
**automatically**. You do not need to wire a CDN purge into your upgrade
pipeline—see the
[upgrade strategy](/pr-preview/pr-21670/installation/production/operational-readiness-checklist/).
</div>


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
[Authentication, Authorization & Identity](/pr-preview/pr-21670/installation/production/authentication-and-identity/)
and the
[environment variables reference](/pr-preview/pr-21670/installation/advanced/environment-variables/).

## Exposing the Broker

The [Meshery Broker](/pr-preview/pr-21670/concepts/architecture/broker/) is reached by Meshery Server
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
[Multi-Cluster & Multi-Cloud](/pr-preview/pr-21670/installation/production/multi-cluster-and-multi-cloud/).

## Egress requirements

Meshery Server typically needs outbound connectivity to:

- **Each managed cluster's Kubernetes API server** (out-of-cluster topologies).
- **The Remote Provider** (e.g. `https://cloud.meshery.io`) for authentication,
  capabilities, and durable state—over `443/tcp`.
- **Container registries and content sources** for images and bundled content
  during install/upgrade. You can reduce or disable some content downloads with
  `SKIP_DOWNLOAD_CONTENT` and `SKIP_DOWNLOAD_EXTENSIONS` in tightly egress-
  restricted environments.
- **The Meshery Helm chart repository**
  (`https://meshery.github.io/meshery.io/charts`) over `443/tcp`, when managing
  clusters in operator mode. Meshery Server reads the repository index to
  resolve the `meshery-operator` chart version before installing it, and
  downloads the chart archive itself. Clusters whose Operator is already
  installed keep reporting status and version without it - only installing or
  upgrading the Operator is withheld, with the reason shown on the connection's
  Operator status card and in its
  [Diagnostics](/pr-preview/pr-21670/guides/infrastructure-management/kubernetes-connection-lifecycle/#diagnostics).
  A blip is not sticky: redeploying the Operator re-resolves the version.

In restricted networks, allow-list the Remote Provider, the chart repository,
and registry endpoints, or run behind an egress proxy. Verify that provider
egress is not silently blocked - failed provider reachability shows up as
authentication problems.

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
[Security Hardening](/pr-preview/pr-21670/installation/production/security-hardening/).

## Connectivity validation

After deployment, validate connectivity end-to-end:

- `mesheryctl system check` runs pre- and post-deployment health checks,
  including connectivity. See the
  [reference](/pr-preview/pr-21670/reference/references/mesheryctl/system/check/).
- In the UI, the connection chip for each cluster reflects live connectivity;
  the Broker and Operator/MeshSync follow the connection's lifecycle.
- For Broker/Operator/MeshSync connectivity problems, use the
  [Operator & MeshSync troubleshooting guide](/pr-preview/pr-21670/guides/troubleshooting/meshery-operator-meshsync/).

<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/design-meshery-mcp-server-architecture-and-registration-interface/7954" target="_blank" rel="noopener noreferrer">
            Design: Meshery MCP Server architecture and registration interface
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/the-meshery-mcp-server-foundation-is-up-lets-agree-on-what-to-build-next/7952" target="_blank" rel="noopener noreferrer">
            The Meshery MCP Server foundation is up, let&#39;s agree on what to build next
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/new-intro-topic/7975" target="_blank" rel="noopener noreferrer">
            New intro topic
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-mcp-server-poc-an-ai-agent-managing-kubernetes-through-mcp/7974" target="_blank" rel="noopener noreferrer">
            Meshery MCP Server POC: an AI agent managing Kubernetes through MCP
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/approach-for-context-window-aware-retrieval-in-the-ai-adapter-issue-20994/7963" target="_blank" rel="noopener noreferrer">
            Approach for context-window-aware retrieval in the AI Adapter (Issue #20994)
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/there-is-some-error-in-running-the-localhost-of-layer5-in-my-server/7818" target="_blank" rel="noopener noreferrer">
            There is some error in running the localhost of layer5 in my server
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/rfc-aligning-the-foundation-for-the-meshery-mcp-server/7913" target="_blank" rel="noopener noreferrer">
            RFC: Aligning the Foundation for the Meshery MCP Server
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-development-meeting-august-12th-2026/7948" target="_blank" rel="noopener noreferrer">
            Meshery Development Meeting | August 12th, 2026
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
