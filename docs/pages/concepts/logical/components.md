---
layout: default
title: Components
permalink: concepts/logical/components
type: concepts
abstract: "Meshery Components identify and characterize infrastructure under management."
language: en
list: include
redirect_from:
- concepts/components
---
In Meshery, a **Component** is a fundamental building block used to represent and define the infrastructure under management. Each component provides granular and specific support for your infrastructure and applications. Each component represents a distinct capability or feature, like a specific service, protocol, or configuration element. Components can be semantically meaningful (orchestratable) or non-semantically meaningful (an annotation).

## Components as building blocks

Meshery Models and their packaged Components represent a standardized and reusable collection of building blocks used to describe both the desired and actual configuration of infrastructure pre and post-deployment. More than this Components are using to describe non-infrastructure concepts, for example, a component might capture a comment or a visual element that you've placed into a Meshery Design. To aid in distinguishing between types of Components, they are categorized into two groups: Semantic and Non-Semantic components. A component is considered semantically meaningful when it behaves as a direct representation of an aspect of your infrastructure. A component is considered non-semantically meaningful when it behaves as an annotation to aid in the comprehension of your designs.

### Semantic Components

Semantic components represent actual infrastructure and application resources that Meshery can manage. These components map directly to real resources in your environment. For example:

- Kubernetes resources like Deployments, Services, and Pods
- Infrastructure configurations like virtual services and gateways
- Cloud provider resources like load balancers and storage volumes

Semantic components can be deployed, configured, and managed by Meshery's lifecycle management capabilities. They have real operational impact when used in Designs.

### Non-Semantic Components

Non-semantic components are visual and organizational elements that help document and organize your Designs but don't represent actual infrastructure. For example:

- Text boxes and comments to add documentation
- Shapes and containers to visually group related components
- Lines and arrows to show logical relationships
- Labels and tags for organization

While non-semantic components help communicate intent and organization in your Designs, they don't result in any actual infrastructure changes when deployed. Meshery's lifecycle management treats them as purely visual/organizational elements.

### Component Properties

Components have several key properties that define their behavior:

- **Model**: The model that the component belongs to, defining its capabilities
- **Kind**: The specific type of resource the component represents
- **Version**: The API version of the component
- **Spec**: The detailed configuration for the component
- **Status**: The current state of the component (for semantic components)

Understanding the distinction between semantic and non-semantic components helps you effectively use them to both manage real infrastructure and document your designs clearly.

Once registered with Meshery Server (in the [Registry](./registry)), components are available for inclusion in [Designs](./designs) that you create. Components can be created and published by anyone, allowing you to share you custom extensions with the community. This fosters a collaborative ecosystem where you can leverage and contribute to the growing collection of Meshery components.

Components having the same `kind`, `apiVersion` and `model.name` attributes are considered duplicates.

<!-- [![Meshery Components]({{ site.baseurl }}/assets/img/architecture/meshery-components.svg
)]({{ site.baseurl }}/assets/img/architecture/meshery-components.svg) -->
<!-- 
 @leecalcote - This is mumbo jumbo to users and needs to be re-written.

 ## Component Status

Components have a status that is represented as a `Connection` object. Both the administrative and real-time status of a component is a normalized representation of the connection's state. The status is represented as a `Connection` object because the status of a component is a *connection* to the component. For example, the status of a Kubernetes cluster is a direct reflection of a Meshery Server's connection to the cluster.

Normalizing and extracting the status of a component as a direct property of the component and putting it into a connection allows multiple systems to share the same component with different states. For example, different Meshery Servers can access the same Kubernetes cluster, but each Meshery Server has its own connection to the cluster with its own status.

Learn more about [Connections](/concepts/logical/connections).
-->
