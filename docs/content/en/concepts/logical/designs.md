---
title: Designs
description: Meshery Designs are descriptive, declarative characterizations of how your Kubernetes infrastructure should be configured.
aliases:
- /concepts/designs/
---

Like a Google Doc, Designs are your primary tool for collaborative authorship of your infrastructure and services. A Design describes all the resources and their properties that you want for a single deployment based on Meshery’s declarative syntax (see [Meshery Schemas repo](https://github.com/meshery/schemas)). By default, Designs are stored in your user account, but can be manually exported, programmatically snapshotted, or automatically synchronized to any OCI-compatible registry (e.g. Docker Hub, AWS ECR, and so on), or Git-based repositories (coming in v0.8). You can share designs and collaborate in real-time on their creation. Designs can be imported, exported, versioned, forked, merged, snapshotted, published, shared, embedded, templatized, and more.

As the primary management unit in Meshery, a Design consists of [Components]({{< ref "concepts/logical/components.md" >}}) and [Relationships]({{< ref "concepts/logical/relationships/index.md" >}}). Configurable and Deployable designs are the deployable unit — used to provision and manage live infrastructure. Annotation-only designs are non-deployable, serving as visual architecture diagrams and documentation. Designs are how you can describe your desired infrastructure state.

## Types of Meshery Designs

Meshery supports two primary classifications of designs to accommodate both architectural planning and infrastructure orchestration:

```text
                      ┌─────────────────────────────────────────┐
                      │             Meshery Design              │
                      └────────────────────┬────────────────────┘
                                           │
                    ┌──────────────────────┴──────────────────────┐
                    ▼                                             ▼
      ┌───────────────────────────┐                 ┌───────────────────────────┐
      │   Annotation-Only Design  │                 │     Deployable Design     │
      ├───────────────────────────┤                 ├───────────────────────────┤
      │ • Visual architecture     │                 │ • Functional workloads    │
      │ • Boundaries & shapes     │                 │ • Kubernetes & Helm specs │
      │ • Non-deployable nodes    │                 │ • Schema-validated props  │
      │ • System documentation    │                 │ • Reconciled to clusters  │
      └───────────────────────────┘                 └───────────────────────────┘
```

### 1. Annotation-Only Designs (Visual & Architecture Diagrams)

**Annotation-only designs** are created for architectural diagramming, conceptual mapping, system topologies, and documentation. They allow platform engineers, architects, and DevOps teams to visualize multi-tier systems, cloud resources, and network boundaries without generating deployable workload manifests.

- **Non-Deployable Components**: Components in an annotation-only design are visual shapes, annotations, note blocks, boundaries, and logical components flagged with `isAnnotation: true` in their metadata.
- **Visual Relationships**: Define directional arrows, dependency annotations, or spatial containment without triggering operational provisioning actions in connected clusters.
- **Catalog & Documentation**: Annotation-only designs can be published to the [Meshery Catalog]({{< ref "concepts/architecture/catalog/index.md" >}}), embedded into documentation, or shared across teams as reference architectures.

### 2. Configurable and Deployable Designs

**Configurable and Deployable designs** represent functional infrastructure and application workloads that can be provisioned, configured, validated, and managed across environments.

- **Model-Backed Components**: Every component is backed by an active [Meshery Model]({{< ref "concepts/logical/models/index.md" >}}) (e.g., Kubernetes `Deployment`, `Service`, `StatefulSet`, Istio `VirtualService`, or Cloud Provider CRDs).
- **Schema-Driven Configuration**: Each component contains a declarative specification that is validated against its model's JSON Schema.
- **Orchestration & Lifecycle**: Deployable designs can be dry-run tested, deployed, updated, or undeployed to connected Kubernetes clusters and cloud providers via Meshery Server's deployment engine.

### Comparison: Annotation-Only vs. Deployable Designs

| Attribute | Annotation-Only Designs | Configurable & Deployable Designs |
| :--- | :--- | :--- |
| **Primary Purpose** | Visual diagrams, topology mapping, documentation | Infrastructure provisioning, GitOps, configuration management |
| **Deployability** | Non-deployable (visual canvas elements) | Fully deployable to connected Kubernetes clusters & clouds |
| **Component Backing** | Visual shapes, annotations (`isAnnotation: true`), sticky notes | Registered Meshery Models & JSON Schemas |
| **Relationships** | Visual pointers, groupings, semantic connections | Structural, hierarchical, and network orchestration relationships |
| **Export Formats** | YAML, OCI image, PNG/SVG image snapshot | YAML, OCI image, Helm chart package (`.tgz`) |
| **Catalog Availability** | Yes (Published as Reference Architecture / Diagram) | Yes (Published as Deployable Pattern / Cloud Native Solution) |


### Constraints on Designs

- Designs belong to only one Workspace at any given time. Designs can be transferred between Workspaces.
- Designs can be shared with other users or teams.
- The user who creates a Design is the Design **Owner**. The Design Owner can grant other users *read* or *write* access to the Design and can *delete* the Design.
  
### Features of Designs

- Designs can be **cloned**. Cloning a Design creates a new Design that is a copy of the original Design. The new Design is owned by the user who cloned it.
- Designs can be **merged**. Merging a Design combines two Designs into a single Design. 
  <!-- - Designs can be forked. Forking a Design creates a new Design that is a copy of the original Design. The new Design is owned by the user who forked it. -->
- Designs can be exported as JSON files or OCI images.
- [Designs can be listed in Artifact Hub](https://artifacthub.io/packages/search?kind=24&sort=relevance&page=1) repos.
- Designs can be imported:
  - as Kubernetes Manifests, Docker Compose, Helm Charts, or Meshery Designs.
  - from individual YAML files, remote HTTP location, local filesystem, or OCI images.
- Designs can be exported or embedded.
  - Designs can be embedded in web pages either as a bundle of HTML or as a React component using the [meshery-design-embed](https://www.npmjs.com/package/@meshery/meshery-design-embed) NPM package.
- Designs can be snapshotted. Snapshots are immutable. Snapshots can be compared for differences between Design versions.
- Designs can be published or unpublished. Published Designs are available to all users of any Meshery instance through the [Catalog]({{< ref "concepts/architecture/catalog/index.md" >}}). Unpublished Designs can still be available to other users if that Design is made public.
- Designs are versioned. Each time a Design is saved, a new version is created.
  <!-- - You can revert to any previous version of a Design. -->
- Designs can be deployed. Deploying a Design involves incorporating one or more components into your Design, configuring their relationships, and deploying them to one or more Environments. By default, any user of a Workspace can deploy a Design. Meshery resolves how to fulfill each component separately - see [Deployment Engine]({{< ref "concepts/architecture/deployment-engine/index.md" >}}).
- Designs can be deleted.
  - Designs can be archived and restored (depending upon Remote Provider)
- Designs can be compared.
- Designs can be validated. Validation involves checking the syntax of the Design and ensuring that all the components and patterns referenced in the Design are available.
- Designs can be dry-run deployed. Use dry-runs (or practice run) to ensure that your design works correctly and will not result in undesired changes.
- Designs can be audited. Auditing a Design involves checking the Design for security vulnerabilities, compliance with best practices, and adherence to the organization's policies.
- Designs can be searched.
- Designs can be filtered and viewed as a live-running deployment in your Environments.
- Designs can be sorted by name, date created, date modified, or by the visibility level (public or private) or Catalog published status.
- Designs can be grouped into Workspaces and shared among teams and deployed to Environment(s).
- Designs can be tagged by Technology and/or by Type (e.g. Deployment, Security, Resiliency, Observability, etc.)
<!-- - While there cannot be two components with the same name within a Design, however, there can be two components with the same name in different Designs. -->
- _(Coming in v0.9)_ Designs can be converted into reusable Patterns. Creating a Pattern involves replacing the values of the variables in the design with the values provided by another user. Patterns can be shared with other users or teams.

### Controlling Access to Designs

When creating a new design by default it's visibility level will be set to **public**. [Remote Providers]({{< ref "reference/extensibility/providers/index.md" >}}) have the option of offering additional visibility levels like **private** and **published**.

## Meshery Designs and Models Explained

Designs are the deployable unit in Meshery. [Models]({{< ref "concepts/logical/models/index.md" >}}) are the unit of packaging for Components. Components are described in Designs. Models are not directly deployed. Designs and their Components are.

### Meshery Models

Meshery Models represent the fundamental building blocks of your infrastructure. Think of them as blueprints or templates that define the structure, components, and configurations of your deployments. These models encapsulate everything from network configurations to service definitions, making them essential for consistent and scalable deployments across environments.

### Meshery Designs

On the other hand, Meshery Designs are the practical implementations based on Meshery Models. They represent declarations of your infrastructure deployments, customized according to specific use cases, environments, and requirements. Meshery Designs allow you to create, manage, and deploy complex architectures seamlessly, leveraging the power and flexibility of Meshery Models as their foundation.

Designs are the blueprints for your deployments, while Meshery Models are the internal components that provide the building blocks and knowledge to fulfill those blueprints.

Because every Component in a Design comes from a Model, and every Model records the registrant that registered it, Meshery decides how to fulfill a Design one Component at a time. A Design whose Components come from different registrants is fulfilled along more than one path in a single deployment. See [Deployment Engine]({{< ref "concepts/architecture/deployment-engine/index.md" >}}).

### Using Designs

See the following tutorials on how to use Meshery Designs for collaboratively managing infrastructure.

<details>
  <summary>
      <a href="{{< ref "guides/tutorials/_index.md" >}}" class="text-black">🧑‍🔬 Tutorials</a>
  </summary>
  <ul class="section-title">
    {{< section-pages section="guides/tutorials" >}}
  </ul>
</details>

Try the [Meshery Playground]({{< ref "installation/playground.md" >}}) for a hands-on experience with Meshery Designs.
