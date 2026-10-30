# Running mesheryctl in a CI Pipeline

> Install mesheryctl on a CI runner, authenticate it without a browser, point it at a Meshery Server, and run commands non-interactively.

---

LLMS index: [llms.txt](/pr-preview/pr-22669/llms.txt)

---

Use this guide to run `mesheryctl` from a CI job, such as a GitHub Actions workflow, a GitLab pipeline, or a Jenkins agent. A CI runner has no browser and nobody to answer prompts, so you install a pinned `mesheryctl`, give it a token, and pass flags that skip every prompt. The examples use GitHub Actions, but each step is a shell command that works on any runner.

## Before you begin

- A runner with `bash`, `curl`, and `jq`. To start Meshery on the runner itself, the runner also needs Docker, or access to a Kubernetes cluster.
- A way to authenticate. Choose one:
  - **A Remote Provider token**, if your pipeline works with your organization's designs, environments, or performance results. Download it from Meshery Cloud as described in [Log in from a machine without a browser](/pr-preview/pr-22669/guides/mesheryctl/authenticate-with-meshery-via-cli/#log-in-from-a-machine-without-a-browser), and store the value of its `token` key as a CI secret. The examples use a secret named `MESHERY_TOKEN`.
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




<ul class="nav nav-tabs" id="tabs-1" role="tablist">
  <li class="nav-item">
      <button class="nav-link active"
          id="tabs-01-00-tab" data-bs-toggle="tab" data-bs-target="#tabs-01-00" role="tab"
          data-td-tp-persist="remote provider token" aria-controls="tabs-01-00" aria-selected="true">
        Remote Provider token
      </button>
    </li><li class="nav-item">
      <button class="nav-link"
          id="tabs-01-01-tab" data-bs-toggle="tab" data-bs-target="#tabs-01-01" role="tab"
          data-td-tp-persist="local provider" aria-controls="tabs-01-01" aria-selected="false">
        Local provider
      </button>
    </li>
</ul>

<div class="tab-content" id="tabs-1-content">
    <div class="tab-body tab-pane fade show active"
        id="tabs-01-00" role="tabpanel" aria-labelled-by="tabs-01-00-tab" tabindex="1">
        <p>Write the token from your secret into a token file, then register the file and attach it to the current context:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">jq -n --arg token <span class="s2">&#34;</span><span class="nv">$MESHERY_TOKEN</span><span class="s2">&#34;</span> <span class="s1">&#39;{&#34;meshery-provider&#34;: &#34;Meshery&#34;, &#34;token&#34;: $token}&#39;</span> &gt; <span class="s2">&#34;</span><span class="nv">$RUNNER_TEMP</span><span class="s2">/auth.json&#34;</span>
</span></span><span class="line"><span class="cl">chmod <span class="m">600</span> <span class="s2">&#34;</span><span class="nv">$RUNNER_TEMP</span><span class="s2">/auth.json&#34;</span>
</span></span><span class="line"><span class="cl">mesheryctl system token create ci -f <span class="s2">&#34;</span><span class="nv">$RUNNER_TEMP</span><span class="s2">/auth.json&#34;</span> --set
</span></span></code></pre></div><p>Keep the token file outside your workspace so that later steps can&rsquo;t upload it as an artifact by accident. <code>$RUNNER_TEMP</code> is the GitHub Actions temporary directory. Use your CI system&rsquo;s equivalent elsewhere.</p>

    </div>
    <div class="tab-body tab-pane fade"
        id="tabs-01-01" role="tabpanel" aria-labelled-by="tabs-01-01-tab" tabindex="1">
        <p>After Meshery Server is running, log in with the Local provider. It signs in without a browser:</p>
<div class="highlight"><pre tabindex="0" class="chroma"><code class="language-bash" data-lang="bash"><span class="line"><span class="cl">mesheryctl system login -p Local
</span></span></code></pre></div>
    </div>
</div>


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

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Deploying designs from CI</div>


`mesheryctl design deploy` and `mesheryctl design apply` do not currently tell Meshery Server which clusters to target, so Meshery Server rejects them. To deploy a design, use Meshery UI as described in [Deploying a Design](/pr-preview/pr-22669/guides/infrastructure-management/deploying-a-design/).
</div>


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

- [Authenticating Meshery via CLI](/pr-preview/pr-22669/guides/mesheryctl/authenticate-with-meshery-via-cli/)
- [Working with mesheryctl](/pr-preview/pr-22669/guides/mesheryctl/working-with-mesheryctl/)
- [GitOps with Meshery](/pr-preview/pr-22669/guides/infrastructure-management/gitops-with-meshery/) for Meshery's own GitHub Actions
- [`mesheryctl` command reference](/pr-preview/pr-22669/reference/mesheryctl/)

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
