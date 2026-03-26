---
title: Compatibility Matrix
aliases:
- /project/compatibility-matrix/
description: Understand which environments, dependencies, and installation paths Meshery supports.
display_title: false
cascade:
  type: compatibility-matrix
---

## Dependency and Platform Compatibility for Installation

Use this page to choose the Meshery installation path that matches your environment. Meshery itself runs on Docker or Kubernetes. `mesheryctl` is the common control point for installing, configuring, and operating Meshery across those targets.

{{% alert color="info" title="What is actually required?" %}}
Docker is enough to run Meshery locally. A Kubernetes cluster is required for most infrastructure lifecycle management features, including deploying and operating Kubernetes-based resources through Meshery.
{{% /alert %}}

## Installation path matrix

| Installation path | Runs on | Required dependencies | Kubernetes required | Best for |
| --- | --- | --- | --- | --- |
| `mesheryctl` only | macOS, Linux, Windows | `mesheryctl` installed via Bash, Homebrew, Scoop, or direct download | No | Installing, upgrading, and managing Meshery deployments |
| Docker deployment (`mesheryctl system start -p docker`) | macOS, Linux, Windows | `mesheryctl`, Docker Engine or Docker Desktop, Docker Compose | No, to start Meshery. Yes, to manage Kubernetes infrastructure through Meshery. | Local evaluation and single-host deployments |
| Kubernetes deployment (`mesheryctl system start`) | Any environment that can reach a cluster | `mesheryctl`, `kubectl`, access to an active Kubernetes cluster | Yes | Full Meshery management workflows on Kubernetes |
| Helm deployment | Any environment that can reach a cluster | Helm v3, `kubectl`, access to an active Kubernetes cluster | Yes | Teams standardizing on Helm-based installation |
| Docker Extension | Docker Desktop | Docker Desktop, Kubernetes enabled in Docker Desktop, Docker Extension support | Yes | Docker Desktop users who want the shortest local setup |
| GitHub Codespaces | GitHub Codespaces | `mesheryctl`, Minikube, `kubectl`, enough Codespaces CPU and memory to run Minikube | Yes | Browser-based evaluation and contributor environments |

## Environment matrix

| Environment | Supported installation paths | Hard dependencies | Recommended additions | Notes |
| --- | --- | --- | --- | --- |
| macOS or Linux | `mesheryctl`, Docker, Kubernetes, Helm | `mesheryctl` for standard installs | Docker Desktop or Docker Engine, `kubectl`, Helm | This is the most direct path for local Docker and Kubernetes workflows. |
| Windows | `mesheryctl`, Docker, Docker Extension, Kubernetes | `mesheryctl`; Docker Desktop for container-based installs | WSL2, Kubernetes enabled in Docker Desktop, `kubectl` | Meshery supports Windows workflows, but Docker-backed paths are generally the smoothest option. |
| Docker Desktop | Docker deployment, Docker Extension | Docker Desktop with Compose support | Kubernetes enabled when you want cluster-backed features | Best fit for local development and evaluation on laptops. |
| Remote or managed Kubernetes cluster | Kubernetes deployment, Helm | Reachable cluster, valid kubeconfig, `kubectl` | LoadBalancer, Ingress, or port-forward access to Meshery | Good fit for shared team environments and production-style installs. |
| GitHub Codespaces | Kubernetes deployment in Minikube | Codespace with enough resources, Minikube, `kubectl`, `mesheryctl` | VS Code desktop or browser-based Codespaces access | Useful when you do not want to manage a local cluster. |

## Feature dependency matrix

| Capability | Docker host only | Reachable Kubernetes cluster |
| --- | --- | --- |
| Start Meshery locally and open the UI | Yes | Optional |
| Run `mesheryctl system check --preflight` | Yes | Optional |
| Use performance management features | Yes | No |
| Deploy and manage Kubernetes infrastructure | No | Yes |
| Use Meshery Operator, MeshSync, and cluster-backed adapters | No | Yes |
| Use the Docker Extension workflow | No | Yes, through Docker Desktop Kubernetes |

## Compatibility stipulations

- `mesheryctl` is the recommended installation and operations interface for Meshery across all supported environments.
- Meshery should be installed onto platform versions that are still maintained by their upstream project or vendor.
- For Kubernetes-based installs, use a Kubernetes release that is still within the Kubernetes project's supported maintenance window.
- Match `kubectl` to the lifecycle of the cluster you target and keep it within the Kubernetes version skew policy.
- For Docker-based installs, use a current Docker Engine or Docker Desktop release that includes Compose support and can run the Meshery container set reliably.
- The Docker Extension path applies when Docker Desktop supports extensions and Kubernetes is enabled in Docker Desktop.
- Standard Meshery installation does not require local Go or Node.js runtimes. Those are contributor dependencies for building Meshery from source.
- If you plan to run Meshery out-of-cluster against Kubernetes, ensure the host running Meshery can reach the Kubernetes API and any broker endpoints exposed by your cluster.

## Release channel guidance

Meshery publishes `stable` and `edge` release channels. Use the release channel that matches how quickly your environment adopts new Kubernetes minors.

| Release channel | Intended use | Kubernetes guidance |
| --- | --- | --- |
| `stable` | Production environments and conservative upgrades | Use for environments that stay on maintained platform releases and adopt Kubernetes upgrades through planned release management. |
| `edge` | Early validation, development, and upcoming-platform testing | Use when validating Meshery against newly introduced platform releases before they become part of your normal stable rollout. |

Stable and edge artifacts are published separately in Meshery's build and release process, so cluster compatibility validation should be considered together with your selected release channel.

## Validation checks

Validate your environment before installing Meshery:

{{< code code="mesheryctl system check --preflight" >}}

If you plan to use Docker, verify that the Docker Engine and Compose plugin are available:

{{< code code="docker version && docker compose version" >}}

If you plan to use Kubernetes, verify that your current context points to the target cluster:

{{< code code="kubectl config current-context" >}}

## Choosing the right path

- Choose Docker when you want to get Meshery running quickly on a single machine.
- Choose Kubernetes or Helm when Meshery needs to manage Kubernetes-backed infrastructure and team environments.
- Choose the Docker Extension when you already standardize on Docker Desktop and want a guided local workflow.
- Choose Codespaces when you want a disposable environment without managing local cluster tooling.

## Related reading

- [Install mesheryctl](/installation/mesheryctl)
- [Install Meshery on Docker](/installation/docker)
- [Install Meshery on Kubernetes](/installation/kubernetes)
- [Install Docker Extension for Meshery](/installation/docker/docker-extension)
- [Running system checks using Meshery CLI](/guides/mesheryctl/running-system-checks-using-mesheryctl)
- [Test status](/project/contributing/test-status) for automated CI and integration test visibility

<!-- 
See also [test status](/project/contributing/test-status), which needs to be combined with the Compatibility Matrix test results to come together under a unified page (set of drillable pages).

{{< compatibility-matrix-kubernetes >}}

## Integration Tests

As a key aspect of Meshery, its integrations with other systems are routinely tested. Unit and integration tests before and after every pull request (before code is to be merged into the project and after code is merged into the project). End-to-end tests are run nightly and automatically posted to the following test matrix.

{{< integration-tests >}} -->
