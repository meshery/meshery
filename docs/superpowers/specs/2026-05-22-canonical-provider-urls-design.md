# Single Canonical Remote-Provider List for Install Artifacts

**Date:** 2026-05-22
**Status:** Approved (design)

## Problem

The remote-provider URL list (`PROVIDER_BASE_URLS`) is duplicated, by hand, across
every install artifact: docker-compose files, the Helm chart values, the raw k8s
deployment manifest, the playground manifest, the playground Makefile, and (as a
compiled-in fallback) the `mesheryctl` Go sources. Each consumer is read by a
different tool - `docker compose`, `helm`, `kubectl apply`, `make`, the Go binary -
and none of them can read a `make` variable, so the list gets copied and drifts.

Two concrete failures motivated this:

1. Production Meshery Playground served a stale `cloud.layer5.io` because its env was
   hand-maintained on the cluster, independent of the repo.
2. The values diverge across files (e.g. `mesheryctl`'s embedded compose ships
   `meshery.digitalocean.com` and `idp.cleverluck.com`, which no other artifact lists).

## Goal

A single, human-readable, name→URL source of truth for the production remote providers,
from which every install artifact derives its `PROVIDER_BASE_URLS`, enforced by a CI
drift check. Removing the duplicated `MESHERY_CLOUD_PROD` / `LAYER5_CLOUD_PROD` variables
from `Makefile.core.mk` is part of this.

### Non-goals

- Governing the live CNCF cluster's env (that is a separate GitOps fix - the cluster
  must be deployed from the checked-in playground manifest, not hand-edited).
- Changing the server's provider-resolution contract.
- Touching the unrelated `MESHERY_CLOUD_PROD` constant in `ui/constants/endpoints.ts`
  (a UI link constant, not the Make variable).
- Restructuring the non-production Make provider vars (`MESHERY_CLOUD_STAGING`,
  `REMOTE_PROVIDER_LOCAL`, `MESHERY_CLOUD_DEV`, `EXOSCALE_*`, `EQUINIX_*`), which serve
  specific dev/staging targets and are out of scope.

## Key constraint discovered

The server has **no fallback** for the remote-provider list. `DefaultProviderURL`
(`server/cmd/main.go:60`) is stamped only on the built-in **Local** provider's
capability paths (`main.go:281`); it is not used for the remote list. If
`PROVIDER_BASE_URLS` is empty, `main.go:358` registers zero remote providers. Therefore
**every install artifact must ship a populated value** - the guarantee that users see
remote providers depends entirely on these files, not on any server default.

### Verified default per install path

| Path | Source file | Required default |
|---|---|---|
| `helm install` | `install/kubernetes/helm/meshery/values.yaml` (`env.PROVIDER_BASE_URLS`, ranged into the container at `templates/deployment.yaml:68-70`) | full |
| `mesheryctl system start` (Docker) | downloads `install/docker/docker-compose.yaml` (`platform.go:413/436`) | full |
| `mesheryctl system start` (Kubernetes) | applies `install/deployment_yamls/k8s/meshery-deployment.yaml` (`platform.go:129/346`) | full |
| Docker Desktop extension | `install/docker-extension/docker-compose.yaml` | full, Meshery pre-selected |
| Playground (hosted) | `install/playground/**` | Meshery only, pre-selected |

## Design (Approach A: generator + drift check)

### 1. Canonical source - `install/providers.env`

`Name=URL` per line, ordered. Production remote providers only (the set composing the
"full" list). The `Name` is the human-friendly provider name and doubles as the valid
value for the `PROVIDER` (pre-selection) env var.

```
Meshery=https://cloud.meshery.io
Layer5=https://cloud.layer5.io
```

### 2. `install/Makefile.core.mk`

- Delete `MESHERY_CLOUD_PROD` and `LAYER5_CLOUD_PROD` (used only on the line below).
- Derive the full list from the canonical file:

```makefile
REMOTE_PROVIDER_URLS := $(shell grep -v '^\#' install/providers.env | cut -d= -f2 | paste -sd, -)
```

Non-production provider vars are left untouched.

### 3. Generator - `install/scripts/sync-provider-urls.py`

A Python 3 script (python3 is present on dev machines and CI runners; multi-line k8s
edits make `sed` fragile). It is the single implementation of the generation logic and
emits **all** derived artifacts - the static install files (in-place value replacement)
and the generated Go file for `mesheryctl` (full file emit).

- Parses `install/providers.env` into an ordered name→URL map.
- Computes `full = join(all urls, ",")` and `playground = url["Meshery"]`.
- Holds a table mapping each managed file to its profile and line format.
- Replaces **only the value** of the managed line in each static file via a
  format-specific regex, preserving all surrounding formatting and comments.
- Emits a generated Go source file (see section 5) carrying the full list.
- Fails loudly if a managed file is missing its expected managed line (that is itself a
  drift signal), rather than silently skipping.
- `--check` mode: compute desired content in memory, diff against on-disk, print the
  diff, exit non-zero if any file (including the generated Go file) is out of sync.
  Writes nothing.

Line formats handled:

| File | Profile | Format |
|---|---|---|
| `install/docker/docker-compose.yaml` | full | compose list item `- "PROVIDER_BASE_URLS=..."` |
| `install/mesheryapp.dockerapp/docker-compose.yml` | full | compose list item (quoted) |
| `install/docker-extension/docker-compose.yaml` | full | compose list item (quoted); `PROVIDER=Meshery` left as-is |
| `install/deployment_yamls/k8s/meshery-deployment.yaml` | full | k8s `name:`/`value:` pair |
| `install/kubernetes/helm/meshery/values.yaml` | full | YAML map `PROVIDER_BASE_URLS: "..."` |
| `install/playground/docker/docker-compose.yaml` | playground | compose list item (unquoted) |
| `install/playground/kubernetes/playground-deployment.yaml` | playground | k8s `name:`/`value:` pair |
| `install/playground/docker/Makefile` | playground | make var `REMOTE_PROVIDER=...` |

### 4. Root `Makefile`

- `generate-install`: runs the script (writes).
- `check-install`: runs the script with `--check`.
- Wire `check-install` into the existing CI/verify target so drift fails the build.

### 5. `mesheryctl` embedded fallback - generated from the canonical list

`mesheryctl/pkg/utils/helpers.go` `Services["meshery"].Environment` currently hardcodes a
stale, non-canonical list (`meshery.digitalocean.com`, `idp.cleverluck.com`). Because this
value is compiled into the binary at `go build`, "inherit from the canonical list" means
**code generation**, not a runtime read.

- The generator emits `mesheryctl/pkg/utils/providers_gen.go` - a `// Code generated ... DO
  NOT EDIT.` file in `package utils` exposing the canonical full list, e.g.:

  ```go
  // Code generated by install/scripts/sync-provider-urls.py from install/providers.env; DO NOT EDIT.
  package utils

  // DefaultProviderBaseURLs is the canonical comma-joined production remote-provider list.
  const DefaultProviderBaseURLs = "https://cloud.meshery.io,https://cloud.layer5.io"
  ```

- `helpers.go` replaces the inline string with `"PROVIDER_BASE_URLS=" + DefaultProviderBaseURLs`.
  The rest of the `Services` map is untouched.
- A `//go:generate` directive documents regeneration; `make generate-install` is the
  canonical entry point and `make check-install` (CI) covers the generated file, so it can
  never drift from `providers.env`.
- The generated file is committed (Go convention), so `go build` needs no extra step.

## Lifecycle: how and when artifacts are populated

Generation is a developer/commit-time step enforced by CI; there is no runtime read of
`providers.env` except by `make` itself. The flow:

1. A human edits `install/providers.env`.
2. `make generate-install` regenerates every committed artifact (static install files +
   `providers_gen.go`).
3. `make check-install` in CI regenerates in memory and fails the build if any committed
   artifact differs - the repo can never merge an out-of-sync artifact.

The value is then frozen into each artifact at the point it ships, and consumed later:

| Artifact | Value frozen at | Consumed at runtime by |
|---|---|---|
| `install/docker/docker-compose.yaml`, `deployment_yamls/k8s/*.yaml` | commit | `mesheryctl system start` downloads the committed file from the repo **at the release tag** |
| Helm `values.yaml` | commit / chart publish | `helm install` reads the packaged chart |
| `docker-extension/docker-compose.yaml` | commit / extension image build | the Docker Desktop extension |
| Playground compose + manifest | commit | applied to the cluster (must deploy from the committed manifest - GitOps gap) |
| `mesheryctl` `providers_gen.go` | `go build` | compiled-in fallback when the download path is not taken |
| `make` `REMOTE_PROVIDER_URLS` | `make` parse time | derived live from `providers.env`; no commit needed |

**Propagation latency:** because `mesheryctl` downloads the install files from a git ref at
runtime, a `providers.env` change reaches existing users only when they target a tag/branch
containing the regenerated files - i.e. at the next release. The same applies to chart
publishes, extension builds, and `mesheryctl` binary releases.

## Behavior changes

- **Docker Desktop extension**: gains `cloud.layer5.io` (now full set); keeps
  `PROVIDER=Meshery` so Meshery stays pre-selected.
- **`mesheryctl` embedded fallback**: drops `digitalocean`/`cleverluck`; now matches the
  canonical full set and is generated (`providers_gen.go`) rather than hand-maintained.
- Playground, helm, mesheryctl-docker, mesheryctl-k8s: no value change (already correct);
  they simply become generated from the canonical source.

## Testing

- `sync-provider-urls.py --check` passes on a freshly generated tree.
- Editing `providers.env` (add/remove/reorder) and running `generate-install` updates all
  full-profile files and `REMOTE_PROVIDER_URLS`; playground files reflect only Meshery.
- `--check` exits non-zero when a managed file is hand-edited out of sync.
- `make` prints the expected `REMOTE_PROVIDER_URLS` (`make -p` or a temporary echo target).
- Script errors when a managed file lacks its expected managed line.
- `helm template install/kubernetes/helm/meshery` renders the full list into the
  container env.
- All managed YAML files remain valid YAML after generation (parse check).
- `providers_gen.go` compiles and `mesheryctl` builds with `DefaultProviderBaseURLs`
  reflecting the canonical full list; `go vet ./...` is clean.
- `--check` exits non-zero when `providers_gen.go` is hand-edited out of sync.

## Risks

- **Regex fragility** if a file's line format changes: mitigated by the "fail loudly on
  missing line" rule and the `--check` CI gate.
- **Python dependency** in install tooling: acceptable (universal on dev/CI); documented.
- **Generated Go file committed:** if a developer hand-edits `providers_gen.go` or forgets
  to regenerate after editing `providers.env`, the build ships a stale list. Mitigated by
  the `make check-install` CI gate and the `DO NOT EDIT` header.
- **Propagation latency:** `providers.env` changes only reach users at the next release /
  chart publish / build, not retroactively (inherent to how artifacts are distributed).
