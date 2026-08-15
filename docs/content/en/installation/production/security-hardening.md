---
title: Security Hardening
linkTitle: Security Hardening
description: >-
  RBAC and least privilege, pod and container security contexts, secret and
  kubeconfig handling, TLS, supply-chain integrity, Broker exposure risk,
  namespace isolation, and extension trust for hardening Meshery in production.
categories: [installation]
weight: 50
aliases:
- /installation/production/security
- /installation/production/hardening
---

Meshery is a management plane with broad visibility into—and, when you use it
for lifecycle management, control over—the infrastructure it connects to. That
makes hardening it a priority. This page covers the controls that reduce
Meshery's attack surface and blast radius in production. Identity and provider
choices are covered separately in
[Authentication, Authorization & Identity]({{< ref "installation/production/authentication-and-identity.md" >}});
read both together.

{{% alert title="Start with identity" color="warning" %}}
The highest-impact security decision is to **preselect a Remote Provider** so
production never runs with the unauthenticated Local Provider, and so you
control which identity providers are accepted. That topic has its own page:
[Authentication, Authorization & Identity]({{< ref "installation/production/authentication-and-identity.md" >}}).
{{% /alert %}}

## Least-privilege RBAC

Meshery and its components act on Kubernetes through their ServiceAccounts.
Grant only what is needed.

- **Scope the ServiceAccount.** The Helm chart uses a `meshery-server`
  ServiceAccount. Bind it to the narrowest roles that still allow your intended
  use. A read-only/discovery deployment needs far less than one performing
  lifecycle management (create/update/delete) of workloads.
- **Enable node watching only where required.** The chart's `rbac.nodes`
  defaults to `false`. Some managed platforms (AKS, AWS, GCP) require permission
  to watch nodes for full discovery; enable `rbac.nodes: true` **only** on the
  clusters that need it rather than globally.
- **Separate duty by cluster.** In out-of-cluster and multi-cloud topologies,
  use a distinct, minimally scoped credential per managed cluster so a single
  compromised credential cannot reach every cluster. See
  [Multi-Cluster & Multi-Cloud]({{< ref "installation/production/multi-cluster-and-multi-cloud.md" >}}).
- **Review adapter permissions.** By default adapters share the Meshery
  ServiceAccount's permissions. If you deploy adapters, consider distinct,
  scoped ServiceAccounts (`serviceAccountNameOverride`) so an adapter does not
  inherit broader rights than it needs.

## Pod and container security context

The Helm chart leaves `podSecurityContext` and `securityContext` empty by
default so Meshery starts in permissive environments. For production, set them
explicitly to run as an unprivileged, locked-down workload:

```yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 2000

securityContext:
  runAsNonRoot: true
  runAsUser: 1000
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities:
    drop:
      - ALL
  seccompProfile:
    type: RuntimeDefault
```

Notes:

- **`readOnlyRootFilesystem: true`** requires that Meshery's writable data
  folder is a mounted, writable volume. Ensure the data directory
  (`USER_DATA_FOLDER`, default under the app user's home) is backed by a
  writable mount so the read-only root does not break the on-disk cache.
- Validate the chosen UID/GID against the image and your `fsGroup` so the cache
  directory is writable.
- Apply equivalent hardening to the **Meshery Operator** (its chart also exposes
  `podSecurityContext`/`securityContext`).

These settings align Meshery with restricted Pod Security Standards. Where you
enforce Pod Security admission, target the **restricted** profile and test
startup.

## Secrets and kubeconfig handling

Meshery's most sensitive inputs are the credentials it uses to reach clusters
and providers.

- **kubeconfig.** For out-of-cluster deployments, Meshery reads Kubernetes
  configuration from `KUBECONFIG_FOLDER` (default `~/.kube`). Mount kubeconfig
  from a Kubernetes Secret rather than baking it into an image, restrict its
  file permissions, and scope each context's credential to least privilege.
  Prefer short-lived or provider-issued credentials where your platform supports
  them.
- **In-cluster credentials.** When Meshery runs in-cluster and manages its own
  cluster, it can use the in-cluster ServiceAccount rather than a kubeconfig
  file—one fewer long-lived secret to manage.
- **Provider and OAuth secrets.** Keep Remote Provider configuration and any
  client secrets in Kubernetes Secrets (or an external secrets manager), not in
  plaintext values files committed to source control.
- **`imagePullSecrets`.** If you mirror images to a private registry, supply
  pull secrets via the chart rather than node-level credentials.

{{% alert title="Don't commit secrets to values.yaml" color="warning" %}}
Keep non-secret configuration in version-controlled Helm values, but source
secrets (kubeconfig, provider client secrets, pull secrets) from Kubernetes
Secrets or an external secrets manager. This keeps your GitOps-friendly
configuration safe to store and share.
{{% /alert %}}

## TLS everywhere

- **User-facing traffic.** Terminate TLS at the ingress so browsers use
  `https://` and `wss://`. See
  [Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}}).
- **Provider egress.** Communication with the Remote Provider is over HTTPS
  (`443`). Do not disable certificate verification or route it through an
  intercepting proxy that breaks trust.
- **Avoid the playground/insecure shortcuts in production.** Development helpers
  that relax security (for example insecure tracing endpoints) belong in
  non-production only.

## Broker exposure risk

When Meshery Server is out-of-cluster, the Broker (`4222/tcp`) is exposed beyond
the cluster so the Server can reach it. Treat that exposure deliberately:

- Restrict the Broker's reachable endpoint to the Meshery Server's network
  origin using load-balancer source ranges, security groups, or firewall rules—
  do not expose `4222` to the public internet broadly.
- Prefer private connectivity (VPC peering, private load balancers, VPN) between
  an out-of-cluster Server and managed-cluster Brokers over public exposure.
- Pair exposure with network policies that allow Broker ingress only from the
  Server. See
  [Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}}).

## Namespace isolation and multi-tenancy

- **Dedicated namespace.** Deploy Meshery into its own namespace (for example
  `meshery`) so RBAC, network policies, quotas, and Pod Security admission apply
  cleanly and Meshery's components are isolated from unrelated workloads.
- **Resource quotas and limits.** Apply `ResourceQuota`/`LimitRange` to the
  namespace so a discovery surge cannot starve neighbors (and pair with the
  sizing guidance in
  [Infrastructure, Sizing & Performance]({{< ref "installation/production/infrastructure-sizing-and-performance.md" >}})).
- **Tenant boundaries.** Use Meshery's organizations, workspaces, and
  environments (provided through a Remote Provider) to separate teams logically,
  and back that with per-cluster credential scoping for hard isolation.

## Supply-chain integrity

- **Pin image provenance.** Pull Meshery images from trusted registries. The
  chart defaults to the `stable-latest` tag with `pullPolicy: Always`; for
  production, pin to a **specific, immutable version tag (or an image digest,
  `sha256:…`)** so deployments are reproducible and auditable, and so an upstream
  tag move cannot silently change what you run. A digest is the strongest
  guarantee, since tags can in principle be reassigned.
- **Mirror to a private registry** in regulated or air-gapped environments and
  scan images as part of your pipeline. Meshery's published
  [security vulnerabilities]({{< ref "project/security-vulnerabilities.md" >}}) and release notes
  help you track fixes.
- **In-cluster components ship pinned, and stay pinnable.** Meshery Server
  installs the `meshery-operator` chart only at a version the chart repository
  publishes - validated against the published index before Helm is called,
  never a moving tag, and never a release candidate chosen automatically - and
  the Operator in turn deploys MeshSync and the Broker at pinned releases with
  pull policies to match, rather than off a `stable-latest` channel. All three
  versions remain yours to choose. For the Operator chart, set
  `operator.version` on the connection; an explicit pin is honored exactly and
  refused loudly if it names an unpublished or moving version (see
  [How Meshery Server manages Meshery Operator]({{< ref "installation/upgrades/index.md#how-meshery-server-manages-meshery-operator" >}})).
  For MeshSync and the Broker, set `spec.version` on the cluster's `MeshSync`
  and `Broker` custom resources - the Operator honors each with a rollout (see
  [Configuring Meshery Operator, MeshSync, and Broker]({{< ref "guides/infrastructure-management/configuring-operator-meshsync-broker.md" >}})).
- **Verify the chart source.** Install from the official
  [Meshery Helm chart]({{< ref "installation/kubernetes/helm.md" >}}) and review values you
  override.

## Trusting an extension

Meshery is extensible by design. [Adapters]({{< ref "extensions/adapters/_index.md" >}}),
[Providers]({{< ref "reference/extensibility/providers/index.md" >}}),
[models and integrations]({{< ref "extensions/models/_index.md" >}}), and
[UI extension points]({{< ref "reference/extensibility/ui.md" >}}) are how a deployment
grows past the core platform. See the
[Extensibility reference]({{< ref "reference/extensibility/_index.md" >}}) for the full set
of extension points and [meshery.io/extensions](https://meshery.io/extensions) for what is
available.

Extensibility is a capability, not a sandbox.

{{% alert title="Enabled extensions are trusted code" color="warning" %}}
**Every extension you enable runs inside your deployment's trust boundary.** Meshery
provides no sandbox, no privilege separation, and no capability restriction between an
extension and itself. Enabling an extension is a trust decision of the same magnitude as
granting Meshery its cluster credentials, and it deserves the same scrutiny.
{{% /alert %}}

### What each extension point can reach

| Extension point | Where it runs | What it can reach |
| :--- | :--- | :--- |
| **Provider extension package, server plugin** | In-process inside Meshery Server, as a native Go plugin | The Meshery datastore, the Broker connection, the MeshSync channel, the Kubernetes connection tracker, and the ability to serve its own HTTP routes under `/api/provider/extension/server/` |
| **Provider extension package, UI components** | In the browser, inside Meshery UI's own origin and JavaScript context | The DOM, the session cookie, and any API the signed-in user is able to call |
| **[Adapters]({{< ref "extensions/adapters/_index.md" >}})** | A separate process, reached from Meshery Server over gRPC | The managed cluster, using the adapter's own credentials. By default adapters share the `meshery-server` ServiceAccount, so an adapter inherits the Server's Kubernetes reach |
| **[Remote Provider]({{< ref "reference/extensibility/providers/index.md" >}})** | An external service you select | Identity, authorization, and durable persistence for the entire deployment. It also names the extension package that Meshery downloads and loads |
| **[Models and integrations]({{< ref "extensions/models/_index.md" >}})** | Registry data inside Meshery Server | The component and relationship definitions Meshery designs and deploys from |
| **[Build-time content]({{< ref "reference/extensibility/build-time.md" >}})** | Baked into your container image | Whatever the image can reach, and it is present before Meshery starts |

The first two rows are the ones most often underestimated. A provider extension package is
not configuration; it is code, and it runs with the privileges of the process that loads
it.

### How a provider extension package is delivered

Understanding the delivery path is what makes the trust decision concrete:

1. **The provider's capabilities document names the package.** After you sign in, the
   Remote Provider returns a capabilities document that carries the package URL and
   version, along with the extension points to wire up.
2. **Meshery downloads and extracts it.** The package is a `.tar.gz`/`.tgz` archive,
   extracted under `~/.meshery/provider/<provider-name>/<package-version>`. Downloads occur
   at login and on capability refresh.
3. **UI components are served to Meshery UI** from that package and load into the same
   browser origin and JavaScript context as Meshery UI itself.
4. **A server plugin, if the package contains one, is loaded into Meshery Server** and its
   entry point is handed the database handler, the Broker connection, the MeshSync channel,
   and the Kubernetes connection tracker.

Two consequences follow directly, and both are operationally relevant:

- **Selecting a Remote Provider is the trust decision.** The package origin comes from the
  provider you chose, so vetting the provider is what vets the code. Prefer first-party and
  maintained providers, and treat adding an unfamiliar provider as you would adding an
  unfamiliar binary to a production host. This is a further reason to preselect a provider
  with `PROVIDER` rather than leaving the choice to whoever signs in first. See
  [Authentication, Authorization & Identity]({{< ref "installation/production/authentication-and-identity.md" >}})
  and
  [Recommended Production Deployment Settings]({{< ref "reference/extensibility/providers/index.md#recommended-production-deployment-settings" >}}).
- **A loaded server plugin stays loaded for the life of the process.** Go plugins cannot be
  unloaded, so removing or rolling back a server-side extension means restarting Meshery
  Server, not just changing a setting.

### Before you enable an extension

Answer these before an extension reaches production, and record the answers where your
change management can find them later:

- [ ] **Who publishes it, and are they someone you would grant cluster credentials?** That
      is the equivalent level of trust.
- [ ] **Is it maintained?** Check recent releases, open security issues, and whether the
      publisher has a disclosure process.
- [ ] **Does it need a server-side plugin, or only UI components?** UI-only extensions carry
      a smaller blast radius than in-process server plugins.
- [ ] **What version are you pinning to,** and can you reproduce that exact package later?
- [ ] **What new egress does it introduce,** and is that egress allowed by your network
      policy? See
      [Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}}).
- [ ] **Has it been reviewed against the Meshery version you run?** Extension and platform
      versions move independently; see
      [Ensuring Extension Compatibility]({{< ref "reference/extensibility/verify-compatibility.md" >}}).
- [ ] **Do you actually need it enabled?** The cheapest hardening is not enabling what you
      will not use.

### Constraining what an extension can do

Meshery cannot restrict an extension's privileges once it is loaded, so the controls
available to you are about **what gets loaded, from where, and what the surrounding
environment permits**:

- **Pin the extension package.** Set `SKIP_DOWNLOAD_EXTENSIONS=true` to stop Meshery
  downloading and refreshing provider extension packages. Existing packages still load, so
  this pins you to a package you have already vetted rather than accepting whatever the
  provider publishes next. See
  [SKIP_DOWNLOAD_EXTENSIONS]({{< ref "reference/extensibility/providers/index.md#skip_download_extensions" >}}).
- **Pin the capabilities document.** `PROVIDER_CAPABILITIES_FILEPATH` loads capabilities
  from a local file instead of the provider's endpoint, which fixes the extension set for a
  deployment. Use it deliberately: it pins capabilities rather than tracking the provider's
  current set. See
  [PROVIDER_CAPABILITIES_FILEPATH]({{< ref "reference/extensibility/providers/index.md#provider_capabilities_filepath" >}}).
- **Pre-package for air-gapped and regulated environments.** Ship a reviewed extension
  package in your own container image with
  [build-time extensibility]({{< ref "reference/extensibility/build-time.md" >}}), so what runs
  is what your pipeline approved.
- **Enable only the adapters you use.** Adapters are opt-in through `ADAPTER_URLS`, and the
  Helm chart ships every adapter subchart disabled by default. Leave it that way for
  adapters you do not need.
- **Give adapters their own ServiceAccount.** By default an adapter shares the
  `meshery-server` ServiceAccount and therefore the Server's Kubernetes reach. Use
  `serviceAccountNameOverride` with a scoped role so enabling an adapter does not extend
  cluster-admin-equivalent rights to it. See
  [Least-privilege RBAC](#least-privilege-rbac) above.
- **Isolate the adapter channel.** Meshery Server reaches adapters over gRPC without mutual
  TLS today, so treat that channel as trusted-network-only: co-locate adapters in the
  Meshery namespace and use NetworkPolicy to allow adapter ports only from Meshery Server.
  See
  [Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}}).
- **Control egress.** An extension inherits the network reach of the process it runs in.
  Apply egress policy so a server-side extension can only reach the destinations your
  deployment requires.
- **Keep the container hardening on.** The
  [pod and container security context](#pod-and-container-security-context) above bounds
  what an in-process extension can do to the host, even though it does not bound what the
  extension can do inside Meshery.
- **Review models and integrations you import.** Registry content shapes what Meshery
  designs and deploys. Prefer the maintained
  [integrations catalog]({{< ref "extensions/models/_index.md" >}}) over ad-hoc registration.

### Removing an extension, and responding to a bad one

- **To remove a UI extension** delivered by a Remote Provider, the provider stops offering
  it in the capabilities document; it is no longer loaded on the next session.
- **To remove a server-side extension**, restart Meshery Server after the package is gone
  or the capability withdrawn. A loaded Go plugin does not unload in place.
- **To pin away from a bad version**, set `SKIP_DOWNLOAD_EXTENSIONS=true` and restore the
  package version you trust under `~/.meshery/provider/<provider-name>/<package-version>`,
  or pin capabilities with `PROVIDER_CAPABILITIES_FILEPATH`.
- **Report a suspected vulnerability in a Meshery extension** the same way you would report
  one in the core platform, through the process in
  [Security Vulnerabilities]({{< ref "project/security-vulnerabilities.md" >}}). Extensions
  maintained in the `meshery-extensions` GitHub organization are handled by Meshery
  maintainers; extensions distributed elsewhere are the responsibility of their publisher,
  and Meshery has no mechanism to disable an extension in a deployment it does not operate.

{{% alert title="The core project does not vet third-party extensions" color="warning" %}}
Meshery's maintainers are responsible for the core platform and the extensions published
under the `meshery-extensions` organization. The security of any other extension is the
responsibility of its author and of the operator who enables it. Review a third-party
extension before enabling it, and enable only what your deployment actually needs.
{{% /alert %}}

## Hardening checklist

- [ ] Remote Provider preselected; Local Provider disabled in production.
- [ ] ServiceAccount RBAC scoped to the minimum required; `rbac.nodes` enabled
      only where needed.
- [ ] Per-cluster, least-privilege credentials in multi-cluster/multi-cloud.
- [ ] `podSecurityContext`/`securityContext` set: non-root, no privilege
      escalation, dropped capabilities, read-only root FS (with a writable data
      volume), `RuntimeDefault` seccomp.
- [ ] Operator hardened with equivalent security context.
- [ ] kubeconfig and provider secrets sourced from Kubernetes Secrets / external
      manager, never committed.
- [ ] TLS terminated at ingress; provider egress over verified HTTPS.
- [ ] Broker exposure restricted to the Server's origin; private connectivity
      preferred.
- [ ] Dedicated namespace with network policies, quotas, and restricted Pod
      Security.
- [ ] Images pinned to immutable version tags from a trusted/mirrored registry
      and scanned.
- [ ] Every enabled extension reviewed and attributed to a trusted publisher;
      unused extensions left disabled.
- [ ] Extension packages pinned (`SKIP_DOWNLOAD_EXTENSIONS`, and
      `PROVIDER_CAPABILITIES_FILEPATH` where the extension set should be fixed).
- [ ] Adapters enabled only where needed, given scoped ServiceAccounts
      (`serviceAccountNameOverride`), and reachable only from Meshery Server.
- [ ] Egress policy accounts for any destination an enabled extension requires.

{{< related-discussions tag="meshery" >}}
