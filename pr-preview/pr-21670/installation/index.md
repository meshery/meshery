# Installation

> Installation procedures for deploying Meshery with mesheryctl.

Source: /pr-preview/pr-21670/installation/

## Supported Platforms

Meshery deploys as a set of Docker containers to either a Docker host or a Kubernetes cluster. A given deployment of Meshery can be described as either an _in-cluster_ or an _out-of-cluster_ deployment. Meshery runs as a standalone management plane on a Docker host (_out-of-cluster_) or within a Kubernetes cluster (_in-cluster_). See the complete list of supported platforms below.

If you are deciding which installation path fits your environment, start with the [Compatibility Matrix](/pr-preview/pr-21670/project/compatibility-matrix/compatibility-matrix/).

For runtime configuration after installation, see [Meshery Server Environment Variables](/pr-preview/pr-21670/installation/advanced/environment-variables/).

If you are deploying Meshery for production use, review the [Production Deployment](/pr-preview/pr-21670/installation/production/) considerations for guidance on high availability, security hardening, performance, multi-cluster and multi-cloud operation, and operational readiness.
