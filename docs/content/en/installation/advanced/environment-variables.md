---
title: Server Environment Variables
description: Installation and runtime configuration reference for the environment variables consumed by Meshery Server.
aliases:
- /reference/environment-variables
- /reference/meshery-server-environment-variables
---

Meshery Server loads configuration from `server/cmd/server-config.env` and from process environment variables via Viper (`viper.AutomaticEnv()`).

This page documents:

- environment variables operators intentionally set for local, Docker, and Kubernetes deployments.
- runtime variables that Meshery reads from the platform when they are present.

It intentionally excludes test-only variables such as `RUN_INTEGRATION_TESTS` and `PATH_TO_SQL_FILE`, and internal-only keys that Meshery sets for itself during startup.

## Core server configuration

| Variable | Default | Description | Example |
| --- | --- | --- | --- |
| `PORT` | `8080` | HTTP listen port for Meshery Server. Local development targets commonly override this to `9081`. | `PORT=9081 make server` |
| `LOG_LEVEL` | `4` | Numeric log level loaded from `server-config.env`. | `LOG_LEVEL=5 make server` |
| `DEBUG` | `false` | Enables debug behavior and forces debug-level logging even when `LOG_LEVEL` is lower. | `DEBUG=true make server-local` |
| `ADAPTER_URLS` | unset | Space- or comma-separated list of Meshery Adapter endpoints to pre-register at startup. | `ADAPTER_URLS='localhost:10000 localhost:10001' make server-local` |
| `USER_DATA_FOLDER` | `$HOME/.meshery/config` | Directory used for Meshery Server state, including the SQLite database. | `USER_DATA_FOLDER=$HOME/.meshery-dev/config make server` |
| `KUBECONFIG_FOLDER` | `$HOME/.kube` | Directory from which Meshery reads Kubernetes configuration files. | `KUBECONFIG_FOLDER=/home/appuser/.kube docker compose up meshery` |
| `KEYS_PATH` | unset | Path to a CSV file whose rows are seeded into the keys store at startup and after database reset. | `KEYS_PATH=/home/appuser/.meshery/keys.csv make server` |
| `OTEL_CONFIG` | unset | Inline YAML configuration for OpenTelemetry tracing. When unset, tracing is disabled. | `OTEL_CONFIG=$'service_name: meshery-server\nservice_version: 1.0.0\nendpoint: localhost:4317\ninsecure: true' make server-local` |

`LOG_LEVEL` values are `0=panic`, `1=fatal`, `2=error`, `3=warn`, `4=info`, `5=debug`, and `6=trace`.

## Provider and authentication configuration

| Variable | Default | Description | Example |
| --- | --- | --- | --- |
| `PROVIDER` | unset | Enforces a single provider and bypasses the provider selection UI. Accepts `Local` (legacy alias `None`) or a registered remote provider name such as `Meshery`. | `PROVIDER=Meshery` |
| `PROVIDER_BASE_URLS` | Canonical remote-provider list from `install/providers.env` | Comma-separated list of remote provider base URLs that Meshery registers at startup. | `PROVIDER_BASE_URLS=https://cloud.meshery.io,https://cloud.layer5.io` |
| `MESHERY_SERVER_CALLBACK_URL` | `http://<request-host>/api/user/token` | Overrides the OAuth callback URL used when Meshery is behind an ingress, reverse proxy, or load balancer. | `MESHERY_SERVER_CALLBACK_URL=https://playground.meshery.io/api/user/token` |
| `PROVIDER_CAPABILITIES_FILEPATH` | unset | Loads provider capabilities from a local JSON file instead of from the remote provider's `/capabilities` endpoint. Useful for offline development and deterministic testing. | `PROVIDER_CAPABILITIES_FILEPATH=/path/to/capabilities.json` |
| `SKIP_DOWNLOAD_EXTENSIONS` | `false` | Skips downloading or refreshing remote provider extension packages. Existing local packages can still be used. | `SKIP_DOWNLOAD_EXTENSIONS=true make server` |
| `PLAYGROUND` | `false` | Enables playground-specific behavior, including provider preselection flows and playground-oriented post-login redirects. | `PLAYGROUND=true PROVIDER=Meshery make server-local` |

The current built-in `PROVIDER_BASE_URLS` default is:

```text
https://cloud.meshery.io,https://perf.smp-spec.io,https://cloud.layer5.io,https://platform.tata-consulting.co.uk,https://collab.eti.cisco.com,https://kickstart.metabit.com,https://provider.od10.in
```

## Feature flags and runtime behavior

| Variable | Default | Description | Example |
| --- | --- | --- | --- |
| `DISABLE_OPERATOR` | `false` | Prevents Meshery Server from automatically deploying Meshery Operator into connected clusters. | `DISABLE_OPERATOR=true make server-without-operator` |
| `MESHSYNC_DEFAULT_DEPLOYMENT_MODE` | `embedded` | Default MeshSync deployment mode for Kubernetes connections that do not already specify one. Accepted values are `embedded` and `operator`. | `MESHSYNC_DEFAULT_DEPLOYMENT_MODE=operator` |
| `SKIP_DOWNLOAD_CONTENT` | `false` | Skips downloading bundled seed content such as patterns and filters before loading local content. | `SKIP_DOWNLOAD_CONTENT=true make server` |
| `SKIP_COMP_GEN` | `false` | Skips background Kubernetes component generation during startup. Explicit API-triggered registration can still run. | `SKIP_COMP_GEN=true make server` |
| `POLICY_EVAL_TIMEOUT` | `3m` | Maximum duration for a single relationship-policy evaluation. Uses Go duration syntax. | `POLICY_EVAL_TIMEOUT=5m make server` |
| `USE_GO_POLICY_ENGINE` | `false` | Uses the Go relationship-policy engine instead of the default Rego/OPA path. The local `make server` target commonly sets this to `true`. | `USE_GO_POLICY_ENGINE=true make server` |

## Build and runtime metadata

These variables are usually injected by the build or deployment pipeline rather than set manually by end users.

| Variable | Default | Description | Example |
| --- | --- | --- | --- |
| `BUILD` | Build-time value; runtime fallback `Not Set` | Meshery Server version string used in server metadata, provider capability requests, and chart selection. | `BUILD=v1.0.0 make server` |
| `COMMITSHA` | Build-time value; runtime fallback `Not Set` | Git commit SHA reported in server metadata. | `COMMITSHA=$(git rev-parse HEAD) make server` |
| `RELEASE_CHANNEL` | Build-time value; runtime fallback `Not Set` | Release channel reported by the server, such as `stable` or `edge`. | `RELEASE_CHANNEL=stable make server` |
| `INSTANCE_ID` | Auto-generated UUID | Unique Meshery instance identifier used in eventing, connection metadata, and system-owned resources. | `INSTANCE_ID=$(uuidgen) make server` |
| `OS` | `meshery` | Value forwarded to remote-provider capability requests so providers can tailor responses to the caller. | `OS=meshery make server` |

## Platform-provided environment variables

Meshery Server reads these values from the hosting environment when they are present. In Kubernetes deployments they are usually injected automatically and typically do not need to be set manually.

| Variable | Default | Description | Example |
| --- | --- | --- | --- |
| `KUBERNETES_SERVICE_HOST` | unset | Kubernetes service host used to detect in-cluster execution and to build in-cluster Kubernetes contexts. | `KUBERNETES_SERVICE_HOST=10.96.0.1` |
| `KUBERNETES_SERVICE_PORT` | unset | Kubernetes service port paired with `KUBERNETES_SERVICE_HOST` for in-cluster execution. | `KUBERNETES_SERVICE_PORT=443` |
| `HOME` | Inherited from the process environment | Home directory used when Meshery resolves paths under `.meshery/` for local registry/model content. | `HOME=/home/appuser` |

## Common deployment examples

### Local development

```bash
PORT=9081 \
DEBUG=true \
USE_GO_POLICY_ENGINE=true \
PROVIDER_BASE_URLS=http://localhost:9876 \
ADAPTER_URLS='localhost:10000 localhost:10001' \
make server-local
```

### Docker deployment with a preselected remote provider

```bash
docker run -p 9081:8080 \
  -e PROVIDER=Meshery \
  -e PROVIDER_BASE_URLS=https://cloud.meshery.io \
  -e KUBECONFIG_FOLDER=/home/appuser/.kube \
  meshery/meshery:stable-latest
```

### Kubernetes or ingress-backed deployment

```bash
MESHERY_SERVER_CALLBACK_URL=https://meshery.example.com/api/user/token \
PROVIDER=Meshery \
MESHSYNC_DEFAULT_DEPLOYMENT_MODE=operator
```
