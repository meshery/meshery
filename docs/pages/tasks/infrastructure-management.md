---
layout: default
title: Infrastructure Management
permalink: tasks/infrastructure-management
redirect_from: tasks/application-management
type: tasks
language: en
list: include
abstract: 'Meshery provides the ability to manage infrastructure for agility, maintainability, diversity, reliability and isolation, security, and speed.'
---

## Using Meshery Designs to Manage Your Infrastructure

Meshery is a versatile platform designed to streamline the lifecycle, configuration, and performance management of infrastructure across Kubernetes clusters.

### Designs in Meshery

Meshery employs the concept of [Designs]({{site.baseurl}}/concepts/logical/designs) as a fundamental construct for managing infrastructure. Designs provide a structured way to organize and deploy various components of your infrastructure. To do so, Meshery utilizes a declarative approach to infrastructure management, similar to Kubernetes manifests. Meshery Designs are written in YAML and are validated against a schema.

#### Importing Existing Infrastructure and Applications

Meshery facilitates the import of infrastructure in multiple formats, including Kubernetes manifests, Helm Charts, and Docker Compose files. Import existing infrastructure definitions by directly from filesystem, via URL, or import directly from a GitHub repository.

- Kubernetes Manifest
- Meshery Design
- Helm Charts
- Docker Compose

<details>
<summary>
<h3>Meshery Internals Provisioning Process</h3>
</summary>
<p>When a request is made to provision a design, it undergoes the following stages:</p>
<h4>1. Import of Referenced Designs</h4>
<p>A Design may reference any number of other Designs, in essence, a Design may import any number of other Designs.  As an editor of a Design, you can make reference to another Design, while following principles of reusing and DRY (Do Not Repeat Yourself). Any referenced Design will subsequently be imported during the provisioning moment. To reference another design, do so by adding the following annotation <pre>type: $(#use \<url-of-remote-pattern\>)</pre> in your Design file. The referenced design will be expanded from the source.</p>
<h4>2. Identification</h4>
<p>Meshery relies on components registered at boot time. Only registered models and components can be managed with Meshery. Currently, models from the ArtifactHub repository are supported.</p>
<h4>3. Validation</h4>
<p>Components in the design are validated against the schema, ensuring consistency, similar to Kubernetes object validation but tailored for Designs.</p>
<h4>4. Dependency Detection and Resolution</h4>
<p>Meshery identifies and resolves dependencies among components. The order of provisioning is critical for successful deployment, and circular dependencies result in the termination of the request.</p>
<h4>5. Provisioning</h4>
<p>A Directed Acyclic Graph (DAG) generated in the previous step is processed. Dependent components are processed sequentially, while others are processed in parallel. Meshery intelligently handles the deployment order to ensure successful deployment.</p>
</details>

<h3>Auto-Deployment of CRDs and Operators</h3>

Meshery automates the deployment of Custom Resource Definitions (CRDs) and operators based on the source from which a particular component was registered. Currently, all components are sourced from the ArtifactHub repository, and the auto-deployment strategy utilizes HelmCharts. (_Support for OCI registries is expected in the near future._)

{% include alert.html title="Upcoming Features" type="info" content="<p>Meshery is actively developing the ability to import custom models and components, expanding the platform's flexibility. Support for OCI registries is expected in the near future.</p>
<p>1. <b>Custom Models and Components:</b>
  Meshery is actively developing the ability to import custom models and components, expanding the platform's flexibility.</p>
<p>2. <b>OCI Registry Support:</b>
  Meshery will soon be adding support for OCI registries in addition to ArtifactHub.
  </p>" %}
