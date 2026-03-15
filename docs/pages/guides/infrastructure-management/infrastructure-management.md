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


See [Importing Designs]({{site.baseurl}}/extensions/import-export-designs) for more information.


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
<p>Meshery identifies and resolves dependencies among components using a dynamic mechanism based on each component’s origin (also known as its <code>host</code> or <code>registrant</code>). Provisioning order is critical—circular dependencies will result in termination of the deployment.</p>

<p>Deployment Mechanism by Source:</p>

<ul>
  <li><b>Artifact Hub:</b> Uses Helm Go client for Kubernetes Operator and CRD deployment via <code>ApplyHelmChart()</code>.</li>
  <li><b>Kubernetes YAML:</b> Direct application with Kubernetes Go client (no auto-dependency handling).</li>
</ul>

<p>This behavior is determined by the component’s host type:</p>

<pre>
if connection.Kind == "artifacthub" ➜ Helm Go client used
if connection.Kind == "kubernetes" ➜ Kubernetes Go client used
</pre>

{% include alert.html title="Implementation Detail" type="info" content="This logic is handled in the <code>NewDependencyHandler()</code> function, which selects the deployment method based on <code>connection.Kind</code>." %}

<h4>5. Provisioning</h4>
<p>A Directed Acyclic Graph (DAG) generated in the previous step is processed. Dependent components are processed sequentially, while others are processed in parallel. Meshery intelligently handles the deployment order to ensure successful deployment.</p>
</details>

<h3>Auto-Deployment of CRDs and Operators</h3>

<p>Meshery automates the deployment of Custom Resource Definitions (CRDs) and operators based on the source from which a particular component was registered. By default, Meshery automatically deploys components that are sourced from Artifact Hub (utilizing Helm Charts). Support for OCI registries is expected in the near future.</p>

<details>
<summary><h4>Understanding CRDs and Why Deployment Order Matters</h4></summary>

<p><b>What is a CRD?</b></p>
<p>Custom Resource Definitions (CRDs) extend Kubernetes to support new resource types (e.g., <code>VirtualService</code>, <code>Gateway</code>). A CRD must be installed <b>before</b> any resource of that type is deployed, or you’ll encounter errors like <i>“no matches for kind…”</i>.</p>

<p><b>Why does Meshery only auto-deploy CRDs for ArtifactHub?</b></p>
<p>Helm Charts (used by ArtifactHub) are bundled with CRDs and operators, making it safe for Meshery to auto-deploy them via Helm Go client.</p>

<p>In contrast, Kubernetes YAML files may be minimal or incomplete. Meshery cannot assume what dependencies you intended. Therefore, Meshery avoids auto-deployment to prevent accidental conflicts or errors.</p>

<p><b>User Responsibility (for non-ArtifactHub components):</b></p>
<p>If you use Kubernetes YAML, ensure that all required CRDs and operators are included. Meshery will apply the YAML as-is using the Kubernetes Go client, but will <b>not auto-deploy</b> missing dependencies for you.</p>

</details>

{% include alert.html title="Upcoming Features" type="info" content="<p>Meshery is actively developing the ability to import custom models and components, expanding the platform's flexibility. Support for OCI registries is expected in the near future.</p><p><b>1. Custom Models and Components:</b> Import your own models for more flexibility.</p><p><b>2. OCI Registry Support:</b> In addition to ArtifactHub, OCI registries will be supported soon.</p>" %}