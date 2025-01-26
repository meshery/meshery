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

<!--
### Semantic Components

Semantic relationships are those that are meaningful in the context of the application or infrastructure. For example, a `Service` in Kubernetes is semantically related to a `Deployment` or a `Pod`. These relationships are meaningful and are managed by Meshery.

### Non-Semantic Components

Non-semantic relationships are those that are meaningful to you as a user and your mental representation of your infrastructure and applications, but are not meaningful in terms of how Meshery evaluates the design or manages these relationships and their associated components. Non-sematic relationships are ignored by Meshery's lifecycle management engine. For example, a `Rectangle` shape that encloses other components (has a parent relationship with other child components) is not semantically meaningful to the way in which Meshery manages these resources. While the `Rectangle` shape might have a parent-child relationship with any number of Meshery-managed components, such a relationship does not implicate any management that Meshery might perform; they are not managed by Meshery.
-->

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
