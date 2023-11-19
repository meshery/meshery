---
layout: default
title: Infrastructure Management
permalink: tasks/infrastructure-management
type: tasks
language: en
list: include
---

## Infrastructure Management with Meshery Designs

Meshery is a versatile platform designed to streamline the storage, management, configuration, and deployment of infrastructure across connected Kubernetes clusters.

### Designs in Meshery
**Overview:**
Meshery employs the concept of "Designs" as a fundamental construct for managing infrastructure. Designs provide a structured way to organize and deploy various components of your infrastructure.

**Importing Infrastructure:**
Meshery facilitates the import of infrastructure in multiple formats, including Kubernetes manifests, Helm Charts, and Docker Compose files. Users can upload directly from the filesystem, provide a URL, or import from a GitHub repository.

### Provisioning Process

When a request is made to provision a design, it undergoes the following stages:

1. **Import of Referenced Designs:**
   Designs support importing remote designs. You can reference a design which will be included at runtime by adding the following annotation `type: $(#use <url-of-remote-pattern>)`. The referenced design will be expanded from the source.

2. **Identification:**
   Meshery relies on components registered at boot time. Only registered models and components can be managed with Meshery. Currently, models from the ArtifactHub repository are supported.

3. **Validation:**
   Components in the design are validated against the schema, ensuring consistency, similar to Kubernetes object validation but tailored for Designs.

4. **Dependency Detection and Resolution:**
   Meshery identifies and resolves dependencies among components. The order of provisioning is critical for successful deployment, and circular dependencies result in the termination of the request.

5. **Provisioning:**
   A Directed Acyclic Graph (DAG) generated in the previous step is processed. Dependent components are processed sequentially, while others are processed in parallel. Meshery intelligently handles the deployment order to ensure successful deployment.

### Auto-Deployment of CRDs and Operators

Meshery automates the deployment of Custom Resource Definitions (CRDs) and operators based on the source from which a particular component was registered. Currently, all components are sourced from the ArtifactHub repository, and the auto-deployment strategy utilizes HelmCharts. (_Support for OCI registries is expected in the near future._)

### Upcoming Features

- **Custom Models and Components:**
  Meshery is actively developing the ability to import custom models and components, expanding the platform's flexibility.

- **OCI Registry Support:**
  Meshery will soon be adding support for OCI registries in addition to ArtifactHub.

---

<img src="{{ site.baseurl }}/assets/img/configuration-management/meshery-applications.png" />
