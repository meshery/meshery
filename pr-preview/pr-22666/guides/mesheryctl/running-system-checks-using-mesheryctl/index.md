# Running system checks using Meshery CLI

> Run pre-flight and post-deployment system health checks using Meshery's CLI, mesheryctl.

---

LLMS index: [llms.txt](/pr-preview/pr-22666/llms.txt)

---

`mesheryctl system check` verifies that your machine is ready to deploy Meshery and that an existing Meshery deployment is healthy. Run it before you deploy Meshery, and again whenever Meshery misbehaves.

Run `mesheryctl system check --help` to see every flag, or see the [`mesheryctl system check` reference](/pr-preview/pr-22666/reference/mesheryctl/).

## Check your machine before you deploy

Run the pre-deployment checks with `--preflight` (or its alias `--pre`):

```bash
mesheryctl system check --preflight
```

`mesheryctl` runs these checks:

- **Docker:** Docker is running and `docker-compose` is available.
- **Kubernetes API:** `mesheryctl` can create a Kubernetes client from your kubeconfig and query the Kubernetes API.
- **Kubernetes version:** your cluster and `kubectl` meet Meshery's minimum versions.

It ends with `Meshery prerequisites met` or `Meshery prerequisites not met`. Only the checks for your context's platform count toward that result. For example, a failed Kubernetes check does not fail the preflight when your context's platform is `docker`.

## Check a running deployment

Run `mesheryctl system check` with no flags after you deploy Meshery:

```bash
mesheryctl system check
```

`mesheryctl` runs the checks that match your context's platform:

- **Docker platform:** the Docker checks.
- **Kubernetes platform:** the Kubernetes API and version checks, plus the Meshery Operator checks.
- **Both platforms:**
  - **Meshery version:** Meshery Server and `mesheryctl` are on the latest release.
  - **Meshery components:** each Meshery Adapter that is registered with Meshery Server is running and reachable.

## Check one part of a deployment

Each of these flags runs only its own check:

| To check | Run |
| --- | --- |
| Every registered Meshery Adapter | `mesheryctl system check --adapters` (or `--components`) |
| One Meshery Adapter | `mesheryctl system check --adapter meshery-istio`, or include the port: `--adapter meshery-istio:10000` |
| Meshery Operator, MeshSync, and Meshery Broker (Kubernetes only) | `mesheryctl system check --operator` |

The adapter checks ask Meshery Server for your session, so log in with `mesheryctl system login` first.

## FAQ

### `mesheryctl system check --preflight` says Kubernetes is missing, but it also says "Meshery prerequisites met". Is that all right?

Yes, if your context's platform is `docker`. Meshery runs on Docker alone. You need a Kubernetes cluster only to manage Kubernetes infrastructure with Meshery.

### The preflight passed, but Meshery doesn't start. What should I do?

Check the installation steps for your platform in [Installation](/pr-preview/pr-22666/installation/), then see [Troubleshooting Meshery Installations](/pr-preview/pr-22666/guides/troubleshooting/installation/).

### Do I need a Kubernetes cluster, or is a Docker host enough?

[Performance management](/pr-preview/pr-22666/guides/performance-management/managing-performance/) works without a Kubernetes cluster. Managing cloud native infrastructure needs at least one Kubernetes cluster connected to Meshery.

### What does Meshery need in production?

See [Production Deployments](/pr-preview/pr-22666/installation/production/).

### Do I need any add-ons to pass the checks?

No. The checks only need Docker or Kubernetes, depending on your platform.

### The adapter check fails with "Auth token not found". What should I do?

Run `mesheryctl system login` to get a token, then run the check again. See [Authenticating Meshery via CLI](/pr-preview/pr-22666/guides/mesheryctl/authenticate-with-meshery-via-cli/).

### My Kubernetes cluster is running, but Meshery can't reach it and the checks fail. What should I do?

Connect the cluster to Meshery. In Meshery UI, click **Add Cluster** in the context switcher, or go to **Lifecycle** > **Connections** and create a Kubernetes connection, then import your kubeconfig. From the CLI, run `mesheryctl connection create --type <aks|eks|gke|minikube>`. See [Registering a Connection](/pr-preview/pr-22666/guides/infrastructure-management/registering-a-connection/).

### The Meshery version check says "CLI is not up-to-date". Should I update mesheryctl often?

Yes. Keep `mesheryctl` on the same release as Meshery Server, and update both together. See [Upgrading Meshery](/pr-preview/pr-22666/installation/upgrades/).

### Should Meshery be running when I run the checks?

Yes, for every check except `--preflight`. The version, adapter, and operator checks all talk to a running Meshery deployment.

### What versions of Kubernetes and kubectl does Meshery need?

Use a Kubernetes minor release that upstream still supports, and keep `kubectl` within the [version skew policy](https://kubernetes.io/releases/version-skew-policy/) of your cluster. See the [Kubernetes release history](https://kubernetes.io/releases/) for the supported releases.

### Only some Meshery Adapters show up in the check. Is that a problem?

No. The check covers only the adapters that are registered with Meshery Server, and the default context deploys none. To deploy more, see [Multiple Adapters](/pr-preview/pr-22666/installation/advanced/multiple-adapters/).

### I started Meshery fresh, and every adapter check fails. What should I do?

Connect a Kubernetes cluster to Meshery as described above, then run the check again. If it still fails, [open an issue](https://github.com/meshery/meshery/issues).

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
