---
title: Authentication, Authorization & Identity
linkTitle: Authentication & Identity
description: >-
  Why to preselect a Remote Provider over the Local Provider, OAuth callback
  configuration, identity providers, capabilities, and keys/permissions for
  production Meshery.
categories: [installation]
weight: 60
aliases:
- /installation/production/authentication
- /installation/production/identity
---

Identity is the first security control to get right in production. Meshery's
authentication, authorization, and durable user state all flow from its
**provider** model. This page explains why production deployments should
preselect a Remote Provider, how to configure the OAuth flow correctly behind an
ingress, and how authorization is expressed. It pairs with
[Security Hardening]({{< ref "installation/production/security-hardening.md" >}}).

## Providers: Local vs. Remote

Meshery supports two kinds of [providers]({{< ref "reference/extensibility/providers/index.md" >}}):

| Provider | Identity | Durable state | Production fit |
| :--- | :--- | :--- | :--- |
| **Local** (built-in; legacy alias `None`) | **None**—sessions are unauthenticated | Only the ephemeral local database | Evaluation and single-operator local use. **Not for shared production.** |
| **Remote** (pluggable, e.g. Meshery Cloud) | Delegated to the provider's identity providers (OAuth) | **Durable, long-term** persistence of users, environments, and saved work | The production choice. |

The Local Provider exists so Meshery can run with zero external dependencies,
but it offers **no authentication** and no durable home for user state. A
Remote Provider adds authentication and authorization, controls which identity
providers are accepted, and persists user data beyond the lifetime of any Server
instance.

{{% alert title="Preselect a Remote Provider in production" color="warning" %}}
Pin a Remote Provider so Meshery never falls back to unauthenticated Local
Provider sessions and so you control accepted identity providers. This is both a
security control (no anonymous access) and a resiliency control (durable state
lives with the provider, not the ephemeral cache).
{{% /alert %}}

## Preselecting the provider

Two environment variables control provider behavior. Set them via the Helm chart
(`env.*`) or your deployment's environment. See the
[environment variables reference]({{< ref "installation/advanced/environment-variables.md" >}}).

- **`PROVIDER`** — enforces a single provider and **bypasses the provider
  selection screen**. Set it to your Remote Provider's registered name (for
  example `Meshery`). Setting `PROVIDER=Local` (or the legacy alias `None`)
  pins the Local Provider—avoid that in production.
- **`PROVIDER_BASE_URLS`** — the comma-separated list of Remote Provider base
  URLs Meshery registers at startup. Restrict this to the provider(s) you
  intend to allow rather than the full default list.

```bash
helm install meshery meshery/meshery --namespace meshery --create-namespace \
  --set env.PROVIDER=Meshery \
  --set env.PROVIDER_BASE_URLS=https://cloud.meshery.io
```

Pinning both `PROVIDER` and a single `PROVIDER_BASE_URLS` entry means users are
taken straight into the chosen Remote Provider's authentication flow, with no
opportunity to select the Local Provider.

`mesheryctl` users can also view and set the provider for CLI-driven workflows
via `mesheryctl system provider` (see the
[reference]({{< ref "reference/references/mesheryctl/system/provider/set.md" >}})), but the server-side
`PROVIDER` setting is what governs the deployment as a whole.

## Configuring the OAuth callback URL

Remote Provider authentication uses an OAuth flow that redirects the user's
browser back to Meshery. Behind an ingress, reverse proxy, or load balancer,
Meshery's internal address differs from the external URL, so you must tell
Meshery its external callback address:

- **`MESHERY_SERVER_CALLBACK_URL`** — overrides the OAuth callback URL. Set it to
  Meshery's external base URL; the effective callback path is
  `/api/user/token`.

```bash
helm install meshery meshery/meshery --namespace meshery --create-namespace \
  --set env.PROVIDER=Meshery \
  --set env.MESHERY_SERVER_CALLBACK_URL=https://meshery.example.com
```

Symptoms of a missing or wrong callback URL include login redirect loops, "redirect
URI mismatch" errors from the identity provider, or landing back on the login
screen after authenticating. When configuring the Remote Provider's OAuth
application, register the same external callback
(`https://meshery.example.com/api/user/token`) as an allowed redirect URI.

{{% alert title="Callback URL and ingress go together" color="info" %}}
If you front Meshery with an ingress (you should), set
`MESHERY_SERVER_CALLBACK_URL` to the public hostname at the same time. See
[Networking & Connectivity]({{< ref "installation/production/networking-and-connectivity.md" >}})
for the ingress, TLS, and WebSocket configuration the OAuth flow depends on.
{{% /alert %}}

## Identity providers

A Remote Provider is what integrates Meshery with identity providers (the OAuth/
OIDC sources that actually authenticate users). By controlling which Remote
Provider(s) you allow—and thus which identity providers are recognized—you
control your authentication surface. Choose a Remote Provider whose supported
identity providers match your organization's SSO, and restrict
`PROVIDER_BASE_URLS` accordingly.

## Capabilities and authorization

After authentication, a Remote Provider returns a **capabilities** document that
determines which features and actions are available to the user—this is how
authorization is expressed through the provider model. Production notes:

- Meshery loads provider capabilities at startup; readiness/health reflect that
  capabilities are loaded (see
  [High Availability & Resiliency]({{< ref "installation/production/high-availability-and-resiliency.md" >}})).
- For tightly egress-restricted or offline environments, capabilities can be
  loaded from a local file via `PROVIDER_CAPABILITIES_FILEPATH` instead of the
  provider's endpoint. Use this deliberately; it pins capabilities rather than
  fetching the provider's current set.
- `SKIP_DOWNLOAD_EXTENSIONS` controls whether provider extension packages are
  downloaded/refreshed; existing local packages can still be used.

## Keys and permissions

Meshery seeds a set of keys/permissions used by its role-based access controls.
The `KEYS_PATH` configuration points Meshery at a CSV of keys that are loaded at
startup (and after a database reset). In the Helm chart this defaults to the
bundled permissions file. If you customize roles/permissions, manage the keys
file as version-controlled configuration and mount it consistently so RBAC is
reproducible across redeploys.

Fine-grained, user-facing authorization (organizations, teams, workspaces, and
roles) is delivered through the Remote Provider, which is another reason a
Remote Provider is required for multi-user production.

## Session and token handling

- Clients (UI, `mesheryctl`, and any API consumer) present a valid **JWT** to
  Meshery's APIs; unauthenticated API access is not available under a Remote
  Provider.
- Serve all authenticated traffic over TLS so tokens are never sent in the
  clear.
- Treat the Remote Provider as part of your availability model—login depends on
  it. See
  [High Availability & Resiliency]({{< ref "installation/production/high-availability-and-resiliency.md" >}}).

## Identity checklist

- [ ] `PROVIDER` set to a Remote Provider name (never `Local`/`None` in
      production).
- [ ] `PROVIDER_BASE_URLS` restricted to the allowed Remote Provider(s).
- [ ] `MESHERY_SERVER_CALLBACK_URL` set to the external URL; matching redirect
      URI registered with the provider's OAuth application.
- [ ] All authenticated traffic served over TLS (`https://`/`wss://`).
- [ ] Custom keys/permissions (if any) version-controlled and mounted
      consistently.
- [ ] Remote Provider included in availability monitoring and runbooks.

{{< related-discussions tag="meshery" >}}
