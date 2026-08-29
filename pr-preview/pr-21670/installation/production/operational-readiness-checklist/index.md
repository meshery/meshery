# Operational Readiness Checklist & Known Caveats

> A consolidated, actionable production-readiness checklist across every dimension, the upgrade strategy for Meshery, and the known caveats and limitations to plan around.

Source: /pr-preview/pr-21670/installation/production/operational-readiness-checklist/

This page brings the whole set together: a single checklist you can work through
before declaring a Meshery deployment production-ready, the upgrade strategy to
adopt, and the known caveats to design around. Each item links back to the page
where it is explained in full.

## Upgrade strategy

Meshery is a composition of components—Server, UI, Operator, MeshSync, Broker,
Adapters, and `mesheryctl`—some of which upgrade together and some
independently. For Kubernetes production deployments:

- **Use Helm and pin versions.** Upgrade with `helm upgrade`, and pin the image
  to an **immutable version tag** rather than tracking `stable-latest`, so
  upgrades are deliberate and reproducible
  ([security hardening](/pr-preview/pr-21670/installation/production/security-hardening/)).
- **Use upgrade-friendly probe values.** Capability reloading can make the
  Server briefly unavailable during an upgrade. The chart provides
  `values-upgrade.yaml` (startup probe, higher failure thresholds) so pods are
  not killed mid-initialization:

  ```bash
  helm upgrade meshery meshery/meshery --namespace meshery \
    -f https://raw.githubusercontent.com/meshery/meshery/master/install/kubernetes/helm/meshery/values-upgrade.yaml \
    --wait --timeout 10m
  ```

  For reproducible production upgrades, download `values-upgrade.yaml`,
  version-control it alongside your own values, and reference your local copy
  (with a pinned chart and image version) rather than fetching from `master` at
  upgrade time.

- **Choose a release channel deliberately.** `stable` for production; `edge`
  only for testing. The channel is reflected by `RELEASE_CHANNEL`.
- **Mind component groupings.** See the
  [Upgrade Guide](/pr-preview/pr-21670/installation/upgrades/) for which components move together
  (e.g., Server/UI/Load Generators/Database) versus independently (Adapters,
  `mesheryctl`). The Operator and its controllers effectively move **with the
  Server**: the chart version a Server deploys to managed clusters tracks its
  own release, and the operator image is pinned inside that chart.
- **Expect a published chart, not necessarily the matching one.** Before
  installing, the Server validates the version it wants against the chart
  repository's published index, so it never hands Helm a chart that does not
  exist. A version that is not published yet resolves to the newest published
  chart, and a version below the minimum deployable chart is raised to the
  oldest published chart at or above it; release candidates are never chosen
  automatically. Every substitution is reported in the events feed. Budget
  outbound access from Meshery Server to the chart repository accordingly
  ([egress requirements](/pr-preview/pr-21670/installation/production/networking-and-connectivity/#egress-requirements)).
- **Let the Server own the Operator.** On Server-managed clusters, do not
  hand-upgrade or hand-configure the `meshery-operator` Helm release — the
  Server's reconciliation re-applies the chart version it resolves and reverts
  manual changes. Upgrading the Server is the supported way to upgrade the
  Operator; the chart's CRD update Job refreshes the `Broker`/`MeshSync` CRD
  schemas on each upgrade (Helm alone never updates CRDs). When you must hold a
  cluster at a specific operator chart, set `operator.version` on that
  connection - an explicit pin is honored and never substituted, and is
  refused with a visible error if it names an unpublished or moving version.
  See
  [How Meshery Server manages Meshery Operator](/pr-preview/pr-21670/installation/upgrades/#how-meshery-server-manages-meshery-operator).
- **Expect CRDs to persist.** The `brokers.meshery.io` and
  `meshsyncs.meshery.io` CRDs (and their objects) deliberately survive
  operator uninstalls and Helm release deletion; include the explicit
  `kubectl delete crd ...` step in decommissioning runbooks only when
  permanent removal is intended.
- **Rehearse and roll back.** Because durable state lives with the Remote
  Provider and the local database is a cache, rolling back the deployment is
  low-risk for data—validate the rollback path anyway. On rollback, managed
  clusters converge back to the operator chart the older Server resolves,
  automatically. The minimum deployable chart version guards that derived
  default, but it is not a blanket guarantee: an explicit `operator.version` is
  honored as written even when it names a chart below the minimum, and if the
  repository publishes nothing at or above the minimum the newest published
  chart is deployed with a warning that the Operator may not become ready.
- **Edge caches need no purge.** If a CDN or caching reverse proxy fronts
  Meshery, its UI cache busts itself on upgrade—hashed asset URLs change and the
  HTML `ETag` follows the build/release version—so a manual cache purge is not
  required in your upgrade pipeline
  ([networking](/pr-preview/pr-21670/installation/production/networking-and-connectivity/)).

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Upgrades and the cache</div>


An upgrade that recreates the Server pod discards the local cache; this is
expected and safe. The Remote Provider preserves durable state and MeshSync
rebuilds the cluster snapshot. Do not block upgrades trying to preserve the
ephemeral database.
</div>


## Known caveats and limitations

Plan around these known characteristics—they are by design and shape good
production architecture:

- **The database is an ephemeral, single-writer cache.** It is SQLite/Bitcask on
  local disk, not a shared multi-writer datastore. Durable state must come from
  a [Remote Provider](/pr-preview/pr-21670/installation/production/authentication-and-identity/);
  this also shapes how Server replication behaves
  ([HA & Resiliency](/pr-preview/pr-21670/installation/production/high-availability-and-resiliency/)).
- **The Broker persists messages in memory only.** No persistent volume is used.
  Brief interruptions are bridged by NATS topic persistence, but a far-behind or
  unavailable Server consumer lets Broker memory climb
  ([sizing](/pr-preview/pr-21670/installation/production/infrastructure-sizing-and-performance/)).
- **The Local Provider has no authentication.** Never run shared production on
  it; preselect a Remote Provider
  ([identity](/pr-preview/pr-21670/installation/production/authentication-and-identity/)).
- **WebSocket support is mandatory at the ingress.** Without it the UI loads but
  never updates live
  ([networking](/pr-preview/pr-21670/installation/production/networking-and-connectivity/)).
- **OAuth requires the external callback URL.** Behind a proxy, an unset or wrong
  `MESHERY_SERVER_CALLBACK_URL` breaks login.
- **Out-of-cluster requires a reachable Broker endpoint.** Cross-cloud
  reachability of the Broker is the most common multi-cloud failure
  ([multi-cloud](/pr-preview/pr-21670/installation/production/multi-cluster-and-multi-cloud/)).
- **Discovery cost scales with cluster size and churn.** Without MeshSync
  blacklisting, very large/churny clusters dominate Server memory and database
  growth.
- **Policy evaluation is time-boxed.** Large designs may hit
  `POLICY_EVAL_TIMEOUT` (default `3m`); tune CPU and the timeout together.
- **The chart ships permissive defaults for portability.** No resource
  requests/limits, empty security contexts, and a `stable-latest`,
  `Always`-pull image are convenient for getting started but are **not**
  production settings—set them explicitly.

## The consolidated production-readiness checklist

Work through these before going live. Each group links to its source page.

### Deployment model

- [ ] Topology chosen deliberately (in-cluster vs. out-of-cluster; Docker vs.
      Kubernetes), with Kubernetes + Helm for production.
      [Deployment Models](/pr-preview/pr-21670/installation/production/deployment-models/)
- [ ] One Operator/MeshSync/Broker per managed cluster (or embedded mode chosen
      deliberately).

### Sizing & performance

- [ ] Explicit CPU/memory **requests and limits** set for Server, Operator,
      MeshSync, and Broker.
      [Sizing](/pr-preview/pr-21670/installation/production/infrastructure-sizing-and-performance/)
- [ ] Server memory headroom for the MeshSync snapshot; Broker memory headroom
      for bursts.
- [ ] MeshSync `informer_config` blacklist tuned for large clusters.
- [ ] Storage provisioned for the Server cache; autoscaling configured where
      appropriate.

### High availability & resiliency

- [ ] Liveness/readiness probes enabled and tuned.
      [HA & Resiliency](/pr-preview/pr-21670/installation/production/high-availability-and-resiliency/)
- [ ] Anti-affinity spreads replicas across nodes/zones.
- [ ] Helm values, env config, and connection definitions in version control
      (your real backup).
- [ ] Redeploy-from-config recovery drill validated.

### Networking & connectivity

- [ ] Ingress fronts Meshery with TLS; **WebSocket upgrades** verified.
      [Networking](/pr-preview/pr-21670/installation/production/networking-and-connectivity/)
- [ ] `MESHERY_SERVER_CALLBACK_URL` set to the external URL.
- [ ] Broker endpoint reachable from the Server and restricted to its origin.
- [ ] Egress to the Remote Provider and registries allowed; network policies
      applied.
- [ ] Any fronting CDN/caching proxy honors origin cache headers (immutable
      assets, `no-cache` HTML) and never caches `/api/*` or the WebSocket.

### Security & identity

- [ ] **Remote Provider preselected** (`PROVIDER`, `PROVIDER_BASE_URLS`); Local
      Provider not used.
      [Identity](/pr-preview/pr-21670/installation/production/authentication-and-identity/)
- [ ] Least-privilege RBAC; `rbac.nodes` only where required; per-cluster
      credentials.
- [ ] Hardened `securityContext` (non-root, dropped caps, read-only root FS with
      writable data volume) on Server and Operator.
      [Security Hardening](/pr-preview/pr-21670/installation/production/security-hardening/)
- [ ] Secrets (kubeconfig, provider, pull) sourced from Secrets/external manager.
- [ ] Images pinned to immutable tags from a trusted/mirrored registry.

### Extensions

- [ ] Every enabled extension attributed to a publisher you trust, and reviewed
      before it reached production.
      [Trusting an extension](/pr-preview/pr-21670/installation/production/security-hardening/#trusting-an-extension)
- [ ] Unused extension points left disabled; adapters enabled only where needed
      (`ADAPTER_URLS`, chart adapter subcharts off by default).
- [ ] Adapters given scoped ServiceAccounts (`serviceAccountNameOverride`) rather
      than sharing `meshery-server`, and reachable only from Meshery Server.
- [ ] Extension packages pinned (`SKIP_DOWNLOAD_EXTENSIONS`), and the capability
      set pinned (`PROVIDER_CAPABILITIES_FILEPATH`) where it should not drift.
      [Providers reference](/pr-preview/pr-21670/reference/extensibility/providers/#runtime-configuration-options)
- [ ] Extension versions validated against the Meshery version you run.
      [Extension compatibility](/pr-preview/pr-21670/reference/extensibility/verify-compatibility/)
- [ ] Egress policy covers every destination an enabled extension requires.
- [ ] Rollback path known: removing a server-side extension requires a Meshery
      Server restart, not just a configuration change.

### Multi-cluster & multi-cloud

- [ ] One least-privilege kubeconfig context per cluster.
      [Multi-Cloud](/pr-preview/pr-21670/installation/production/multi-cluster-and-multi-cloud/)
- [ ] Per-cloud node-watch RBAC and Broker endpoint reachability validated.
- [ ] Private connectivity for private clusters; cross-region latency assessed.

### Monitoring & observability

- [ ] External uptime check on `/healthz/ready`; workload metrics scraped.
      [Monitoring & KPIs](/pr-preview/pr-21670/installation/production/monitoring-observability-and-kpis/)
- [ ] Dedicated alerts on Server memory and Broker memory; per-cluster
      connection health and provider reachability monitored.
- [ ] OpenTelemetry tracing to a real collector; logs centralized at
      `LOG_LEVEL=4`.
- [ ] `mesheryctl system check` scheduled as a synthetic check.

### Operations

- [ ] Upgrade procedure uses `values-upgrade.yaml`; `stable` release channel.
      [Upgrade Guide](/pr-preview/pr-21670/installation/upgrades/)
- [ ] Runbooks reference the
      [troubleshooting guides](/pr-preview/pr-21670/guides/troubleshooting/meshery-operator-meshsync/).
- [ ] Ownership, on-call, and escalation defined for the Meshery deployment.

## Where to go next

- Revisit any area above via its linked page in this set.
- For component internals, see the
  [Meshery Architecture](/pr-preview/pr-21670/concepts/architecture/) reference.
- For runtime configuration, see the
  [Server Environment Variables](/pr-preview/pr-21670/installation/advanced/environment-variables/)
  reference.

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
