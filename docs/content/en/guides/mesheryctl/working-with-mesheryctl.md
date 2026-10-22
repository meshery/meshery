---
title: Using Meshery CLI
categories: [mesheryctl]
description: Guides for common tasks while using Meshery's CLI, mesheryctl.
---

`mesheryctl` is Meshery's command line interface. Use it to manage the lifecycle of Meshery itself and to call Meshery's infrastructure management features from your terminal. Its commands fall into these groups:

- `mesheryctl` - global flags (`--config`, `--verbose`), `version`, and `completion`
- `mesheryctl system` - Meshery lifecycle, contexts, providers, tokens, and troubleshooting
- `mesheryctl design` - design import, export, deployment, and evaluation
- `mesheryctl model` - model management and registry operations
- `mesheryctl component` - component lookup
- `mesheryctl relationship` - relationship management
- `mesheryctl registry` - capability registry generation and publishing
- `mesheryctl connection` - connection management
- `mesheryctl environment` - environment management
- `mesheryctl workspace` - workspace management
- `mesheryctl organization` - organization management
- `mesheryctl filter` - WebAssembly filter management
- `mesheryctl perf` - performance management
- `mesheryctl adapter` - Meshery Adapter operations
- `mesheryctl exp` - experimental commands

## Configure Meshery deployments with meshconfig

_Meshconfig_ is the `mesheryctl` configuration file at `~/.meshery/config.yaml`. It holds one or more _contexts_, and each context describes one Meshery deployment. `mesheryctl system` commands such as `start`, `stop`, `status`, and `reset` act on the current context.

## Meshery CLI FAQ

#### What is the meshconfig?

Meshconfig is to `mesheryctl` what kubeconfig is to `kubectl`. `mesheryctl` reads it from `~/.meshery/config.yaml` by default. Pass the global `--config <path>` flag to any `mesheryctl` command to use a different file.

#### What is a context?

A context describes a single Meshery deployment: where Meshery Server runs, which platform it runs on, which release channel and version it uses, which provider it enforces, and which Meshery Adapters it deploys. Contexts let you keep several Meshery deployments, for example a local Docker one and a shared Kubernetes one, and switch between them.

#### What is `current-context`?

`current-context` names the context that `mesheryctl` commands use. Only one context is current at a time. Change it with `mesheryctl system context switch <name>`, or override it for a single `mesheryctl system` command with `-c <name>`.

#### What's the difference between contexts and environments?

A context configures a Meshery deployment (Meshery Server, its adapters, and its operator settings) and lives in your local meshconfig. An [environment]({{< ref "concepts/logical/environments.md" >}}) is a group of connections, such as Kubernetes clusters, that Meshery manages, and it lives in Meshery Server.

#### What does a meshconfig look like?

`mesheryctl` creates `~/.meshery/config.yaml` the first time you run it. The file has this shape. Most fields are optional.

```yaml
contexts:
  local:
    endpoint: http://localhost:9081   # URL of Meshery Server
    token: default                    # name of an entry under tokens
    platform: kubernetes              # "docker" or "kubernetes"
    channel: stable                   # "stable" or "edge"
    version: latest                   # "latest" or a release tag such as v0.9.2
    provider: Meshery                 # optional: enforce a provider, for example "Meshery" or "Local"
    components:                       # optional: Meshery Adapters to deploy
      - meshery-istio
      - meshery-linkerd
    env:                              # optional: extra environment variables for Meshery Server
      meshsync_default_deployment_mode: embedded

current-context: local

tokens:
  - name: default
    location: /home/me/.meshery/auth.json
```

`components` accepts `meshery-istio`, `meshery-linkerd`, `meshery-consul`, `meshery-nsm`, `meshery-traefik-mesh`, `meshery-kuma`, and `meshery-cilium`. `mesheryctl` uppercases each key under `env` and passes it to Meshery Server as an environment variable, so `meshsync_default_deployment_mode` becomes `MESHSYNC_DEFAULT_DEPLOYMENT_MODE`. See [Environment Variables]({{< ref "installation/advanced/environment-variables.md" >}}) for the variables Meshery Server reads.

To see a generated context, run:

```bash
mesheryctl system context create test
mesheryctl system context view test
```

#### What does the `--config` flag do?

`--config` is a global flag that every `mesheryctl` command accepts. It points `mesheryctl` at a meshconfig other than `~/.meshery/config.yaml`, which is useful when you keep a separate meshconfig per project or per CI job.

#### How do endpoints work in meshconfig?

The `endpoint` is the URL where `mesheryctl` reaches Meshery Server and where it opens Meshery UI.

- **Docker:** `mesheryctl system start` publishes Meshery on the port in the context's endpoint, for example `http://localhost:9081`.
- **Kubernetes:** `mesheryctl system dashboard` looks up the Meshery service through the Kubernetes API and writes the address it finds back to the context's endpoint.

#### Can I get an API token using mesheryctl?

Yes. Run `mesheryctl system login`. See [Authenticating Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}).

## Run a specific Meshery version

Every time you run `mesheryctl system start`, `mesheryctl` regenerates the services in `~/.meshery/meshery.yaml` from the current context. Edits you make to the image tags in that file do not last. To pin Meshery to a release, set it on the context instead:

```bash
mesheryctl system channel set stable-v0.9.2
mesheryctl system restart
```

To go back to the newest stable release, run `mesheryctl system channel set stable`.

## Suggested reading

For every `mesheryctl` command and flag, see the [`mesheryctl` command reference]({{< ref "reference/references/mesheryctl/_index.md" >}}).

{{< related-discussions tag="mesheryctl" >}}

{{< discuss >}}
