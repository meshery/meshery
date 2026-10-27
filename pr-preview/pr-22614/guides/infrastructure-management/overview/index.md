# Overview

> Meshery provides the ability to manage infrastructure for agility, maintainability, diversity, reliability and isolation, security, and speed.

---

LLMS index: [llms.txt](/pr-preview/pr-22614/llms.txt)

---

## Using Meshery Designs to Manage Your Infrastructure

Meshery is a versatile platform designed to streamline the lifecycle, configuration, and performance management of infrastructure across Kubernetes clusters.

### Designs in Meshery

Meshery employs the concept of [Designs](/pr-preview/pr-22614/concepts/logical/designs/) as a fundamental construct for managing infrastructure. Designs provide a structured way to organize and deploy various components of your infrastructure. To do so, Meshery utilizes a declarative approach to infrastructure management, similar to Kubernetes manifests. Meshery Designs are written in YAML and are validated against a schema.

#### Importing Existing Infrastructure and Applications

Meshery imports infrastructure definitions in these formats, from a file on your machine or from a direct URL to the file:

- Kubernetes Manifest
- Meshery Design
- Helm Charts
- Docker Compose
- Kustomize packages


See [Importing Designs](/pr-preview/pr-22614/guides/configuration-management/import-export-designs/) for more information.


<details>
<summary>
<h3>Meshery Internals Provisioning Process</h3>
</summary>
<p>When a request is made to provision a design, it undergoes the following stages:</p>
<h4>1. Identification</h4>
<p>Meshery relies on components registered in its <a href="/pr-preview/pr-22614/concepts/logical/registry/">Registry</a>. Only registered models and components can be managed with Meshery. Every registered model records the <b>registrant</b> that registered it - a connected Kubernetes cluster, Artifact Hub, GitHub, Meshery itself, or a deployed Meshery Adapter.</p>
<h4>2. Validation</h4>
<p>Components in the design are validated against the schema, ensuring consistency, similar to Kubernetes object validation but tailored for Designs.</p>

<h4>3. Dependency Detection and Resolution</h4>
<p>Deployment order comes from what the design declares. A component's <code>dependsOn</code> entries name other components of the same design, and each entry becomes an edge of the graph provisioning walks. Before anything is deployed, Meshery rejects a design whose dependencies name a component the design does not contain, name a component whose name is shared by more than one component, or form a cycle. The registrant plays no part in this: it decides who fulfills a component and what can be installed on the component's behalf, not what a component depends on.</p>

<p>Installing what a component needs before applying it - its Operators and CRDs - is a separate, registrant-specific behavior. Meshery Server does it unless the deploy request sets <code>skipCRD=true</code>, and what it can install depends on the source of the model. In Meshery UI, the <b>Include Dependencies</b> option in the deploy dialog currently applies only to the dry run; the deployment itself always installs dependencies:</p>

<ul>
  <li><b>Artifact Hub:</b> Uses Helm Go client for Kubernetes Operator and CRD deployment via <code>ApplyHelmChart()</code>.</li>
  <li><b>Kubernetes YAML:</b> Direct application with Kubernetes Go client (no auto-dependency handling).</li>
</ul>

<p>Installing dependencies is best-effort: the component is applied whether or not the install succeeded, and any failure is reported next to the component in the deployment summary.</p>

<h4>4. Provisioning</h4>
<p>A Directed Acyclic Graph (DAG) generated in the previous step is processed. Dependent components are processed sequentially, while others are processed in parallel. Meshery intelligently handles the deployment order to ensure successful deployment.</p>

<p>A component is deployed only once every component it declares a dependency on has been deployed successfully. If one of those fails to apply, the component that depends on it is withheld rather than deployed, and is reported as such along with the dependency that failed; components that declared no dependency on the failed one are unaffected.</p>

<p>Each component is fulfilled by whoever registered its model: Meshery Server applies the component itself when that registrant is only a source of definitions, and delegates over gRPC to a Meshery Adapter when the registrant advertises a network endpoint. One design can use both paths in a single deployment. See <a href="/pr-preview/pr-22614/concepts/architecture/deployment-engine/">Deployment Engine</a> for both paths in full, and for what withholding does when a component's own prerequisites fail, when several clusters are selected, and when a design is undeployed.</p>
</details>

<h3 id="auto-deployment-of-crds-and-operators">Auto-Deployment of CRDs and Operators</h3>

<p>Meshery can deploy the Custom Resource Definitions (CRDs) and operators a component needs, based on the source its model was registered from. For components whose models come from Artifact Hub, Meshery installs the Helm chart that carries those CRDs and operators before it applies the component.</p>

<div class="alert alert-info" role="alert"><div class="h4 alert-heading" role="heading">Registrants decide more than dependencies</div>


The registrant of a component's model also decides <i>who applies the component itself</i> - Meshery Server, or a Meshery Adapter reached over gRPC. See <a href="/pr-preview/pr-22614/concepts/architecture/deployment-engine/">Deployment Engine</a>.
</div>


<details>
<summary><h4>Understanding CRDs and Why Deployment Order Matters</h4></summary>

<p><b>What is a CRD?</b></p>
<p>Custom Resource Definitions (CRDs) extend Kubernetes to support new resource types (e.g., <code>VirtualService</code>, <code>Gateway</code>). A CRD must be installed <b>before</b> any resource of that type is deployed, or you'll encounter errors like <i>“no matches for kind…”</i>.</p>

<p><b>Why does Meshery only auto-deploy CRDs for ArtifactHub?</b></p>
<p>Helm Charts (used by ArtifactHub) are bundled with CRDs and operators, making it safe for Meshery to auto-deploy them via Helm Go client.</p>

<p>In contrast, Kubernetes YAML files may be minimal or incomplete. Meshery cannot assume what dependencies you intended. Therefore, Meshery avoids auto-deployment to prevent accidental conflicts or errors.</p>

<p><b>User Responsibility (for non-ArtifactHub components):</b></p>
<p>If you use Kubernetes YAML, ensure that all required CRDs and operators are included. Meshery will apply the YAML as-is using the Kubernetes Go client, but will <b>not auto-deploy</b> missing dependencies for you.</p>

</details>
