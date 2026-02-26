---
title: Compatibility Matrix
aliases: 
 - /project/compatibility-matrix/
description: An installation compatibility matrix and project test status dashboard.
display_title: false
cascade:
  type: compatibility-matrix
---

# Compatibility Matrix

Meshery Server and Meshery Adapters are tested daily for their compatibility with the infrastructure they manage and the platforms Meshery deploys on (Kubernetes and Docker). End-to-end test results are automatically posted to the following compatibility matrix.

See also https://docs.meshery.io/project/contributing/test-status, which needs to be combined with the Compatibility Matrix test results to come together under a unified page (set of drillable pages).

{{< compatibility-matrix-kubernetes >}}

## Integration Tests

As a key aspect of Meshery, its integrations with other systems are routinely tested. Unit and integration tests before and after every pull request (before code is to be merged into the project and after code is merged into the project). End-to-end tests are run nightly and automatically posted to the following test matrix.

{{< integration-tests >}}
