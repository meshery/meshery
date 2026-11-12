# Installation

> Choose where to run Meshery, on Kubernetes or a Docker host, and find the installation guide for your platform.

---

LLMS index: [llms.txt](/docs-preview/pr-preview/pr-23341/llms.txt)

---

## Supported Platforms

Meshery runs as a set of containers, either on a Docker host or in a Kubernetes cluster. On a Docker host, Meshery runs as a standalone management plane outside the clusters it manages, which is an _out-of-cluster_ deployment. Inside a Kubernetes cluster, it is an _in-cluster_ deployment. The supported platforms are listed below.

To decide which installation path fits your environment, start with the [Compatibility Matrix](/docs-preview/pr-preview/pr-23341/project/compatibility-matrix/compatibility-matrix/).

For runtime configuration after installation, see [Meshery Server Environment Variables](/docs-preview/pr-preview/pr-23341/installation/advanced/environment-variables/).

Before you deploy Meshery for production use, review the [Production Deployment](/docs-preview/pr-preview/pr-23341/installation/production/) guidance on high availability, security hardening, performance, multi-cluster and multi-cloud operation, and operational readiness.

---

Section pages:

- [Overview](/docs-preview/pr-preview/pr-23341/installation/overview/): Meshery is the self-service engineering platform, enabling collaborative design and operation of cloud and cloud native infrastructure.
- [Quick Start Guide](/docs-preview/pr-preview/pr-23341/installation/quick-start/): Install mesheryctl and start Meshery on a Kubernetes cluster or a Docker host with a single command, then connect your clusters and start designing.
- [Using Meshery Playground](/docs-preview/pr-preview/pr-23341/installation/playground/): Try Meshery in your browser, connected to live Kubernetes clusters, without installing anything.
- [Install using mesheryctl](/docs-preview/pr-preview/pr-23341/installation/mesheryctl/): Use Meshery CLI to install Meshery on supported platforms.
- [Install on Docker](/docs-preview/pr-preview/pr-23341/installation/docker/): Install Meshery on Docker
- [Install on Kubernetes](/docs-preview/pr-preview/pr-23341/installation/kubernetes/): Install Meshery on Kubernetes. Deploy Meshery in Kubernetes in-cluster or outside of Kubernetes out-of-cluster.
- [Production Deployment](/docs-preview/pr-preview/pr-23341/installation/production/): Considerations, best practices, security hardening, performance bounds, and operational readiness guidance for deploying Meshery in production.
- [Advanced Configuration](/docs-preview/pr-preview/pr-23341/installation/advanced/): Advanced installation and runtime configuration topics for Meshery.
- [Upgrading Meshery](/docs-preview/pr-preview/pr-23341/installation/upgrades/): How Meshery's components version together, and the step-by-step procedure for upgrading the CLI, Server, Operator, MeshSync, and Broker on Docker and Kubernetes deployments.
- [Accessing Meshery UI](/docs-preview/pr-preview/pr-23341/installation/accessing-meshery-ui/): Accessing Meshery's web-based user interface
