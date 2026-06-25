---
title: Compatibility Matrix
aliases:
- /project/compatibility-matrix/
description: Understand which environments, dependencies, and installation paths Meshery supports.
display_title: false
cascade:
  type: compatibility-matrix
draft: true
weight: 100
---

## Dependency and Platform Compatibility for Installation

Use this page to choose the Meshery installation path that matches your environment. Meshery itself runs on [Docker]({{< ref "installation/docker/_index.md" >}}) or [Kubernetes]({{< ref "installation/kubernetes/_index.md" >}}). [mesheryctl]({{< ref "installation/mesheryctl/_index.md" >}}) is the common control point for installing, configuring, and operating Meshery across those targets.

{{% alert color="info" title="What is actually required?" %}}
Docker is enough to run Meshery locally and use many of its features. A Kubernetes cluster is required for most infrastructure lifecycle management features, including deploying and operating Cloud-based resources through Meshery. See the [Feature dependency matrix](#feature-dependency-matrix) below for specifics.
{{% /alert %}}

## Installation Path Matrix

| Installation path | Runs on | Required dependencies | Kubernetes required | Best for |
| --- | --- | --- | --- | --- |
| [`mesheryctl`]({{< ref "installation/mesheryctl/_index.md" >}}) only | macOS, Linux, Windows | `mesheryctl` installed via Bash, Homebrew, Scoop, or direct download | No | Installing, upgrading, and managing Meshery deployments |
| [Docker deployment]({{< ref "installation/docker/_index.md" >}}) (`mesheryctl system start -p docker`) | macOS, Linux, Windows | `mesheryctl`, Docker Engine or Docker Desktop, Docker Compose | No, to start Meshery. Yes, to manage Kubernetes infrastructure through Meshery. | Local evaluation and single-host deployments |
| [Kubernetes deployment]({{< ref "installation/kubernetes/_index.md" >}}) (`mesheryctl system start`) | Any environment that can reach a cluster | `mesheryctl`, `kubectl`, access to an active Kubernetes cluster | Yes | Full Meshery management workflows on Kubernetes |
| [Helm deployment]({{< ref "installation/kubernetes/helm.md" >}}) | Any environment that can reach a cluster | Helm v3, `kubectl`, access to an active Kubernetes cluster | Yes | Teams standardizing on Helm-based installation |
| [Docker Extension]({{< ref "installation/docker/docker-extension/index.md" >}}) | Docker Desktop | Docker Desktop, Kubernetes enabled in Docker Desktop, Docker Extension support | Yes | Docker Desktop users who want the shortest local setup |
| [GitHub Codespaces]({{< ref "installation/advanced/codespaces/index.md" >}}) | GitHub Codespaces | `mesheryctl`, Minikube, `kubectl`, enough Codespaces CPU and memory to run Minikube | Yes | Browser-based evaluation and contributor environments |

## Environment Matrix

| Environment | Supported installation paths | Hard dependencies | Recommended additions | Notes |
| --- | --- | --- | --- | --- |
| [Linux or macOS]({{< ref "installation/mesheryctl/linux-mac/_index.md" >}}) | `mesheryctl`, Docker, Kubernetes, Helm | `mesheryctl` for standard installs | Docker Desktop or Docker Engine, `kubectl`, Helm | This is the most direct path for local Docker and Kubernetes workflows. |
| [Windows]({{< ref "installation/mesheryctl/windows/_index.md" >}}) | `mesheryctl`, Docker, Docker Extension, Kubernetes | `mesheryctl`; Docker Desktop for container-based installs | WSL2, Kubernetes enabled in Docker Desktop, `kubectl` | Meshery supports Windows workflows, but Docker-backed paths are generally the smoothest option. |
| [Docker Desktop]({{< ref "installation/docker/docker-extension/index.md" >}}) | Docker deployment, Docker Extension | Docker Desktop with Compose support | Kubernetes enabled when you want cluster-backed features | Best fit for local development and evaluation on laptops. |
| Remote or managed Kubernetes cluster | Kubernetes deployment, Helm | Reachable cluster, valid kubeconfig, `kubectl` | LoadBalancer, Ingress, or port-forward access to Meshery | Good fit for shared team environments and production-style installs. See also [Upgrading Meshery]({{< ref "installation/upgrades/index.md" >}}). |
| [GitHub Codespaces]({{< ref "installation/advanced/codespaces/index.md" >}}) | Kubernetes deployment in Minikube | Codespace with enough resources, Minikube, `kubectl`, `mesheryctl` | VS Code desktop or browser-based Codespaces access | Useful when you do not want to manage a local cluster. |

## Feature Dependency Matrix

| Capability | Docker host only | Reachable Kubernetes cluster |
| --- | --- | --- |
| Start Meshery locally and open the UI | Yes | Optional |
| Run [`mesheryctl system check --preflight`]({{< ref "guides/mesheryctl/running-system-checks-using-mesheryctl.md" >}}) | Yes | Optional |
| Use performance management features | Yes | Yes |
| Deploy and manage Kubernetes infrastructure | No | Yes |
| Use [Meshery Operator and MeshSync]({{< ref "guides/troubleshooting/meshery-operator-meshsync.md" >}}) and cluster-backed adapters | No | Yes |
| Use the Docker Extension workflow | No | Yes, through Docker Desktop Kubernetes |

## Compatibility Stipulations

- Meshery should be installed onto platform versions that are still maintained by their upstream project or vendor.
- For Kubernetes-based installs, use a Kubernetes release that is still within the Kubernetes project's supported maintenance window and follow the operational guidance in [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}}).
- Match `kubectl` to the lifecycle of the cluster you target and keep it within the Kubernetes version skew policy.
- For Docker-based installs, use a current Docker Engine or Docker Desktop release that includes Compose support and can run the Meshery container set reliably. See [Install Meshery on Docker]({{< ref "installation/docker/_index.md" >}}).
- The Docker Extension path applies when Docker Desktop supports extensions and Kubernetes is enabled in Docker Desktop. See [Install Docker Extension for Meshery]({{< ref "installation/docker/docker-extension/index.md" >}}).
- Standard Meshery installation does not require local Go or Node.js runtimes. Those are contributor dependencies for [building Meshery from source]({{< ref "project/contributing/_index.md" >}}).
- If you plan to run Meshery out-of-cluster against Kubernetes, ensure the host running Meshery can reach the Kubernetes API and any broker endpoints exposed by your cluster. If connectivity fails, start with [Troubleshooting Meshery Installations]({{< ref "guides/troubleshooting/installation.md" >}}).

## Release Channel Guidance

Meshery publishes `stable` and `edge` release channels. Use the release channel that matches how quickly your environment adopts new Kubernetes minors. See [Build & Release (CI)]({{< ref "project/contributing/build-and-release.md" >}}) for how these channels are produced and [Upgrading Meshery]({{< ref "installation/upgrades/index.md" >}}) for changing deployed components over time.

| Release channel | Intended use | Kubernetes guidance |
| --- | --- | --- |
| [`stable`]({{< ref "project/contributing/build-and-release.md#stable-channel" >}}) | Production environments and conservative upgrades | Use for environments that stay on maintained platform releases and adopt Kubernetes upgrades through planned release management. |
| [`edge`]({{< ref "project/contributing/build-and-release.md#edge-channel" >}}) | Early validation, development, and upcoming-platform testing | Use when validating Meshery against newly introduced platform releases before they become part of your normal stable rollout. |

Stable and edge artifacts are published separately in Meshery's build and release process, so cluster compatibility validation should be considered together with your selected release channel and your [upgrade path]({{< ref "installation/upgrades/index.md" >}}).

## Validation Checks

Validate your environment before installing Meshery:

{{< code code="mesheryctl system check --preflight" >}}

For a deeper walkthrough of what these checks validate, see [Running system checks using Meshery CLI]({{< ref "guides/mesheryctl/running-system-checks-using-mesheryctl.md" >}}) and the [`mesheryctl system check` reference]({{< ref "reference/references/mesheryctl/system/check.md" >}}).

If you plan to use Docker, verify that the Docker Engine and Compose plugin are available:

{{< code code="docker version && docker compose version" >}}

If you plan to use Kubernetes, verify that your current context points to the target cluster:

{{< code code="kubectl config current-context" >}}

## Choosing the Right Path

- Choose the Docker Extension when you already standardize on Docker Desktop and want a guided local workflow.
- Choose Docker when you want to get Meshery running quickly on a single machine.
- Choose Kubernetes or Helm when Meshery needs to manage Kubernetes-backed infrastructure and team environments.
- Choose the Meshery Playground when you want to try Meshery without installing anything locally. See [Try Meshery Playground]({{< ref "installation/playground.md" >}}).

For rollout planning after initial installation, continue with [Upgrading Meshery]({{< ref "installation/upgrades/index.md" >}}).

## Related Reading

- [Install mesheryctl]({{< ref "installation/mesheryctl/_index.md" >}})
- [Install Meshery on Docker]({{< ref "installation/docker/_index.md" >}})
- [Install Meshery on Kubernetes]({{< ref "installation/kubernetes/_index.md" >}})
- [Install Docker Extension for Meshery]({{< ref "installation/docker/docker-extension/index.md" >}})
- [Build & Release (CI)]({{< ref "project/contributing/build-and-release.md" >}})
- [Upgrading Meshery]({{< ref "installation/upgrades/index.md" >}})
- [Running system checks using Meshery CLI]({{< ref "guides/mesheryctl/running-system-checks-using-mesheryctl.md" >}})
- [Troubleshooting Meshery Installations]({{< ref "guides/troubleshooting/installation.md" >}})
