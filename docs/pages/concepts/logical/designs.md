---
layout: default
title: Designs
permalink: concepts/designs
type: concepts
abstract: "Meshery Designs are descriptive, declarative characterizations of how your Kubernetes infrastructure should be configured."
language: en
list: include
---

## Designs

Like a Google Doc, Designs are your primary tool for collaborative authorship of your infrastructure and services. A Design describes All the resources and their properties that you want for a single deployment written inAll the resources and their properties that you want for a single deployment written in YAML based on Meshery's declarative syntax (see [Meshery Schemas repo]((https://github.com/meshery/schemas)). By default, Designs are stored in your user account, but can be manually exported, programmatically snapshotted, or automatically synchronized to any OCI-compatible registry (e.g. Docker Hub, AWS ECR, and so on), or Git-based repositories (coming in v0.8). You can share designs and collaborate in real-time on their creation. Designs can be imported, exported, versioned, forked, merged, snapshotted, published, shared, embedded, templatized, and more.

<!-- ### Using Designs -->

A Design consists of **Components** and/or **Patterns**. A Design is the deployable unit in Meshery. Designs are how the users can describe the desired infrastructure state.
There cannot be two components with the same name within a Design. However, there can be two components with the same name in different Designs.

### Design Relationships and Restrictions

- Designs belong to one an only one Workspace at any given time. Designs may be transferred between Workspaces.
- Designs may be shared with other users or teams.
- Designs may be cloned or merged.
- Designs may be exported or imported.
- Designs may be published or unpublished. Published Designs are available to all users of any Meshery instance through the Meshery Catalog. Unpublished Designs may still be available to other users if that Design is made public.
- Designs may be versioned.
- Designs may be cloned.
- Designs may be snapshotted. Snapshots are immutable. Snapshots may be compared for differences between Design versions.
  - Snapshots may be exported or embedded.
- Designs may be embedded.
- Designs may be converted into reusable Patterns. Creating a Pattern. involves replacing the values of the variables in the design with the values provided by another user.
  - The user who creates a Pattern is called the **pattern owner**.
- Designs may be deployed. Deploying a Design involves incorporating one or more components into your Design, configuring their relationships, and deploying them to one or more Environments. By default, any user of a Workspace can deploy a Design.
- Designs may be deleted.
  - Designs may be archived (depending upon Remote Provider)
  - Designs may be restored (depending upon Remote Provider).
- Designs may be compared.
- Designs may be validated. Validation involves checking the syntax of the Design and ensuring that all the components and patterns referenced in the Design are available.
- Designs may be dry-run deployed.
- Designs may be audited.
- Designs may be tagged.
- Designs may be searched.
- Designs may be filtered.
- Designs may be sorted.
- Designs may be grouped into Workspaces and shared among teams and deployed to Environment(s).
- Designs may be grouped by Technology and/or by Type (e.g. Deployment, Security, Resiliency, Observability, etc.)

### Controlling Access to Designs

When creating a new design by default it's visibility level will be set to **public**. Remove providers have the option of offering additional visibility levels like **private** and **published**.

