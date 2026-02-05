---
title: "Compatibility Matrix"
description: "Meshery Server and Adapter compatibility with infrastructure platforms"
weight: 55
aliases:
  - /installation/compatibility-matrix/
display_title: "false"
---

# Compatibility Matrix

Meshery Server and Meshery Adapters are tested daily for their compatibility with the infrastructure they manage and the platforms Meshery deploys on (Kubernetes and Docker). Integration test results are automatically posted to the following compatibility matrix.

## Kubernetes

Meshery is compatible with all versions of Kubernetes. The following table shows the tested versions:

| Kubernetes Version | Status |
|-------------------|--------|
| 1.29.x | ✅ Tested |
| 1.28.x | ✅ Tested |
| 1.27.x | ✅ Tested |
| 1.26.x | ✅ Tested |
| 1.25.x | ✅ Tested |

## Docker

The following minimum Docker build versions are required:

| Name | Version |
|------|---------|
| [Docker Engine](/installation/docker) | **19.x** and above |
| [Docker Desktop](/installation/docker/docker-extension) (via Docker Extension) | **2.0.x** and above |

## Integration Tests

For a complete overview of the latest integration tests and their status, please visit the [Meshery Integration Tests](https://github.com/meshery/meshery/actions) page.

{{< related-discussions tag="meshery" >}}
