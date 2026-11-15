---
title: Installation
linkTitle: "🗺️ Overview & Installation"
description: Choose where to run Meshery, on Kubernetes or a Docker host, and find the installation guide for your platform.
aliases: 
- /platforms
- /platforms/
- /installation/platforms
- /installation/platforms/
- /installation/
weight: 1
---

## Supported Platforms

Meshery runs as a set of containers, either on a Docker host or in a Kubernetes cluster. On a Docker host, Meshery runs as a standalone management plane outside the clusters it manages, which is an _out-of-cluster_ deployment. Inside a Kubernetes cluster, it is an _in-cluster_ deployment. The supported platforms are listed below.

To decide which installation path fits your environment, start with the [Compatibility Matrix]({{< ref "project/compatibility-matrix/compatibility-matrix.md" >}}).

For runtime configuration after installation, see [Meshery Server Environment Variables]({{< ref "installation/advanced/environment-variables.md" >}}).

Before you deploy Meshery for production use, review the [Production Deployment]({{< ref "installation/production/_index.md" >}}) guidance on high availability, security hardening, performance, multi-cluster and multi-cloud operation, and operational readiness.
