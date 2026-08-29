# Multi-Cluster & Multi-Cloud Operations

> Managed vs. unmanaged cluster connections, one Operator per cluster, kubeconfig and context management, MeshSync deployment modes, and cloud-specific guidance for operating Meshery across many clusters and clouds.

Source: /pr-preview/pr-21670/installation/production/multi-cluster-and-multi-cloud/

A core reason to run Meshery as a central management plane is to operate many
clusters—often across multiple clouds—from one place. This page covers how
Meshery connects to clusters, the difference between managed and unmanaged
connections, how discovery is deployed per cluster, and the cloud-specific
details that decide whether a fleet "just works."

## How Meshery connects to a cluster

A Meshery–cluster connection is established by providing Meshery Server access to
the cluster's Kubernetes API:

- **In-cluster**, Meshery uses its in-cluster ServiceAccount for the cluster it
  runs in.
- **Out-of-cluster**, Meshery uses a **kubeconfig context** per cluster.
  `mesheryctl system start` or uploading a kubeconfig in the UI creates the
  connection; each context becomes a managed connection.

On connection, Meshery deploys one [Operator](/pr-preview/pr-21670/concepts/architecture/operator/)
per cluster (unless you use embedded MeshSync—see below), which manages that
cluster's [MeshSync](/pr-preview/pr-21670/concepts/architecture/meshsync/) and
[Broker](/pr-preview/pr-21670/concepts/architecture/broker/). On disconnect, those components are
removed. The Operator is also controllable from the UI's on/off switch
independently of the connection.

## Managed vs. unmanaged connections

"Managed" and "unmanaged" can mean two different—and both relevant—things in a
multi-cluster context. Be clear about which you mean:

### Meshery-managed discovery vs. library (embedded) discovery

| Connection style | What's deployed into the cluster | Trade-offs |
| :--- | :--- | :--- |
| **Operator-managed** | Operator + MeshSync + Broker run in the managed cluster. | Full, event-driven discovery with in-cluster components; the Operator self-heals them. Requires permission to deploy into the cluster. |
| **Embedded (library, default)** | Nothing—MeshSync runs as a library inside Meshery Server for that connection. | No in-cluster footprint; useful where you can't or won't deploy the Operator. Shifts discovery work into the Server process. |

The mode for new connections is governed by
`MESHSYNC_DEFAULT_DEPLOYMENT_MODE` (`operator` or `embedded`), which **defaults
to `embedded`**, and the mode can be switched per connection on the connections
page. Switching from operator to
embedded undeploys the in-cluster components and starts the in-Server routine;
switching back redeploys them. Choose per cluster based on whether in-cluster
deployment is acceptable and on the Server's capacity to absorb embedded
discovery (see
[Infrastructure, Sizing & Performance](/pr-preview/pr-21670/installation/production/infrastructure-sizing-and-performance/)).

### Cloud-managed vs. self-managed Kubernetes

Separately, the clusters themselves may be **cloud-managed** (EKS, GKE, AKS,
and similar) or **self-managed**. Meshery connects to both the same way, but
cloud-managed clusters differ in how they surface node permissions and
load-balanced endpoints, covered under [cloud-specific notes](#cloud-specific-guidance).

You can also disable Operator deployment entirely for a deployment with
`DISABLE_OPERATOR=true`, which prevents Meshery from automatically deploying the
Operator into connected clusters—useful when discovery is handled in embedded
mode or by policy.

## Managing kubeconfig and contexts at scale

For a fleet, kubeconfig/context management is the operational backbone:

- **One context per cluster**, each scoped to a **least-privilege credential**
  for that cluster (see
  [Security Hardening](/pr-preview/pr-21670/installation/production/security-hardening/)). Avoid a
  single all-powerful credential spanning the fleet.
- Mount kubeconfig from a Secret and point Meshery at it with
  `KUBECONFIG_FOLDER` (default `~/.kube`). Keep context names stable and
  meaningful.
- Treat the **set of connections as version-controlled configuration** so the
  fleet can be reconstructed during recovery (see
  [High Availability & Resiliency](/pr-preview/pr-21670/installation/production/high-availability-and-resiliency/)).
- Prefer short-lived or provider-issued credentials where the cloud supports
  them, and rotate per-cluster credentials independently.

## Reaching each cluster's Broker

In multi-cluster/out-of-cluster operation, Meshery Server must reach each
cluster's Broker on `4222/tcp`. The Operator publishes a reachable `external`
endpoint into the Broker custom resource `status`, selecting (in order of
preference) the LoadBalancer hostname, the LoadBalancer IP, the kubeconfig host
with NodePort, the ClusterIP with cluster port, or a worker node IP with
NodePort.

The practical implications across a fleet:

- Ensure each cluster's Broker Service type (`LoadBalancer` or `NodePort`) is
  supported and that the resulting endpoint is reachable from the Server.
- Restrict that exposure to the Server's network origin (security groups,
  load-balancer source ranges, private connectivity). See
  [Networking & Connectivity](/pr-preview/pr-21670/installation/production/networking-and-connectivity/).

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">The #1 multi-cloud pitfall</div>


The most common cross-cloud failure is an unreachable Broker endpoint: the
cluster publishes a LoadBalancer hostname or NodePort the central Server cannot
reach (blocked by security groups, private subnets, or missing routes).
Validate Broker reachability from the Server for **every** cluster you add.
</div>


## Cloud-specific guidance

The architecture is identical across clouds; these are the per-provider details
that matter:

- **Node-watch RBAC.** Full discovery on AKS, AWS, and GCP may require
  permission to watch nodes. Enable `rbac.nodes: true` on those clusters (it
  defaults to `false`). Grant it only where needed.
- **Load-balanced Broker endpoints.** Clouds differ in whether a `LoadBalancer`
  Service surfaces a **hostname** (commonly AWS ELB) or an **IP** (commonly GCP/
  Azure). The Operator's endpoint selection handles both, but your firewall and
  reachability checks must account for the form your cloud uses.
- **Private clusters.** For private API servers or nodes (private EKS/GKE/AKS),
  the central Server needs network reachability to both the API server and the
  Broker—via VPC/VNet peering, private load balancers, transit gateways, or a
  VPN. Public exposure is discouraged (see
  [Security Hardening](/pr-preview/pr-21670/installation/production/security-hardening/)).
- **Cross-region latency.** A central Server managing distant clusters incurs
  latency on discovery streaming and API calls. Keep it acceptable, and consider
  regional management planes if a single central plane spans high-latency links.

## Operating a fleet

- **Per-cluster connection health.** Each connection's chip in the UI reflects
  live connectivity; Broker/Operator/MeshSync follow the connection lifecycle.
  Watch these as fleet KPIs (see
  [Monitoring, Observability & Health KPIs](/pr-preview/pr-21670/installation/production/monitoring-observability-and-kpis/)).
- **Per-cluster discovery scope.** Tune MeshSync's `informer_config` blacklist
  per cluster to control footprint on large clusters
  ([sizing](/pr-preview/pr-21670/installation/production/infrastructure-sizing-and-performance/)).
- **Blast-radius isolation.** Distinct per-cluster credentials and network
  policies mean a problem in one cluster's connection does not cascade across
  the fleet.
- **Consistent lifecycle.** Keep the Operator/MeshSync mode consistent with your
  policy across clusters, and codify it via `MESHSYNC_DEFAULT_DEPLOYMENT_MODE`
  and your connection configuration.

## Multi-cluster checklist

- [ ] One kubeconfig context per cluster, each with a least-privilege
      credential.
- [ ] Broker endpoint reachable from the Server and locked to the Server's
      origin for every cluster.
- [ ] `rbac.nodes` enabled only on clusters that require node watching.
- [ ] MeshSync mode (operator vs. embedded) chosen deliberately per cluster.
- [ ] Private connectivity (peering/VPN) for private clusters; no broad public
      Broker exposure.
- [ ] Cross-region latency assessed; regional planes considered if needed.
- [ ] Connection set and per-cluster config under version control for recovery.

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
