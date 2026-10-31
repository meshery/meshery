# Authentication, Authorization & Identity

> Why to preselect a Remote Provider over the Local Provider, OAuth callback configuration, identity providers, capabilities, and keys/permissions for production Meshery.

---

LLMS index: [llms.txt](/pr-preview/pr-22696/llms.txt)

---

Identity is the first security control to get right in production. Meshery's
authentication, authorization, and durable user state all flow from its
**provider** model. This page explains why production deployments should
preselect a Remote Provider, how to configure the OAuth flow correctly behind an
ingress, and how authorization is expressed. It pairs with
[Security Hardening](/pr-preview/pr-22696/installation/production/security-hardening/).

## Providers: Local vs. Remote

Meshery supports two kinds of [providers](/pr-preview/pr-22696/reference/extensibility/providers/):

| Provider | Identity | Durable state | Production fit |
| :--- | :--- | :--- | :--- |
| **Local** (built-in; legacy alias `None`) | **None** - sessions are unauthenticated | Only the ephemeral local database | Evaluation and single-operator local use. **Not for shared production.** |
| **Remote** (pluggable, e.g. Meshery Cloud) | Delegated to the provider's identity providers (OAuth) | **Durable, long-term** persistence of users, environments, and saved work | The production choice. |

The Local Provider exists so Meshery can run with zero external dependencies,
but it offers **no authentication** and no durable home for user state. A
Remote Provider adds authentication and authorization, controls which identity
providers are accepted, and persists user data beyond the lifetime of any Server
instance.

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Preselect a Remote Provider in production</div>


Pin a Remote Provider so Meshery never falls back to unauthenticated Local
Provider sessions and so you control accepted identity providers. This is both a
security control (no anonymous access) and a resiliency control (durable state
lives with the provider, not the ephemeral cache).
</div>


## Preselecting the provider

Two environment variables control provider behavior. Set them via the Helm chart
(`env.*`) or your deployment's environment. See the
[environment variables reference](/pr-preview/pr-22696/installation/advanced/environment-variables/).

- **`PROVIDER`** - hard-enforces a single provider at **server boot**. Set it to your Remote Provider's registered name (for example `Meshery`). When it is set to a remote, the Local Provider is not registered and cannot be selected via the chooser, cookie, header, or `?provider=`. Setting `PROVIDER=Local` (or the legacy alias `None`) pins the Local Provider; avoid that in production. If `PROVIDER` is set but does not match a registered provider, Meshery **refuses to start** rather than falling back to the chooser.
- **`PROVIDER_BASE_URLS`** - the comma-separated list of Remote Provider base URLs Meshery registers at startup. Restrict this to the provider you intend to pin rather than the full default list. This setting is read at process start; change it and restart (for example with `helm upgrade`) to re-point an existing deployment. Provider-bound data does not migrate. `mesheryctl system provider switch` changes `PROVIDER` only, never this list - see [below](#preselecting-the-provider).

```bash
helm install meshery meshery/meshery --namespace meshery --create-namespace \
  --set env.PROVIDER=Meshery \
  --set env.PROVIDER_BASE_URLS=https://cloud.meshery.io
```

Pinning both `PROVIDER` and a single `PROVIDER_BASE_URLS` entry means users are
taken straight into the chosen Remote Provider's authentication flow. The Local Provider is not registered on that server, and cookies or query parameters cannot switch to it.

`mesheryctl` users can also view and set the provider for CLI-driven workflows
via `mesheryctl system provider` (see the
[reference](/pr-preview/pr-22696/reference/mesheryctl/system/provider/set/)), but the server-side
`PROVIDER` setting is what governs the deployment as a whole.

`mesheryctl system provider switch` sets the provider on the CLI's current
context and restarts the deployment, which is how that value reaches the server
as `PROVIDER`. It does **not** add anything to `PROVIDER_BASE_URLS` - that list
is read separately at server start - so the target provider's base URL must
already be present there. Two consequences follow on a deployment that is
already pinned:

- `mesheryctl system provider set` validates against the running server's
  `/api/providers`, which reports only the pinned provider, so selecting a
  different one requires `--force`.
- If the value you force is not resolvable when the server restarts, the server
  refuses to start rather than falling back to the chooser. Add the provider's
  URL to `PROVIDER_BASE_URLS` first.

## Configuring the OAuth callback URL

Remote Provider authentication uses an OAuth flow that redirects the user's
browser back to Meshery. Behind an ingress, reverse proxy, or load balancer,
Meshery's internal address differs from the external URL, so you must tell
Meshery its external callback address:

- **`MESHERY_SERVER_CALLBACK_URL`** - overrides the OAuth callback URL. Set it to
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

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Callback URL and ingress go together</div>


If you front Meshery with an ingress (you should), set
`MESHERY_SERVER_CALLBACK_URL` to the public hostname at the same time. See
[Networking & Connectivity](/pr-preview/pr-22696/installation/production/networking-and-connectivity/)
for the ingress, TLS, and WebSocket configuration the OAuth flow depends on.
</div>


## Identity providers

A Remote Provider is what integrates Meshery with identity providers (the OAuth/
OIDC sources that actually authenticate users). By controlling which Remote
Provider(s) you allow - and thus which identity providers are recognized - you
control your authentication surface. Choose a Remote Provider whose supported
identity providers match your organization's SSO, and restrict
`PROVIDER_BASE_URLS` accordingly.

## Capabilities and authorization

After authentication, a Remote Provider returns a **capabilities** document that
determines which features and actions are available to the user - this is how
authorization is expressed through the provider model. Production notes:

- Meshery loads provider capabilities at startup; readiness/health reflect that
  capabilities are loaded (see
  [High Availability & Resiliency](/pr-preview/pr-22696/installation/production/high-availability-and-resiliency/)).
- For tightly egress-restricted or offline environments, capabilities can be
  loaded from a local file via `PROVIDER_CAPABILITIES_FILEPATH` instead of the
  provider's endpoint. Use this deliberately; it pins capabilities rather than
  fetching the provider's current set.
- `SKIP_DOWNLOAD_EXTENSIONS` controls whether provider extension packages are
  downloaded/refreshed; existing local packages can still be used.

<div class="alert alert-warning" role="alert"><div class="h4 alert-heading" role="heading">Choosing a provider is also an extension-trust decision</div>


The capabilities document does more than enable features: it names the **extension
package** Meshery downloads and loads, whose UI components run in Meshery UI's own browser
context and whose server-side plugin, if present, runs in-process inside Meshery Server.
Selecting a Remote Provider therefore decides whose code runs inside your deployment. See
[Trusting an extension](/pr-preview/pr-22696/installation/production/security-hardening/#trusting-an-extension).
</div>


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
- Treat the Remote Provider as part of your availability model - login depends on
  it. See
  [High Availability & Resiliency](/pr-preview/pr-22696/installation/production/high-availability-and-resiliency/).

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

<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-seeking-an-experienced-meshery-contributor-as-a-guide/8166" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate: Seeking an Experienced Meshery Contributor as a Guide
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-to-reactivate-my-slack-account/8168" target="_blank" rel="noopener noreferrer">
            Request to Reactivate My Slack Account
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/bug-new-catalog-publish-request-email-renders-with-oversized-social-icons-and-squashed-content-on-mobile/8138" target="_blank" rel="noopener noreferrer">
            [Bug] &#34;New Catalog publish request&#34; email renders with oversized social icons and squashed content on mobile
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/new-contributor-intro-vedit-agrawal-bughunterx2101-looking-to-pair-with-a-meshmate-for-meshery-ui-ai-adapter/8172" target="_blank" rel="noopener noreferrer">
            New Contributor Intro: Vedit Agrawal (@BugHunterX2101) | Looking to pair with a MeshMate for Meshery UI &amp; AI Adapter
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-frontend-backend-cloud-native-development/8161" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate — Frontend, Backend &amp; Cloud-Native Development
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-for-a-meshmate-active-meshery-contributor/8158" target="_blank" rel="noopener noreferrer">
            Request for a MeshMate — Active Meshery Contributor
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/how-to-present-an-issue-pr-in-a-development-meeting-feedback-wanted/8107" target="_blank" rel="noopener noreferrer">
            How to Present an Issue/PR in a Development Meeting: feedback wanted
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newbie-here-excited-to-join/8052" target="_blank" rel="noopener noreferrer">
            Newbie here -- Excited to join
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
