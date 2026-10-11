---
title: Overview
aliases: 
- /tasks/application-management
- /tasks/infrastructure-management
categories: [infrastructure]
description: Meshery provides the ability to manage infrastructure for agility, maintainability, diversity, reliability and isolation, security, and speed.
weight: -10
---

## Using Meshery Designs to Manage Your Infrastructure

Meshery is a versatile platform designed to streamline the lifecycle, configuration, and performance management of infrastructure across Kubernetes clusters.

### Designs in Meshery

Meshery employs the concept of [Designs]({{< ref "concepts/logical/designs.md" >}}) as a fundamental construct for managing infrastructure. Designs provide a structured way to organize and deploy various components of your infrastructure. To do so, Meshery utilizes a declarative approach to infrastructure management, similar to Kubernetes manifests. Meshery Designs are written in YAML and are validated against a schema.

#### Importing Existing Infrastructure and Applications

Meshery facilitates the import of infrastructure in multiple formats, including Kubernetes manifests, Helm Charts, and Docker Compose files. Import existing infrastructure definitions by directly from filesystem, via URL, or import directly from a GitHub repository.

- Kubernetes Manifest
- Meshery Design
- Helm Charts
- Docker Compose


See [Importing Designs]({{< ref "guides/configuration-management/import-export-designs.md" >}}) for more information.


<details>
<summary>
<h3>Meshery Internals Provisioning Process</h3>
</summary>
<p>When a request is made to provision a design, it undergoes the following stages:</p>
<h4>1. Import of Referenced Designs</h4>
<p>A Design may reference any number of other Designs, in essence, a Design may import any number of other Designs.  As an editor of a Design, you can make reference to another Design, while following principles of reusing and DRY (Do Not Repeat Yourself). Any referenced Design will subsequently be imported during the provisioning moment. To reference another design, do so by adding the following annotation <pre>type: $(#use \<url-of-remote-pattern\>)</pre> in your Design file. The referenced design will be expanded from the source.</p>
<h4>2. Identification</h4>
<p>Meshery relies on components registered in its <a href="{{< ref "concepts/logical/registry.md" >}}">Registry</a>. Only registered models and components can be managed with Meshery. Every registered model records the <b>registrant</b> that registered it - a connected Kubernetes cluster, Artifact Hub, GitHub, Meshery itself, or a deployed Meshery Adapter.</p>
<h4>3. Validation</h4>
<p>Components in the design are validated against the schema, ensuring consistency, similar to Kubernetes object validation but tailored for Designs.</p>

<h4>4. Dependency Detection and Resolution</h4>
<p>Deployment order comes from what the design declares. A component's <code>dependsOn</code> entries name other components of the same design, and each entry becomes an edge of the graph provisioning walks. Before anything is deployed, Meshery rejects a design whose dependencies name a component the design does not contain, name a component whose name is shared by more than one component, or form a cycle. The registrant plays no part in this: it decides who fulfills a component and what can be installed on the component's behalf, not what a component depends on.</p>

<p>Installing what a component needs before applying it - its Operators and CRDs - is that separate, registrant-specific behavior, and what Meshery is able to install depends on the source of the model:</p>

<ul>
  <li><b>Artifact Hub:</b> Uses Helm Go client for Kubernetes Operator and CRD deployment via <code>ApplyHelmChart()</code>.</li>
  <li><b>Kubernetes YAML:</b> Direct application with Kubernetes Go client (no auto-dependency handling).</li>
</ul>

<p>This behavior is determined by the component’s host type:</p>

<pre>
if connection.Kind == "artifacthub" ➜ Helm Go client used
if connection.Kind == "kubernetes" ➜ Kubernetes Go client used
</pre>

{{% alert color="info" title="Implementation Detail" %}}
This logic is handled in the <code>NewDependencyHandler()</code> function, which selects the deployment method based on <code>connection.Kind</code>.
{{% /alert %}}

<h4>5. Provisioning</h4>
<p>A Directed Acyclic Graph (DAG) generated in the previous step is processed. Dependent components are processed sequentially, while others are processed in parallel. Meshery intelligently handles the deployment order to ensure successful deployment.</p>

<p>A component is deployed only once every component it declares a dependency on has been deployed successfully. If one of those fails to apply, the component that depends on it is withheld rather than deployed, and is reported as such along with the dependency that failed; components that declared no dependency on the failed one are unaffected.</p>

<p>Each component is fulfilled by whoever registered its model: Meshery Server applies the component itself when that registrant is only a source of definitions, and delegates over gRPC to a Meshery Adapter when the registrant advertises a network endpoint. One design can use both paths in a single deployment. See <a href="{{< ref "concepts/architecture/deployment-engine/index.md" >}}">Deployment Engine</a> for both paths in full, and for what withholding does when a component's own prerequisites fail, when several clusters are selected, and when a design is undeployed.</p>
</details>

<h3 id="auto-deployment-of-crds-and-operators">Auto-Deployment of CRDs and Operators</h3>

<p>Meshery automates the deployment of Custom Resource Definitions (CRDs) and operators based on the source from which a particular component was registered. By default, Meshery automatically deploys components that are sourced from Artifact Hub (utilizing Helm Charts). Support for OCI registries is expected in the near future.</p>

{{% alert color="info" title="Registrants decide more than dependencies" %}}
The registrant of a component's model also decides <i>who applies the component itself</i> - Meshery Server, or a Meshery Adapter reached over gRPC. See <a href="{{< ref "concepts/architecture/deployment-engine/index.md" >}}">Deployment Engine</a>.
{{% /alert %}}

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

{{% alert color="info" title="Upcoming Features" %}}
<p>Meshery is actively developing the ability to import custom models and components, expanding the platform's flexibility. Support for OCI registries is expected in the near future.</p><p><b>1. Custom Models and Components:</b> Import your own models for more flexibility.</p><p><b>2. OCI Registry Support:</b> In addition to ArtifactHub, OCI registries will be supported soon.</p>
{{% /alert %}}
