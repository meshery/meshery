---
title: Infrastructure Management
description: Meshery provides the ability to manage infrastructure for agility, maintainability, diversity, reliability and isolation, security, and speed.
weight: 35
aliases:
  - /guides/infrastructure-management/overview
  - /tasks/application-management
  - /tasks/infrastructure-management
---

## Using Meshery Designs to Manage Your Infrastructure

Meshery is a versatile platform designed to streamline the lifecycle, configuration, and performance management of infrastructure across Kubernetes clusters.

### Designs in Meshery

Meshery employs the concept of [Designs](/concepts/logical/designs) as a fundamental construct for managing infrastructure. Designs provide a structured way to organize and deploy various components of your infrastructure. To do so, Meshery utilizes a declarative approach to infrastructure management, similar to Kubernetes manifests. Meshery Designs are written in YAML and are validated against a schema.

#### Importing Existing Infrastructure and Applications

Meshery facilitates the import of infrastructure in multiple formats, including Kubernetes manifests, Helm Charts, and Docker Compose files. Import existing infrastructure definitions by directly from filesystem, via URL, or import directly from a GitHub repository.

- Kubernetes Manifest
- Meshery Design
- Helm Charts
- Docker Compose

See [Importing Designs](/extensions/kanvas/import-export-designs) for more information.

{{% details summary="### Meshery Internals Provisioning Process" %}}
When a request is made to provision a design, it undergoes the following stages:
#### 1. Import of Referenced Designs
A Design may reference any number of other Designs, in essence, a Design may import any number of other Designs.  As an editor of a Design, you can make reference to another Design, while following principles of reusing and DRY (Do Not Repeat Yourself). Any referenced Design will subsequently be imported during the provisioning moment. To reference another design, do so by adding the following annotation <pre>type: $(#use \<url-of-remote-pattern\>)</pre> in your Design file. The referenced design will be expanded from the source.

#### 2. Identification
Meshery relies on components registered at boot time. Only registered models and components can be managed with Meshery. Currently, models from the ArtifactHub repository are supported.

#### 3. Validation
Components in the design are validated against the schema, ensuring consistency, similar to Kubernetes object validation but tailored for Designs.

#### 4. Dependency Detection and Resolution
Meshery identifies and resolves dependencies among components using a dynamic mechanism based on each component’s origin (also known as its `host` or `registrant`). Provisioning order is critical—circular dependencies will result in termination of the deployment.

Deployment Mechanism by Source:
- **Artifact Hub:** Uses Helm Go client for Kubernetes Operator and CRD deployment via `ApplyHelmChart()`.
- **Kubernetes YAML:** Direct application with Kubernetes Go client (no auto-dependency handling).

This behavior is determined by the component’s host type:
<pre>
if connection.Kind == "artifacthub" ➜ Helm Go client used
if connection.Kind == "kubernetes" ➜ Kubernetes Go client used
</pre>

{{< alert title="Implementation Detail" type="info" >}}
This logic is handled in the `NewDependencyHandler()` function, which selects the deployment method based on `connection.Kind`."
{{< /alert >}}

#### 5. Provisioning
A Directed Acyclic Graph (DAG) generated in the previous step is processed. Dependent components are processed sequentially, while others are processed in parallel. Meshery intelligently handles the deployment order to ensure successful deployment.
{{% /details %}}

### Auto-Deployment of CRDs and Operators

Meshery automates the deployment of Custom Resource Definitions (CRDs) and operators based on the source from which a particular component was registered. By default, Meshery automatically deploys components that are sourced from Artifact Hub (utilizing Helm Charts). Support for OCI registries is expected in the near future.

{{% details summary="#### Understanding CRDs and Why Deployment Order Matters" %}}
**What is a CRD?**

Custom Resource Definitions (CRDs) extend Kubernetes to support new resource types (e.g., `VirtualService`, `Gateway`). A CRD must be installed **before** any resource of that type is deployed, or you’ll encounter errors like _“no matches for kind…”_.

**Why does Meshery only auto-deploy CRDs for ArtifactHub?**

Helm Charts (used by ArtifactHub) are bundled with CRDs and operators, making it safe for Meshery to auto-deploy them via Helm Go client.

In contrast, Kubernetes YAML files may be minimal or incomplete. Meshery cannot assume what dependencies you intended. Therefore, Meshery avoids auto-deployment to prevent accidental conflicts or errors.

**User Responsibility (for non-ArtifactHub components):**

If you use Kubernetes YAML, ensure that all required CRDs and operators are included. Meshery will apply the YAML as-is using the Kubernetes Go client, but will **not auto-deploy** missing dependencies for you.
{{% /details %}}

{{< alert title="Upcoming Features" type="info" >}}
Meshery is actively developing the ability to import custom models and components, expanding the platform's flexibility. Support for OCI registries is expected in the near future.

1. **Custom Models and Components:** Import your own models for more flexibility.
2. **OCI Registry Support:** In addition to ArtifactHub, OCI registries will be supported soon.
{{< /alert >}}
