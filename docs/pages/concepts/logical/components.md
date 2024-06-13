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
In Meshery, a **Component** is a fundamental building block used to represent and define the infrastructure under management. Each component provides granular and specific support for your infrastructure and applications.

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