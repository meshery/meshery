---
layout: default
title: Designs
permalink: concepts/logical/designs
type: concepts
abstract: "Meshery Designs are descriptive, declarative characterizations of how your Kubernetes infrastructure should be configured."
language: en
list: include
redirect_from:
- concepts/designs
---

Like a Google Doc, Designs are your primary tool for collaborative authorship of your infrastructure and services. A Design describes all the resources and their properties that you want for a single deployment based on Meshery’s declarative syntax (see [Meshery Schemas repo](https://github.com/meshery/schemas)). By default, Designs are stored in your user account, but can be manually exported, programmatically snapshotted, or automatically synchronized to any OCI-compatible registry (e.g. Docker Hub, AWS ECR, and so on), or Git-based repositories (coming in v0.8). You can share designs and collaborate in real-time on their creation. Designs can be imported, exported, versioned, forked, merged, snapshotted, published, shared, embedded, templatized, and more.

As the deployable unit in Meshery, a Design consists of [Components]({{site.baseurl}}/concepts/logical/components) and [Relationships]({{site.baseurl}}/concepts/logical/relationships). Designs are how you can describe your desired infrastructure state.

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
  - Designs can be embedded in web pages either as a bundle of HTML or as a React component using the [meshery-design-embed](https://www.npmjs.com/package/@layer5/meshery-design-embed) NPM package.
- Designs can be [snapshotted](https://docs.meshery.io/extensions/snapshot). Snapshots are immutable. Snapshots can be compared for differences between Design versions.
- Designs can be published or unpublished. Published Designs are available to all users of any Meshery instance through the [Catalog]({{site.baseurl}}/concepts/catalog). Unpublished Designs can still be available to other users if that Design is made public.
- Designs are versioned. Each time a Design is saved, a new version is created.
  <!-- - You can revert to any previous version of a Design. -->
- Designs can be deployed. Deploying a Design involves incorporating one or more components into your Design, configuring their relationships, and deploying them to one or more Environments. By default, any user of a Workspace can deploy a Design.
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

When creating a new design by default it's visibility level will be set to **public**. Remove providers have the option of offering additional visibility levels like **private** and **published**.

## Meshery Designs and Models Explained

Designs are the deployable unit in Meshery. [Models]({{site.baseurl}}/concepts/logical/models) are the unit of packaging for Components. Components are described in Designs. Models are not directly deployed. Designs and their Components are.

### Meshery Models

Meshery Models represent the fundamental building blocks of your infrastructure. Think of them as blueprints or templates that define the structure, components, and configurations of your deployments. These models encapsulate everything from network configurations to service definitions, making them essential for consistent and scalable deployments across environments.

### Meshery Designs

On the other hand, Meshery Designs are the practical implementations based on Meshery Models. They represent declarations of your infrastructure deployments, customized according to specific use cases, environments, and requirements. Meshery Designs allow you to create, manage, and deploy complex architectures seamlessly, leveraging the power and flexibility of Meshery Models as their foundation.

Designs are the blueprints for your deployments, while Meshery Models are the internal components that provide the building blocks and knowledge to fulfill those blueprints.

### Using Designs

See the following tutorials on how to use Meshery Designs for collaboratively managing infrastructure.

  <details>
  <summary>
    <p style="display:inline">
      <a href="{{ site.baseurl }}/guides/infrastructure-management" class="text-black">🧑‍🔬 Tutorials</a>
    </p>
  </summary>
  <ul class="section-title">
      {% assign tutorials = site.pages | where: "category","tutorials" %}
      {% for item in tutorials %}
      {% if item.type=="guides" and item.category=="tutorials" and item.language=="en" -%}
        <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
        {% if item.abstract != " " %}
          -  {{ item.abstract }}
        {% endif %}
        </li>
        {% endif %}
      {% endfor %}
  </ul>
</details>

Try the [Meshery Playground](/installation/playground) for a hands-on experience with Meshery Designs.
