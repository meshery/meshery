---
layout: default
title: Infrastructure Management
permalink: guides/infrastructure-management/overview
redirect_from: 
- tasks/application-management
- tasks/infrastructure-management
type: guides
category: infrastructure
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

See [Importing Applications]({{site.baseurl}}guides/configuration-management/importing-apps) for more information.

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
<p>Meshery uses <a href="https://github.com/meshery/meshkit/blob/bd00372a4645ff28abe11dae2442f6a627f8c3f9/models/meshmodel/core/v1beta1/host.go">Meshkit</a> to efficiently manage and resolve component dependencies. By utilizing the Helm Go client, it ensures that the provisioning sequence is maintained for successful deployment as a Helm chart. Circular dependencies are promptly detected, resulting in the termination of the request to prevent deployment issues.</p>
<h4>5. Provisioning</h4>
<p>A Directed Acyclic Graph (DAG) generated in the previous step is processed. Dependent components are processed sequentially, while others are processed in parallel. Meshery intelligently handles the deployment order to ensure successful deployment.</p>
</details>

<h3>Auto-Deployment of CRDs and Operators</h3>
Meshery automates the deployment of Custom Resource Definitions (CRDs) and operators based on the source from which a particular component was registered. When a deployment request is made for a Helm Chart without the upgrade flag, Meshery does not check for the existence of the Helm Chart and always tries to install it. In cases where the chart/release already exists, the installation fails (from the Helm side), and Meshery errors out. Currently, by default, Meshery automatically deploys components that are sourced from Artifact Hub (utilizing Helm Charts).

Meshery gives users control over installing dependencies through a simple checkbox option. Users can decide whether to install dependencies during deployment. By default, Meshery tries to install dependencies without checking for an existing release, making deployment quicker. (Support for OCI registries is coming soon.)
{% include alert.html title="Upcoming Features" type="info" content="<p>Meshery is actively developing the ability to import custom models and components, expanding the platform's flexibility. Support for OCI registries is expected in the near future.</p>
<p>1. <b>Custom Models and Components:</b>
  Meshery is actively developing the ability to import custom models and components, expanding the platform's flexibility.</p>
<p>2. <b>OCI Registry Support:</b>
  Meshery will soon be adding support for OCI registries in addition to ArtifactHub.
  </p>" %}

<h3>Frequently Asked Questions</h3>

**Q: How does Meshery determine where to get the required dependency for provisioning components?**

**A:** Meshery determines the source based on the metadata of models and components. This metadata includes the actual location of the components, whether it's a Helm chart, Kubernetes manifests, or any other format. Depending on the type of source (Helm chart or Kubernetes manifest), Meshery deploys the components accordingly.

**Q: How does Meshery know what Operator to provision when provisioning components from a given Model?**

**A:** Meshery determines which Operator to provision based on the `metadata.sourceuri` field, which captures the location or source of the models and components. This process depends on the registrant, such as ArtifactHub or GitHub.
