# Installation

> Choose where to run Meshery, on Kubernetes or a Docker host, and find the installation guide for your platform.

---

LLMS index: [llms.txt](/llms.txt)

---

## Supported Platforms

Meshery runs as a set of containers, either on a Docker host or in a Kubernetes cluster. On a Docker host, Meshery runs as a standalone management plane outside the clusters it manages, which is an _out-of-cluster_ deployment. Inside a Kubernetes cluster, it is an _in-cluster_ deployment. The supported platforms are listed below.

To decide which installation path fits your environment, start with the [Compatibility Matrix](/project/compatibility-matrix/compatibility-matrix/).

For runtime configuration after installation, see [Meshery Server Environment Variables](/installation/advanced/environment-variables/).

Before you deploy Meshery for production use, review the [Production Deployment](/installation/production/) guidance on high availability, security hardening, performance, multi-cluster and multi-cloud operation, and operational readiness.

---

Section pages:

- [Overview](/installation/overview/): Meshery is the self-service engineering platform, enabling collaborative design and operation of cloud and cloud native infrastructure.
- [Quick Start Guide](/installation/quick-start/): Install mesheryctl and start Meshery on a Kubernetes cluster or a Docker host with a single command, then connect your clusters and start designing.
- [Using Meshery Playground](/installation/playground/): Try Meshery in your browser, connected to live Kubernetes clusters, without installing anything.
- [Install using mesheryctl](/installation/mesheryctl/): Use Meshery CLI to install Meshery on supported platforms.
- [Install on Docker](/installation/docker/): Install Meshery on Docker
- [Install on Kubernetes](/installation/kubernetes/): Install Meshery on Kubernetes. Deploy Meshery in Kubernetes in-cluster or outside of Kubernetes out-of-cluster.
- [Production Deployment](/installation/production/): Considerations, best practices, security hardening, performance bounds, and operational readiness guidance for deploying Meshery in production.
- [Advanced Configuration](/installation/advanced/): Advanced installation and runtime configuration topics for Meshery.
- [Upgrading Meshery](/installation/upgrades/): How Meshery's components version together, and the step-by-step procedure for upgrading the CLI, Server, Operator, MeshSync, and Broker on Docker and Kubernetes deployments.
- [Accessing Meshery UI](/installation/accessing-meshery-ui/): Accessing Meshery's web-based user interface
