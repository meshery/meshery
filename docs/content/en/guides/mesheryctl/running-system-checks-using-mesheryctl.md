---
title: Running system checks using Meshery CLI
categories: [mesheryctl]
description: Run pre-flight and post-deployment system health checks using Meshery's CLI, mesheryctl.
---

`mesheryctl system check` verifies that your machine is ready to deploy Meshery and that an existing Meshery deployment is healthy. Run it before you deploy Meshery, and again whenever Meshery misbehaves.

Run `mesheryctl system check --help` to see every flag, or see the [`mesheryctl system check` reference]({{< ref "reference/references/mesheryctl/_index.md" >}}).

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

Check the installation steps for your platform in [Installation]({{< ref "installation/_index.md" >}}), then see [Troubleshooting Meshery Installations]({{< ref "guides/troubleshooting/installation/index.md" >}}).

### Do I need a Kubernetes cluster, or is a Docker host enough?

[Performance management]({{< ref "guides/performance-management/managing-performance/index.md" >}}) works without a Kubernetes cluster. Managing cloud native infrastructure needs at least one Kubernetes cluster connected to Meshery.

### What does Meshery need in production?

See [Production Deployments]({{< ref "installation/production/_index.md" >}}).

### Do I need any add-ons to pass the checks?

No. The checks only need Docker or Kubernetes, depending on your platform.

### The adapter check fails with "Auth token not found". What should I do?

Run `mesheryctl system login` to get a token, then run the check again. See [Authenticating Meshery via CLI]({{< ref "guides/mesheryctl/authenticate-with-meshery-via-cli/index.md" >}}).

### My Kubernetes cluster is running, but Meshery can't reach it and the checks fail. What should I do?

Connect the cluster to Meshery. In Meshery UI, click **Add Cluster** in the context switcher, or go to **Lifecycle** > **Connections** and create a Kubernetes connection, then import your kubeconfig. From the CLI, run `mesheryctl connection create --type <aks|eks|gke|minikube>`. See [Registering a Connection]({{< ref "guides/infrastructure-management/registering-a-connection.md" >}}).

### The Meshery version check says "CLI is not up-to-date". Should I update mesheryctl often?

Yes. Keep `mesheryctl` on the same release as Meshery Server, and update both together. See [Upgrading Meshery]({{< ref "installation/upgrades/index.md" >}}).

### Should Meshery be running when I run the checks?

Yes, for every check except `--preflight`. The version, adapter, and operator checks all talk to a running Meshery deployment.

### What versions of Kubernetes and kubectl does Meshery need?

Use a Kubernetes minor release that upstream still supports, and keep `kubectl` within the [version skew policy](https://kubernetes.io/releases/version-skew-policy/) of your cluster. See the [Kubernetes release history](https://kubernetes.io/releases/) for the supported releases.

### Only some Meshery Adapters show up in the check. Is that a problem?

No. The check covers only the adapters that are registered with Meshery Server, and the default context deploys none. To deploy more, see [Multiple Adapters]({{< ref "installation/advanced/multiple-adapters.md" >}}).

### I started Meshery fresh, and every adapter check fails. What should I do?

Connect a Kubernetes cluster to Meshery as described above, then run the check again. If it still fails, [open an issue](https://github.com/meshery/meshery/issues).

{{< related-discussions tag="mesheryctl" >}}
