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

Meshery Models and their packaged Components represent a standardized and reusable collection of building blocks used to describe both the desired and actual configuration of infrastructure pre and post-deployment. More than this, Components are used to describe non-infrastructure concepts, for example, a component might capture a comment or a visual element that you've placed into a Meshery Design. To aid in distinguishing between types of Components, they are categorized into two groups: Semantic and Non-Semantic components. A component is considered semantically meaningful when it behaves as a direct representation of an aspect of your infrastructure. A component is considered non-semantically meaningful when it behaves as an annotation to aid in the comprehension of your designs.

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

Components have several key properties that define their structure and behavior:

#### Core Identity Properties

- **schemaVersion**: The schema version of the component definition (e.g., `components.meshery.io/v1beta1`)
- **version**: The version of the component definition itself (e.g., `v1.0.0`)
- **displayName**: Human-readable name displayed in the UI
- **description**: Detailed description of the component's purpose
- **format**: Data format of the component definition (typically JSON)

#### Component Specification

- **component**: Nested object containing the actual resource specification:
  - **kind**: The specific type of resource the component represents (e.g., `Deployment`, `Service`, `NamespacesTopic`)
  - **version**: The API version of the underlying resource (e.g., `apps/v1`, `servicebus.azure.com/v1api20210101preview`)
  - **schema**: JSON schema defining the structure and validation rules for the component's configuration

#### Model Association

- **model**: Reference to the parent model that packages this component, including:
  - **name**: The model identifier (e.g., `kubernetes`, `azure-service-bus`)
  - **displayName**: Human-readable model name
  - **version**: The model version
  - **category**: Classification of the model (e.g., `Orchestration and Management`, `Observability and Analysis`)
  - **registrant**: Information about who/what registered the model

#### Visual and Behavioral Properties

- **styles**: Visual styling properties for UI representation:
  - **primaryColor**: Primary color for the component icon
  - **secondaryColor**: Secondary accent color
  - **shape**: Visual shape in diagrams (e.g., `circle`, `rectangle`)
  - **svgColor**: Colored SVG icon
  - **svgWhite**: Monochrome SVG icon

- **capabilities**: Array of capability objects defining what operations can be performed with the component:
  - **Performance testing**: Initiate load tests and collect metrics
  - **Configuration management**: Modify workload settings, labels, and annotations
  - **Relationship visualization**: View connections with other components
  - **Style customization**: Change visual appearance and shape
  - **Interactive operations**: Drag-and-drop, compound operations

#### Metadata

- **metadata**: Additional component information:
  - **isAnnotation**: Boolean indicating if this is a non-semantic (visual/organizational) component
  - **isNamespaced**: Boolean indicating whether the component is namespace-scoped (for Kubernetes resources)
  - **published**: Publication status
  - **source_uri**: Original source location of the component definition
  - **genealogy**: Component lineage information
  - **capabilities**: Extended capability definitions

#### Configuration

- **configuration**: The actual configuration values for the component instance (populated when used in a Design)
- **status**: The current state of deployed semantic components (managed through Connections)

Understanding these properties helps you effectively work with components, whether you're creating custom components, managing infrastructure, or designing cloud-native architectures. The distinction between semantic and non-semantic components (indicated by `metadata.isAnnotation`) is particularly important for understanding which components have real operational impact versus those used for documentation and organization.

Once registered with Meshery Server (in the [Registry](./registry)), components are available for inclusion in [Designs](./designs) that you create. Components can be created and published by anyone, allowing you to share your custom extensions with the community. This fosters a collaborative ecosystem where you can leverage and contribute to the growing collection of Meshery components.

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
