---
title: Security Hardening
linkTitle: Security Hardening
description: >-
  RBAC and least privilege, pod and container security contexts, secret and
  kubeconfig handling, TLS, supply-chain integrity, Broker exposure risk, and
  namespace isolation for hardening Meshery in production.
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
- **Verify the chart source.** Install from the official
  [Meshery Helm chart]({{< ref "installation/kubernetes/helm.md" >}}) and review values you
  override.

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

{{< related-discussions tag="meshery" >}}
