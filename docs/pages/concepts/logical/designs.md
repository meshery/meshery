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

### Design Relationships and Restrictions

- Access to Designs may be granted to one or more users.

### Controlling Access to Designs

When creating a new design by default it's visibility level will be set to **public**. Remove providers have the option of offering additional visibility levels like **private** and **published**.

