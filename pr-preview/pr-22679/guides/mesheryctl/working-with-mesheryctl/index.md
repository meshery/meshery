# Using Meshery CLI

> Guides for common tasks while using Meshery's CLI, mesheryctl.

---

LLMS index: [llms.txt](/pr-preview/pr-22679/llms.txt)

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

### What is the meshconfig?

Meshconfig is to `mesheryctl` what kubeconfig is to `kubectl`. `mesheryctl` reads it from `~/.meshery/config.yaml` by default. Pass the global `--config <path>` flag to any `mesheryctl` command to use a different file.

### What is a context?

A context describes a single Meshery deployment: where Meshery Server runs, which platform it runs on, which release channel and version it uses, which provider it enforces, and which Meshery Adapters it deploys. Contexts let you keep several Meshery deployments, for example a local Docker one and a shared Kubernetes one, and switch between them.

### What is `current-context`?

`current-context` names the context that `mesheryctl` commands use. Only one context is current at a time. Change it with `mesheryctl system context switch <name>`, or override it for a single `mesheryctl system` command with `-c <name>`.

### What's the difference between contexts and environments?

A context configures a Meshery deployment (Meshery Server, its adapters, and its operator settings) and lives in your local meshconfig. An [environment](/pr-preview/pr-22679/concepts/logical/environments/) is a group of connections, such as Kubernetes clusters, that Meshery manages, and it lives in Meshery Server.

### What does a meshconfig look like?

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

`components` accepts `meshery-istio`, `meshery-linkerd`, `meshery-consul`, `meshery-nsm`, `meshery-traefik-mesh`, `meshery-kuma`, and `meshery-cilium`. `mesheryctl` uppercases each key under `env` and passes it to Meshery Server as an environment variable, so `meshsync_default_deployment_mode` becomes `MESHSYNC_DEFAULT_DEPLOYMENT_MODE`. See [Environment Variables](/pr-preview/pr-22679/installation/advanced/environment-variables/) for the variables Meshery Server reads.

To see a generated context, run:

```bash
mesheryctl system context create test
mesheryctl system context view test
```

### What does the `--config` flag do?

`--config` is a global flag that every `mesheryctl` command accepts. It points `mesheryctl` at a meshconfig other than `~/.meshery/config.yaml`, which is useful when you keep a separate meshconfig per project or per CI job.

### How do endpoints work in meshconfig?

The `endpoint` is the URL where `mesheryctl` reaches Meshery Server and where it opens Meshery UI.

- **Docker:** `mesheryctl system start` publishes Meshery on the port in the context's endpoint, for example `http://localhost:9081`.
- **Kubernetes:** `mesheryctl system dashboard` looks up the Meshery service through the Kubernetes API and writes the address it finds back to the context's endpoint.

### Can I get an API token using mesheryctl?

Yes. Run `mesheryctl system login`. See [Authenticating Meshery via CLI](/pr-preview/pr-22679/guides/mesheryctl/authenticate-with-meshery-via-cli/).

## Run a specific Meshery version

Every time you run `mesheryctl system start`, `mesheryctl` regenerates the services in `~/.meshery/meshery.yaml` from the current context. Edits you make to the image tags in that file do not last. To pin Meshery to a release, set it on the context instead:

```bash
mesheryctl system channel set stable-v0.9.2
mesheryctl system restart
```

To go back to the newest stable release, run `mesheryctl system channel set stable`.

## Suggested reading

For every `mesheryctl` command and flag, see the [`mesheryctl` command reference](/pr-preview/pr-22679/reference/mesheryctl/).

<div class="related-discussions">
  <h3>Recent Discussions with "mesheryctl" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/pairing-up-with-a-meshmate-to-gacefully-start-contributing-to-meshery-and-its-project-ecosystem/6957" target="_blank" rel="noopener noreferrer">
            Pairing up with a Meshmate to gacefully start contributing to Meshery and its project ecosystem
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/development-meeting-contributing-to-meshery-cli-april-30-2025/6927" target="_blank" rel="noopener noreferrer">
            [Development Meeting] Contributing to Meshery CLI - April 30, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/meshery-development-meeting-april-30th-2025/6926" target="_blank" rel="noopener noreferrer">
            Meshery Development Meeting | April 30th, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newcomers-meeting-end-to-end-testing-in-meshery-cli-using-bats-april-17-2025/6897" target="_blank" rel="noopener noreferrer">
            [Newcomers’ Meeting] End-to-End Testing in Meshery CLI using BATs – April 17, 2025
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/unable-to-lint-mesheryctl/6854" target="_blank" rel="noopener noreferrer">
            unable to lint mesheryctl 
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/mesheryctl-system-login-problem/6687" target="_blank" rel="noopener noreferrer">
            Mesheryctl system login problem
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/error-while-launching-the-meshery-dashboard/6600" target="_blank" rel="noopener noreferrer">
            Error while launching the Meshery Dashboard
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/cant-find-the-file-path-for-meshery-designs/6319" target="_blank" rel="noopener noreferrer">
            Can&#39;t find the file path for meshery designs
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/mesheryctl" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>mesheryctl</code>
    </a>
  </p>
</div>


<div class="alert alert-dark" role="alert">
  <h4 class="alert-heading">Discussion Forum</h4>
  <p>Don't find an answer to your question here? Ask on the <a href="https://discuss.meshery.io/">Discussion Forum</a>.</p>
</div>
