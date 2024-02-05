---
layout: enhanced
title: Designs
permalink: concepts/logical/designs
type: concepts
abstract: "Meshery Designs are descriptive, declarative characterizations of how your Kubernetes infrastructure should be configured."
language: en
list: include
redirect_from:
- concepts/designs
---

Like a Google Doc, Designs are your primary tool for collaborative authorship of your infrastructure and services. A Design describes all the resources and their properties that you want for a single deployment written in YAML based on Meshery’s declarative syntax (see [Meshery Schemas repo](<(https://github.com/meshery/schemas)>)). By default, Designs are stored in your user account, but can be manually exported, programmatically snapshotted, or automatically synchronized to any OCI-compatible registry (e.g. Docker Hub, AWS ECR, and so on), or Git-based repositories (coming in v0.8). You can share designs and collaborate in real-time on their creation. Designs can be imported, exported, versioned, forked, merged, snapshotted, published, shared, embedded, templatized, and more.

<!-- ### Using Designs -->

A Design consists of [Components]({{site.baseurl}}/concepts/logical/components) and/or [Patterns](../concepts/logical/patterns). A Design is the deployable unit in Meshery. Designs are how the users can describe the desired infrastructure state.
There cannot be two components with the same name within a Design. However, there can be two components with the same name in different Designs.

### Design Relationships and Restrictions

- Designs belong to only one Workspace at any given time. Designs can be transferred between Workspaces.
- Designs can be shared with other users or teams.
- Designs can be cloned or merged.
- Designs can be exported or imported.
- Designs can be published or unpublished. Published Designs are available to all users of any Meshery instance through the Meshery Catalog. Unpublished Designs can still be available to other users if that Design is made public.
- Designs can be versioned.
- Designs can be cloned.
- Designs can be snapshotted. Snapshots are immutable. Snapshots can be compared for differences between Design versions.
  - Snapshots can be exported or embedded.
- Designs can be embedded.
- Designs can be converted into reusable Patterns. Creating a Pattern involves replacing the values of the variables in the design with the values provided by another user.
  - The user who creates a Pattern is called the **pattern owner**.
- Designs can be deployed. Deploying a Design involves incorporating one or more components into your Design, configuring their relationships, and deploying them to one or more Environments. By default, any user of a Workspace can deploy a Design.
- Designs can be deleted.
  - Designs can be archived (depending upon Remote Provider)
  - Designs can be restored (depending upon Remote Provider).
- Designs can be compared.
- Designs can be validated. Validation involves checking the syntax of the Design and ensuring that all the components and patterns referenced in the Design are available.
- Designs can be dry-run deployed.
- Designs can be audited.
- Designs can be tagged.
- Designs can be searched.
- Designs can be filtered.
- Designs can be sorted.
- Designs can be grouped into Workspaces and shared among teams and deployed to Environment(s).
- Designs can be grouped by Technology and/or by Type (e.g. Deployment, Security, Resiliency, Observability, etc.)

### Controlling Access to Designs

When creating a new design by default it's visibility level will be set to **public**. Remove providers have the option of offering additional visibility levels like **private** and **published**.

