---
title: Running mesheryctl in a CI Pipeline
categories: [mesheryctl]
description: Install mesheryctl on a CI runner, authenticate it without a browser, point it at a Meshery Server, and run commands non-interactively.
---

Use this guide to run `mesheryctl` from a CI job, such as a GitHub Actions workflow, a GitLab pipeline, or a Jenkins agent. A CI runner has no browser and nobody to answer prompts, so you install a pinned `mesheryctl`, give it a token, and pass flags that skip every prompt. The examples use GitHub Actions, but each step is a shell command that works on any runner.

## Before you begin

- A runner with `bash`, `curl`, and `jq`. To start Meshery on the runner itself, the runner also needs Docker, or access to a Kubernetes cluster.
- A way to authenticate. Choose one:
  - **A Remote Provider token**, if your pipeline works with your organization's designs, environments, or performance results. Download it from Meshery Cloud as described in [Log in from a machine without a browser]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md#log-in-from-a-machine-without-a-browser" >}}), and store the value of its `token` key as a CI secret. The examples use a secret named `MESHERY_TOKEN`.
  - **The Local provider**, if the pipeline starts a throwaway Meshery Server and needs nothing from your organization. It needs no secret.

## Install a pinned mesheryctl

Pin the version so that a new release doesn't change your pipeline's behavior. Pick a version from the [Meshery releases](https://github.com/meshery/meshery/releases), then install only the CLI with `-n`:

```bash
curl -L https://meshery.io/install | MESHERY_VERSION=<version> bash -s -- -n
```

The script installs `mesheryctl` to `/usr/local/bin`. Set `MESHERY_INSTALL_DIR` to install it somewhere else, for example a directory that you add to `PATH`.

## Point mesheryctl at a Meshery Server

Choose the branch that matches your pipeline.

### Use an existing Meshery Server

Create a context for the server and make it current. `--yes` answers the confirmation prompt:

```bash
mesheryctl system context create ci --url https://meshery.example.com --platform kubernetes --set --yes
```

### Start Meshery on the runner

Start Meshery Server in Docker. `--skip-browser` keeps `mesheryctl` from trying to open Meshery UI, and `--skip-update` skips the image update check. `mesheryctl` returns before the server is ready, so wait for the version endpoint to answer:

```bash
mesheryctl system start --platform docker --skip-browser --skip-update --yes
for attempt in $(seq 1 30); do
  curl --silent --fail http://localhost:9081/api/system/version > /dev/null && break
  sleep 10
done
curl --silent --show-error --fail http://localhost:9081/api/system/version
```

The final `curl` fails the job if the server never came up.

## Authenticate without a browser

{{< tabpane text=true >}}
{{% tab header="Remote Provider token" %}}
Write the token from your secret into a token file, then register the file and attach it to the current context:

```bash
jq -n --arg token "$MESHERY_TOKEN" '{"meshery-provider": "Meshery", "token": $token}' > "$RUNNER_TEMP/auth.json"
chmod 600 "$RUNNER_TEMP/auth.json"
mesheryctl system token create ci -f "$RUNNER_TEMP/auth.json" --set
```

Keep the token file outside your workspace so that later steps can't upload it as an artifact by accident. `$RUNNER_TEMP` is the GitHub Actions temporary directory. Use your CI system's equivalent elsewhere.
{{% /tab %}}
{{% tab header="Local provider" %}}
After Meshery Server is running, log in with the Local provider. It signs in without a browser:

```bash
mesheryctl system login -p Local
```
{{% /tab %}}
{{< /tabpane >}}

Confirm that `mesheryctl` can reach the server with your token:

```bash
mesheryctl version
mesheryctl connection list
```

## Run commands without prompts

- Pass `--yes` (`-y`) to commands that ask for confirmation, such as `system` subcommands and `perf apply`.
- Refer to designs, profiles, and other resources by ID, or by a name that matches only one of them. Some commands prompt you to choose when a name matches more than one resource, and a prompt stalls the job.
- Use `--output-format json` (`-o json`) on `view` commands, and parse the result with `jq`, instead of scraping table output.
- Add `--verbose` when a step fails, to see the requests `mesheryctl` makes.

Common pipeline steps:

```bash
# Import the design stored in your repository, so your team sees the latest version in Meshery
mesheryctl design import -f deploy/app.yaml -n "my-app"

# Package and import a model you maintain
mesheryctl model build my-model/v1.0.0
mesheryctl model import -f my-model-v1-0-0.tar

# Run a saved performance profile and fail the job if mesheryctl fails
mesheryctl perf apply "checkout latency" -y

# Export the current design for an artifact
mesheryctl design export "my-app" --type oci
```

{{% alert color="info" title="Deploying designs from CI" %}}
`mesheryctl design deploy` and `mesheryctl design apply` do not currently tell Meshery Server which clusters to target, so Meshery Server rejects them. To deploy a design, use Meshery UI as described in [Deploying a Design]({{< ref "guides/infrastructure-management/deploying-a-design.md" >}}).
{{% /alert %}}

## A complete GitHub Actions job

This job installs `mesheryctl`, authenticates with a Remote Provider token against an existing Meshery Server, and imports the design kept in the repository:

```yaml
name: Sync design to Meshery
on:
  push:
    branches: [main]
    paths: ["deploy/**"]

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install mesheryctl
        run: curl -L https://meshery.io/install | MESHERY_VERSION=<version> bash -s -- -n

      - name: Configure mesheryctl
        env:
          MESHERY_URL: ${{ vars.MESHERY_URL }}
          MESHERY_TOKEN: ${{ secrets.MESHERY_TOKEN }}
        run: |
          mesheryctl system context create ci --url "$MESHERY_URL" --platform kubernetes --set --yes
          jq -n --arg token "$MESHERY_TOKEN" '{"meshery-provider": "Meshery", "token": $token}' > "$RUNNER_TEMP/auth.json"
          chmod 600 "$RUNNER_TEMP/auth.json"
          mesheryctl system token create ci -f "$RUNNER_TEMP/auth.json" --set

      - name: Import design
        run: mesheryctl design import -f deploy/app.yaml -n "my-app"
```

Store the server URL as a repository variable named `MESHERY_URL`, and the token as a secret named `MESHERY_TOKEN`. The job passes both to the shell as environment variables, rather than expanding `${{ vars.MESHERY_URL }}` inside the script, so a crafted variable value can't inject commands into a step that holds your token.

## Related

- [Authenticating Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}})
- [Working with mesheryctl]({{< ref "guides/mesheryctl/working-with-mesheryctl.md" >}})
- [GitOps with Meshery]({{< ref "guides/infrastructure-management/gitops-with-meshery/index.md" >}}) for Meshery's own GitHub Actions
- [`mesheryctl` command reference]({{< ref "reference/mesheryctl/_index.md" >}})

{{< related-discussions tag="mesheryctl" >}}
