---
layout: default
title: mesheryctl-system-provider-switch
permalink: reference/mesheryctl/system/provider/switch
redirect_from: reference/mesheryctl/system/provider/switch/
type: reference
display-title: "false"
language: en
command: system
subcommand: provider
---

# mesheryctl system provider switch

Switch provider and redeploy

## Synopsis

Switch provider of context in focus and redeploy Meshery. Run `mesheryctl system provider list` to see the available providers.

<pre class='codeblock-pre'>
<div class='codeblock'>
mesheryctl system provider switch [provider] [flags]

</div>
</pre>

## Description

A **provider** determines how Meshery handles user identity and authentication. Switching the provider updates the active Meshery context and triggers a redeployment of Meshery.

### What changes when you switch providers?

Switching providers affects several core components of your Meshery instance:
* **Authentication Flow:** Determines if you log in via a remote identity provider (e.g., Meshery Cloud) or bypass authentication (Local mode).
* **Persistence:** Where your designs, filters, and profiles are stored (locally on your machine or in a remote database).
* **UI Experience:** After switching, you may be redirected to a different login or landing page.

### Local (None) Provider Behavior

The **Local** provider (often specified as `None`) is designed for offline or local-only deployments. 
* **No Authentication:** It does not require an external login.
* **Direct Access:** The UI should skip the login screen and redirect you directly to the dashboard.

> **Note:** If you are using the Local provider and are stuck at a `/user/login` screen, this is likely a configuration mismatch. Ensure your `meshconfig` is updated correctly using the troubleshooting steps below.

### Using a provider preset in meshconfig

The active provider is stored in the Meshery configuration file (meshconfig), typically located at:

- `~/.meshery/config.yaml`

If the provider is preset in the meshconfig, `mesheryctl` will use that provider when starting or redeploying Meshery.

---

## Troubleshooting

### I see a blank white page after switching providers

If you see a blank page in the browser after switching providers (commonly at a route like `/user/login`), verify which provider is currently active:

```bash
mesheryctl system provider view
