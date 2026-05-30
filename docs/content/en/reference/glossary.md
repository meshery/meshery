---
title: "Glossary"
description: >
  A reference guide to terminology and concepts used across the Meshery ecosystem.
categories: [Reference]
tags: [concepts, terminology]
weight: 5
---

Meshery introduces several project-specific terms that are used consistently across
its documentation, codebase, and community. This page serves as a single reference
for those terms.

## Adapter

A Meshery Adapter is a gRPC-based plugin that extends Meshery's management capabilities
to a specific cloud native technology (e.g., Istio, Linkerd, Consul). Adapters translate
Meshery's API calls into technology-specific operations and report operational data back
to Meshery Server.

See: [Adapters]({{< ref "/concepts/architecture/adapters" >}})

## Broker

The Meshery Broker (based on NATS) enables asynchronous, event-driven communication
between Meshery Server and its components (such as MeshSync and Adapters). It acts as
a message bus for internal system events and telemetry.

See: [Broker]({{< ref "/concepts/architecture/broker" >}})

## Catalog

The Meshery Catalog is a library of curated, reusable cloud native design templates,
patterns, and configurations. Users can publish their own designs to the catalog or
clone published designs to use as starting points.

See: [Catalog](https://meshery.io/catalog)

## Component

A Component is the atomic unit of a Meshery Model. It represents a specific Kubernetes
resource type (e.g., `Deployment`, `Service`, `VirtualService`). Components carry
schema-validated configuration and are used within Designs.

See: [Components]({{< ref "/concepts/logical/components" >}})

## Connection

A Connection represents an authenticated link between Meshery and an external system
(such as a Kubernetes cluster, Prometheus instance, or cloud provider). Connections
are managed and can be associated with Environments.

See: [Connections]({{< ref "/concepts/logical/connections" >}})

## Credential

A Credential stores the authentication data (tokens, kubeconfigs, API keys) required
to establish a Connection to an external system.

See: [Credentials]({{< ref "/concepts/logical/credentials" >}})

## Design

A Meshery Design is a collection of configured Components and their Relationships,
representing a desired state of cloud native infrastructure. Designs are the primary
artifact for GitOps workflows in Meshery and can be deployed, undeployed, exported,
or shared via the Catalog.

See: [Designs]({{< ref "/concepts/logical/designs" >}})

## Environment

An Environment is a named grouping of Connections and Credentials. Environments make
it easier to manage and share sets of infrastructure resources (e.g., "production",
"staging") across teams and Workspaces.

See: [Environments]({{< ref "/concepts/logical/environments" >}})

## Kanvas

Kanvas is the visual design and management extension for Meshery. It provides a
drag-and-drop canvas interface for creating, editing, and collaborating on Meshery
Designs in real time.

See: [Kanvas]({{< ref "/extensions/kanvas" >}})

## MeshSync

MeshSync is a Kubernetes operator component of Meshery that continuously discovers
and synchronizes the state of Kubernetes resources into Meshery's internal database.
It acts as the "eyes" of Meshery inside a cluster.

See: [MeshSync]({{< ref "/concepts/architecture/meshsync" >}})

## Model

A Meshery Model is a versioned, schema-validated package that describes a cloud native
technology (e.g., Kubernetes, Istio, AWS). A Model contains one or more Components and
defines the Relationships between them. Models are registered in the Meshery Registry.

See: [Models]({{< ref "/concepts/logical/models" >}})

## Operator

The Meshery Operator is a Kubernetes operator that manages the lifecycle of Meshery's
in-cluster components, including MeshSync and the Broker.

See: [Operator]({{< ref "/concepts/architecture/operator" >}})

## Pattern

A Pattern is a reusable, composable unit of cloud native infrastructure configuration.
Patterns are derived from industry best practices and are surfaced through the Meshery
Catalog.

See: [Patterns]({{< ref "/concepts/logical/patterns" >}})

## Policy

A Policy in Meshery is an OPA (Open Policy Agent) Rego rule used to evaluate and
enforce configuration constraints on Designs. Policies can be context-aware and operate
across components and their relationships.

See: [Policies]({{< ref "/concepts/logical/policies" >}})

## Registry

The Meshery Registry is the internal database that stores all registered Models,
Components, and Relationships. The `mesheryctl registry` commands interact with it.

See: [Registry]({{< ref "/concepts/logical/registry" >}})

## Relationship

A Relationship describes how two or more Components interact with each other within
a Design (e.g., a Pod mounting a PersistentVolume, a Service selecting a Deployment).
Relationships drive visual connections in Kanvas and power policy evaluation.

See: [Relationships]({{< ref "/concepts/logical/relationships" >}})

## Workspace

A Workspace is the top-level organizational unit in Meshery for teams. It groups
related Designs, Environments, and Connections together and controls access for
team members.

See: [Workspaces]({{< ref "/concepts/logical/workspaces" >}})

---

{{% alert color="info" title="Missing a term?" %}}
If you encounter a term not listed here, please
[open a documentation issue](https://github.com/meshery/meshery/issues/new?template=documentation.md)
or submit a pull request to add it.
{{% /alert %}}